import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { PendingButtonWrapper } from "../../utils/PendingButton/test-components";

describe("PendingButton", () => {
  test("displays normal text when pending is false", async () => {
    const screen = await render(<PendingButtonWrapper pending={false} />);
    const button = screen.getByRole("button", { name: "Submit" });
    await expect(button).toMatchScreenshot();
  });

  test("displays pending text when pending is true", async () => {
    const screen = await render(<PendingButtonWrapper pending={true} />);
    const button = screen.getByRole("button", { name: "Submit" });
    await expect(button).toMatchScreenshot();
  });
});
