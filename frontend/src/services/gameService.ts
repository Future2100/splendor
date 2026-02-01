import api from './api'
import type { Game, CreateGameRequest, GameListResponse } from '../types'

export const gameService = {
  async createGame(data: CreateGameRequest): Promise<{ game: Game; room_code: string }> {
    const response = await api.post('/games', data)
    return response.data
  },

  async listGames(status?: string, limit = 20, offset = 0): Promise<GameListResponse> {
    const params = new URLSearchParams()
    if (status) params.append('status', status)
    params.append('limit', limit.toString())
    params.append('offset', offset.toString())

    const response = await api.get<GameListResponse>(`/games?${params.toString()}`)
    return response.data
  },

  async getGame(gameId: number): Promise<Game> {
    const response = await api.get<{ game: Game }>(`/games/${gameId}`)
    return response.data.game
  },

  async joinGame(roomCode: string): Promise<Game> {
    const response = await api.post<{ game: Game }>('/games/join', { room_code: roomCode })
    return response.data.game
  },

  async leaveGame(gameId: number): Promise<void> {
    await api.post(`/games/${gameId}/leave`)
  },

  async startGame(gameId: number): Promise<Game> {
    const response = await api.post<{ game: Game }>(`/games/${gameId}/start`)
    return response.data.game
  },

  async getGameState(gameId: number): Promise<any> {
    const response = await api.get(`/games/${gameId}/state`)
    return response.data.state
  },

  async takeGems(gameId: number, gems: Record<string, number>): Promise<void> {
    await api.post(`/games/${gameId}/take-gems`, { gems })
  },

  async purchaseCard(gameId: number, cardId: number, fromReserve: boolean): Promise<void> {
    await api.post(`/games/${gameId}/purchase-card`, { card_id: cardId, from_reserve: fromReserve })
  },

  async reserveCard(gameId: number, cardId: number, tier: number): Promise<void> {
    await api.post(`/games/${gameId}/reserve-card`, { card_id: cardId, tier })
  },
}
