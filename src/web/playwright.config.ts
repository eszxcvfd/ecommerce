import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  retries: 0,
  use: {
    baseURL: 'http://localhost:3000',
    headless: true,
  },
  webServer: [
    {
      command: 'go run ./cmd/test/',
      port: 8080,
      cwd: '../api',
      timeout: 15000,
      reuseExistingServer: false,
    },
    {
      command: 'npx nuxt dev --port 3000',
      port: 3000,
      cwd: '.',
      timeout: 60000,
      reuseExistingServer: false,
      env: {
        API_URL: 'http://localhost:8080',
      },
    },
  ],
})
