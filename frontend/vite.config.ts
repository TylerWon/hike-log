/// <reference types="vitest/config" />
import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), babel({ presets: [reactCompilerPreset()] }), tailwindcss()],
  server: {
    // Settings for Vite dev server to work with Docker Compose
    host: true,
    port: 5173,
    watch: {
      usePolling: true,
    },
  },
  test: {
    // Vitest Browser Mode settings
    browser: {
      enabled: true,
      headless: true,
      instances: [{ browser: "chromium" }],
      provider: playwright(),
    },
    // Vitest settings
    restoreMocks: true,
    unstubGlobals: true,
  },
});
