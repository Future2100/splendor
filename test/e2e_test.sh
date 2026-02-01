#!/bin/bash

# Splendor End-to-End Test Script
# This script tests the complete game flow

set -e

API_URL="http://localhost:8080/api/v1"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "======================================"
echo "  Splendor E2E Test"
echo "======================================"
echo ""

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

test_endpoint() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local expected_status=$5
    local token=$6

    echo -n "Testing $name... "

    if [ -z "$token" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$url" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$url" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$data")
    fi

    status=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$status" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((TESTS_PASSED++))
        echo "$body"
    else
        echo -e "${RED}✗ FAILED${NC} (Expected $expected_status, got $status)"
        echo "$body"
        ((TESTS_FAILED++))
    fi

    echo ""
}

# Health Check
echo "=== 1. Health Check ==="
test_endpoint "Health check" "GET" "/health" "" 200

# User Registration
echo "=== 2. User Registration ==="
USER1_DATA='{"username":"alice","email":"alice@example.com","password":"password123"}'
USER2_DATA='{"username":"bob","email":"bob@example.com","password":"password123"}'

USER1_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "$USER1_DATA")
USER1_TOKEN=$(echo $USER1_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "User 1 Token: ${USER1_TOKEN:0:20}..."

USER2_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "$USER2_DATA")
USER2_TOKEN=$(echo $USER2_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "User 2 Token: ${USER2_TOKEN:0:20}..."
echo ""

# Login Test
echo "=== 3. Login Test ==="
LOGIN_DATA='{"email":"alice@example.com","password":"password123"}'
test_endpoint "User login" "POST" "/auth/login" "$LOGIN_DATA" 200

# Get Current User
echo "=== 4. Get Current User ==="
test_endpoint "Get current user" "GET" "/auth/me" "" 200 "$USER1_TOKEN"

# Create Game
echo "=== 5. Create Game ==="
CREATE_GAME_DATA='{"num_players":2}'
GAME_RESPONSE=$(curl -s -X POST "$API_URL/games" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $USER1_TOKEN" \
    -d "$CREATE_GAME_DATA")
GAME_ID=$(echo $GAME_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
ROOM_CODE=$(echo $GAME_RESPONSE | grep -o '"room_code":"[^"]*' | cut -d'"' -f4)
echo "Game ID: $GAME_ID"
echo "Room Code: $ROOM_CODE"
echo ""

# List Games
echo "=== 6. List Games ==="
test_endpoint "List games" "GET" "/games" "" 200

# Join Game
echo "=== 7. Join Game ==="
JOIN_DATA="{\"room_code\":\"$ROOM_CODE\"}"
test_endpoint "User 2 joins game" "POST" "/games/join" "$JOIN_DATA" 200 "$USER2_TOKEN"

# Get Game
echo "=== 8. Get Game Details ==="
test_endpoint "Get game" "GET" "/games/$GAME_ID" "" 200

# Start Game
echo "=== 9. Start Game ==="
test_endpoint "Start game" "POST" "/games/$GAME_ID/start" "" 200 "$USER1_TOKEN"

# Get Game State
echo "=== 10. Get Game State ==="
test_endpoint "Get game state" "GET" "/games/$GAME_ID/state" "" 200

# Statistics
echo "=== 11. Statistics ==="
test_endpoint "Get leaderboard" "GET" "/stats/leaderboard" "" 200

# Summary
echo "======================================"
echo "  Test Summary"
echo "======================================"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"
echo "======================================"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
