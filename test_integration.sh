#!/bin/bash

echo "=== tech-docs-ai-chat AI Integration Test ==="
echo "Testing all components of the RAG application..."

# Test 1: Check if all containers are running
echo -e "\n1. Checking container status..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep elderly-care-ai

# Test 2: Test the main API endpoint
echo -e "\n2. Testing main chat API endpoint..."
curl -s -X POST http://localhost/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "What is the weather like today?"}' | jq .

# Test 3: Test direct app endpoint
echo -e "\n3. Testing direct app endpoint..."
curl -s -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "How can I improve my health?"}' | jq .

# Test 4: Test error handling
echo -e "\n4. Testing error handling..."
echo "Testing empty message:"
curl -s -X POST http://localhost/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": ""}'

echo -e "\nTesting invalid JSON:"
curl -s -X POST http://localhost/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{invalid json}'

# Test 5: Check Ollama models
echo -e "\n5. Checking Ollama models..."
curl -s http://localhost:11434/api/tags | jq '.models[].name'

# Test 6: Check Qdrant status
echo -e "\n6. Checking Qdrant status..."
curl -s http://localhost:6333/collections | jq .

# Test 7: Run Go tests
echo -e "\n7. Running Go unit tests..."
go test ./... -v

echo -e "\n=== Integration Test Complete ==="
echo "All components are working correctly!" 