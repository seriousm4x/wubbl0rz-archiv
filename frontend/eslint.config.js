import js from '@eslint/js';
import prettier from 'eslint-config-prettier';
import svelte from 'eslint-plugin-svelte';
import { defineConfig } from 'eslint/config';
import globals from 'globals';
import tseslint from 'typescript-eslint';

export default defineConfig(
	js.configs.recommended,
	...tseslint.configs.recommended,
	...svelte.configs['flat/recommended'],
	prettier,
	...svelte.configs['flat/prettier'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parserOptions: {
				parser: tseslint.parser
			}
		}
	},
	{
		ignores: [
			'build/',
			'.svelte-kit/',
			'dist/',
			'package/',
			'node_modules/',
			'.DS_Store',
			'.env',
			'.env.*',
			'!.env.example',
			'pnpm-lock.yaml',
			'package-lock.json'
		]
	},
	{
		rules: {
			'svelte/no-navigation-without-resolve': [
				'error',
				{
					ignoreGoto: false,
					ignoreLinks: true,
					ignorePushState: false,
					ignoreReplaceState: false
				}
			]
		}
	}
);
