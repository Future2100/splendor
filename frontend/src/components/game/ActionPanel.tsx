import { useState } from 'react'
import { motion } from 'framer-motion'
import type { PlayerState } from '../../types'

interface ActionPanelProps {
  playerState: PlayerState
  availableGems: Record<string, number>
  onTakeGems: (gems: Record<string, number>) => void
}

const gemIcons: Record<string, string> = {
  diamond: '💎',
  sapphire: '🔷',
  emerald: '💚',
  ruby: '❤️',
  onyx: '⚫',
}

const gemNames: Record<string, string> = {
  diamond: 'Diamond',
  sapphire: 'Sapphire',
  emerald: 'Emerald',
  ruby: 'Ruby',
  onyx: 'Onyx',
}

export default function ActionPanel({ playerState, availableGems, onTakeGems }: ActionPanelProps) {
  const [selectedGems, setSelectedGems] = useState<Record<string, number>>({})

  const handleGemClick = (gemType: string) => {
    const current = selectedGems[gemType] || 0
    const available = availableGems[gemType] || 0
    const totalSelected = Object.values(selectedGems).reduce((a, b) => a + b, 0)

    // Can take up to 2 of the same type
    if (current < available && current < 2 && totalSelected < 3) {
      setSelectedGems({ ...selectedGems, [gemType]: current + 1 })
    } else if (current > 0) {
      // Click again to deselect
      const newSelected = { ...selectedGems }
      delete newSelected[gemType]
      setSelectedGems(newSelected)
    }
  }

  const handleSubmit = () => {
    onTakeGems(selectedGems)
    setSelectedGems({})
  }

  const totalSelected = Object.values(selectedGems).reduce((a, b) => a + b, 0)
  const selectedTypes = Object.keys(selectedGems).length

  // Can take 3 different or 2 same (if 4+ available)
  const canSubmit = (totalSelected === 3 && selectedTypes === 3) ||
                     (totalSelected === 2 && selectedTypes === 1 && availableGems[Object.keys(selectedGems)[0]] >= 4)

  const totalPlayerGems = Object.values(playerState.gems).reduce((a, b) => a + b, 0)

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="card bg-green-900/40 border-2 border-green-500"
    >
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-xl font-bold text-white">🎯 Your Turn</h3>
        <div className="text-white/70 text-sm">
          Your gems: {totalPlayerGems}/10
        </div>
      </div>

      {/* Gem Selection */}
      <div className="grid grid-cols-5 gap-2 mb-4">
        {['diamond', 'sapphire', 'emerald', 'ruby', 'onyx'].map((gemType) => {
          const available = availableGems[gemType] || 0
          const selected = selectedGems[gemType] || 0
          const isDisabled = available === 0 || totalPlayerGems >= 10

          return (
            <button
              key={gemType}
              onClick={() => handleGemClick(gemType)}
              disabled={isDisabled}
              className={`relative rounded-xl p-3 text-center transition-all ${
                selected > 0
                  ? 'bg-green-500 ring-2 ring-green-300 shadow-lg scale-105'
                  : 'bg-white/10 hover:bg-white/20'
              } ${isDisabled ? 'opacity-30 cursor-not-allowed' : 'cursor-pointer hover:scale-105'}`}
            >
              <div className="text-3xl mb-1">{gemIcons[gemType]}</div>
              <div className="text-white text-xs font-semibold">{gemNames[gemType]}</div>
              <div className="text-white/60 text-[10px] mt-1">
                Bank: {available}
              </div>
              {selected > 0 && (
                <div className="absolute -top-2 -right-2 bg-yellow-400 text-black font-bold rounded-full w-6 h-6 flex items-center justify-center text-sm shadow-lg">
                  {selected}
                </div>
              )}
            </button>
          )
        })}
      </div>

      {/* Submit */}
      <div className="flex items-center justify-between bg-white/10 rounded-lg p-3">
        <div className="text-white">
          <div className="font-bold">Selected: {totalSelected}</div>
          {!canSubmit && totalSelected > 0 && (
            <div className="text-xs text-yellow-300">
              {totalSelected === 2 && availableGems[Object.keys(selectedGems)[0]] < 4
                ? '⚠️ Need 4+ in bank for 2 same'
                : '⚠️ Select 3 different or 2 same'}
            </div>
          )}
        </div>
        <button
          onClick={handleSubmit}
          disabled={!canSubmit}
          className="bg-green-500 hover:bg-green-600 disabled:bg-gray-500 disabled:cursor-not-allowed text-white font-bold py-2 px-6 rounded-lg transition-all shadow-lg disabled:opacity-50"
        >
          Take Gems
        </button>
      </div>
    </motion.div>
  )
}
