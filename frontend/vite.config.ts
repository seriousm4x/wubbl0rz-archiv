import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import { Features } from 'lightningcss';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import adapter from '@sveltejs/adapter-node';

export default defineConfig({
	css: {
		lightningcss: { exclude: Features.LightDark | Features.DirSelector }
	},
	plugins: [
		tailwindcss(),
		sveltekit({
			preprocess: vitePreprocess(),
			adapter: adapter(),
			env: {
				dir: '..'
			}
		})
	]
});
