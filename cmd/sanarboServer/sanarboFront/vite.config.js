import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
import vuetify from 'vite-plugin-vuetify'

import { fileURLToPath, URL } from 'node:url'
import fs from 'node:fs'

const actualRoot = fs.realpathSync(fileURLToPath(new URL('.', import.meta.url)))

// https://vitejs.dev/config/
export default defineConfig({
	root: actualRoot,

	build: {
		chunkSizeWarningLimit: '750kb',
		rollupOptions: {
			output: {
				manualChunks: {
					'mdi-icons': ['@mdi/font/css/materialdesignicons.css'],
				}
			}
		}
	},
  plugins: [
		vue(),
		vuetify({ autoImport: true }),
	],
	define: { "process.env": {} },
	resolve: {
		alias: {
			'@': fs.realpathSync(fileURLToPath(new URL('./src', import.meta.url))),
		},
	},
	test: {
		global: true,
		environment: 'jsdom',
		deps: {
			inline: ['vuetify'],
		},
	},
})
