import { defineProject, mergeConfig } from "vitest/config";

import { COMMON_TEST_CONFIG } from "./common-config";

type NodeTestConfigArgs = {
  name: string;
};

// Creates a Vitest Project config to run tests in Node.
export function defineNodeTestConfig({ name }: NodeTestConfigArgs) {
  return mergeConfig(
    COMMON_TEST_CONFIG,
    defineProject({
      test: {
        environment: "node",
        name,
      },
    }),
  );
}
