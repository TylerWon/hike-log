import type { UserWorkspaceConfig } from "vitest/config";

import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";
import path from "node:path";
import { defineProject } from "vitest/config";

type BrowserModeProjectArgs = {
  browser?: BrowserOptions;
  name: string;
};

type BrowserOptions = NonNullable<NonNullable<UserWorkspaceConfig["test"]>["browser"]>;

// Creates a Vitest Project config to run tests in Vitest Browser Mode.
export function defineBrowserModeProject({ browser, name }: BrowserModeProjectArgs) {
  return defineProject({
    plugins: [react(), babel({ presets: [reactCompilerPreset()] }), tailwindcss()],
    test: {
      browser: {
        enabled: true,
        headless: true,
        instances: [{ browser: "chromium", viewport: { height: 720, width: 1280 } }],
        locators: {
          exact: true,
        },
        provider: playwright(),
        ...browser,
      },
      name,
      // Construct an absolute path to the test setup file. setupFiles paths are otherwise relative to each project
      // config which causes issues when defineBrowserModeProject() is used in different projects.
      setupFiles: [path.join(import.meta.dirname, "browser-mode-test-setup.ts")],
    },
  });
}
