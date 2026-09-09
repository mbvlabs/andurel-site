import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import inertia from '@inertiajs/vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig(({ command }) => ({
  plugins: [react(), inertia({ ssr: true }), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./resources/js', import.meta.url)),
    },
  },
  base: '/assets/dist/',
  build: {
    manifest: "vite/manifest.json",
    assetsDir: "",
    outDir: 'assets/dist',
    rollupOptions: {
      input: 'resources/js/app.tsx',
    },
  },
  // CNB drops node_modules after the frontend build. cmd/ssr only has the
  // bundle, so production SSR must inline npm imports instead of resolving them.
  ssr: command === 'build' ? { noExternal: true } : {},
  server: {
    port: 5173,
    strictPort: true,
  },
}))
