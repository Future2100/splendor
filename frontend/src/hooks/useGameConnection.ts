import { useEffect, useState } from 'react'
import { useWebSocket } from '../context/WebSocketContext'
import { gameService } from '../services/gameService'
import type { Game } from '../types'

export function useGameConnection(gameId: number) {
  const { connect, disconnect, isConnected, lastMessage } = useWebSocket()
  const [game, setGame] = useState<Game | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadGame()
  }, [gameId])

  useEffect(() => {
    if (game) {
      // Get access token
      const token = localStorage.getItem('access_token')
      if (token) {
        connect(gameId, token)
      }
    }

    return () => {
      disconnect()
    }
  }, [game])

  // Handle WebSocket messages
  useEffect(() => {
    if (!lastMessage) return

    switch (lastMessage.type) {
      case 'game_update':
        // Reload game when state changes
        loadGame()
        break

      case 'player_event':
        // Handle player events (join, leave)
        loadGame()
        break

      case 'game_end':
        // Handle game end
        loadGame()
        break

      case 'error':
        setError(lastMessage.payload.message || 'An error occurred')
        break
    }
  }, [lastMessage])

  const loadGame = async () => {
    try {
      setLoading(true)
      const gameData = await gameService.getGame(gameId)
      setGame(gameData)
      setError(null)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load game')
    } finally {
      setLoading(false)
    }
  }

  return {
    game,
    loading,
    error,
    isConnected,
    refreshGame: loadGame,
  }
}
