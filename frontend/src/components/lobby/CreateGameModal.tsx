import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

interface CreateGameModalProps {
  isOpen: boolean
  onClose: () => void
  onCreate: (numPlayers: number) => void
  loading: boolean
}

export default function CreateGameModal({ isOpen, onClose, onCreate, loading }: CreateGameModalProps) {
  const [numPlayers, setNumPlayers] = useState(4)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onCreate(numPlayers)
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            onClick={onClose}
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            className="relative card max-w-md w-full"
          >
            <h2 className="text-2xl font-bold text-white mb-6">Create New Game</h2>

            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                <label className="block text-white mb-4">Number of Players</label>
                <div className="grid grid-cols-3 gap-3">
                  {[2, 3, 4].map((num) => (
                    <button
                      key={num}
                      type="button"
                      onClick={() => setNumPlayers(num)}
                      className={`py-4 px-6 rounded-lg font-semibold text-lg transition-all ${
                        numPlayers === num
                          ? 'bg-purple-600 text-white shadow-lg scale-105'
                          : 'bg-white/10 text-white/70 hover:bg-white/20'
                      }`}
                    >
                      {num}
                    </button>
                  ))}
                </div>
              </div>

              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={onClose}
                  className="btn-secondary flex-1"
                  disabled={loading}
                >
                  Cancel
                </button>
                <button type="submit" className="btn-primary flex-1" disabled={loading}>
                  {loading ? 'Creating...' : 'Create Game'}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  )
}
