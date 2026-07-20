import { test, expect } from '@playwright/test'

test.describe('Public catalog page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('shows page title and subtitle', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Sàn Sản phẩm số')
    await expect(page.locator('.catalog__subtitle')).toBeVisible()
  })

  test('displays all six categories as filter buttons', async ({ page }) => {
    const categoryBtns = page.locator('.catalog__category-grid .catalog__filter-btn')
    await expect(categoryBtns).toHaveCount(7) // "Tất cả" + 6 categories

    const expected = ['Tất cả', 'kiến trúc', 'cơ khí', 'điện tử', 'đồ họa', 'đồ án', 'luận văn']
    const texts = await categoryBtns.allTextContents()
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
    const ftTitle = page.locator('.product-card__title', { hasText: 'Mẫu vách CNC đồng tiền hiện đại' })
    await expect(ftTitle).toBeVisible()

    const ftTitle2 = page.locator('.product-card__title', { hasText: 'Mẫu vách cổng CNC cây nghệ thuật' })
    await expect(ftTitle2).toBeVisible()

    const ftCard = page.locator('.product-card').filter({ hasText: 'Mẫu vách CNC đồng tiền hiện đại' })
    await expect(ftCard.locator('img')).toBeAttached()
  })

  test('filters products by selected category via backend query', async ({ page }) => {
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Click "cơ khí" category filter
    await page.locator('.catalog__filter-btn', { hasText: 'cơ khí' }).click()

    // Wait for API response — only 2 products in "cơ khí"
    await expect(cards).toHaveCount(2)

    // Every visible card shows "cơ khí" as its category
    const visibleCategories = await page.locator('.product-card__category').allTextContents()
    for (const cat of visibleCategories) {
      expect(cat.trim()).toBe('cơ khí')
    }

    // Active category button class
    const activeBtn = page.locator('.catalog__filter-btn--active', { hasText: 'cơ khí' })
    await expect(activeBtn).toBeVisible()

    // Click "Tất cả" restores all products
    await page.locator('.catalog__filter-btn', { hasText: 'Tất cả' }).click()
    await expect(cards).toHaveCount(12)
  })

  test('clicking a product card navigates to its detail page', async ({ page }) => {
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Get the first card's link href before clicking
    const firstCard = cards.first()
    const href = await firstCard.getAttribute('href')
    expect(href).toMatch(/^\/san-pham\//)

    // Click the first product card
    await firstCard.click()

    // Should navigate to the detail page
    await expect(page).toHaveURL(new RegExp(href!))

    // Product detail page should show the product title
    await expect(page.locator('.product-detail__title')).toBeVisible()
  })
})

test.describe('Search and sort', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('search input is visible and filters by name', async ({ page }) => {
    const searchInput = page.locator('.catalog__search-input')
    await expect(searchInput).toBeVisible()

    await searchInput.fill('CNC')

    // Wait for debounced search + API response
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(2)

    const titles = await page.locator('.product-card__title').allTextContents()
    expect(titles.join(' ')).toContain('CNC')
  })

  test('search empty state shows reset suggestion', async ({ page }) => {
    const searchInput = page.locator('.catalog__search-input')
    await searchInput.fill('zzzznotfound')

    // Wait for empty state
    const empty = page.locator('.catalog__empty')
    await expect(empty).toBeVisible()
    await expect(empty).toContainText('Không tìm thấy')

    // Reset button should be visible
    await expect(page.locator('.catalog__reset-btn')).toBeVisible()
  })

  test('reset button clears search and filters', async ({ page }) => {
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Apply category filter
    await page.locator('.catalog__filter-btn', { hasText: 'cơ khí' }).click()
    await expect(cards).toHaveCount(2)

    // Click reset
    await page.locator('.catalog__reset-btn').click()
    await expect(cards).toHaveCount(12)
  })

  test('sort dropdown changes product order', async ({ page }) => {
    // Wait for products to load
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Select sort by price ascending
    await page.locator('.catalog__select[aria-label="Sắp xếp"]').selectOption('price_asc')

    // Wait for re-fetch and check that first card is free
    await expect(cards.first().locator('.product-card__price--free')).toBeVisible()
  })

  test('format filter narrows results', async ({ page }) => {
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(12)

    // Select "dxf" format
    await page.locator('.catalog__select[aria-label="Lọc theo định dạng"]').selectOption('dxf')

    // Should show only DXF products
    await expect(cards).toHaveCount(2)
  })

  test('product card shows description excerpt and format badges', async ({ page }) => {
    const firstCard = page.locator('.product-card').first()
    await expect(firstCard.locator('.product-card__description')).toBeVisible()
    // At least one format badge
    await expect(firstCard.locator('.product-card__format-badge').first()).toBeVisible()
  })

  test('URL is synced with search query', async ({ page }) => {
    await page.locator('.catalog__search-input').fill('CNC')

    // Wait for debounce + API
    const cards = page.locator('.product-card')
    await expect(cards).toHaveCount(2)

    // URL should contain query param
    await expect(page).toHaveURL(/q=CNC/)
  })
})

test.describe('Product detail page', () => {
  test('shows product details for a free product', async ({ page }) => {
    await page.goto('/san-pham/sp-001')

    // Product title should be visible
    await expect(page.locator('.product-detail__title')).toHaveText('Bản vẽ nhà phố 3 tầng')

    // Description should be visible
    await expect(page.locator('.product-detail__description')).toBeVisible()

    // Category should be visible
    await expect(page.locator('.product-detail__category')).toContainText('kiến trúc')

    // Free label should be visible
    await expect(page.locator('.product-detail__price')).toContainText('Miễn phí')

    // Formats should be visible
    const formats = page.locator('.product-detail__formats .format-badge')
    await expect(formats).toHaveCount(2)
  })

  test('shows paid product price in Xu', async ({ page }) => {
    await page.goto('/san-pham/sp-017')

    await expect(page.locator('.product-detail__title')).toContainText('Mẫu vách CNC')
    await expect(page.locator('.product-detail__price')).toContainText('100')
    await expect(page.locator('.product-detail__price')).toContainText('Xu')
  })

  test('shows demo image for product', async ({ page }) => {
    await page.goto('/san-pham/sp-001')
    await expect(page.locator('.product-detail__image img')).toBeAttached()
  })

  test('shows 404 for non-existent product', async ({ page }) => {
    await page.goto('/san-pham/nonexistent')
    await expect(page.locator('.product-detail__not-found')).toBeVisible()
  })

  test('does not show detail for hidden product', async ({ page }) => {
    await page.goto('/san-pham/sp-010')
    await expect(page.locator('.product-detail__not-found')).toBeVisible()
  })

  test('shows star rating on detail page', async ({ page }) => {
    await page.goto('/san-pham/sp-001')
    await expect(page.locator('.product-detail__rating')).toContainText('★')
  })
})
