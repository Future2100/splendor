import { useParams } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useFullGameState } from '../hooks/useFullGameState'
import { gameService } from '../services/gameService'
import ConnectionStatus from '../components/common/ConnectionStatus'
import GameBoard from '../components/game/GameBoard'

export default function GamePage() {
  const { id } = useParams<{ id: string }>()
  const gameId = parseInt(id || '0')
  const { user } = useAuth()

  const { gameState, loading, error, isConnected, refreshGameState } = useFullGameState(gameId)

  const handleTakeGems = async (gems: Record<string, number>) => {
    try {
      await gameService.takeGems(gameId, gems)
      setTimeout(refreshGameState, 500)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to take gems')
    }
  }

  const handlePurchaseCard = async (cardId: number, fromReserve: boolean) => {
    try {
      await gameService.purchaseCard(gameId, cardId, fromReserve)
      setTimeout(refreshGameState, 500)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to purchase card')
    }
  }

  const handleReserveCard = async (cardId: number, tier: number) => {
    try {
      await gameService.reserveCard(gameId, cardId, tier)
      setTimeout(refreshGameState, 500)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to reserve card')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-white text-xl">Loading game...</div>
      </div>
    )
  }

  if (error || !gameState || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="card max-w-md">
          <h2 className="text-2xl font-bold text-white mb-4">Error</h2>
          <p className="text-red-400">{error || 'Game not found'}</p>
        </div>
      </div>
    )
  }

  if (gameState.game.status === 'waiting') {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="card max-w-md text-center">
          <h2 className="text-2xl font-bold text-white mb-4">Waiting for Game to Start</h2>
          <p className="text-white/60">Room Code: {gameState.game.room_code}</p>
          <p className="text-white/60 mt-2">The host will start the game soon...</p>
        </div>
      </div>
    )
  }

  if (gameState.game.status === 'completed') {
    const winner = gameState.players.find(p => p.user_id === gameState.game.winner_id)
    const isWinner = gameState.game.winner_id === user.id

    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="card max-w-2xl text-center">
          <div className="text-6xl mb-4">{isWinner ? '🎉' : '🏆'}</div>
          <h2 className="text-4xl font-bold text-white mb-6">
            {isWinner ? 'You Won!' : 'Game Over'}
          </h2>

          <div className="bg-white/10 backdrop-blur-sm rounded-lg p-6 mb-6">
            <p className="text-2xl text-white mb-4">
              Winner: <span className="text-yellow-400 font-bold">{winner?.user?.username}</span>
            </p>

            <div className="space-y-3">
              <h3 className="text-lg font-bold text-white mb-3">Final Scores</h3>
              {gameState.players
                .sort((a, b) => b.victory_points - a.victory_points)
                .map((player, index) => (
                  <div
                    key={player.id}
                    className={`flex items-center justify-between p-3 rounded-lg ${
                      player.user_id === gameState.game.winner_id
                        ? 'bg-yellow-500/20 border-2 border-yellow-400'
                        : 'bg-white/5'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-2xl font-bold text-white/60">#{index + 1}</span>
                      <span className="text-white font-semibold">{player.user?.username}</span>
                      {player.user_id === gameState.game.winner_id && (
                        <span className="text-yellow-400">👑</span>
                      )}
                    </div>
                    <span className="text-2xl font-bold text-white">{player.victory_points} pts</span>
                  </div>
                ))}
            </div>
          </div>

          <button
            onClick={() => window.location.href = '/lobby'}
            className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 text-white font-bold py-3 px-8 rounded-lg transition-all shadow-lg"
          >
            Back to Lobby
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen p-4 md:p-8">
      <ConnectionStatus isConnected={isConnected} />

      <div className="max-w-7xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-3xl md:text-4xl font-bold text-white">Splendor</h1>
            <p className="text-white/60 text-sm">Room: {gameState.game.room_code}</p>
            <p className="text-purple-300 text-xs mt-1">
              Playing as <span className="font-semibold">{user.username}</span>
            </p>
          </div>
          <div className="flex items-center gap-3">
            <span className="bg-green-500/20 text-green-300 px-4 py-2 rounded-full font-semibold text-sm">
              In Progress
            </span>
          </div>
        </div>

        <GameBoard
          gameState={gameState}
          currentUserId={user.id}
          onTakeGems={handleTakeGems}
          onPurchaseCard={handlePurchaseCard}
          onReserveCard={handleReserveCard}
        />
      </div>
    </div>
  )
}
