import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import * as child_process from 'node:child_process';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	// omit warning on build: https://svelte.dev/docs/svelte/compiler-warnings#state_referenced_locally
	onwarn: (warning, handler) => {
		if (warning.code === 'state_referenced_locally') return;
		handler(warning);
	},

	kit: {
		// adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
		// If your environment is not supported or you settled on a specific environment, switch out the adapter.
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		adapter: adapter(),
		version: {
			name: child_process.execSync('git rev-parse HEAD').toString().trim(),
			pollInterval: 10 * 1000 // 10 sec
		}
	}
};

export default config;
