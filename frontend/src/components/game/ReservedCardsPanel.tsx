import { motion } from 'framer-motion'
import DevelopmentCard from './DevelopmentCard'
import type { DevelopmentCard as CardType } from '../../types'

interface ReservedCardsPanelProps {
  reservedCards: CardType[]
  onPurchase: (cardId: number) => void
  disabled: boolean
  canAffordCard: (cost: Record<string, number>) => boolean
}

export default function ReservedCardsPanel({
  reservedCards,
  onPurchase,
  disabled,
  canAffordCard
}: ReservedCardsPanelProps) {
  if (reservedCards.length === 0) {
    return null
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="card border-l-4 border-yellow-500 p-3"
    >
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-base font-bold text-white">Your Reserved Cards</h3>
        <div className="text-white/60 text-xs">
          {reservedCards.length}/3
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        {reservedCards.map((card) => (
          <DevelopmentCard
            key={card.id}
            card={card}
            onPurchase={() => onPurchase(card.id)}
            disabled={disabled}
            canAfford={!disabled && canAffordCard(card.cost)}
          />
        ))}
      </div>
    </motion.div>
  )
}
