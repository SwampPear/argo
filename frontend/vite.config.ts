import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  resolve: {
    alias: {
      // (optional) if you don’t want the plugin:
      // '@wails': '/src/wailsjs',
      // '@go': '/src/wailsjs/go',
      // '@state': '/src/state',
      // '@components': '/src/components'
    }
  }
})
