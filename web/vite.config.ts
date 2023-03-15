import { VitePWA } from 'vite-plugin-pwa';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import svgrPlugin from 'vite-plugin-svgr';

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: 'dist/web',
  },
  plugins: [
    react(),
    svgrPlugin({
      svgrOptions: {
        icon: true,
        // ...svgr options (https://react-svgr.com/docs/options/)
      },
    }),
    VitePWA({
      registerType: 'autoUpdate',
      workbox: {
        globPatterns: [
          'index.html',
          'assets/*.{js,css,html,ico,png,svg,jpeg,jpg}',
          'locales/**/*.json',
        ],
        cacheId: 'shinpuru-v1',
        navigateFallbackDenylist: [/^\/api\/(auth|public)\/.*/, /^\/invite/],
        runtimeCaching: [
          {
            urlPattern: /\/api\/(?:v\d\/)?guilds\/\d+\/(members|\d+)/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'guild-members-cache',
              expiration: {
                maxAgeSeconds: 60 * 5,
              },
              cacheableResponse: {
                statuses: [200],
              },
            },
          },
          {
            urlPattern: /\/api\/me/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'self-cache',
              expiration: {
                maxAgeSeconds: 60 * 5,
              },
              cacheableResponse: {
                statuses: [200],
              },
            },
          },
          {
            urlPattern: /\/api\/guilds\/(?:\d+)\/(?:\d+)\/permissions\/allowed/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'permissions-cache',
              expiration: {
                maxAgeSeconds: 60 * 5,
              },
              cacheableResponse: {
                statuses: [200],
              },
            },
          },
          {
            urlPattern: /\/api\/allpermissions/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'all-permissions-cache',
              expiration: {
                maxAgeSeconds: 24 * 3600,
              },
              cacheableResponse: {
                statuses: [200],
              },
            },
          },
          {
            urlPattern: /\/api\/.*/,
            handler: 'NetworkOnly',
          },
          {
            urlPattern: /^https:\/\/cdn\.discordapp\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'discord-cdn-cache',
              expiration: {
                maxEntries: 5000,
                maxAgeSeconds: 60 * 60 * 24 * 365, // <== 365 days
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
          {
            urlPattern: /^https:\/\/fonts\.googleapis\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'google-fonts-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 365, // <== 365 days
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
          {
            urlPattern: /^https:\/\/fonts\.gstatic\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'gstatic-fonts-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 365, // <== 365 days
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
        ],
      },
      includeAssets: ['favicon.ico', 'logo192.png', 'logo512.png'],
      manifest: {
        short_name: 'shinpuru',
        name: 'shinpuru web interface',
        description: 'The web interface of the shinpuru Discord bot.',
        icons: [
          {
            src: 'favicon.ico',
            sizes: '64x64 32x32 24x24 16x16',
            type: 'image/x-icon',
          },
          {
            src: 'logo192.png',
            type: 'image/png',
            sizes: '192x192',
          },
          {
            src: 'logo512.png',
            type: 'image/png',
            sizes: '512x512',
          },
        ],
        start_url: '.',
        display: 'standalone',
        orientation: 'any',
        theme_color: '#000000',
        background_color: '#ffffff',
      },
    }),
  ],
});
