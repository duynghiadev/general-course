# Performance Benchmark: Node.js vs Go

## Overview

A detailed comparison of performance characteristics between Node.js and Go, based on the article from [ITNEXT](https://itnext.io/performance-benchmark-node-js-vs-go-9dbad158c3b0).

## Key Points

### Test Environment

- Hardware specifications
- Software versions used
- Testing methodology

### Benchmark Results

#### 1. Response Time

- Average response times
- Response time under load
- Latency distribution

#### 2. Resource Usage

- CPU utilization
- Memory consumption
- Thread/goroutine management

#### 3. Concurrency Handling

- Performance under concurrent requests
- Maximum throughput
- Connection handling

## Analysis

### Strengths and Weaknesses

#### Node.js

- Pros:
  - Event-driven architecture
  - Rich ecosystem
  - Easy to get started
- Cons:
  - Single-threaded nature
  - CPU-bound task limitations

#### Go

- Pros:
  - Built-in concurrency
  - Low memory footprint
  - Fast execution
- Cons:
  - Steeper learning curve
  - Less extensive package ecosystem

## Conclusions

- Summary of findings
- Use case recommendations
- Performance optimization tips

## References

1. [Original Article - Performance Benchmark: Node.js vs Go](https://itnext.io/performance-benchmark-node-js-vs-go-9dbad158c3b0)
2. Additional references and resources

## Notes

- Date of benchmark
- Any specific conditions or limitations
- Suggestions for further testing

---

# Go vs Node.js Performance Benchmark

This repo compares performance between Node.js and Go using three benchmark types:

1. Loop Performance
2. Concurrency Test
3. HTTP Benchmark (`/ping` endpoint using wrk)

## Structure

- `node_test/`: Node.js benchmarks using Express
- `go_test/`: Go benchmarks using Fiber
- `wrk_scripts/`: Shell scripts for HTTP benchmarking

## Requirements

- Node.js v22+
- Go 1.24+
- wrk (install via `brew install wrk` or build manually)

## Usage

### Node.js

```bash
cd node_test
npx esbuild index.js --bundle --platform=node --outfile=dist/index.js
node dist/index.js
```

## Go

```bash
cd go_test
go build -o benchmark-test main.go
./benchmark-test
```

## Run HTTP Benchmark

```bash
wrk_scripts/run_benchmark.sh
```
