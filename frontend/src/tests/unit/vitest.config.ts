import { defineProject } from "vitest/config";

// Unit test configuration
export default defineProject({
  test: {
    environment: "node",
    name: "Unit Tests",
  },
});
