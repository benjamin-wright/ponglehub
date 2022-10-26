import { fileURLToPath } from 'url'
import { defineConfig } from 'vite'

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        appIndex: fileURLToPath(new URL('index.html', import.meta.url)),
      },
    },
  },
})
