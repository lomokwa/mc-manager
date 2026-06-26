import { defineConfig } from 'vite'

// In dev, proxy the API (REST + the /api/console WebSocket) to the Go server
// so the app can use same-origin requests and skip CORS entirely. Override
// the target with VITE_API_PROXY when the API runs elsewhere.
const target = process.env.VITE_API_PROXY ?? 'http://localhost:8080'

export default defineConfig({
  server: {
    proxy: {
      '/api': { target, changeOrigin: true, ws: true },
    },
  },
})
