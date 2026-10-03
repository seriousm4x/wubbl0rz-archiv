import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import { Features } from 'lightningcss';

export default defineConfig({
	css: {
		lightningcss: { exclude: Features.LightDark | Features.DirSelector }
	},
	plugins: [tailwindcss(), sveltekit()]
});
