import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter(),
		experimental: {
			remoteFunctions: true,
			tracing: {
				server: true
			},
			instrumentation: {
				server: true
			}
		},
		alias: {
			$api: 'api',
			$lib: 'src/lib',
			$src: 'src'
		}
	},
	compilerOptions: {
		runes: true,
		experimental: {
			async: true
		}
	}
};

export default config;
