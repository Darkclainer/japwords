/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    fontFamily: {
      sans: ['Inter'],
    },
    colors: {
      gray: 'var(--gray)',
      'dark-gray': 'var(--dark-gray)',
      'mid-gray': 'var(--mid-gray)',
      blue: 'var(--blue)',
      'dark-blue': 'var(--dark-blue)',
      black: 'var(--black)',
      white: 'var(--white)',
      green: 'var(--green)',
      'dark-green': 'var(--dark-green)',
    },
  },
  plugins: [],
};
