import { motion } from 'framer-motion'
import type { DevelopmentCard as CardType } from '../../types'

interface DevelopmentCardProps {
  card: CardType
  onPurchase?: () => void
  onReserve?: () => void
  disabled?: boolean
  canAfford?: boolean
  canReserve?: boolean
}

// 宝石类型的主题色
const gemThemes: Record<string, { bg: string; accent: string; border: string }> = {
  diamond: {
    bg: 'bg-gradient-to-br from-cyan-400 via-blue-500 to-indigo-600',
    accent: 'bg-blue-700',
    border: 'border-cyan-300'
  },
  sapphire: {
    bg: 'bg-gradient-to-br from-blue-500 via-indigo-600 to-blue-700',
    accent: 'bg-blue-900',
    border: 'border-blue-300'
  },
  emerald: {
    bg: 'bg-gradient-to-br from-emerald-400 via-green-500 to-emerald-600',
    accent: 'bg-emerald-800',
    border: 'border-emerald-300'
  },
  ruby: {
    bg: 'bg-gradient-to-br from-rose-400 via-red-500 to-rose-600',
    accent: 'bg-red-800',
    border: 'border-red-300'
  },
  onyx: {
    bg: 'bg-gradient-to-br from-gray-600 via-slate-700 to-gray-800',
    accent: 'bg-gray-900',
    border: 'border-gray-400'
  },
}

// 宝石图标（使用emoji，和Gem Bank保持一致）
const GemIcon = ({ type, size = 'large' }: { type: string; size?: 'small' | 'large' | 'medium' }) => {
  const sizeClass = size === 'large' ? 'text-2xl' : size === 'medium' ? 'text-base' : 'text-sm'
  const gemIcons: Record<string, string> = {
    diamond: '💎',
    sapphire: '🔷',
    emerald: '💚',
    ruby: '❤️',
    onyx: '⚫',
  }

  return (
    <span className={sizeClass}>
      {gemIcons[type]}
    </span>
  )
}

export default function DevelopmentCard({ card, onPurchase, onReserve, disabled, canAfford = true, canReserve = true }: DevelopmentCardProps) {
  console.log(`[DevelopmentCard ${card.id}] Props:`, { disabled, canAfford, canReserve, hasOnPurchase: !!onPurchase })

  const theme = gemThemes[card.gem_type]
  const tierStyles = {
    1: 'border-2',
    2: 'border-[3px] shadow-xl',
    3: 'border-4 shadow-2xl'
  }

  return (
    <motion.div
      whileHover={{ scale: disabled ? 1 : 1.03 }}
      className={`relative ${theme.bg} rounded-xl ${tierStyles[card.tier]} ${theme.border} overflow-hidden flex flex-col ${
        disabled ? 'opacity-60' : ''
      }`}
      style={{
        minHeight: '200px',
        maxHeight: '200px',
        boxShadow: disabled ? undefined : '0 4px 8px rgba(0,0,0,0.15)'
      }}
    >
      {/* 装饰性图案 */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-0 right-0 w-20 h-20 bg-white rounded-full -translate-y-10 translate-x-10" />
        <div className="absolute bottom-0 left-0 w-16 h-16 bg-white rounded-full translate-y-8 -translate-x-8" />
      </div>

      {/* Victory Points Badge */}
      {card.victory_points > 0 && (
        <div className="absolute top-1.5 right-1.5 bg-gradient-to-br from-yellow-300 to-yellow-500 text-gray-900 font-black rounded-full w-8 h-8 flex items-center justify-center text-base shadow-md border-2 border-yellow-200 z-10">
          {card.victory_points}
        </div>
      )}

      {/* Gem Icon - Centered at very top, not overlapping white area */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 z-20" style={{ filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.3))' }}>
        <GemIcon type={card.gem_type} size="large" />
      </div>

      {/* Card Content */}
      <div className="relative flex flex-col h-full p-2 pt-7">
        {/* Cost Display - Fixed 130px to guarantee all 4 cost items visible */}
        <div className="flex-none h-[130px]">
          <div className="bg-white/95 backdrop-blur-sm rounded-lg p-1.5 shadow-inner border border-white/60 h-full overflow-hidden">
            <div className="space-y-0.5 h-full overflow-y-auto scrollbar-thin scrollbar-thumb-gray-400 scrollbar-track-transparent pr-0.5">
              {Object.entries(card.cost)
                .filter(([, cost]) => cost > 0)
                .map(([gemType, cost]) => (
                  <div key={gemType} className="flex items-center justify-between bg-white/80 rounded px-1.5 py-0.5">
                    <div className="flex items-center gap-1">
                      <GemIcon type={gemType} size="medium" />
                      <span className="capitalize text-[9px] font-semibold text-gray-800">
                        {gemType}
                      </span>
                    </div>
                    <span className="font-bold text-[10px] text-gray-900 bg-white px-1.5 py-0.5 rounded border border-gray-300 min-w-[20px] text-center shadow-sm">
                      {cost}
                    </span>
                  </div>
                ))}
            </div>
          </div>
        </div>

        {/* Bottom: Action Buttons - Fixed at Bottom */}
        {!disabled && (onPurchase || onReserve) && (
          <div className="flex-none mt-1 flex gap-1">
            {onPurchase && (
              <button
                onClick={canAfford ? onPurchase : undefined}
                disabled={!canAfford}
                className={`flex-1 font-bold text-[10px] py-1.5 px-2 rounded-lg transition-all duration-200 border ${
                  canAfford
                    ? `${theme.accent} text-white border-transparent hover:scale-105 hover:shadow-md active:scale-95 cursor-pointer`
                    : 'bg-gray-400 text-gray-200 border-gray-300 cursor-not-allowed opacity-50 pointer-events-none'
                }`}
              >
                Buy
              </button>
            )}
            {onReserve && (
              <button
                onClick={canReserve ? onReserve : undefined}
                disabled={!canReserve}
                className={`flex-1 font-bold text-[10px] py-1.5 px-2 rounded-lg transition-all duration-200 border ${
                  canReserve
                    ? 'bg-amber-600 text-white border-transparent hover:scale-105 hover:shadow-md active:scale-95 cursor-pointer'
                    : 'bg-gray-400 text-gray-200 border-gray-300 cursor-not-allowed opacity-50 pointer-events-none'
                }`}
              >
                Reserve
              </button>
            )}
          </div>
        )}
      </div>
    </motion.div>
  )
}
