# Splendor Frontend

A beautiful React + TypeScript frontend for the Splendor board game.

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Framer Motion** - Animations
- **React Router** - Routing
- **Axios** - HTTP client

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

The app will be available at http://localhost:5173

### Build

```bash
npm run build
```

## Project Structure

```
src/
├── components/       # React components
│   ├── common/      # Shared components
│   ├── auth/        # Authentication components
│   ├── lobby/       # Game lobby components
│   ├── game/        # Game board components
│   └── stats/       # Statistics components
├── hooks/           # Custom React hooks
├── context/         # React context providers
├── services/        # API services
├── types/           # TypeScript type definitions
├── utils/           # Utility functions
├── pages/           # Page components
└── styles/          # Global styles
```

## Features

### Phase 1 (Complete)
- ✅ Project setup
- ✅ Routing
- ✅ Basic UI with Tailwind CSS
- ✅ Page templates

### Upcoming Phases
- Phase 2: Authentication system
- Phase 3: Game lobby
- Phase 4: WebSocket integration
- Phase 5-9: Game mechanics
- Phase 10: Statistics and leaderboard
- Phase 11: UI polish and animations

## Environment Variables

Copy `.env.example` to `.env` and configure:

```
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```
