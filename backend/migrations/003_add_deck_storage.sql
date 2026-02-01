-- Migration: Add deck array storage to game_state table
-- This allows proper card replacement when cards are purchased or reserved

ALTER TABLE game_state
ADD COLUMN IF NOT EXISTS deck_tier1 JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS deck_tier2 JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS deck_tier3 JSONB DEFAULT '[]';

-- Add index for better performance (only if it doesn't exist)
CREATE INDEX IF NOT EXISTS idx_game_state_game_id ON game_state(game_id);
