import { motion } from 'framer-motion'

export default function GameRules() {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="card bg-gradient-to-br from-blue-900/30 to-indigo-900/30 border-2 border-blue-500/30"
    >
      <div className="flex items-center gap-2 mb-3">
        <span className="text-2xl">📋</span>
        <h3 className="text-lg font-bold text-white">Game Rules</h3>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm">
        {/* Take Gems */}
        <div className="bg-white/5 backdrop-blur-sm rounded-lg p-3">
          <h4 className="text-white font-bold mb-2 flex items-center gap-1">
            <span>💎</span>
            <span>Take Gems</span>
          </h4>
          <ul className="text-white/80 space-y-1 text-xs">
            <li>• Take 3 different colored gems, OR</li>
            <li>• Take 2 same colored gems (if 4+ available)</li>
            <li>• Maximum 10 gems total in hand</li>
            <li>• Gold can only be obtained by reserving</li>
          </ul>
        </div>

        {/* Purchase Card */}
        <div className="bg-white/5 backdrop-blur-sm rounded-lg p-3">
          <h4 className="text-white font-bold mb-2 flex items-center gap-1">
            <span>🎴</span>
            <span>Purchase Card</span>
          </h4>
          <ul className="text-white/80 space-y-1 text-xs">
            <li>• Pay gems shown on card</li>
            <li>• Permanent gems count as payment</li>
            <li>• Gold acts as any color</li>
            <li>• Gain permanent gem & victory points</li>
          </ul>
        </div>

        {/* Reserve Card */}
        <div className="bg-white/5 backdrop-blur-sm rounded-lg p-3">
          <h4 className="text-white font-bold mb-2 flex items-center gap-1">
            <span>🃏</span>
            <span>Reserve Card</span>
          </h4>
          <ul className="text-white/80 space-y-1 text-xs">
            <li>• Maximum 3 reserved cards</li>
            <li>• Get 1 gold token (if available)</li>
            <li>• Can reserve from visible or deck</li>
            <li>• Purchase later from your reserve</li>
          </ul>
        </div>

        {/* Victory Conditions */}
        <div className="bg-white/5 backdrop-blur-sm rounded-lg p-3">
          <h4 className="text-white font-bold mb-2 flex items-center gap-1">
            <span>🏆</span>
            <span>Victory</span>
          </h4>
          <ul className="text-white/80 space-y-1 text-xs">
            <li>• First to reach 15 points triggers end</li>
            <li>• All players get equal turns</li>
            <li>• Highest points wins</li>
            <li>• Tiebreaker: fewer cards wins</li>
            <li>• Nobles visit automatically (3 pts)</li>
          </ul>
        </div>
      </div>
    </motion.div>
  )
}
