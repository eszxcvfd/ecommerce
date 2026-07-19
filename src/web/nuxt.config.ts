// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  compatibilityDate: '2026-07-19',

  app: {
    head: {
      title: 'Sàn Sản phẩm số — Thiết kế và Kỹ thuật',
      meta: [
        { name: 'description', content: 'Marketplace Sản phẩm số thiết kế và kỹ thuật' }
      ]
    }
  },
  // Vite dev proxy to Go backend
  vite: {
    server: {
      proxy: {
        '/api': {
          target: process.env.API_URL || 'http://localhost:8080',
          changeOrigin: true
        }
      }
    }
  },

  // CSS files
  css: ['~/assets/styles/main.css'],

  // Type-check in dev and build
  typescript: {
    strict: true
  }
})
