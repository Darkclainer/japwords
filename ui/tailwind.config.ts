import type { Config } from 'tailwindcss';

const defaultTheme = require('tailwindcss/defaultTheme');

export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    fontFamily: {
      sans: ['Inter', ...defaultTheme.fontFamily.sans],
      mono: ['"Roboto mono"', ...defaultTheme.fontFamily.mono],
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
      gray: 'rgb(var(--gray) / <alpha-value>)',
      'dark-gray': 'rgb(var(--dark-gray) / <alpha-value>)',
      'mid-gray': 'rgb(var(--mid-gray) / <alpha-value>)',
      blue: 'rgb(var(--blue) / <alpha-value>)',
      'dark-blue': 'rgb(var(--dark-blue) / <alpha-value>)',
      black: 'rgb(var(--black) / <alpha-value>)',
      white: 'rgb(var(--white) / <alpha-value>)',
      green: 'rgb(var(--green) / <alpha-value>)',
      'dark-green': 'rgb(var(--dark-green) / <alpha-value>)',
      red: 'rgb(var(--red) / <alpha-value>)',
      'light-red': 'rgb(var(--light-red) / <alpha-value>)',
      'dark-red': 'rgb(var(--dark-red) / <alpha-value>)',
      'error-red': 'rgb(var(--error-red) / <alpha-value>)',
      'warning-yellow': 'rgb(var(--warning-yellow) / <alpha-value>)',
    },
    extend: {
      boxShadow: {
        'tab-inner': 'inset 2px -2px 4px 0 rgb(0 0 0 / 0.05)',
      },
      animation: {
        rspin: 'rspin 1s linear infinite',
        reversePing: 'reversePing 250ms cubic-bezier(0, 0, 0.2, 1)',

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
        reversePing: {
          from: { opacity: '0', transform: 'scale(1.5)' },
          to: { opacity: '1', transform: 'scale(1)' },
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
