# Quick Start Guide

This guide will help you get the Splendor game system up and running in less than 5 minutes.

## Prerequisites

- Docker and Docker Compose installed
- OR: Go 1.21+, Node.js 18+, and PostgreSQL 15

## Method 1: Docker (Recommended - Easiest)

### Step 1: Start All Services

```bash
cd /Users/shanks/go/src/splendor
docker-compose up --build -d
```

This command will:
- Build the backend Go application
- Build the frontend React application
- Start PostgreSQL database with seed data
- Configure networking between services

### Step 2: Wait for Services to Start

```bash
# Check service status
docker-compose ps

# Watch logs (Ctrl+C to exit)
docker-compose logs -f

# Wait for "Server starting on port 8080" message
```

Services should be ready in 30-60 seconds.

### Step 3: Access the Application

Open your browser and navigate to:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/health

### Step 4: Test the Application

1. **Register a User**
   - Click "Register" on the homepage
   - Enter username, email, and password
   - Submit the form

2. **Create a Game**
   - After login, you'll see the game lobby
   - Click "Create Game"
   - Select number of players (2-4)
   - Note the room code

3. **Join Game (Optional - Open in Incognito Window)**
   - Open another browser window (incognito mode)
   - Register a different user
   - Enter the room code
   - Click "Join Game"

4. **Start Playing**
   - The game creator can start the game once players have joined
   - Take turns performing actions:
     - Take gems (3 different or 2 same)
     - Purchase cards
     - Reserve cards
   - First player to reach 15 points triggers game end

### Step 5: Stop Services

```bash
docker-compose down

# To remove volumes and reset database:
docker-compose down -v
```

## Method 2: Local Development

### Step 1: Start Database

```bash
docker-compose up postgres -d
```

### Step 2: Run Database Migrations

```bash
cd backend
export DATABASE_URL="postgres://splendor:splendor_password@localhost:5432/splendor?sslmode=disable"

# Apply migrations
psql $DATABASE_URL -f migrations/001_initial_schema.sql
psql $DATABASE_URL -f migrations/002_seed_cards_and_nobles.sql
```

### Step 3: Start Backend

```bash
cd backend

# Set environment variables
export DATABASE_URL="postgres://splendor:splendor_password@localhost:5432/splendor?sslmode=disable"
export JWT_SECRET="your-secret-key-for-development"
export PORT="8080"
export ENVIRONMENT="development"
export FRONTEND_URL="http://localhost:5173"

# Run the server (already compiled)
./splendor-server

# Or rebuild and run:
go run cmd/server/main.go
```

### Step 4: Start Frontend

```bash
# In a new terminal
cd frontend

# Create .env file
cat > .env << EOF
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
EOF

# Start dev server
npm run dev
```

The frontend will be available at http://localhost:5173

## Method 3: Run E2E Tests

### Automated API Testing

```bash
# Make sure services are running
docker-compose up -d

# Wait for services to be ready
sleep 10

# Run the test suite
cd test
chmod +x e2e_test.sh
./e2e_test.sh
```

The test script will:
1. Check health endpoint
2. Register two users
3. Login both users
4. Create a game
5. Join the game
6. Start the game
7. Fetch game state
8. Retrieve statistics
9. Check leaderboard

### Expected Output

```
✅ Health check passed
✅ User 1 registration successful
✅ User 2 registration successful
✅ User 1 login successful
✅ User 2 login successful
✅ Game creation successful - Room Code: ABCD1234
✅ User 2 joined game
✅ Game started successfully
✅ Game state retrieved
✅ User statistics retrieved
✅ Leaderboard retrieved

🎉 All E2E tests passed!
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :5432  # PostgreSQL

# Kill the process or change ports in docker-compose.yml
```

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### Frontend Not Loading

```bash
# Check backend is accessible
curl http://localhost:8080/health

# Check nginx logs
docker-compose logs frontend

# Rebuild frontend
docker-compose up --build frontend
```

### WebSocket Connection Failed

```bash
# Check WebSocket endpoint
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  http://localhost:8080/api/v1/ws/games/1?token=YOUR_JWT_TOKEN

# Verify nginx WebSocket proxy config
docker-compose exec frontend cat /etc/nginx/conf.d/default.conf
```

### Reset Everything

```bash
# Stop and remove all containers, networks, and volumes
docker-compose down -v

# Remove built images
docker-compose down --rmi all

# Start fresh
docker-compose up --build -d
```

## API Endpoints Quick Reference

### Authentication
- POST `/api/v1/auth/register` - Register new user
- POST `/api/v1/auth/login` - Login user
- POST `/api/v1/auth/refresh` - Refresh JWT token
- GET `/api/v1/auth/me` - Get current user

### Games
- GET `/api/v1/games` - List all games
- POST `/api/v1/games` - Create new game
- GET `/api/v1/games/:id` - Get game details
- POST `/api/v1/games/join` - Join game by room code
- POST `/api/v1/games/:id/leave` - Leave game
- POST `/api/v1/games/:id/start` - Start game
- GET `/api/v1/games/:id/state` - Get full game state

### Gameplay
- POST `/api/v1/games/:id/take-gems` - Take gems action
- POST `/api/v1/games/:id/purchase-card` - Purchase card
- POST `/api/v1/games/:id/reserve-card` - Reserve card

### Statistics
- GET `/api/v1/stats/users/:id` - User statistics
- GET `/api/v1/stats/leaderboard` - Global leaderboard

### WebSocket
- WS `/api/v1/ws/games/:id?token=<jwt>` - Real-time game updates

## Default Credentials (Development)

For quick testing, you can use these credentials after running migrations:

```
Database:
  Host: localhost
  Port: 5432
  User: splendor
  Password: splendor_password
  Database: splendor

JWT Secret: your-production-secret-key-change-this (change in production!)
```

## Production Deployment Notes

Before deploying to production:

1. **Change JWT Secret**: Generate a strong random secret
   ```bash
   openssl rand -base64 32
   ```

2. **Use Strong Database Password**: Update in docker-compose.yml

3. **Enable HTTPS**: Use a reverse proxy (nginx, Caddy, Traefik)

4. **Set Proper CORS Origins**: Update `FRONTEND_URL` environment variable

5. **Enable Rate Limiting**: Add middleware for API protection

6. **Monitor Logs**: Set up logging aggregation (ELK, Datadog, etc.)

7. **Backup Database**: Schedule regular PostgreSQL backups

## Next Steps

- Read `FINAL_SUMMARY.md` for complete project documentation
- Check `BUILD_VERIFICATION.md` for build details
- Review `test/README.md` for testing guide
- See `README.md` for project overview

## Support

If you encounter any issues:
1. Check the logs: `docker-compose logs -f`
2. Verify environment variables
3. Ensure ports are not in use
4. Try resetting: `docker-compose down -v && docker-compose up --build -d`

Happy gaming! 🎮✨
