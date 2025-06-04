#!/bin/bash

echo "Benchmarking Node.js..."
wrk -t12 -c400 -d30s http://localhost:3000/ping

echo "=========================="

echo "Benchmarking Go..."
wrk -t12 -c400 -d30s http://localhost:3001/ping
