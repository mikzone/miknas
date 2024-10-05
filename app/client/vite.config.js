import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: { transformAssetUrls }
    }),
    vueDevTools(),

    // @quasar/plugin-vite options list:
    // https://github.com/quasarframework/quasar/blob/dev/vite-plugin/index.d.ts
    quasar({
      sassVariables: 'src/quasar-variables.scss'
    })
  ],
  server: {
    proxy: {
      // /s/ 前缀的都代表server，也就是服务端的接口
      '/s/': {
        target: 'http://127.0.0.1:2020',
        secure: false, // 不合法的证书需要用这个跳过
        changeOrigin: false,
      },
    },
    host: '0.0.0.0',
    https: false,
    port: 8010,
  },
  build: {
    rollupOptions: {
      output: {
        // sourceMap: true,
        manualChunks(id) {
          const mapChunks = {
            '@quasar': 'quasar',
            'quasar': 'quasar',
            'axios': 'quasar',
            'pinia': 'vue',
            'vue': 'vue',
            'vue-router': 'vue',
          }
          if (id.includes('/node_modules/')) {
            for (const [chunkName, toChunkName] of Object.entries(
              mapChunks
            )) {
              if (id.includes(`/node_modules/${chunkName}/`)) {
                return toChunkName;
              }
            }
          }
        },
      },
    },
  },
})
