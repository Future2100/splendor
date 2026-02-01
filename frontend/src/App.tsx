import { Routes, Route } from 'react-router-dom'
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import LobbyPage from './pages/LobbyPage'
import GamePage from './pages/GamePage'
import StatsPage from './pages/StatsPage'

function App() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900">
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/lobby" element={<LobbyPage />} />
        <Route path="/game/:id" element={<GamePage />} />
        <Route path="/stats" element={<StatsPage />} />
      </Routes>
    </div>
  )
}

export default App
