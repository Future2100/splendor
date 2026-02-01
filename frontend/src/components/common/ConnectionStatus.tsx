import { motion, AnimatePresence } from 'framer-motion'

interface ConnectionStatusProps {
  isConnected: boolean
}

export default function ConnectionStatus({ isConnected }: ConnectionStatusProps) {
  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -20 }}
        className="fixed top-4 right-4 z-50"
      >
        <div
          className={`flex items-center gap-2 px-4 py-2 rounded-lg shadow-lg ${
            isConnected
              ? 'bg-green-500/20 border border-green-500'
              : 'bg-red-500/20 border border-red-500'
          }`}
        >
          <div
            className={`w-2 h-2 rounded-full ${
              isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
            }`}
          />
          <span className={`text-sm font-semibold ${isConnected ? 'text-green-200' : 'text-red-200'}`}>
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
        </div>
      </motion.div>
    </AnimatePresence>
  )
}
