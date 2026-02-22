import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/go-api': {
        target: 'http://localhost:8080',
        rewrite: (path) => path.replace(/^\/go-api/, ''),
      },
      '/rust-api': {
        target: 'http://localhost:8081',
        rewrite: (path) => path.replace(/^\/rust-api/, ''),
      },
    },
  },
})
