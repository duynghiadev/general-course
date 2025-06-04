package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Address represents a physical address
type Address struct {
	City string `comparer:"ignore"` // Ignored via struct tag
	Zip  string
}

// Person represents an individual with personal details
type Person struct {
	Name    string `comparer:"custom=case_insensitive"` // Custom comparison rule
	Age     int    `comparer:"ignore"`
	Address Address
	Hobbies []string
}

// ComparisonResult captures differences between two structs
type ComparisonResult struct {
	Path      string
	DiffType  string
	Value1    any
	Value2    any
	FieldType reflect.Type
	Timestamp time.Time
}

// CustomCompareFunc defines a custom comparison function
type CustomCompareFunc func(path string, v1, v2 any) bool

// StructComparerConfig holds configuration for struct comparison
type StructComparerConfig struct {
	MaxGoroutines    int
	IgnoreFields     []string
	CompareTimeout   time.Duration
	MaxDepth         int
	CustomComparers  map[string]CustomCompareFunc
	EnableLogging    bool
	Logger           *log.Logger
	IgnoreTag        string
	CustomCompareTag string
}

// StructComparer manages concurrent struct comparison operations
type StructComparer struct {
	config       StructComparerConfig
	ignoreFields map[string]struct{}
	customFields map[string]CustomCompareFunc
	results      []ComparisonResult
	cache        sync.Map
	semaphore    chan struct{}
	mu           sync.RWMutex
	logger       *log.Logger
}

// NewStructComparer initializes a new StructComparer with configuration
func NewStructComparer(config StructComparerConfig) *StructComparer {
	if config.MaxGoroutines <= 0 {
		config.MaxGoroutines = 10
	}
	if config.CompareTimeout == 0 {
		config.CompareTimeout = 5 * time.Second
	}
	if config.MaxDepth <= 0 {
		config.MaxDepth = 100
	}
	if config.IgnoreTag == "" {
		config.IgnoreTag = "ignore"
	}
	if config.CustomCompareTag == "" {
		config.CustomCompareTag = "custom"
	}
	if config.Logger == nil && config.EnableLogging {
		config.Logger = log.New(os.Stdout, "StructComparer: ", log.LstdFlags)
	}

	ignoreFields := make(map[string]struct{}, len(config.IgnoreFields))
	for _, field := range config.IgnoreFields {
		ignoreFields[field] = struct{}{}
	}

	customFields := make(map[string]CustomCompareFunc)
	for k, v := range config.CustomComparers {
		customFields[k] = v
	}

	return &StructComparer{
		config:       config,
		ignoreFields: ignoreFields,
		customFields: customFields,
		semaphore:    make(chan struct{}, config.MaxGoroutines),
		results:      make([]ComparisonResult, 0),
		logger:       config.Logger,
	}
}

// Reset clears the comparer's results and cache
func (sc *StructComparer) Reset() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.results = make([]ComparisonResult, 0)
	sc.cache = sync.Map{}
	sc.log("Reset comparer state")
}

// Compare compares two structs and returns differences
func (sc *StructComparer) Compare(ctx context.Context, v1, v2 any) ([]ComparisonResult, error) {
	ctx, cancel := context.WithTimeout(ctx, sc.config.CompareTimeout)
	defer cancel()

	if err := sc.compareStructs(ctx, "root", reflect.ValueOf(v1), reflect.ValueOf(v2), 0); err != nil {
		return nil, fmt.Errorf("failed to compare structs: %w", err)
	}

	return sc.GetResults(), nil
}

