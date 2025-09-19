import prettier from 'eslint-config-prettier';
import { includeIgnoreFile } from '@eslint/compat';
import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import { fileURLToPath } from 'node:url';
import ts from 'typescript-eslint';
import svelteConfig from './svelte.config.js';

const gitignorePath = fileURLToPath(new URL('./.gitignore', import.meta.url));
const tsconfigRootDir = fileURLToPath(new URL('.', import.meta.url));

export default ts.config(
	{ ignores: ['api/**', 'src/lib/components/ui/**'] },
	includeIgnoreFile(gitignorePath),
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs.recommended,
	prettier,
	...svelte.configs.prettier,
	{
		languageOptions: {
			globals: { ...globals.browser, ...globals.node }
		},
		rules: {
			// typescript-eslint strongly recommend that you do not use the no-undef lint rule on TypeScript projects.
			// see: https://typescript-eslint.io/troubleshooting/faqs/eslint/#i-get-errors-from-the-no-undef-rule-about-global-variables-not-being-defined-even-though-there-are-no-typescript-errors
			'no-undef': 'off',
			'svelte/no-navigation-without-resolve': ['error', { ignoreGoto: true, ignoreLinks: true }]
		}
	},
	{
		files: ['src/**/*.svelte', 'src/**/*.svelte.ts', 'src/**/*.svelte.js'],
		languageOptions: {
			parserOptions: {
				project: [fileURLToPath(new URL('./tsconfig.eslint.json', import.meta.url))],
				tsconfigRootDir,
				extraFileExtensions: ['.svelte'],
				parser: ts.parser,
				svelteConfig
			}
		},
		rules: {
			'@typescript-eslint/no-deprecated': 'error'
		}
	},
	{
		files: ['src/**/*.{ts,tsx}'],
		languageOptions: {
			parserOptions: {
				project: [fileURLToPath(new URL('./tsconfig.eslint.json', import.meta.url))],
				tsconfigRootDir
			}
		},
		rules: {
			'@typescript-eslint/no-deprecated': 'error'
		}
	}
);
