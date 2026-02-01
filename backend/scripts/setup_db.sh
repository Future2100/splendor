#!/bin/bash

# Database setup script for Splendor
# This script creates the database and runs all migrations

set -e

echo "🎮 Splendor Database Setup"
echo "=========================="

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default database URL
DB_URL="${DATABASE_URL:-postgres://localhost:5432/splendor?sslmode=disable}"

echo "📦 Creating database..."
createdb splendor 2>/dev/null || echo "Database already exists"

echo "📝 Running migrations..."
echo "  → 001_initial_schema.sql"
psql "$DB_URL" -f migrations/001_initial_schema.sql

echo "  → 002_seed_cards_and_nobles.sql"
psql "$DB_URL" -f migrations/002_seed_cards_and_nobles.sql

echo ""
echo "✅ Database setup complete!"
echo ""
echo "📊 Verification:"
psql "$DB_URL" -c "SELECT tier, COUNT(*) as count FROM development_cards GROUP BY tier ORDER BY tier;"
psql "$DB_URL" -c "SELECT COUNT(*) as nobles FROM nobles;"
echo ""
echo "🚀 Ready to start the server!"
