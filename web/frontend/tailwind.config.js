/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Sedai.io-inspired color palette
        primary: {
          50: '#f5f3ff',
          100: '#ede9fe',
          200: '#ddd6fe',
          300: '#c4b5fd',
          400: '#a78bfa',
          500: '#8A43FF', // Primary purple
          600: '#7c3aed',
          700: '#6d28d9',
          800: '#5b21b6',
          900: '#4c1d95',
        },
        accent: {
          cyan: '#41E9E0',
          teal: '#26BDFF',
          magenta: '#FF83D1',
          pink: '#FF18DB',
          blue: '#7376FF',
        },
        background: {
          navy: '#0A0E29',
          slate: '#2A3947',
          dark: '#1a1f3a',
          card: '#1e2438',
        },
        text: {
          primary: '#FFFFFF',
          secondary: '#BABABA',
          muted: '#6B7280',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        heading: ['Outfit', 'Inter', 'system-ui', 'sans-serif'],
      },
      backgroundImage: {
        'gradient-purple': 'linear-gradient(135deg, #8A43FF 0%, #FF18DB 100%)',
        'gradient-cyan': 'linear-gradient(135deg, #41E9E0 0%, #26BDFF 100%)',
        'gradient-radial': 'radial-gradient(circle, var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
      },
      boxShadow: {
        'glow-purple': '0 0 20px rgba(138, 67, 255, 0.4)',
        'glow-cyan': '0 0 20px rgba(65, 233, 224, 0.4)',
        'card': '0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -1px rgba(0, 0, 0, 0.2)',
        'card-hover': '0 10px 15px -3px rgba(138, 67, 255, 0.3), 0 4px 6px -2px rgba(138, 67, 255, 0.2)',
      },
      borderRadius: {
        'pill': '36px',
      },
      animation: {
        'gradient': 'gradient 6s ease infinite',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'shimmer': 'shimmer 2s linear infinite',
      },
      keyframes: {
        gradient: {
          '0%, 100%': {
            'background-size': '200% 200%',
            'background-position': 'left center'
          },
          '50%': {
            'background-size': '200% 200%',
            'background-position': 'right center'
          },
        },
        shimmer: {
          '0%': { transform: 'translateX(-100%)' },
          '100%': { transform: 'translateX(100%)' },
        },
      },
    },
  },
  plugins: [],
}
