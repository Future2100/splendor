import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import DevelopmentCard from './DevelopmentCard'
import type { DevelopmentCard as CardType } from '../../types'

describe('DevelopmentCard', () => {
  it('should render all cost items even with 4+ gem types', () => {
    // 创建一个有4个cost项目的卡片（最复杂的情况）
    const card: CardType = {
      id: 1,
      tier: 1,
      gem_type: 'diamond',
      victory_points: 0,
      cost: {
        diamond: 1,
        sapphire: 1,
        emerald: 1,
        ruby: 3,
        onyx: 0,
        gold: 0
      }
    }

    render(
      <DevelopmentCard
        card={card}
        onPurchase={() => {}}
        onReserve={() => {}}
        canAfford={true}
        canReserve={true}
      />
    )

    // 验证所有4个cost项目都被渲染（首字母大写）
    expect(screen.getByText(/diamond/i)).toBeInTheDocument()
    expect(screen.getByText(/sapphire/i)).toBeInTheDocument()
    expect(screen.getByText(/emerald/i)).toBeInTheDocument()
    expect(screen.getByText(/ruby/i)).toBeInTheDocument()

    // 验证cost值都显示（使用getAllByText因为有多个"1"）
    const ones = screen.getAllByText('1')
    expect(ones.length).toBeGreaterThanOrEqual(3) // diamond, sapphire, emerald都是1
    expect(screen.getByText('3')).toBeInTheDocument() // ruby是3
  })

  it('should render Buy and Reserve buttons', () => {
    const card: CardType = {
      id: 1,
      tier: 1,
      gem_type: 'diamond',
      victory_points: 1,
      cost: {
        diamond: 3,
        sapphire: 0,
        emerald: 0,
        ruby: 0,
        onyx: 0,
        gold: 0
      }
    }

    render(
      <DevelopmentCard
        card={card}
        onPurchase={() => {}}
        onReserve={() => {}}
        canAfford={true}
        canReserve={true}
      />
    )

    expect(screen.getByText('Buy')).toBeInTheDocument()
    expect(screen.getByText('Reserve')).toBeInTheDocument()
  })

  it('should show victory points badge when > 0', () => {
    const card: CardType = {
      id: 1,
      tier: 1,
      gem_type: 'diamond',
      victory_points: 5,
      cost: {
        diamond: 3,
        sapphire: 0,
        emerald: 0,
        ruby: 0,
        onyx: 0,
        gold: 0
      }
    }

    render(
      <DevelopmentCard
        card={card}
        onPurchase={() => {}}
        canAfford={true}
      />
    )

    expect(screen.getByText('5')).toBeInTheDocument()
  })

  it('should disable Buy button when cannot afford', () => {
    const card: CardType = {
      id: 1,
      tier: 1,
      gem_type: 'diamond',
      victory_points: 0,
      cost: {
        diamond: 3,
        sapphire: 0,
        emerald: 0,
        ruby: 0,
        onyx: 0,
        gold: 0
      }
    }

    render(
      <DevelopmentCard
        card={card}
        onPurchase={() => {}}
        canAfford={false}
      />
    )

    const buyButton = screen.getByText('Buy')
    expect(buyButton).toBeDisabled()
  })

  it('should have consistent height of 200px', () => {
    const card: CardType = {
      id: 1,
      tier: 1,
      gem_type: 'diamond',
      victory_points: 0,
      cost: {
        diamond: 1,
        sapphire: 0,
        emerald: 0,
        ruby: 0,
        onyx: 0,
        gold: 0
      }
    }

    const { container } = render(
      <DevelopmentCard
        card={card}
        onPurchase={() => {}}
        canAfford={true}
      />
    )

    const cardElement = container.firstChild as HTMLElement
    const style = window.getComputedStyle(cardElement)
    expect(style.minHeight).toBe('200px')
    expect(style.maxHeight).toBe('200px')
  })
})
