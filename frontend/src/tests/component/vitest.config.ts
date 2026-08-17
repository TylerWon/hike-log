import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";
import { defineProject } from "vitest/config";

// Component test configuration
export default defineProject({
  plugins: [react(), babel({ presets: [reactCompilerPreset()] }), tailwindcss()],
  test: {
    // Configure Vitest Browser Mode
    browser: {
      enabled: true,
      headless: true,
      instances: [{ browser: "chromium" }],
      locators: {
        exact: true,
      },
      provider: playwright(),
    },
    name: "Component Tests",
    setupFiles: ["./setup.ts"], // Runs before every test file
  },
});
