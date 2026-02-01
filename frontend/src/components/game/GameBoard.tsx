import { motion } from 'framer-motion'
import type { FullGameState, DevelopmentCard as CardType } from '../../types'
import GemBank from './GemBank'
import CardTiers from './CardTiers'
import NobleDisplay from './NobleDisplay'
import PlayerPanels from './PlayerPanels'
import ActionPanel from './ActionPanel'
import ReservedCardsPanel from './ReservedCardsPanel'
import GameRules from './GameRules'

interface GameBoardProps {
  gameState: FullGameState
  currentUserId: number
  onTakeGems: (gems: Record<string, number>) => void
  onPurchaseCard: (cardId: number, fromReserve: boolean) => void
  onReserveCard: (cardId: number, tier: number) => void
}

export default function GameBoard({
  gameState,
  currentUserId,
  onTakeGems,
  onPurchaseCard,
  onReserveCard
}: GameBoardProps) {
  const isMyTurn = gameState.game.current_turn_player_id === currentUserId
  const myPlayerState = gameState.player_states[currentUserId]

  // Find a card by ID across all tiers and reserves
  const findCard = (cardId: number): CardType | null => {
    const allCards = [
      ...gameState.game_state.visible_cards_tier1,
      ...gameState.game_state.visible_cards_tier2,
      ...gameState.game_state.visible_cards_tier3,
      ...(myPlayerState?.reserved_cards || [])
    ]
    return allCards.find(c => c.id === cardId) || null
  }

  // Calculate if player can afford a card
  const canAffordCard = (cardCost: Record<string, number>): boolean => {
    if (!myPlayerState) {
      console.log('[canAffordCard] No player state')
      return false
    }

    let goldNeeded = 0
    const breakdown: Record<string, { cost: number; permanent: number; owned: number; total: number; shortage: number }> = {}

    for (const [gemType, cost] of Object.entries(cardCost)) {
      if (cost > 0) {
        const permanent = myPlayerState.permanent_gems[gemType] || 0
        const owned = myPlayerState.gems[gemType] || 0
        const total = permanent + owned
        const shortage = Math.max(0, cost - total)

        breakdown[gemType] = { cost, permanent, owned, total, shortage }

        if (total < cost) {
          goldNeeded += cost - total
        }
      }
    }

    const availableGold = myPlayerState.gems.gold || 0
    const canAfford = goldNeeded <= availableGold

    console.log('[canAffordCard]', {
      cardCost,
      breakdown,
      goldNeeded,
      availableGold,
      canAfford
    })

    return canAfford
  }

  // Check if player can reserve more cards
  const canReserveMore = myPlayerState ? myPlayerState.reserved_cards.length < 3 : false

  // Validated purchase handler
  const handlePurchaseCard = (cardId: number, fromReserve: boolean) => {
    if (!isMyTurn) {
      alert('It is not your turn!')
      return
    }

    const card = findCard(cardId)
    if (!card) {
      alert('Card not found!')
      return
    }

    if (!canAffordCard(card.cost)) {
      alert('You cannot afford this card!')
      return
    }

    onPurchaseCard(cardId, fromReserve)
  }

  // Validated reserve handler
  const handleReserveCard = (cardId: number, tier: number) => {
    if (!isMyTurn) {
      alert('It is not your turn!')
      return
    }

    if (!canReserveMore) {
      alert('You already have 3 reserved cards!')
      return
    }

    onReserveCard(cardId, tier)
  }

  return (
    <div className="space-y-2">
      {/* Top Bar: Turn Indicator */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex items-center justify-center gap-3"
      >
        <div className={`px-4 py-2 rounded-lg font-bold text-sm shadow-md ${
          isMyTurn
            ? 'bg-gradient-to-r from-green-500 to-emerald-600 text-white animate-pulse'
            : 'bg-gradient-to-r from-gray-600 to-gray-700 text-gray-300'
        }`}>
          {isMyTurn ? "🎮 YOUR TURN" : "⏳ WAITING..."}
        </div>
        <div className="text-white/70 text-xs font-medium">
          Turn {gameState.game.turn_number}
        </div>
      </motion.div>

      {/* Main Game Area */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-2">
        {/* Left Sidebar: Gem Bank & Nobles */}
        <div className="lg:col-span-3 flex flex-col gap-2" style={{ height: '664px' }}>
          <div className="flex-none">
            <GemBank gems={gameState.game_state.available_gems} />
          </div>
          <div className="flex-1 min-h-0">
            <NobleDisplay nobles={gameState.game_state.available_nobles} />
          </div>
        </div>

        {/* Center: Card Tiers */}
        <div className="lg:col-span-9">
          <CardTiers
            tier1={gameState.game_state.visible_cards_tier1}
            tier2={gameState.game_state.visible_cards_tier2}
            tier3={gameState.game_state.visible_cards_tier3}
            deckCounts={{
              tier1: gameState.game_state.deck_tier1_count,
              tier2: gameState.game_state.deck_tier2_count,
              tier3: gameState.game_state.deck_tier3_count,
            }}
            onPurchase={handlePurchaseCard}
            onReserve={handleReserveCard}
            disabled={!isMyTurn}
            canAffordCard={canAffordCard}
            canReserveMore={canReserveMore}
          />
        </div>
      </div>

      {/* Player Panels */}
      <PlayerPanels
        players={gameState.players}
        playerStates={gameState.player_states}
        currentUserId={currentUserId}
        currentTurnUserId={gameState.game.current_turn_player_id || 0}
      />

      {/* Reserved Cards Panel */}
      {myPlayerState && myPlayerState.reserved_cards.length > 0 && (
        <ReservedCardsPanel
          reservedCards={myPlayerState.reserved_cards}
          onPurchase={(cardId) => handlePurchaseCard(cardId, true)}
          disabled={!isMyTurn}
          canAffordCard={canAffordCard}
        />
      )}

      {/* Action Panel (when it's your turn) */}
      {isMyTurn && myPlayerState && (
        <ActionPanel
          playerState={myPlayerState}
          availableGems={gameState.game_state.available_gems}
          onTakeGems={onTakeGems}
        />
      )}

      {/* Game Rules - Always visible at bottom */}
      <GameRules />
    </div>
  )
}
