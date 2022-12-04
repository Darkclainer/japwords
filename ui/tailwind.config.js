/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.{html,svelte}'],
	theme: {
		fontFamily: {
			sans: ['Inter']
		},
		colors: {
			gray: '#F4F4F4',
			'dark-gray': '#5C5C5C',
			blue: '#617DB6',
			black: '#000000',
			white: '#FFFFFF',
			green: '#77B661',
			'dark-green': '#35512B'
		},
		extend: {}
	},
	plugins: []
};
