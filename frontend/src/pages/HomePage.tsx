import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useAuth } from '../context/AuthContext'

export default function HomePage() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="text-center"
      >
        <h1 className="text-6xl font-bold text-white mb-4 text-glow">
          SPLENDOR
        </h1>
        <p className="text-xl text-white/80 mb-12">
          A Renaissance merchant trading game
        </p>

        {user ? (
          <>
            <p className="text-white mb-6">
              Welcome back, <span className="font-bold text-purple-400">{user.username}</span>!
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link to="/lobby">
                <button className="btn-primary w-full sm:w-auto">
                  Go to Lobby
                </button>
              </Link>
              <Link to="/stats">
                <button className="btn-secondary w-full sm:w-auto">
                  My Stats
                </button>
              </Link>
              <button onClick={handleLogout} className="btn-secondary w-full sm:w-auto">
                Logout
              </button>
            </div>
          </>
        ) : (
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link to="/login">
              <button className="btn-primary w-full sm:w-auto">
                Login
              </button>
            </Link>
            <Link to="/register">
              <button className="btn-secondary w-full sm:w-auto">
                Register
              </button>
            </Link>
          </div>
        )}

        <div className="mt-12 text-white/60">
          <p className="font-semibold">✅ Phase 1-2 Complete:</p>
          <p className="text-sm mt-2">
            Project Setup ✓ | Database ✓ | Authentication ✓
          </p>
          <p className="text-xs mt-4">Next: Game Lobby & Real-time Communication</p>
        </div>
      </motion.div>
    </div>
  )
}
