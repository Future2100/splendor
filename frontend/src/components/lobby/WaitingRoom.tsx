import { useState } from 'react'
import { motion } from 'framer-motion'
import type { Game } from '../../types'

interface WaitingRoomProps {
  game: Game
  currentUserId: number
  onLeave: () => void
  onStart: () => void
  loading: boolean
}

export default function WaitingRoom({ game, currentUserId, onLeave, onStart, loading }: WaitingRoomProps) {
  const isCreator = currentUserId === game.created_by
  const playerCount = game.players?.length || 0
  const canStart = playerCount >= 2
  const [copiedText, setCopiedText] = useState<string>('')

  const joinUrl = `${window.location.origin}/lobby?join=${game.room_code}`

  const copyToClipboard = async (text: string, label: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedText(label)
      setTimeout(() => setCopiedText(''), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  return (
    <div className="max-w-4xl mx-auto">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="card"
      >
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">Waiting Room</h1>
          <p className="text-2xl text-purple-400 font-mono">Room Code: {game.room_code}</p>
          <p className="text-white/60 mt-2">Share this code with your friends to join!</p>

          {/* Share Section */}
          <div className="mt-6 bg-white/5 backdrop-blur-sm rounded-lg p-4">
            <h3 className="text-sm font-semibold text-white/80 mb-3">📤 Share with friends</h3>
            <div className="flex flex-col sm:flex-row gap-2">
              <button
                onClick={() => copyToClipboard(game.room_code, 'code')}
                className="flex-1 bg-purple-600 hover:bg-purple-700 text-white font-semibold py-2 px-4 rounded-lg transition-all flex items-center justify-center gap-2"
              >
                {copiedText === 'code' ? (
                  <>
                    <span>✓</span>
                    <span>Code Copied!</span>
                  </>
                ) : (
                  <>
                    <span>📋</span>
                    <span>Copy Room Code</span>
                  </>
                )}
              </button>
              <button
                onClick={() => copyToClipboard(joinUrl, 'link')}
                className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition-all flex items-center justify-center gap-2"
              >
                {copiedText === 'link' ? (
                  <>
                    <span>✓</span>
                    <span>Link Copied!</span>
                  </>
                ) : (
                  <>
                    <span>🔗</span>
                    <span>Copy Join Link</span>
                  </>
                )}
              </button>
            </div>
          </div>
        </div>

        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-white">Players ({playerCount}/{game.num_players})</h2>
            {isCreator && (
              <span className="bg-purple-600 text-white text-sm px-3 py-1 rounded-full">
                You are the host
              </span>
            )}
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {game.players?.map((player, index) => (
              <motion.div
                key={player.id}
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: index * 0.1 }}
                className="bg-white/10 rounded-lg p-4 text-center"
              >
                <div className="w-16 h-16 mx-auto mb-2 bg-gradient-to-br from-purple-500 to-blue-500 rounded-full flex items-center justify-center text-white text-2xl font-bold">
                  {player.user?.username.charAt(0).toUpperCase()}
                </div>
                <p className="text-white font-semibold">{player.user?.username}</p>
                {player.user_id === game.created_by && (
                  <p className="text-purple-400 text-xs mt-1">Host</p>
                )}
              </motion.div>
            ))}

            {/* Empty slots */}
            {Array.from({ length: game.num_players - playerCount }).map((_, index) => (
              <div
                key={`empty-${index}`}
                className="bg-white/5 rounded-lg p-4 text-center border-2 border-dashed border-white/20"
              >
                <div className="w-16 h-16 mx-auto mb-2 bg-white/5 rounded-full flex items-center justify-center text-white/30 text-2xl">
                  ?
                </div>
                <p className="text-white/40 text-sm">Waiting...</p>
              </div>
            ))}
          </div>
        </div>

        <div className="flex gap-4">
          <button onClick={onLeave} className="btn-secondary flex-1" disabled={loading}>
            Leave Game
          </button>
          {isCreator && (
            <button
              onClick={onStart}
              className="btn-primary flex-1"
              disabled={loading || !canStart}
            >
              {loading ? 'Starting...' : canStart ? 'Start Game' : `Need ${2 - playerCount} more players`}
            </button>
          )}
        </div>

        {!isCreator && (
          <p className="text-center text-white/60 text-sm mt-4">
            Waiting for host to start the game...
          </p>
        )}
      </motion.div>
    </div>
  )
}
