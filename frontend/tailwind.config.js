/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        diamond: '#E8F4F8',
        sapphire: '#2E5EAA',
        emerald: '#2D8659',
        ruby: '#C41E3A',
        onyx: '#2B2B2B',
        gold: '#FFD700',
      },
      animation: {
        'gem-float': 'float 3s ease-in-out infinite',
        'card-flip': 'flip 0.6s ease-in-out',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        flip: {
          '0%': { transform: 'rotateY(0deg)' },
          '100%': { transform: 'rotateY(360deg)' },
        },
      },
    },
  },
  plugins: [],
}
