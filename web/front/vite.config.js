import { defineConfig } from "vite";
import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

const gateTarget = process.env.VITE_GATE_URL ?? "http://127.0.0.1:8080";

/** @type {import('vite').UserConfig} */
export default defineConfig({
  plugins: [sveltekit(), tailwindcss()],
  server: {
    port: 5173,
    strictPort: false,
    proxy: {
      "/auth": { target: gateTarget, changeOrigin: true },
      "/spaces": { target: gateTarget, changeOrigin: true },
      "/authorize": { target: gateTarget, changeOrigin: true },
      "/invites": { target: gateTarget, changeOrigin: true },
      "/agents": { target: gateTarget, changeOrigin: true },
      "/tasks": { target: gateTarget, changeOrigin: true },
      "/schedules": { target: gateTarget, changeOrigin: true },
      "/memory": { target: gateTarget, changeOrigin: true },
      "/mcp": { target: gateTarget, changeOrigin: true },
      "/skills": { target: gateTarget, changeOrigin: true },
      "/tools": { target: gateTarget, changeOrigin: true },
      "/identity": { target: gateTarget, changeOrigin: true },
      "/dashboard": { target: gateTarget, changeOrigin: true },
      "/kits": { target: gateTarget, changeOrigin: true },
      "/projects": { target: gateTarget, changeOrigin: true },
      "/files": { target: gateTarget, changeOrigin: true },
      "/health": { target: gateTarget, changeOrigin: true },
      "/metrics": { target: gateTarget, changeOrigin: true },
      "/ws": { target: gateTarget, changeOrigin: true, ws: true },
    },
  },
  optimizeDeps: {
    include: ["pdfjs-dist", "three"],
  },
  resolve: {
    alias: {
      $lib: path.resolve("./src/lib"),
    },
  },
});