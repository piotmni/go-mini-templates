#!/bin/bash

# just some LLM stuff

# API Test Script
# Tests all endpoints for the minimal Go application

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="$BASE_URL/api/v1"
TOKEN="${API_TOKEN:-pass}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "API Test Script"
echo "Base URL: $BASE_URL"
echo "=========================================="
echo ""

# -----------------------------------------------------------------------------
# Health Check
# -----------------------------------------------------------------------------
echo -e "${YELLOW}=== Health Check ===${NC}"
echo "GET /health"
curl -s -X GET "$BASE_URL/health" | jq .
echo ""

# -----------------------------------------------------------------------------
# Categories CRUD
# -----------------------------------------------------------------------------
echo -e "${YELLOW}=== Categories ===${NC}"

# Create Category
echo -e "${GREEN}POST /api/v1/categories${NC} - Create a category"
CATEGORY_RESPONSE=$(curl -s -X POST "$API_URL/categories" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Work"}')
echo "$CATEGORY_RESPONSE" | jq .
CATEGORY_ID=$(echo "$CATEGORY_RESPONSE" | jq -r '.id')
echo "Created category ID: $CATEGORY_ID"
echo ""

# Create another category for testing
echo -e "${GREEN}POST /api/v1/categories${NC} - Create another category"
CATEGORY2_RESPONSE=$(curl -s -X POST "$API_URL/categories" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Personal"}')
echo "$CATEGORY2_RESPONSE" | jq .
CATEGORY2_ID=$(echo "$CATEGORY2_RESPONSE" | jq -r '.id')
echo ""

# Get All Categories
echo -e "${GREEN}GET /api/v1/categories${NC} - List all categories"
curl -s -X GET "$API_URL/categories" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Get Category by ID
echo -e "${GREEN}GET /api/v1/categories/:id${NC} - Get category by ID"
curl -s -X GET "$API_URL/categories/$CATEGORY_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Update Category
echo -e "${GREEN}PUT /api/v1/categories/:id${NC} - Update category"
curl -s -X PUT "$API_URL/categories/$CATEGORY_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Work Updated"}' | jq .
echo ""

# -----------------------------------------------------------------------------
# Notes CRUD
# -----------------------------------------------------------------------------
echo -e "${YELLOW}=== Notes ===${NC}"

# Create Note
echo -e "${GREEN}POST /api/v1/notes${NC} - Create a note"
NOTE_RESPONSE=$(curl -s -X POST "$API_URL/notes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"category_id\": \"$CATEGORY_ID\", \"title\": \"My First Note\", \"content\": \"This is the content of my first note.\"}")
echo "$NOTE_RESPONSE" | jq .
NOTE_ID=$(echo "$NOTE_RESPONSE" | jq -r '.id')
echo "Created note ID: $NOTE_ID"
echo ""

# Create another note
echo -e "${GREEN}POST /api/v1/notes${NC} - Create another note"
NOTE2_RESPONSE=$(curl -s -X POST "$API_URL/notes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"category_id\": \"$CATEGORY2_ID\", \"title\": \"Personal Note\", \"content\": \"Personal stuff here.\"}")
echo "$NOTE2_RESPONSE" | jq .
NOTE2_ID=$(echo "$NOTE2_RESPONSE" | jq -r '.id')
echo ""

# Get All Notes
echo -e "${GREEN}GET /api/v1/notes${NC} - List all notes"
curl -s -X GET "$API_URL/notes" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Get Notes by Category
echo -e "${GREEN}GET /api/v1/notes?category_id=:id${NC} - List notes by category"
curl -s -X GET "$API_URL/notes?category_id=$CATEGORY_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Get Note by ID
echo -e "${GREEN}GET /api/v1/notes/:id${NC} - Get note by ID"
curl -s -X GET "$API_URL/notes/$NOTE_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Update Note
echo -e "${GREEN}PUT /api/v1/notes/:id${NC} - Update note"
curl -s -X PUT "$API_URL/notes/$NOTE_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"category_id\": \"$CATEGORY_ID\", \"title\": \"Updated Note Title\", \"content\": \"Updated content here.\"}" | jq .
echo ""

# -----------------------------------------------------------------------------
# Error Cases
# -----------------------------------------------------------------------------
echo -e "${YELLOW}=== Error Cases ===${NC}"

# Missing Authorization
echo -e "${RED}GET /api/v1/categories${NC} - Missing auth header (should fail)"
curl -s -X GET "$API_URL/categories" | jq .
echo ""

# Invalid Token
echo -e "${RED}GET /api/v1/categories${NC} - Invalid token (should fail)"
curl -s -X GET "$API_URL/categories" \
  -H "Authorization: Bearer invalid-token" | jq .
echo ""

# Get non-existent category
echo -e "${RED}GET /api/v1/categories/:id${NC} - Non-existent category (should fail)"
curl -s -X GET "$API_URL/categories/00000000-0000-0000-0000-000000000000" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Create category with missing name
echo -e "${RED}POST /api/v1/categories${NC} - Missing name (should fail)"
curl -s -X POST "$API_URL/categories" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}' | jq .
echo ""

# Create note with invalid category_id
echo -e "${RED}POST /api/v1/notes${NC} - Invalid category_id (should fail)"
curl -s -X POST "$API_URL/notes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"category_id": "invalid", "title": "Test"}' | jq .
echo ""

# -----------------------------------------------------------------------------
# Cleanup (Delete)
# -----------------------------------------------------------------------------
echo -e "${YELLOW}=== Cleanup ===${NC}"

# Delete Notes
echo -e "${GREEN}DELETE /api/v1/notes/:id${NC} - Delete note"
curl -s -w "HTTP Status: %{http_code}\n" -X DELETE "$API_URL/notes/$NOTE_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""

echo -e "${GREEN}DELETE /api/v1/notes/:id${NC} - Delete second note"
curl -s -w "HTTP Status: %{http_code}\n" -X DELETE "$API_URL/notes/$NOTE2_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""

# Delete Categories
echo -e "${GREEN}DELETE /api/v1/categories/:id${NC} - Delete category"
curl -s -w "HTTP Status: %{http_code}\n" -X DELETE "$API_URL/categories/$CATEGORY_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""

echo -e "${GREEN}DELETE /api/v1/categories/:id${NC} - Delete second category"
curl -s -w "HTTP Status: %{http_code}\n" -X DELETE "$API_URL/categories/$CATEGORY2_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""

echo "=========================================="
echo "Tests completed!"
echo "=========================================="
