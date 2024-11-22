import { defineConfig } from "vite";
import laravel from "laravel-vite-plugin";
import react from "@vitejs/plugin-react";
// import vue from '@vitejs/plugin-vue';


export default defineConfig({
  plugins: [
    laravel({
      // input: ["resources/js/app.js", "resources/css/app.css"],
      // ssr: "resources/js/ssr.js",
      input: ["resources/js/app.tsx", "resources/css/app.css"],
      ssr: "resources/js/ssr.tsx",
      publicDirectory: "public",
      buildDirectory: "build",
      refresh: true,
    }),
    // vue({
    //   include: [/\.vue$/],
    // }),
    react({}),
  ],
  optimizeDeps: {
    force: true,
    esbuildOptions: {
      loader: {
        ".js": "jsx",
        ".ts": "tsx",
      },
    },
  },
});
