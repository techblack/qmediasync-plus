import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [vue(), Components({ resolvers: [AntDesignVueResolver({ importStyle: false })] })],
  root: '.',
  server: {
    proxy: {
      '/api': process.env.VITE_API_PROXY || 'http://127.0.0.1:12333',
    },
  },
  build: { outDir: '../web_statics', emptyOutDir: false }
})
