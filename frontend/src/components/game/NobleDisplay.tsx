import { motion } from 'framer-motion'
import type { Noble } from '../../types'

interface NobleDisplayProps {
  nobles: Noble[]
}

// 宝石小图标组件
const GemIcon = ({ type, size = 'small' }: { type: string; size?: 'small' | 'tiny' }) => {
  const sizeClass = size === 'tiny' ? 'text-xs' : 'text-sm'
  const gemIcons: Record<string, string> = {
    diamond: '💎',
    sapphire: '🔷',
    emerald: '💚',
    ruby: '❤️',
    onyx: '⚫',
  }

  return (
    <span className={sizeClass}>{gemIcons[type]}</span>
  )
}

export default function NobleDisplay({ nobles }: NobleDisplayProps) {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      className="card h-full flex flex-col"
    >
      <h3 className="text-sm font-bold text-white mb-1.5 flex items-center gap-1 flex-none">
        <span className="text-base">👑</span>
        <span>Nobles</span>
      </h3>
      <div className="space-y-1.5 flex-1 overflow-y-auto">
        {nobles.map((noble) => (
          <div
            key={noble.id}
            className="bg-gradient-to-br from-purple-600/90 to-purple-900/90 rounded-lg p-1.5 shadow-md border border-purple-400/30 hover:border-purple-400/60 transition-all"
          >
            <div className="flex items-center justify-between mb-1">
              <span className="text-white font-semibold text-[10px] truncate pr-1">{noble.name}</span>
              <span className="bg-gradient-to-br from-yellow-300 to-yellow-500 text-gray-900 font-black px-1.5 py-0.5 rounded-full text-[10px] shadow-md flex-shrink-0">
                {noble.victory_points}
              </span>
            </div>

            <div className="bg-black/20 rounded-md p-1">
              <div className="grid grid-cols-2 gap-1">
                {Object.entries(noble.required).map(([gem, count]) => (
                  count > 0 && (
                    <div key={gem} className="flex items-center gap-1 bg-white/10 px-1 py-0.5 rounded-md">
                      <GemIcon type={gem} size="tiny" />
                      <span className="text-white font-bold text-[10px]">{count}</span>
                    </div>
                  )
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </motion.div>
  )
}
