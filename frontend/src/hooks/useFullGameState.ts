import { useEffect, useState } from 'react'
import { useWebSocket } from '../context/WebSocketContext'
import { gameService } from '../services/gameService'
import type { FullGameState } from '../types'

export function useFullGameState(gameId: number) {
  const { connect, disconnect, isConnected, lastMessage } = useWebSocket()
  const [gameState, setGameState] = useState<FullGameState | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadGameState()
  }, [gameId])

  useEffect(() => {
    if (gameState) {
      const token = localStorage.getItem('access_token')
      if (token) {
        connect(gameId, token)
      }
    }

    return () => {
      disconnect()
    }
  }, [gameState])

  // Handle WebSocket messages
  useEffect(() => {
    if (!lastMessage) return

    switch (lastMessage.type) {
      case 'game_update':
      case 'player_event':
      case 'game_end':
        loadGameState()
        break

      case 'error':
        setError(lastMessage.payload.message || 'An error occurred')
        break
    }
  }, [lastMessage])

  const loadGameState = async () => {
    try {
      setLoading(true)
      const state = await gameService.getGameState(gameId)
      setGameState(state)
      setError(null)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load game state')
    } finally {
      setLoading(false)
    }
  }

  return {
    gameState,
    loading,
    error,
    isConnected,
    refreshGameState: loadGameState,
  }
}
