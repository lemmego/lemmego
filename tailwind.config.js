/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.{gohtml,templ,html,js}",
    "./resources/views/**/*.html",
    "./resources/js/**/*.vue",
    "./resources/js/**/*.jsx",
    "./resources/js/**/*.tsx",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
};
