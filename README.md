# Splendor - Board Game Web Application

A complete implementation of the Splendor board game with real-time multiplayer functionality.

## Features

- **Full Authentication System** - User registration, login, JWT tokens
- **Real-time Multiplayer** - WebSocket-based game synchronization
- **Complete Game Rules** - All Splendor mechanics implemented
- **Beautiful UI** - Animations, transitions, responsive design
- **Statistics & Leaderboard** - Track your progress and compete

## Tech Stack

### Backend
- **Go** - Backend language
- **Gin** - HTTP framework
- **PostgreSQL** - Database
- **WebSocket** - Real-time communication
- **JWT** - Authentication

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Framer Motion** - Animations

## Project Structure

```
splendor/
├── backend/          # Go backend
│   ├── cmd/         # Application entrypoint
│   ├── internal/    # Private application code
│   ├── pkg/         # Public libraries
│   └── migrations/  # Database migrations
└── frontend/         # React frontend
    └── src/         # Source code
```

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### Backend Setup

```bash
cd backend
cp .env.example .env
# Edit .env with your database credentials
go mod download
go run cmd/server/main.go
```

### Frontend Setup

```bash
cd frontend
npm install
cp .env.example .env
npm run dev
```

### Database Setup

```bash
# Create database
createdb splendor

# Run migrations (coming in Phase 1)
psql splendor < migrations/001_initial_schema.sql
psql splendor < migrations/002_seed_cards_and_nobles.sql
```

## Development Status

### ✅ Phase 1: Project Initialization (Complete)
- Backend Go project structure
- Frontend React project structure
- Configuration files
- Basic routing and pages

### 🚧 Phase 2: Authentication (Next)
- User registration and login
- JWT token management
- Password hashing
- Protected routes

### 📋 Phase 3-13: Remaining Features
See the full implementation plan in the project documentation.

## Game Rules

Splendor is a game of chip-collecting and card development. Players are merchants of the Renaissance trying to buy gem mines, means of transportation, and shops—all in order to acquire the most prestige points.

- **Players**: 2-4
- **Goal**: First to 15 points triggers the end game
- **Actions**: Take gems, Reserve cards, Purchase cards
- **Nobles**: Visit players who meet requirements

## API Endpoints

```
POST   /api/v1/auth/register       # Register user
POST   /api/v1/auth/login          # Login
GET    /api/v1/games               # List games
POST   /api/v1/games               # Create game
WS     /api/v1/ws/games/:id        # WebSocket connection
GET    /api/v1/stats/leaderboard   # Leaderboard
```

## Contributing

This is a learning project. Feel free to fork and modify!

## License

MIT
