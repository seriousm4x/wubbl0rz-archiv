import themes from 'daisyui/src/theming/themes';

/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/daisyui/**/*.{js,jsx,ts,tsx}'],
	theme: {
		extend: {}
	},
	plugins: [require('daisyui')],
	daisyui: {
		themes: [
			{
				light: {
					...themes['[data-theme=light]'],
					primary: '#9147ff',
					'primary-content': '#ffffff',
					secondary: '#1f1b27'
				}
			},
			{
				dark: {
					...themes['[data-theme=dark]'],
					primary: '#9147ff',
					'primary-content': '#ffffff',
					secondary: '#1f1b27',
					'base-100': '#1f1f23',
					'base-200': '#18181b',
					'base-300': '#0e0e10',
					'base-content': '#E5E5E5'
				}
			},
			'cupcake',
			'bumblebee',
			'emerald',
			'corporate',
			'synthwave',
			'retro',
			'cyberpunk',
			'valentine',
			'halloween',
			'garden',
			'forest',
			'aqua',
			'lofi',
			'pastel',
			'fantasy',
			'wireframe',
			'black',
			'luxury',
			'dracula',
			'cmyk',
			'autumn',
			'business',
			'acid',
			'lemonade',
			'night',
			'coffee',
			'winter'
		]
	},
	fontSize: {
		sm: '0.750rem',
		base: '1rem',
		xl: '1.333rem',
		'2xl': '1.777rem',
		'3xl': '2.369rem',
		'4xl': '3.158rem',
		'5xl': '4.210rem'
	},
	fontFamily: {
		heading: 'Poppins',
		body: 'Poppins'
	},
	fontWeight: {
		normal: '400',
		bold: '700'
	}
};