// compareStructs compares two struct values recursively
func (sc *StructComparer) compareStructs(ctx context.Context, path string, v1, v2 reflect.Value, depth int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error at path %s: %w", path, err)
	}

	if depth > sc.config.MaxDepth {
		return fmt.Errorf("maximum recursion depth exceeded at path %s", path)
	}

	if err := sc.acquireSemaphore(); err != nil {
		return fmt.Errorf("semaphore acquisition failed at path %s: %w", path, err)
	}
	defer sc.releaseSemaphore()

	if cached := sc.checkCache(v1, v2); cached {
		sc.log("Cache hit for path %s", path)
		return nil
	}

	if err := sc.validateValues(v1, v2, path); err != nil {
		return err
	}

	if v1.Type() != v2.Type() {
		sc.addResult(path, "TypeMismatch", v1.Type(), v2.Type(), v1.Type())
		return nil
	}

	switch v1.Kind() {
	case reflect.Struct:
		return sc.compareStructFields(ctx, path, v1, v2, depth+1)
	case reflect.Slice, reflect.Array:
		return sc.compareArrayElements(ctx, path, v1, v2, depth+1)
	default:
		return sc.comparePrimitiveValues(path, v1, v2)
	}
}

// acquireSemaphore controls concurrent goroutine execution
func (sc *StructComparer) acquireSemaphore() error {
	select {
	case sc.semaphore <- struct{}{}:
		return nil
	default:
		return fmt.Errorf("max concurrent operations reached")
	}
}

// releaseSemaphore releases the semaphore
func (sc *StructComparer) releaseSemaphore() {
	<-sc.semaphore
}

// checkCache verifies if comparison result is cached
func (sc *StructComparer) checkCache(v1, v2 reflect.Value) bool {
	cacheKey := fmt.Sprintf("%p-%p", v1.Interface(), v2.Interface())
	if _, exists := sc.cache.LoadOrStore(cacheKey, true); exists {
		return true
	}
	return false
}

// validateValues ensures values are valid for comparison
func (sc *StructComparer) validateValues(v1, v2 reflect.Value, path string) error {
	if !v1.IsValid() || !v2.IsValid() {
		return fmt.Errorf("invalid value at path %s", path)
	}
	return nil
}

// compareStructFields compares all fields of two structs
func (sc *StructComparer) compareStructFields(ctx context.Context, path string, v1, v2 reflect.Value, depth int) error {
	var wg sync.WaitGroup
	errChan := make(chan error, v1.NumField())

	for i := 0; i < v1.NumField(); i++ {
		field := v1.Type().Field(i)
		fieldName := field.Name

		// Check struct tags for ignore
		if tag := field.Tag.Get("comparer"); strings.Contains(tag, sc.config.IgnoreTag) {
			sc.log("Ignoring field %s due to tag", fieldName)
			continue
		}

		// Check config-based ignore fields
		if _, ignore := sc.ignoreFields[fieldName]; ignore {
			sc.log("Ignoring field %s due to config", fieldName)
			continue
		}

		// Check for custom comparison
		if custom, exists := sc.customFields[fieldName]; exists {
			if !custom(path+"."+fieldName, v1.Field(i).Interface(), v2.Field(i).Interface()) {
				sc.addResult(path+"."+fieldName, "CustomMismatch", v1.Field(i).Interface(), v2.Field(i).Interface(), v1.Field(i).Type())
			}
			continue
		}

		// Check struct tags for custom comparison
		if tag := field.Tag.Get("comparer"); strings.Contains(tag, sc.config.CustomCompareTag) {
			if custom, exists := sc.customFields[tag[strings.Index(tag, "=")+1:]]; exists {
				if !custom(path+"."+fieldName, v1.Field(i).Interface(), v2.Field(i).Interface()) {
					sc.addResult(path+"."+fieldName, "CustomMismatch", v1.Field(i).Interface(), v2.Field(i).Interface(), v1.Field(i).Type())
				}
				continue
			}
		}

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			newPath := fmt.Sprintf("%s.%s", path, v1.Type().Field(index).Name)
			if err := sc.compareStructs(ctx, newPath, v1.Field(index), v2.Field(index), depth); err != nil {
				errChan <- fmt.Errorf("field comparison failed at %s: %w", newPath, err)
			}
		}(i)
	}

	return sc.collectGoroutineErrors(&wg, errChan)
}

