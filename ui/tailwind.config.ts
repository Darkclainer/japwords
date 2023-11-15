import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    fontFamily: {
      sans: ['Inter'],
    },
    fontSize: {
      sm: ['0.875rem', '1.25rem'],
      base: ['1rem', '1.5rem'],
      lg: ['1.125rem', '1.5rem'],
      xl: ['1.25rem', '1.75rem'],
      '2xl': ['1.5rem', '1.8125rem'],
      '3xl': ['1.875rem', '2.25rem'],
      '4xl': ['2.25rem', '2.5rem'],
      '5xl': ['3rem', '1'],
      '6xl': ['3.75rem', '1'],
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
      red: 'var(--red)',
      'light-red': 'var(--light-red)',
      'dark-red': 'var(--dark-red)',
      'error-red': 'var(--error-red)',
    },
    extend: {
      boxShadow: {
        'tab-inner': 'inset 2px -2px 4px 0 rgb(0 0 0 / 0.05)',
      },
      animation: {
        rspin: 'rspin 1s linear infinite',

        slideDownAndFade: 'slideDownAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)',
        slideLeftAndFade: 'slideLeftAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)',
        slideUpAndFade: 'slideUpAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)',
        slideRightAndFade: 'slideRightAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)',
      },
      keyframes: {
        rspin: {
          from: { transform: 'rotate(0deg)' },
          to: { transform: 'rotate(-360deg)' },
        },
        slideDownAndFade: {
          from: { opacity: '0', transform: 'translateY(-2px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        slideLeftAndFade: {
          from: { opacity: '0', transform: 'translateX(2px)' },
          to: { opacity: '1', transform: 'translateX(0)' },
        },
        slideUpAndFade: {
          from: { opacity: '0', transform: 'translateY(2px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        slideRightAndFade: {
          from: { opacity: '0', transform: 'translateX(-2px)' },
          to: { opacity: '1', transform: 'translateX(0)' },
        },
      },
    },
  },
  plugins: [
    require('tailwindcss-radix')({
      variantPrefix: 'radix',
    }),
  ],
} satisfies Config;
