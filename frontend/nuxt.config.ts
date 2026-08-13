import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  app: {
    head: {
      title: 'KyokusuLib',
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      link: [
        { rel: 'icon', type: 'image/png', href: '/images/kyokusulib_logo_3.png' },
      ],
    },
  },
  devtools: { enabled: true },

  modules: ['@pinia/nuxt', '@vueuse/nuxt', '@nuxtjs/color-mode', '@nuxt/icon'],

  routeRules: {
    '/api/**': {
      proxy: (process.env.NUXT_API_INTERNAL_URL || 'http://localhost:8080') + '/api/**',
    },
    '/uploads/**': {
      proxy: (process.env.NUXT_API_INTERNAL_URL || 'http://localhost:8080') + '/uploads/**',
    },
  },

  colorMode: {
    preference: 'dark',
    fallback: 'dark',
    classSuffix: '',
    storageKey: 'theme',
  },

  icon: {
    clientBundle: {
      scan: true,
      includeCustomCollections: true,
    },
    serverBundle: 'local',
  },

  vite: {
    plugins: [tailwindcss()],
    optimizeDeps: {
      include: ['@vue/devtools-core', '@vue/devtools-kit'],
    },
    resolve: {
      dedupe: ['vue'],
    },
  },

  build: {
    transpile: ['@kyokusu-ui/vue'],
  },

  css: ['~/assets/css/global.css'],

  postcss: {
    plugins: {
      autoprefixer: {},
    },
  },

  sourcemap: {
    client: 'hidden',
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '',
    },
  },

  compatibilityDate: 'latest',
})
