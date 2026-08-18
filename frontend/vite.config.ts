/// <reference types="vitest/config" />
import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
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
    // Vitest Projects. Projects allow tests to run in different environments. Each folder listed will be treated as
    // their own Project and have their own configuration file.
    projects: ["src/tests/component", "src/tests/unit", "src/tests/visual"],

    // Vitest settings that apply to all projects
    restoreMocks: true,
    unstubGlobals: true,
  },
});
