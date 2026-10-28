import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { PendingButtonWrapper } from "../../utils/PendingButton/test-components";

describe("PendingButton", () => {
  test("displays normal text when pending is false", async () => {
    const screen = await render(<PendingButtonWrapper pending={false} />);

    const text = screen.getByText("Submit");
    await expect.element(text).toBeInTheDocument();

    const pendingText = screen.getByText("Adding...");
    await expect.element(pendingText).not.toBeInTheDocument();
  });

  test("displays pending text when pending is true", async () => {
    const screen = await render(<PendingButtonWrapper pending={true} />);

    const text = screen.getByText("Submit");
    await expect.element(text).not.toBeInTheDocument();

    const pendingText = screen.getByText("Adding...");
    await expect.element(pendingText).toBeInTheDocument();

    const spinner = screen.getByRole("img");
    await expect.element(spinner).toBeInTheDocument();
  });
});
