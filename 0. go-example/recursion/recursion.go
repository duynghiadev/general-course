package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Address struct {
	City string
	Zip  string
}

type Person struct {
	Name    string
	Age     int
	Address Address
	Hobbies []string
}

// ComparisonResult holds the details of struct comparison
type ComparisonResult struct {
	Path      string
	DiffType  string
	Value1    interface{}
	Value2    interface{}
	Timestamp time.Time
}

// StructComparer handles struct comparison with advanced features
type StructComparer struct {
	mu             sync.RWMutex
	results        []ComparisonResult
	cache          sync.Map
	maxGoroutines  int
	semaphore      chan struct{}
	ignoreFields   map[string]bool
	compareTimeout time.Duration
}

// NewStructComparer creates a new instance of StructComparer with options
func NewStructComparer(opts ...func(*StructComparer)) *StructComparer {
	sc := &StructComparer{
		maxGoroutines:  10,
		semaphore:      make(chan struct{}, 10),
		ignoreFields:   make(map[string]bool),
		compareTimeout: 5 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(sc)
	}

	return sc
}

// WithMaxGoroutines sets the maximum number of concurrent comparisons
func WithMaxGoroutines(max int) func(*StructComparer) {
	return func(sc *StructComparer) {
		sc.maxGoroutines = max
		sc.semaphore = make(chan struct{}, max)
	}
}

// WithIgnoreFields specifies fields to ignore during comparison
func WithIgnoreFields(fields []string) func(*StructComparer) {
	return func(sc *StructComparer) {
		for _, field := range fields {
			sc.ignoreFields[field] = true
		}
	}
}

// CompareStructs compares two structs with context and advanced features
func (sc *StructComparer) CompareStructs(ctx context.Context, path string, v1, v2 reflect.Value) error {
	// Context timeout handling
	select {
	case <-ctx.Done():
		return fmt.Errorf("comparison timed out at path %s: %w", path, ctx.Err())
	default:
	}

	// Acquire semaphore for concurrent control
	sc.semaphore <- struct{}{}
	defer func() { <-sc.semaphore }()

	// Cache check
	cacheKey := fmt.Sprintf("%p-%p", v1.Interface(), v2.Interface())
	if _, exists := sc.cache.Load(cacheKey); exists {
		return nil
	}
	sc.cache.Store(cacheKey, true)

	// Validation
	if err := sc.validateValues(v1, v2, path); err != nil {
		return err
	}

	// Type comparison
	if v1.Type() != v2.Type() {
		sc.addResult(path, "TypeMismatch", v1.Type(), v2.Type())
		return nil
	}

	// Handle different kinds
	switch v1.Kind() {
	case reflect.Struct:
		return sc.compareStructs(ctx, path, v1, v2)
	case reflect.Slice, reflect.Array:
		return sc.compareArrays(ctx, path, v1, v2)
	default:
		return sc.comparePrimitives(path, v1, v2)
	}
}

// validateValues checks if values are valid for comparison
func (sc *StructComparer) validateValues(v1, v2 reflect.Value, path string) error {
	if !v1.IsValid() || !v2.IsValid() {
		return fmt.Errorf("invalid value at path %s", path)
	}
	return nil
}

// compareStructs handles struct comparison with goroutines
func (sc *StructComparer) compareStructs(ctx context.Context, path string, v1, v2 reflect.Value) error {
	var wg sync.WaitGroup
	errChan := make(chan error, v1.NumField())

	for i := 0; i < v1.NumField(); i++ {
		fieldName := v1.Type().Field(i).Name

		// Skip ignored fields
		if sc.ignoreFields[fieldName] {
			continue
		}

		wg.Add(1)
		go func(fieldIndex int) {
			defer wg.Done()

			fieldName := v1.Type().Field(fieldIndex).Name
			newPath := path + "." + fieldName

			if err := sc.CompareStructs(ctx, newPath, v1.Field(fieldIndex), v2.Field(fieldIndex)); err != nil {
				errChan <- fmt.Errorf("error comparing field %s: %w", newPath, err)
			}
		}(i)
	}

	// Wait for all comparisons to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// compareArrays handles slice/array comparison
func (sc *StructComparer) compareArrays(ctx context.Context, path string, v1, v2 reflect.Value) error {
	if v1.Len() != v2.Len() {
		sc.addResult(path, "LengthMismatch", v1.Len(), v2.Len())
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, v1.Len())

	for i := 0; i < v1.Len(); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			newPath := fmt.Sprintf("%s[%d]", path, index)
			if err := sc.CompareStructs(ctx, newPath, v1.Index(index), v2.Index(index)); err != nil {
				errChan <- err
			}
		}(i)
	}

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

// comparePrimitives handles primitive type comparison
func (sc *StructComparer) comparePrimitives(path string, v1, v2 reflect.Value) error {
	if !reflect.DeepEqual(v1.Interface(), v2.Interface()) {
		sc.addResult(path, "ValueMismatch", v1.Interface(), v2.Interface())
	}
	return nil
}

// addResult thread-safely adds a comparison result
func (sc *StructComparer) addResult(path, diffType string, val1, val2 interface{}) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.results = append(sc.results, ComparisonResult{
		Path:      path,
		DiffType:  diffType,
		Value1:    val1,
		Value2:    val2,
		Timestamp: time.Now(),
	})
}

// GetResults returns all comparison results
func (sc *StructComparer) GetResults() []ComparisonResult {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	results := make([]ComparisonResult, len(sc.results))
	copy(results, sc.results)
	return results
}

func main() {
	// Create test structures
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
		Name: "Alice",
		Age:  31,
		Address: Address{
			City: "New York",
			Zip:  "10002",
		},
		Hobbies: []string{"reading", "swimming"},
	}

	// Create comparer with options
	comparer := NewStructComparer(
		WithMaxGoroutines(5),
		WithIgnoreFields([]string{"Age"}),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Perform comparison
	if err := comparer.CompareStructs(ctx, "Person", reflect.ValueOf(p1), reflect.ValueOf(p2)); err != nil {
		fmt.Printf("Comparison error: %v\n", err)
		return
	}

	// Print results
	for _, result := range comparer.GetResults() {
		fmt.Printf("Difference at %s (%s): %v != %v\n",
			result.Path, result.DiffType, result.Value1, result.Value2)
	}
}
