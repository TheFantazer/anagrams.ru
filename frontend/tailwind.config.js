/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'bg-dark': '#0a0a0f',
                'accent': '#63e6be',
                'accent-hover': '#4ecdc4',
            },
            fontFamily: {
                'mono': ['"Space Mono"', 'monospace'],
                'sans': ['"Outfit"', 'sans-serif'],
            },
            borderRadius: {
                'card': '16px',
                'button': '12px',
            },
        },
    },
    plugins: [],
}