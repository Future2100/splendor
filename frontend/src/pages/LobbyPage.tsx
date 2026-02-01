import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { gameService } from '../services/gameService'
import CreateGameModal from '../components/lobby/CreateGameModal'
import GameCard from '../components/lobby/GameCard'
import WaitingRoom from '../components/lobby/WaitingRoom'
import type { Game } from '../types'

export default function LobbyPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [games, setGames] = useState<Game[]>([])
  const [currentGame, setCurrentGame] = useState<Game | null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('waiting')

  useEffect(() => {
    loadGames()
    // Auto-refresh every 3 seconds
    const interval = setInterval(loadGames, 3000)
    return () => clearInterval(interval)
  }, [statusFilter])

  // Auto-refresh current game in waiting room
  useEffect(() => {
    if (currentGame && currentGame.status === 'waiting') {
      const interval = setInterval(async () => {
        try {
          const updated = await gameService.getGame(currentGame.id)
          setCurrentGame(updated)

          // If game started, navigate to game page
          if (updated.status === 'in_progress') {
            navigate(`/game/${updated.id}`)
          }
        } catch (err) {
          console.error('Failed to refresh game:', err)
        }
      }, 2000)
      return () => clearInterval(interval)
    }
  }, [currentGame?.id, currentGame?.status, navigate])

  useEffect(() => {
    // Check if we have a room code in URL for auto-join
    const joinRoomCode = searchParams.get('join')
    if (joinRoomCode && user && !currentGame) {
      handleJoinGame(joinRoomCode)
      return
    }

    // Check if we have a room code in URL
    const roomCode = searchParams.get('room')
    if (roomCode && user) {
      checkCurrentGame()
    }
  }, [searchParams, user])

  const loadGames = async () => {
    try {
      const response = await gameService.listGames(statusFilter || undefined)
      setGames(response.games)
    } catch (err) {
      console.error('Failed to load games:', err)
    }
  }

  const checkCurrentGame = async () => {
    // Check if user is in any waiting game
    try {
      const response = await gameService.listGames('waiting')
      const myGame = response.games.find((g) =>
        g.players?.some((p) => p.user_id === user?.id)
      )
      if (myGame) {
        const fullGame = await gameService.getGame(myGame.id)
        setCurrentGame(fullGame)
      }
    } catch (err) {
      console.error('Failed to check current game:', err)
    }
  }

  const handleCreateGame = async (numPlayers: number) => {
    setLoading(true)
    setError('')
    try {
      const response = await gameService.createGame({ num_players: numPlayers })
      setCurrentGame(response.game)
      setShowCreateModal(false)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create game')
    } finally {
      setLoading(false)
    }
  }

  const handleJoinGame = async (roomCode: string) => {
    setLoading(true)
    setError('')
    try {
      const game = await gameService.joinGame(roomCode)
      setCurrentGame(game)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to join game')
    } finally {
      setLoading(false)
    }
  }

  const handleLeaveGame = async () => {
    if (!currentGame) return
    setLoading(true)
    try {
      await gameService.leaveGame(currentGame.id)
      setCurrentGame(null)
      loadGames()
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to leave game')
    } finally {
      setLoading(false)
    }
  }

  const handleStartGame = async () => {
    if (!currentGame) return
    setLoading(true)
    try {
      const game = await gameService.startGame(currentGame.id)
      navigate(`/game/${game.id}`)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to start game')
    } finally {
      setLoading(false)
    }
  }

  const handleViewGame = async (gameId: number) => {
    try {
      const game = await gameService.getGame(gameId)
      // If game is waiting, show waiting room instead of game page
      if (game.status === 'waiting') {
        setCurrentGame(game)
      } else {
        navigate(`/game/${gameId}`)
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load game')
    }
  }

  // If user is in a waiting game, show waiting room
  if (currentGame && currentGame.status === 'waiting' && user) {
    return (
      <div className="min-h-screen p-8">
        <WaitingRoom
          game={currentGame}
          currentUserId={user.id}
          onLeave={handleLeaveGame}
          onStart={handleStartGame}
          loading={loading}
        />
      </div>
    )
  }

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold text-white">Game Lobby</h1>
            {user && (
              <p className="text-purple-300 text-sm mt-1">
                Logged in as <span className="font-semibold">{user.username}</span>
              </p>
            )}
          </div>
          <button onClick={() => setShowCreateModal(true)} className="btn-primary">
            Create Game
          </button>
        </div>

        {error && (
          <div className="bg-red-500/20 border border-red-500 text-red-200 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {/* Filter tabs */}
        <div className="flex gap-2 mb-6">
          {['waiting', 'in_progress', 'completed'].map((status) => (
            <button
              key={status}
              onClick={() => setStatusFilter(status)}
              className={`px-6 py-2 rounded-lg font-semibold transition-all ${
                statusFilter === status
                  ? 'bg-purple-600 text-white'
                  : 'bg-white/10 text-white/70 hover:bg-white/20'
              }`}
            >
              {status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ')}
            </button>
          ))}
        </div>

        {games.length === 0 ? (
          <div className="card text-center py-12">
            <p className="text-white/60 text-lg">
              No {statusFilter} games found. Create a new game to get started!
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {games.map((game) => (
              <GameCard
                key={game.id}
                game={game}
                onJoin={handleJoinGame}
                onView={handleViewGame}
                currentUserId={user?.id}
              />
            ))}
          </div>
        )}

        <CreateGameModal
          isOpen={showCreateModal}
          onClose={() => setShowCreateModal(false)}
          onCreate={handleCreateGame}
          loading={loading}
        />
      </div>
    </div>
  )
}
