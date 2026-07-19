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
    await expect(cards).toHaveCount(6)

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
    await expect(cards).toHaveCount(6)

    // Products sp-001 and sp-003 are free ("Miễn phí")
    const freePrices = page.locator('.product-card__price--free')
    expect(await freePrices.count()).toBeGreaterThanOrEqual(2)

    // Products sp-002, sp-004, sp-006 are paid (show "Xu")
    const allPrices = await page.locator('.product-card__price').allTextContents()
    const paidPrices = allPrices.filter(t => t.includes('Xu'))
    expect(paidPrices.length).toBeGreaterThanOrEqual(3)
  })

  test('does NOT show non-approved products', async ({ page }) => {
    const cardTitles = await page.locator('.product-card__title').allTextContents()
    expect(cardTitles).not.toContain('Bản nháp chưa duyệt')
    expect(cardTitles).not.toContain('Đang chờ duyệt')
  })
})
