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
      ignored: process.env.VITEST ? undefined : ["**/src/tests/**"], // prevent hot reload when test files are changed, unless Vitest is running
      usePolling: true,
    },
  },
  test: {
    // Vitest projects - allows tests to be run in different environments
    projects: [
      {
        test: {
          environment: "node",
          include: ["src/tests/**/*.test.ts"],
          name: "Unit Tests",
        },
      },
      {
        test: {
          browser: {
            enabled: true,
            headless: true,
            instances: [{ browser: "chromium" }],
            provider: playwright(),
          },
          include: ["src/tests/**/*.test.tsx"],
          name: "Component Tests",
        },
      },
    ],
    // Vitest settings (applies to all projects)
    restoreMocks: true,
    unstubGlobals: true,
  },
});
