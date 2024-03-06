import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

// If is production

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        // Use .env to get ip and port
        target: 'http://10.0.2.249',
        changeOrigin: true,
      },

      '/ws': {
        target: 'ws://10.0.2.249',
        changeOrigin: true,
        ws: true,
      },
    },
  },
});
