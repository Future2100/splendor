import { motion } from 'framer-motion'
import type { Game } from '../../types'

interface GameCardProps {
  game: Game
  onJoin?: (roomCode: string) => void
  onView?: (gameId: number) => void
  currentUserId?: number
}

export default function GameCard({ game, onJoin, onView, currentUserId }: GameCardProps) {
  const playerCount = game.players?.length || 0
  const isCreator = currentUserId === game.created_by
  const isPlayer = game.players?.some(p => p.user_id === currentUserId) || false
  const canJoin = game.status === 'waiting' && playerCount < game.num_players && !isPlayer

  const statusColors = {
    waiting: 'bg-yellow-500',
    in_progress: 'bg-green-500',
    completed: 'bg-gray-500',
  }

  const statusLabels = {
    waiting: 'Waiting',
    in_progress: 'In Progress',
    completed: 'Completed',
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="card hover:shadow-2xl transition-shadow"
    >
      <div className="flex items-start justify-between mb-4">
        <div>
          <div className="flex items-center gap-3">
            <h3 className="text-xl font-bold text-white">Room {game.room_code}</h3>
            <span className={`${statusColors[game.status]} text-white text-xs px-2 py-1 rounded-full`}>
              {statusLabels[game.status]}
            </span>
          </div>
          <p className="text-white/60 text-sm mt-1">
            Created by {game.creator?.username || 'Unknown'}
          </p>
        </div>
      </div>

      <div className="space-y-3">
        <div className="flex items-center justify-between text-white/80">
          <span>Players:</span>
          <span className="font-semibold">
            {playerCount} / {game.num_players}
          </span>
        </div>

        {game.players && game.players.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {game.players.map((player) => (
              <span
                key={player.id}
                className="bg-white/10 text-white text-xs px-3 py-1 rounded-full"
              >
                {player.user?.username}
              </span>
            ))}
          </div>
        )}

        <div className="flex gap-2 pt-2">
          {canJoin && onJoin && (
            <button
              onClick={() => onJoin(game.room_code)}
              className="btn-primary flex-1"
            >
              Join Game
            </button>
          )}
          {isPlayer && game.status === 'waiting' && onView && (
            <button
              onClick={() => onView(game.id)}
              className="btn-primary flex-1"
            >
              {isCreator ? 'Enter Lobby' : 'Return to Lobby'}
            </button>
          )}
          {game.status === 'in_progress' && onView && (
            <button
              onClick={() => onView(game.id)}
              className={isPlayer ? "btn-primary flex-1" : "btn-secondary flex-1"}
            >
              {isPlayer ? 'Continue Game' : 'View Game'}
            </button>
          )}
          {game.status === 'completed' && onView && (
            <button
              onClick={() => onView(game.id)}
              className="btn-secondary flex-1"
            >
              View Results
            </button>
          )}
        </div>
      </div>
    </motion.div>
  )
}
