import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
import vuetify from 'vite-plugin-vuetify'

// https://vitejs.dev/config/
export default defineConfig({
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
	test: {
		global: true,
		environment: 'jsdom',
		deps: {
			inline: ['vuetify'],
		},
	},
})
