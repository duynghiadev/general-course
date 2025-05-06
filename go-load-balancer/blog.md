[https://hackernoon.com/this-150-line-go-script-is-actually-a-full-on-load-balancer](https://hackernoon.com/this-150-line-go-script-is-actually-a-full-on-load-balancer)

# Go Load Balancer Project

## Overview

This project implements a simple but robust HTTP load balancer in Go, featuring round-robin load balancing with health checks for backend servers.

## Components

### Load Balancer

The load balancer implements the following features:

- Round-robin load balancing
- Health checks for backend servers
- Reverse proxy functionality
- Automatic failover for dead backends
- Custom error handling

### Backend Server

A simple HTTP server that:

- Responds to HTTP requests
- Provides health check endpoint
- Simulates processing time
- Returns detailed request information

## How to Run

1. Build the backend server:

```bash
cd simple-backend
go build -o simplebackend simplebackend.go
```

2. Start multiple backend instances:

```bash
./simplebackend -port 8082 &
./simplebackend -port 8083 &
./simplebackend -port 8084 &
```

3. Build and run the load balancer:

```bash
cd load-balancer
go build -o loadbalancer loadbalancer.go

./loadbalancer -port 8081
```

## Testing

You can test the load balancer using curl:

```bash
curl http://localhost:8081
```

Multiple requests will be distributed across the backend servers.

## Key Features

### Backend Structure

```go
type Backend struct {
    URL          *url.URL
    Alive        bool
    mux          sync.RWMutex
    ReverseProxy *httputil.ReverseProxy
}
```

### Load Balancer Structure

```go
type LoadBalancer struct {
    backends []*Backend
    current  uint64
}
```

### Health Checking

- Periodic health checks of backend servers
- Configurable check interval
- TCP connection verification
- Automatic backend status updates

### Error Handling

- Graceful handling of backend failures
- Automatic removal of dead backends
- Custom error responses
- Request logging

## Configuration

- Default load balancer port: 8081
- Default health check interval: 1 minute
- Backend ports: 8082, 8083, 8084
- All ports configurable via command-line flags

## Stopping the Services

To stop all running services:

```bash
pkill -f simplebackend
pkill -f loadbalancer
```

## Future Improvements

1. Configuration file support
2. Different load balancing algorithms
3. Dynamic backend registration
4. Metrics and monitoring
5. TLS support
6. Rate limiting
7. Session persistence
