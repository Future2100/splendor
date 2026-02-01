-- Splendor Game Data
-- Development Cards (90 total) and Nobles (10 total)

-- TIER 1 CARDS (40 cards) - 0-1 victory points
-- Diamond cards (8 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(1, 'diamond', 0, 0, 3, 0, 0, 0),
(1, 'diamond', 0, 0, 0, 0, 2, 1),
(1, 'diamond', 0, 0, 1, 1, 1, 1),
(1, 'diamond', 0, 0, 2, 0, 2, 0),
(1, 'diamond', 0, 0, 0, 4, 0, 0),
(1, 'diamond', 0, 0, 1, 2, 1, 1),
(1, 'diamond', 0, 0, 2, 2, 0, 1),
(1, 'diamond', 1, 0, 0, 0, 3, 0);

-- Sapphire cards (8 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(1, 'sapphire', 0, 0, 0, 0, 0, 3),
(1, 'sapphire', 0, 1, 0, 0, 0, 2),
(1, 'sapphire', 0, 1, 0, 1, 1, 1),
(1, 'sapphire', 0, 0, 0, 2, 0, 2),
(1, 'sapphire', 0, 0, 0, 0, 4, 0),
(1, 'sapphire', 0, 1, 0, 1, 2, 1),
(1, 'sapphire', 0, 1, 0, 2, 2, 0),
(1, 'sapphire', 1, 0, 0, 0, 0, 3);

-- Emerald cards (8 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(1, 'emerald', 0, 0, 0, 0, 3, 0),
(1, 'emerald', 0, 2, 1, 0, 0, 0),
(1, 'emerald', 0, 1, 1, 0, 1, 1),
(1, 'emerald', 0, 2, 0, 0, 2, 0),
(1, 'emerald', 0, 0, 4, 0, 0, 0),
(1, 'emerald', 0, 1, 2, 0, 1, 1),
(1, 'emerald', 0, 0, 2, 0, 2, 1),
(1, 'emerald', 1, 3, 0, 0, 0, 0);

-- Ruby cards (8 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(1, 'ruby', 0, 3, 0, 0, 0, 0),
(1, 'ruby', 0, 0, 0, 1, 0, 2),
(1, 'ruby', 0, 1, 1, 1, 0, 1),
(1, 'ruby', 0, 2, 0, 2, 0, 0),
(1, 'ruby', 0, 0, 0, 4, 0, 0),
(1, 'ruby', 0, 2, 1, 1, 0, 1),
(1, 'ruby', 0, 2, 0, 1, 0, 2),
(1, 'ruby', 1, 0, 3, 0, 0, 0);

-- Onyx cards (8 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(1, 'onyx', 0, 0, 0, 3, 0, 0),
(1, 'onyx', 0, 2, 0, 0, 1, 0),
(1, 'onyx', 0, 1, 1, 1, 1, 0),
(1, 'onyx', 0, 0, 2, 0, 0, 2),
(1, 'onyx', 0, 4, 0, 0, 0, 0),
(1, 'onyx', 0, 1, 1, 1, 2, 0),
(1, 'onyx', 0, 1, 2, 0, 0, 2),
(1, 'onyx', 1, 0, 0, 3, 0, 0);

-- TIER 2 CARDS (30 cards) - 1-3 victory points
-- Diamond cards (6 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(2, 'diamond', 1, 0, 0, 3, 2, 2),
(2, 'diamond', 1, 0, 0, 0, 5, 0),
(2, 'diamond', 2, 0, 0, 0, 5, 3),
(2, 'diamond', 2, 0, 5, 0, 0, 0),
(2, 'diamond', 2, 6, 0, 0, 0, 0),
(2, 'diamond', 3, 0, 0, 0, 0, 6);

-- Sapphire cards (6 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(2, 'sapphire', 1, 2, 0, 2, 3, 0),
(2, 'sapphire', 1, 5, 0, 0, 0, 0),
(2, 'sapphire', 2, 5, 0, 3, 0, 0),
(2, 'sapphire', 2, 0, 0, 0, 0, 5),
(2, 'sapphire', 2, 0, 0, 6, 0, 0),
(2, 'sapphire', 3, 0, 0, 0, 6, 0);

-- Emerald cards (6 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(2, 'emerald', 1, 3, 0, 2, 0, 2),
(2, 'emerald', 1, 0, 0, 5, 0, 0),
(2, 'emerald', 2, 0, 0, 5, 3, 0),
(2, 'emerald', 2, 0, 5, 0, 0, 0),
(2, 'emerald', 2, 0, 6, 0, 0, 0),
(2, 'emerald', 3, 6, 0, 0, 0, 0);

-- Ruby cards (6 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(2, 'ruby', 1, 2, 3, 0, 0, 2),
(2, 'ruby', 1, 0, 5, 0, 0, 0),
(2, 'ruby', 2, 3, 0, 0, 0, 5),
(2, 'ruby', 2, 0, 0, 0, 5, 0),
(2, 'ruby', 2, 0, 0, 0, 6, 0),
(2, 'ruby', 3, 0, 6, 0, 0, 0);

-- Onyx cards (6 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(2, 'onyx', 1, 0, 2, 2, 0, 3),
(2, 'onyx', 1, 0, 0, 0, 0, 5),
(2, 'onyx', 2, 0, 3, 0, 5, 0),
(2, 'onyx', 2, 5, 0, 0, 0, 0),
(2, 'onyx', 2, 0, 0, 0, 0, 6),
(2, 'onyx', 3, 0, 0, 6, 0, 0);

-- TIER 3 CARDS (20 cards) - 3-5 victory points
-- Diamond cards (4 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(3, 'diamond', 3, 0, 3, 3, 5, 3),
(3, 'diamond', 4, 0, 0, 0, 7, 0),
(3, 'diamond', 4, 0, 0, 3, 6, 3),
(3, 'diamond', 5, 0, 0, 0, 7, 3);

-- Sapphire cards (4 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(3, 'sapphire', 3, 3, 0, 5, 3, 3),
(3, 'sapphire', 4, 7, 0, 0, 0, 0),
(3, 'sapphire', 4, 6, 0, 3, 3, 0),
(3, 'sapphire', 5, 7, 0, 3, 0, 0);

-- Emerald cards (4 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(3, 'emerald', 3, 5, 3, 0, 3, 3),
(3, 'emerald', 4, 0, 7, 0, 0, 0),
(3, 'emerald', 4, 3, 6, 0, 3, 0),
(3, 'emerald', 5, 3, 7, 0, 0, 0);

-- Ruby cards (4 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(3, 'ruby', 3, 3, 5, 3, 0, 3),
(3, 'ruby', 4, 0, 0, 7, 0, 0),
(3, 'ruby', 4, 3, 0, 6, 0, 3),
(3, 'ruby', 5, 0, 3, 7, 0, 0);

-- Onyx cards (4 cards)
INSERT INTO development_cards (tier, gem_type, victory_points, cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx) VALUES
(3, 'onyx', 3, 3, 3, 5, 3, 0),
(3, 'onyx', 4, 0, 0, 0, 0, 7),
(3, 'onyx', 4, 0, 3, 6, 3, 0),
(3, 'onyx', 5, 0, 0, 7, 3, 0);

-- NOBLES (10 total) - Each worth 3 victory points
INSERT INTO nobles (name, victory_points, required_diamond, required_sapphire, required_emerald, required_ruby, required_onyx) VALUES
('Catherine de Medici', 3, 0, 0, 4, 4, 0),
('Charles V', 3, 0, 3, 3, 3, 0),
('Francis I', 3, 0, 4, 0, 0, 4),
('Henry VIII', 3, 4, 0, 0, 0, 4),
('Mary Stuart', 3, 0, 0, 3, 3, 3),
('Anne of Brittany', 3, 3, 0, 0, 3, 3),
('Isabella I', 3, 0, 4, 4, 0, 0),
('Cosimo de Medici', 3, 4, 0, 4, 0, 0),
('Suleiman the Magnificent', 3, 3, 3, 0, 0, 3),
('Elisabeth of Austria', 3, 4, 4, 0, 0, 0);

-- Verify counts
SELECT tier, COUNT(*) as count FROM development_cards GROUP BY tier ORDER BY tier;
SELECT COUNT(*) as noble_count FROM nobles;
