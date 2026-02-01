import DevelopmentCard from './DevelopmentCard'
import type { DevelopmentCard as CardType } from '../../types'

interface CardTiersProps {
  tier1: CardType[]
  tier2: CardType[]
  tier3: CardType[]
  deckCounts: { tier1: number; tier2: number; tier3: number }
  onPurchase: (cardId: number, fromReserve: boolean) => void
  onReserve: (cardId: number, tier: number) => void
  disabled: boolean
  canAffordCard: (cost: Record<string, number>) => boolean
  canReserveMore: boolean
}

// 牌堆背面卡片
const DeckCard = ({ tier, count }: { tier: number; count: number }) => {
  const tierColors = {
    1: 'from-green-600 to-emerald-700 border-green-400',
    2: 'from-blue-600 to-indigo-700 border-blue-400',
    3: 'from-purple-600 to-violet-700 border-purple-400'
  }

  return (
    <div className={`relative bg-gradient-to-br ${tierColors[tier as keyof typeof tierColors]} rounded-xl border-3 flex flex-col items-center justify-center shadow-lg w-full`}
         style={{ minHeight: '200px', maxHeight: '200px' }}>
      <div className="absolute inset-0 opacity-20">
        <div className="absolute inset-3 border-3 border-white rounded-lg" />
        <div className="absolute inset-6 border-2 border-white rounded-md" />
      </div>
      <div className="relative text-white text-center">
        <div className="text-5xl font-black mb-1">{count}</div>
        <div className="text-sm font-bold uppercase tracking-wide">Tier {tier}</div>
        <div className="text-[10px] opacity-80">Left</div>
      </div>
    </div>
  )
}

export default function CardTiers({
  tier1,
  tier2,
  tier3,
  deckCounts,
  onPurchase,
  onReserve,
  disabled,
  canAffordCard,
  canReserveMore
}: CardTiersProps) {
  const tiers = [
    { cards: tier3, tier: 3, deckCount: deckCounts.tier3, label: 'Tier 3', bgColor: 'bg-gradient-to-r from-purple-900/30 to-violet-900/30' },
    { cards: tier2, tier: 2, deckCount: deckCounts.tier2, label: 'Tier 2', bgColor: 'bg-gradient-to-r from-blue-900/30 to-indigo-900/30' },
    { cards: tier1, tier: 1, deckCount: deckCounts.tier1, label: 'Tier 1', bgColor: 'bg-gradient-to-r from-green-900/30 to-emerald-900/30' },
  ]

  return (
    <div className="space-y-2">
      {tiers.map(({ cards, tier, deckCount, bgColor }) => (
        <div key={tier} className={`${bgColor} rounded-xl p-2 backdrop-blur-sm border border-white/10`}>
          {/* Cards Row - Grid Layout with smaller deck card */}
          <div className="grid gap-2" style={{ gridTemplateColumns: '120px repeat(4, 1fr)' }}>
            {/* Deck Card */}
            <div>
              <DeckCard tier={tier} count={deckCount} />
            </div>

            {/* Visible Cards */}
            {cards.map((card) => {
              const affordability = !disabled && canAffordCard(card.cost)
              console.log(`[CardTiers] Card ${card.id} (tier ${tier}): disabled=${disabled}, canAfford=${affordability}`)
              return (
                <div key={card.id}>
                  <DevelopmentCard
                    card={card}
                    onPurchase={() => onPurchase(card.id, false)}
                    onReserve={() => onReserve(card.id, tier)}
                    disabled={disabled}
                    canAfford={affordability}
                    canReserve={!disabled && canReserveMore}
                  />
                </div>
              )
            })}
          </div>
        </div>
      ))}
    </div>
  )
}
