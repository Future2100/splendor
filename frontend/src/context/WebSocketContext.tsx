import { createContext, useContext, useState, useEffect, useCallback, ReactNode, useRef } from 'react'
import type { WSMessage } from '../types'

interface WebSocketContextType {
  isConnected: boolean
  connect: (gameId: number, token: string) => void
  disconnect: () => void
  sendMessage: (message: WSMessage) => void
  lastMessage: WSMessage | null
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined)

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080'

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const [isConnected, setIsConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<number | null>(null)
  const gameIdRef = useRef<number | null>(null)
  const tokenRef = useRef<string | null>(null)

  const connect = useCallback((gameId: number, token: string) => {
    // Store connection params for reconnection
    gameIdRef.current = gameId
    tokenRef.current = token

    // Close existing connection
    if (wsRef.current) {
      wsRef.current.close()
    }

    try {
      const ws = new WebSocket(`${WS_URL}/api/v1/ws/games/${gameId}?token=${token}`)

      ws.onopen = () => {
        console.log('WebSocket connected')
        setIsConnected(true)
        // Clear reconnect timeout if connection successful
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current)
          reconnectTimeoutRef.current = null
        }
      }

      ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data)
          console.log('WebSocket message received:', message)
          setLastMessage(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      ws.onclose = () => {
        console.log('WebSocket disconnected')
        setIsConnected(false)

        // Attempt to reconnect after 3 seconds
        if (gameIdRef.current && tokenRef.current) {
          reconnectTimeoutRef.current = window.setTimeout(() => {
            console.log('Attempting to reconnect...')
            connect(gameIdRef.current!, tokenRef.current!)
          }, 3000)
        }
      }

      wsRef.current = ws
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }, [])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    gameIdRef.current = null
    tokenRef.current = null

    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }

    setIsConnected(false)
  }, [])

  const sendMessage = useCallback((message: WSMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not connected')
    }
  }, [])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  const value = {
    isConnected,
    connect,
    disconnect,
    sendMessage,
    lastMessage,
  }

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>
}

export function useWebSocket() {
  const context = useContext(WebSocketContext)
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }
  return context
}