// compareArrayElements compares elements of arrays or slices
func (sc *StructComparer) compareArrayElements(ctx context.Context, path string, v1, v2 reflect.Value, depth int) error {
	if v1.Len() != v2.Len() {
		sc.addResult(path, "LengthMismatch", v1.Len(), v2.Len(), v1.Type())
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, v1.Len())

	for i := 0; i < v1.Len(); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			newPath := fmt.Sprintf("%s[%d]", path, index)
			if err := sc.compareStructs(ctx, newPath, v1.Index(index), v2.Index(index), depth); err != nil {
				errChan <- fmt.Errorf("array comparison failed at %s: %w", newPath, err)
			}
		}(i)
	}

	return sc.collectGoroutineErrors(&wg, errChan)
}

// comparePrimitiveValues compares primitive type values
func (sc *StructComparer) comparePrimitiveValues(path string, v1, v2 reflect.Value) error {
	if !reflect.DeepEqual(v1.Interface(), v2.Interface()) {
		sc.addResult(path, "ValueMismatch", v1.Interface(), v2.Interface(), v1.Type())
	}
	return nil
}

// collectGoroutineErrors collects errors from concurrent operations
func (sc *StructComparer) collectGoroutineErrors(wg *sync.WaitGroup, errChan chan error) error {
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// addResult safely records a comparison result
func (sc *StructComparer) addResult(path, diffType string, val1, val2 any, fieldType reflect.Type) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.results = append(sc.results, ComparisonResult{
		Path:      path,
		DiffType:  diffType,
		Value1:    val1,
		Value2:    val2,
		FieldType: fieldType,
		Timestamp: time.Now(),
	})
	sc.log("Recorded difference at %s: %s", path, diffType)
}

// GetResults returns a copy of all comparison results
func (sc *StructComparer) GetResults() []ComparisonResult {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	results := make([]ComparisonResult, len(sc.results))
	copy(results, sc.results)
	return results
}

// log records a message if logging is enabled
func (sc *StructComparer) log(format string, args ...any) {
	if sc.config.EnableLogging && sc.logger != nil {
		sc.logger.Printf(format, args...)
	}
}

func main() {
	// Define custom comparison function
	caseInsensitiveCompare := func(path string, v1, v2 any) bool {
		s1, ok1 := v1.(string)
		s2, ok2 := v2.(string)
		if !ok1 || !ok2 {
			return false
		}
		return strings.EqualFold(s1, s2)
	}

	// Initialize test data
	p1 := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			City: "New York",
			Zip:  "10001",
		},
		Hobbies: []string{"reading", "hiking"},
	}

	p2 := Person{
		Name: "ALICE",
		Age:  31,
		Address: Address{
			City: "New York",
			Zip:  "10002",
		},
		Hobbies: []string{"reading", "swimming"},
	}

	// Configure and create comparer
	comparer := NewStructComparer(StructComparerConfig{
		MaxGoroutines:    5,
		IgnoreFields:     []string{"Age"},
		CompareTimeout:   5 * time.Second,
		MaxDepth:         10,
		CustomComparers:  map[string]CustomCompareFunc{"case_insensitive": caseInsensitiveCompare},
		EnableLogging:    true,
		IgnoreTag:        "ignore",
		CustomCompareTag: "custom",
	})

	// Perform comparison
	results, err := comparer.Compare(context.Background(), p1, p2)
	if err != nil {
		fmt.Printf("Comparison error: %v\n", err)
		return
	}

	// Print results
	for _, result := range results {
		fmt.Printf("Difference at %s (%s, Type: %v): %v != %v\n",
			result.Path, result.DiffType, result.FieldType, result.Value1, result.Value2)
	}

	// Reset comparer and reuse
	comparer.Reset()
	fmt.Println("\nAfter reset, results count:", len(comparer.GetResults()))
}
