// Game Types
export type GemType = 'diamond' | 'sapphire' | 'emerald' | 'ruby' | 'onyx' | 'gold'

export type GameStatus = 'waiting' | 'in_progress' | 'completed'

export type MoveType = 'take_gems' | 'reserve_card' | 'purchase_card'

export interface User {
  id: number
  username: string
  email: string
  created_at: string
}

export interface Game {
  id: number
  room_code: string
  status: GameStatus
  current_turn_player_id?: number
  turn_number: number
  winner_id?: number
  created_by: number
  num_players: number
  created_at: string
  started_at?: string
  completed_at?: string
  players?: GamePlayer[]
  creator?: User
}

export interface GamePlayer {
  id: number
  game_id: number
  user_id: number
  player_position: number
  victory_points: number
  is_active: boolean
  joined_at: string
  user?: User
}

export interface DevelopmentCard {
  id: number
  tier: 1 | 2 | 3
  gem_type: GemType
  victory_points: number
  cost: Record<string, number>
}

export interface Noble {
  id: number
  name: string
  victory_points: number
  required: Record<string, number>
}

export interface GameState {
  id: number
  game_id: number
  available_gems: Record<string, number>
  visible_cards_tier1: DevelopmentCard[]
  visible_cards_tier2: DevelopmentCard[]
  visible_cards_tier3: DevelopmentCard[]
  available_nobles: Noble[]
  deck_tier1_count: number
  deck_tier2_count: number
  deck_tier3_count: number
  updated_at: string
}

export interface PlayerState {
  id: number
  game_player_id: number
  gems: Record<string, number>
  permanent_gems: Record<string, number>
  purchased_cards: DevelopmentCard[]
  reserved_cards: DevelopmentCard[]
  nobles: Noble[]
  updated_at: string
}

export interface FullGameState {
  game: Game
  players: GamePlayer[]
  game_state: GameState
  player_states: Record<number, PlayerState>
}

export interface GameMove {
  id: number
  game_id: number
  game_player_id: number
  move_number: number
  move_type: MoveType
  move_data: any
  created_at: string
}

export interface GameStatistics {
  id: number
  user_id: number
  total_games: number
  total_wins: number
  total_losses: number
  average_points: number
  average_moves_per_game: number
  favorite_gem_type: GemType
  total_nobles_earned: number
  total_cards_purchased: number
  updated_at: string
}

// WebSocket Message Types
export interface WSMessage {
  type: 'game_update' | 'player_event' | 'game_end' | 'error' | 'move'
  payload: any
}

// Auth Types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: User
}

// Game API Types
export interface CreateGameRequest {
  num_players: number
}

export interface JoinGameRequest {
  room_code: string
}

export interface GameListResponse {
  games: Game[]
  total: number
}
