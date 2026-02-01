import { motion } from 'framer-motion'
import type { GamePlayer, PlayerState } from '../../types'

interface PlayerPanelsProps {
  players: GamePlayer[]
  playerStates: Record<number, PlayerState>
  currentUserId: number
  currentTurnUserId: number
}

const gemIcons: Record<string, string> = {
  diamond: '💎',
  sapphire: '🔷',
  emerald: '💚',
  ruby: '❤️',
  onyx: '⚫',
  gold: '🪙',
}

export default function PlayerPanels({
  players,
  playerStates,
  currentUserId,
  currentTurnUserId
}: PlayerPanelsProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {players.map((player) => {
        const state = playerStates[player.user_id]
        const isCurrentUser = player.user_id === currentUserId
        const isCurrentTurn = player.user_id === currentTurnUserId

        const totalGems = state
          ? Object.values(state.gems).reduce((sum, count) => sum + (count || 0), 0)
          : 0

        return (
          <motion.div
            key={player.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className={`card relative ${isCurrentTurn ? 'ring-4 ring-green-500 shadow-lg shadow-green-500/50' : ''} ${
              isCurrentUser ? 'bg-purple-900/40 border-2 border-purple-500' : ''
            }`}
          >
            {/* Turn Indicator */}
            {isCurrentTurn && (
              <div className="absolute -top-2 left-1/2 transform -translate-x-1/2 bg-green-500 text-white text-xs font-bold px-3 py-1 rounded-full shadow-lg">
                TURN
              </div>
            )}

            {/* Player Info */}
            <div className="text-center mb-3">
              <div className="w-14 h-14 mx-auto mb-2 bg-gradient-to-br from-purple-500 to-blue-500 rounded-full flex items-center justify-center text-white text-xl font-bold shadow-lg">
                {player.user?.username.charAt(0).toUpperCase()}
              </div>
              <p className="text-white font-bold text-sm">{player.user?.username}</p>
              {isCurrentUser && (
                <p className="text-purple-300 text-xs font-semibold">You</p>
              )}
              <div className="mt-1 inline-block bg-yellow-400 text-black font-bold px-3 py-1 rounded-full text-base shadow-md">
                {player.victory_points} pts
              </div>
            </div>

            {state && (
              <div className="space-y-2">
                {/* Gems Owned */}
                <div className="bg-white/10 backdrop-blur-sm rounded-lg p-2">
                  <div className="text-xs text-white/80 font-semibold mb-1 flex items-center justify-between">
                    <span>Gems ({totalGems}/10)</span>
                  </div>
                  <div className="grid grid-cols-3 gap-1">
                    {Object.entries(state.gems).map(([gemType, count]) => (
                      count > 0 && (
                        <div key={gemType} className="flex items-center gap-1 text-white">
                          <span className="text-sm">{gemIcons[gemType]}</span>
                          <span className="text-xs font-bold">{count}</span>
                        </div>
                      )
                    ))}
                  </div>
                  {totalGems === 0 && (
                    <div className="text-white/50 text-xs text-center">No gems</div>
                  )}
                </div>

                {/* Permanent Gems (Bonuses) */}
                <div className="bg-white/10 backdrop-blur-sm rounded-lg p-2">
                  <div className="text-xs text-white/80 font-semibold mb-1">
                    Bonuses ({state.purchased_cards.length} cards)
                  </div>
                  <div className="grid grid-cols-3 gap-1">
                    {Object.entries(state.permanent_gems).map(([gemType, count]) => (
                      count > 0 && (
                        <div key={gemType} className="flex items-center gap-1 text-white">
                          <span className="text-sm">{gemIcons[gemType]}</span>
                          <span className="text-xs font-bold">{count}</span>
                        </div>
                      )
                    ))}
                  </div>
                  {state.purchased_cards.length === 0 && (
                    <div className="text-white/50 text-xs text-center">No bonuses</div>
                  )}
                </div>

                {/* Other Stats */}
                <div className="flex justify-between text-xs text-white/70">
                  <span>Reserved: {state.reserved_cards.length}/3</span>
                  <span>Nobles: {state.nobles.length}</span>
                </div>
              </div>
            )}
          </motion.div>
        )
      })}
    </div>
  )
}
