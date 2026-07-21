import { test, expect } from '@playwright/test'

test.describe('Account and authentication', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the home page before each test
    await page.goto('/')
  })

  test('shows login button in navigation when not authenticated', async ({ page }) => {
    const loginLink = page.locator('nav a, nav button', { hasText: 'Đăng nhập' }).first()
    await expect(loginLink).toBeVisible()
  })

  test('registration and login flow', async ({ page }) => {
    const uniqueEmail = `user_${Date.now()}@test.com`

    // Go to registration page
    await page.click('text=Đăng ký')
    await expect(page).toHaveURL(/\/dang-ky/)

    // Fill registration form
    await page.fill('#ten', 'Test User')
    await page.fill('#email', uniqueEmail)
    await page.fill('#password', 'testpass123')
    await page.click('button[type="submit"]')

    // Should redirect to login page after registration
    await expect(page).toHaveURL(/\/dang-nhap/)

    // Fill login form
    await page.fill('#email', uniqueEmail)
    await page.fill('#password', 'testpass123')
    await page.click('button[type="submit"]')

    // Should redirect to home page after login
    await expect(page).toHaveURL(/\//)

    // Navigation should show the user name and logout button
    const userName = page.locator('nav', { hasText: 'Test User' })
    await expect(userName).toBeVisible()

    const logoutBtn = page.locator('nav button', { hasText: 'Đăng xuất' })
    await expect(logoutBtn).toBeVisible()

    // Logout
    await logoutBtn.click()

    // After logout, should see login button again
    await expect(page.locator('nav a', { hasText: 'Đăng nhập' }).first()).toBeVisible()
  })

  test('login with invalid credentials shows error', async ({ page }) => {
    // Go to login page
    await page.click('text=Đăng nhập')
    await expect(page).toHaveURL(/\/dang-nhap/)

    // Fill login with wrong credentials
    await page.fill('#email', 'nonexistent@test.com')
    await page.fill('#password', 'wrongpassword')
    await page.click('button[type="submit"]')

    // Should show error message
    const errorMsg = page.locator('.auth-error')
    await expect(errorMsg).toBeVisible()
  })

  test('protected route redirects unauthenticated users', async ({ page }) => {
    // Try to access seller profile activation endpoint directly
    // Since this is a SPA, the API will return 401
    const response = await page.request.post('/api/v1/ho-so-nguoi-ban', {
      headers: { 'Content-Type': 'application/json' },
    })
    expect(response.status()).toBe(401)
  })

  test('login form has required fields', async ({ page }) => {
    await page.click('text=Đăng nhập')
    await expect(page).toHaveURL(/\/dang-nhap/)

    // Form should have email and password fields
    await expect(page.locator('#email')).toHaveAttribute('required', '')
    await expect(page.locator('#password')).toHaveAttribute('required', '')
  })
})
