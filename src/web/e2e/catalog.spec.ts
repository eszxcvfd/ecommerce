import { test, expect } from '@playwright/test'

test.describe('Public catalog page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('shows page title and subtitle', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Sàn Sản phẩm số')
    await expect(page.locator('.catalog__subtitle')).toBeVisible()
  })

  test('displays all six categories', async ({ page }) => {
    const categories = page.locator('.catalog__category-card')
    await expect(categories).toHaveCount(6)

    const expected = ['kiến trúc', 'cơ khí', 'điện tử', 'đồ họa', 'đồ án', 'luận văn']
    const texts = await categories.allTextContents()
    for (const name of expected) {
      expect(texts.some(t => t.trim().includes(name))).toBeTruthy()
    }
  })

  test('shows product cards with required fields', async ({ page }) => {
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Each card must show name, price, category, and rating
    const firstCard = cards.first()
    await expect(firstCard.locator('.product-card__title')).toBeVisible()
    await expect(firstCard.locator('.product-card__category')).toBeVisible()
    await expect(firstCard.locator('.product-card__price')).toBeVisible()
    // Rating might be present or show "Chưa có đánh giá"
    await expect(firstCard.locator('.product-card__rating')).toBeVisible()
  })

  test('includes free and paid price labels', async ({ page }) => {
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)
    // Free products include "Miễn phí" label
    const freePrices = page.locator('.product-card__price--free')
    expect(await freePrices.count()).toBeGreaterThanOrEqual(2)
    // Paid products show "Xu" amount
    const allPrices = await page.locator('.product-card__price').allTextContents()
    const paidPrices = allPrices.filter(t => t.includes('Xu'))
    expect(paidPrices.length).toBeGreaterThanOrEqual(3)
  })

  test('does NOT show non-approved products', async ({ page }) => {
    const cardTitles = await page.locator('.product-card__title').allTextContents()
    expect(cardTitles).not.toContain('Bản nháp chưa duyệt')
    expect(cardTitles).not.toContain('Đang chờ duyệt')
  })

  test('shows filethietke.vn-backed products with thumbnails', async ({ page }) => {
    // Auto-retrying locators confirm title and thumbnail
    const ftTitle = page.locator('.product-card__title', { hasText: 'Mẫu vách CNC đồng tiền hiện đại' })
    await expect(ftTitle).toBeVisible()

    const ftTitle2 = page.locator('.product-card__title', { hasText: 'Mẫu vách cổng CNC cây nghệ thuật' })
    await expect(ftTitle2).toBeVisible()

    // Card with source-backed thumbnail renders an img element (avoid network assertion)
    const ftCard = page.locator('.product-card').filter({ hasText: 'Mẫu vách CNC đồng tiền hiện đại' })
    await expect(ftCard.locator('img')).toBeAttached()
  })

  test('filters products by selected category', async ({ page }) => {
    // All 12 approved products are initially visible
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Click "cơ khí" category filter
    await page.locator('.catalog__category-card', { hasText: 'cơ khí' }).click()

    // Only 2 products in "cơ khí" remain visible
    await expect(cards).toHaveCount(2)

    // Every visible card shows "cơ khí" as its category
    const visibleCategories = await page.locator('.product-card__category').allTextContents()
    for (const cat of visibleCategories) {
      expect(cat.trim()).toBe('cơ khí')
    }

    // Active category button exposes selected state
    await expect(page.locator('.catalog__category-card', { hasText: 'cơ khí' })).toHaveAttribute('aria-pressed', 'true')

    // Click "Tất cả" restores all products
    await page.locator('.catalog__filter-btn').click()
    await expect(cards).toHaveCount(12)
  })
})
