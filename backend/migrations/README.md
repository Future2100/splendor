# Database Migrations

## Running Migrations

### Manual Migration

```bash
# Create database
createdb splendor

# Run migrations in order
psql splendor < migrations/001_initial_schema.sql
psql splendor < migrations/002_seed_cards_and_nobles.sql
```

### Using Environment Variable

```bash
# Set DATABASE_URL in .env file
export DATABASE_URL="postgres://username:password@localhost:5432/splendor?sslmode=disable"

# Run migrations
psql $DATABASE_URL -f migrations/001_initial_schema.sql
psql $DATABASE_URL -f migrations/002_seed_cards_and_nobles.sql
```

## Migration Files

### 001_initial_schema.sql
Creates all database tables:
- `users` - User accounts
- `games` - Game instances
- `game_players` - Players in each game
- `development_cards` - Reference data for cards (90 cards)
- `nobles` - Reference data for nobles (10 nobles)
- `game_state` - Current board state for each game
- `player_state` - Each player's resources and cards
- `game_moves` - Move history
- `game_statistics` - Aggregated player statistics

### 002_seed_cards_and_nobles.sql
Populates reference data:
- **40 Tier 1 cards** (8 of each gem type, 0-1 points)
- **30 Tier 2 cards** (6 of each gem type, 1-3 points)
- **20 Tier 3 cards** (4 of each gem type, 3-5 points)
- **10 Nobles** (historical figures, 3 points each)

## Verify Installation

```sql
-- Check card distribution
SELECT tier, COUNT(*) as count
FROM development_cards
GROUP BY tier
ORDER BY tier;
-- Should show: Tier 1: 40, Tier 2: 30, Tier 3: 20

-- Check nobles
SELECT COUNT(*) FROM nobles;
-- Should show: 10

-- Check all tables
\dt
```

## Rollback

To completely reset the database:

```bash
dropdb splendor
createdb splendor
# Then re-run migrations
```

## Future Migrations

When adding new migrations:
1. Name them sequentially: `003_description.sql`, `004_description.sql`
2. Include both UP and DOWN migrations if possible
3. Test thoroughly before committing
4. Document changes in this README
