import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  root: "./frontend",
  server: {
    proxy: {
      "/backend": {
        target: "http://127.0.0.1:7827",
      },
    },
  },
});
