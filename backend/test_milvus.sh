#!/bin/bash

# Test script for Milvus client connection
echo "Testing Milvus client connection to 0.0.0.0:19530..."

# Set environment variables for testing
export MILVUS_HOST=0.0.0.0
export MILVUS_PORT=19530

# Change to backend directory
cd "$(dirname "$0")"

# Run the test
go test -v ./internal/milvus/... -run TestMilvusConnection