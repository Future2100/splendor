-- Splendor Database Schema
-- Phase 1: Initial Schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- Games table
CREATE TABLE IF NOT EXISTS games (
    id BIGSERIAL PRIMARY KEY,
    room_code VARCHAR(10) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting', -- waiting, in_progress, completed
    current_turn_player_id BIGINT,
    turn_number INT NOT NULL DEFAULT 0,
    winner_id BIGINT,
    created_by BIGINT NOT NULL REFERENCES users(id),
    num_players INT NOT NULL DEFAULT 4,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    CONSTRAINT chk_num_players CHECK (num_players BETWEEN 2 AND 4),
    CONSTRAINT chk_status CHECK (status IN ('waiting', 'in_progress', 'completed'))
);

CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_games_room_code ON games(room_code);
CREATE INDEX idx_games_created_by ON games(created_by);

-- Game players table (junction table for many-to-many)
CREATE TABLE IF NOT EXISTS game_players (
    id BIGSERIAL PRIMARY KEY,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    player_position INT NOT NULL, -- 0-3
    victory_points INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_player_position CHECK (player_position BETWEEN 0 AND 3),
    UNIQUE(game_id, user_id),
    UNIQUE(game_id, player_position)
);

CREATE INDEX idx_game_players_game_id ON game_players(game_id);
CREATE INDEX idx_game_players_user_id ON game_players(user_id);

-- Development cards reference data
CREATE TABLE IF NOT EXISTS development_cards (
    id BIGSERIAL PRIMARY KEY,
    tier INT NOT NULL, -- 1, 2, or 3
    gem_type VARCHAR(20) NOT NULL, -- diamond, sapphire, emerald, ruby, onyx
    victory_points INT NOT NULL DEFAULT 0,
    cost_diamond INT NOT NULL DEFAULT 0,
    cost_sapphire INT NOT NULL DEFAULT 0,
    cost_emerald INT NOT NULL DEFAULT 0,
    cost_ruby INT NOT NULL DEFAULT 0,
    cost_onyx INT NOT NULL DEFAULT 0,
    CONSTRAINT chk_tier CHECK (tier IN (1, 2, 3)),
    CONSTRAINT chk_gem_type CHECK (gem_type IN ('diamond', 'sapphire', 'emerald', 'ruby', 'onyx'))
);

CREATE INDEX idx_development_cards_tier ON development_cards(tier);

-- Nobles reference data
CREATE TABLE IF NOT EXISTS nobles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    victory_points INT NOT NULL DEFAULT 3,
    required_diamond INT NOT NULL DEFAULT 0,
    required_sapphire INT NOT NULL DEFAULT 0,
    required_emerald INT NOT NULL DEFAULT 0,
    required_ruby INT NOT NULL DEFAULT 0,
    required_onyx INT NOT NULL DEFAULT 0
);

-- Game state table (current board state)
CREATE TABLE IF NOT EXISTS game_state (
    id BIGSERIAL PRIMARY KEY,
    game_id BIGINT UNIQUE NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    available_gems JSONB NOT NULL, -- {diamond: 4, sapphire: 4, ...}
    visible_cards_tier1 JSONB NOT NULL, -- Array of card IDs
    visible_cards_tier2 JSONB NOT NULL,
    visible_cards_tier3 JSONB NOT NULL,
    available_nobles JSONB NOT NULL, -- Array of noble IDs
    deck_tier1_count INT NOT NULL DEFAULT 0,
    deck_tier2_count INT NOT NULL DEFAULT 0,
    deck_tier3_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_game_state_game_id ON game_state(game_id);

-- Player state table (each player's current state)
CREATE TABLE IF NOT EXISTS player_state (
    id BIGSERIAL PRIMARY KEY,
    game_player_id BIGINT UNIQUE NOT NULL REFERENCES game_players(id) ON DELETE CASCADE,

    -- Temporary gems (tokens in hand)
    gems_diamond INT NOT NULL DEFAULT 0,
    gems_sapphire INT NOT NULL DEFAULT 0,
    gems_emerald INT NOT NULL DEFAULT 0,
    gems_ruby INT NOT NULL DEFAULT 0,
    gems_onyx INT NOT NULL DEFAULT 0,
    gems_gold INT NOT NULL DEFAULT 0,

    -- Permanent gems (from purchased cards)
    permanent_diamond INT NOT NULL DEFAULT 0,
    permanent_sapphire INT NOT NULL DEFAULT 0,
    permanent_emerald INT NOT NULL DEFAULT 0,
    permanent_ruby INT NOT NULL DEFAULT 0,
    permanent_onyx INT NOT NULL DEFAULT 0,

    -- Cards and nobles
    purchased_cards JSONB NOT NULL DEFAULT '[]', -- Array of card IDs
    reserved_cards JSONB NOT NULL DEFAULT '[]', -- Array of card IDs (max 3)
    nobles JSONB NOT NULL DEFAULT '[]', -- Array of noble IDs

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_max_gems CHECK (
        gems_diamond + gems_sapphire + gems_emerald +
        gems_ruby + gems_onyx + gems_gold <= 10
    )
);

CREATE INDEX idx_player_state_game_player_id ON player_state(game_player_id);

-- Game moves history
CREATE TABLE IF NOT EXISTS game_moves (
    id BIGSERIAL PRIMARY KEY,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    game_player_id BIGINT NOT NULL REFERENCES game_players(id),
    move_number INT NOT NULL,
    move_type VARCHAR(50) NOT NULL, -- take_gems, reserve_card, purchase_card
    move_data JSONB NOT NULL, -- Details of the move
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_move_type CHECK (move_type IN ('take_gems', 'reserve_card', 'purchase_card'))
);

CREATE INDEX idx_game_moves_game_id ON game_moves(game_id);
CREATE INDEX idx_game_moves_game_player_id ON game_moves(game_player_id);
CREATE INDEX idx_game_moves_move_number ON game_moves(game_id, move_number);

-- Game statistics (aggregated data per user)
CREATE TABLE IF NOT EXISTS game_statistics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_games INT NOT NULL DEFAULT 0,
    total_wins INT NOT NULL DEFAULT 0,
    total_losses INT NOT NULL DEFAULT 0,
    average_points DECIMAL(5,2) NOT NULL DEFAULT 0,
    average_moves_per_game DECIMAL(5,2) NOT NULL DEFAULT 0,
    favorite_gem_type VARCHAR(20),
    total_nobles_earned INT NOT NULL DEFAULT 0,
    total_cards_purchased INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_favorite_gem CHECK (
        favorite_gem_type IS NULL OR
        favorite_gem_type IN ('diamond', 'sapphire', 'emerald', 'ruby', 'onyx')
    )
);

CREATE INDEX idx_game_statistics_user_id ON game_statistics(user_id);
CREATE INDEX idx_game_statistics_total_wins ON game_statistics(total_wins DESC);

-- Update timestamp trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply update timestamp triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_state_updated_at BEFORE UPDATE ON game_state
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_player_state_updated_at BEFORE UPDATE ON player_state
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_game_statistics_updated_at BEFORE UPDATE ON game_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
