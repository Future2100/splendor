# Splendor Testing Guide

## End-to-End Testing

### Prerequisites

1. **PostgreSQL running** with database created:
```bash
createdb splendor
```

2. **Database migrations applied**:
```bash
cd backend
./scripts/setup_db.sh
```

3. **Backend server running**:
```bash
cd backend
go run cmd/server/main.go
```

### Running E2E Tests

```bash
cd test
./e2e_test.sh
```

### What the E2E Test Covers

1. ✓ Health check endpoint
2. ✓ User registration (2 users)
3. ✓ User login
4. ✓ Get current user (authenticated)
5. ✓ Create game
6. ✓ List games
7. ✓ Join game (second user)
8. ✓ Get game details
9. ✓ Start game
10. ✓ Get game state (with cards, gems, nobles)
11. ✓ Get leaderboard

### Manual Testing

#### Complete Game Flow

1. **Start backend**:
```bash
cd backend
go run cmd/server/main.go
```

2. **Start frontend** (when Node.js available):
```bash
cd frontend
npm install
npm run dev
```

3. **Test in browser**:
   - Open http://localhost:5173
   - Register two users in different browsers/tabs
   - Create a game with user 1
   - Join with user 2 using room code
   - Start game
   - Play through turns:
     - Take gems
     - Purchase cards
     - Reserve cards
     - Get nobles
     - Reach 15 points

### API Testing with curl

#### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

#### Create Game
```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"num_players":2}'
```

#### Join Game
```bash
curl -X POST http://localhost:8080/api/v1/games/join \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"room_code":"ROOM_CODE"}'
```

#### Start Game
```bash
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/start \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Get Game State
```bash
curl http://localhost:8080/api/v1/games/GAME_ID/state
```

#### Take Gems
```bash
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/take-gems \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"gems":{"diamond":1,"sapphire":1,"emerald":1}}'
```

#### Purchase Card
```bash
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/purchase-card \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"card_id":1,"from_reserve":false}'
```

#### Reserve Card
```bash
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/reserve-card \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"card_id":5,"tier":1}'
```

### WebSocket Testing

Connect to WebSocket:
```
ws://localhost:8080/api/v1/ws/games/GAME_ID?token=YOUR_JWT_TOKEN
```

Send message:
```json
{
  "type": "move",
  "payload": {
    "action": "take_gems",
    "data": {...}
  }
}
```

### Expected Results

All E2E tests should pass with:
- ✓ All HTTP endpoints returning correct status codes
- ✓ Valid JSON responses
- ✓ Proper authentication and authorization
- ✓ Game state correctly initialized
- ✓ WebSocket connections established

### Troubleshooting

**Database connection issues**:
```bash
# Check PostgreSQL is running
psql -l

# Recreate database
dropdb splendor
createdb splendor
cd backend && ./scripts/setup_db.sh
```

**Port already in use**:
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

**Token expired**:
- Tokens expire after 15 minutes
- Re-login to get new tokens
- Or use the refresh endpoint

### Performance Testing

Test concurrent games:
```bash
# Run multiple test scripts simultaneously
for i in {1..10}; do
    ./e2e_test.sh &
done
wait
```

### Load Testing with Apache Bench

```bash
# Test API endpoint performance
ab -n 1000 -c 10 http://localhost:8080/health

# Test with authentication
ab -n 100 -c 5 -H "Authorization: Bearer TOKEN" \
  http://localhost:8080/api/v1/games
```
