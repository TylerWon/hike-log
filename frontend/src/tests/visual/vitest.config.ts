import { defineBrowserModeTestConfig } from "../configs/browser-mode-config";

// Visual test configuration
export default defineBrowserModeTestConfig({
  browser: {
    expect: {
      toMatchScreenshot: {
        comparatorOptions: {
          allowedMismatchedPixelRatio: 0.01, // 1% of pixels can differ
        },
      },
    },
  },
  name: "Visual Tests",
});
