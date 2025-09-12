#!/bin/bash

BASE_URL="http://localhost:8081"

echo "=============================="
echo "1️⃣  Testing /api/healthz"
curl -s "${BASE_URL}/api/healthz"
echo -e "\n"

echo "=============================="
echo "2️⃣  Testing /api/version"
curl -s "${BASE_URL}/api/version"
echo -e "\n"

echo "=============================="
echo "3️⃣  Testing /api/register"
curl -s -X POST "${BASE_URL}/api/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
echo -e "\n"

echo "=============================="
echo "4️⃣  Testing /api/login"
curl -s -X POST "${BASE_URL}/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
echo -e "\n"

echo "=============================="
echo "5️⃣  Testing /api/profile"
curl -s -X GET "${BASE_URL}/api/profile"
echo -e "\n"

echo "=============================="
echo "6️⃣  Testing /api/projects (list)"
curl -s -X GET "${BASE_URL}/api/projects"
echo -e "\n"

echo "=============================="
echo "7️⃣  Testing /api/projects (create)"
curl -s -X POST "${BASE_URL}/api/projects" \
  -H "Content-Type: application/json" \
  -d '{"title":"New Project","description":"Project description"}'
echo -e "\n"

echo "=============================="
echo "✅ All tests completed!"
