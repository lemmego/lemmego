/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{tmpl,html,js}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms"), require("daisyui")],
};
