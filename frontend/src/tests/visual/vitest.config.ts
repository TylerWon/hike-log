import { defineBrowserModeProject } from "../configs/browser-mode-config";

// Visual test configuration
export default defineBrowserModeProject({
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
