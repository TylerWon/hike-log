import { defineBrowserModeTestConfig } from "../configs/browser-mode-config";

// Configuration for unit tests that need to run in a browser environment
export default defineBrowserModeTestConfig({
  browser: {
    screenshotFailures: false,
  },
  name: "Unit Tests - Browser",
});
