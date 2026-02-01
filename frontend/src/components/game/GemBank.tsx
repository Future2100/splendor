import { motion } from 'framer-motion'

interface GemBankProps {
  gems: Record<string, number>
}

const gemData: Record<string, { icon: string; color: string; name: string }> = {
  diamond: { icon: '💎', color: 'from-blue-400 to-blue-600', name: 'Diamond' },
  sapphire: { icon: '🔷', color: 'from-blue-600 to-blue-800', name: 'Sapphire' },
  emerald: { icon: '💚', color: 'from-green-500 to-green-700', name: 'Emerald' },
  ruby: { icon: '❤️', color: 'from-red-500 to-red-700', name: 'Ruby' },
  onyx: { icon: '⚫', color: 'from-gray-600 to-gray-800', name: 'Onyx' },
  gold: { icon: '🪙', color: 'from-yellow-400 to-yellow-600', name: 'Gold' },
}

export default function GemBank({ gems }: GemBankProps) {
  const gemTypes = ['diamond', 'sapphire', 'emerald', 'ruby', 'onyx', 'gold']

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      className="card"
    >
      <h3 className="text-sm font-bold text-white mb-1.5 flex items-center gap-1">
        <span className="text-base">💰</span>
        <span>Gem Bank</span>
      </h3>
      <div className="grid grid-cols-2 gap-1.5">
        {gemTypes.map((gemType) => {
          const data = gemData[gemType]
          const count = gems[gemType] || 0

          return (
            <div
              key={gemType}
              className={`bg-gradient-to-br ${data.color} rounded-lg p-1.5 shadow-md border-2 border-white/20`}
            >
              <div className="text-center">
                <div className="text-xl mb-0.5">{data.icon}</div>
                <div className="text-lg font-bold text-white">{count}</div>
                <div className="text-[9px] text-white/90 font-semibold">
                  {data.name}
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </motion.div>
  )
}
