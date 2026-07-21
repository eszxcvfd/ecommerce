import { test, expect } from '@playwright/test'

test.describe('Seller draft management', () => {
  const uniqueEmail = `seller_${Date.now()}@test.com`
  const password = 'testpass123'

  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('registration, seller activation, and draft creation flow', async ({ page }) => {
    // Register a new account
    await page.click('text=Đăng ký')
    await expect(page).toHaveURL(/\/dang-ky/)

    await page.fill('#ten', 'Người Bán E2E')
    await page.fill('#email', uniqueEmail)
    await page.fill('#password', password)
    await page.click('button[type="submit"]')

    // Should redirect to login page
    await expect(page).toHaveURL(/\/dang-nhap/)

    // Login
    await page.fill('#email', uniqueEmail)
    await page.fill('#password', password)
    await page.click('button[type="submit"]')

    // Should be on home page, authenticated
    await expect(page.locator('nav', { hasText: 'Người Bán E2E' })).toBeVisible()

    // Activate seller profile by calling the API directly (no UI for this yet)
    const apiUrl = process.env.API_URL || 'http://localhost:8080'
    const response = await page.request.post(`${apiUrl}/api/v1/ho-so-nguoi-ban`, {
      headers: { 'Content-Type': 'application/json' },
    })
    expect(response.status()).toBe(201)

    // Navigate to seller drafts page
    await page.click('text=Bản nháp')
    await expect(page).toHaveURL(/\/seller/)

    // Should see empty state
    await expect(page.locator('text=Chưa có bản nháp nào')).toBeVisible()

    // Navigate to create draft page
    await page.click('text=Tạo bản nháp mới')
    await expect(page).toHaveURL(/\/seller\/tao-moi/)

    // Fill in draft form
    await page.fill('#ten', 'Bản vẽ kiến trúc E2E Test')
    await page.fill('#mo_ta', 'Mô tả ngắn cho bản vẽ kiểm thử')
    await page.fill('#mo_ta_chi_tiet', 'Mô tả chi tiết đầy đủ hơn cho bản vẽ kiểm thử E2E')
    await page.selectOption('#danh_muc', 'kiến trúc')
    await page.fill('#so_xu', '15000')
    await page.fill('#giay_phep', 'Giấy phép tiêu chuẩn')
    await page.fill('#anh_demo', 'https://example.com/demo.jpg')

    // Add file entries
    await page.click('text=+ Thêm tệp')
    await page.fill('input[id^="tep_ten_0"]', 'facade.dwg')
    await page.selectOption('select[id^="tep_dd_0"]', 'dwg')
    await page.fill('input[id^="tep_size_0"]', '3072000')

    await page.click('text=+ Thêm tệp')
    await page.fill('input[id^="tep_ten_1"]', 'model.skp')
    await page.selectOption('select[id^="tep_dd_1"]', 'skp')
    await page.fill('input[id^="tep_size_1"]', '5120000')

    // Submit
    await page.click('button[type="submit"]')

    // Should show success message and redirect to drafts list
    await expect(page.locator('text=Đã tạo bản nháp thành công')).toBeVisible()
    await expect(page).toHaveURL(/\/seller/, { timeout: 5000 })

    // Drafts list should show our draft
    await expect(page.locator('text=Bản vẽ kiến trúc E2E Test')).toBeVisible()
    await expect(page.locator('text=15,000 Xu')).toBeVisible()
  })

  test('unauthenticated user cannot access seller pages', async ({ page }) => {
    // Try to access seller API directly
    const apiUrl = process.env.API_URL || 'http://localhost:8080'

    const response = await page.request.post(`${apiUrl}/api/v1/seller/san-pham`, {
      headers: { 'Content-Type': 'application/json' },
      data: { ten: 'test' },
    })
    expect(response.status()).toBe(401)
  })

  test('seller without profile cannot create drafts', async ({ page }) => {
    // Register and login (without activating seller profile)
    const buyerEmail = `buyer_${Date.now()}@test.com`
    await page.click('text=Đăng ký')
    await expect(page).toHaveURL(/\/dang-ky/)

    await page.fill('#ten', 'Người Mua E2E')
    await page.fill('#email', buyerEmail)
    await page.fill('#password', password)
    await page.click('button[type="submit"]')

    await expect(page).toHaveURL(/\/dang-nhap/)

    await page.fill('#email', buyerEmail)
    await page.fill('#password', password)
    await page.click('button[type="submit"]')

    // Try to create a draft via API
    const apiUrl = process.env.API_URL || 'http://localhost:8080'
    const response = await page.request.post(`${apiUrl}/api/v1/seller/san-pham`, {
      headers: { 'Content-Type': 'application/json' },
      data: {
        ten: 'Test Product',
        danh_muc: 'kiến trúc',
        tep: [{ ten_tep: 'test.dwg', dinh_dang: 'dwg', dung_luong_bytes: 1000 }],
      },
    })
    expect(response.status()).toBe(403)
  })
})
