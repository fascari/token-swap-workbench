import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ command }) => ({
  base: command === "serve" ? "/" : "/static/",
  plugins: [react()],
  root: "web/app",
  build: {
    outDir: "../static",
    emptyOutDir: true
  },
  server: {
    port: 5173,
    proxy: {
      "/health": "http://127.0.0.1:8080",
      "/v1": "http://127.0.0.1:8080"
    }
  }
}));
