import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import Toast from "../../../components/Toast";
import { clickButton } from "../../utils/button";

describe("Toast", () => {
  test("displays toast for 1 second then invokes onDismiss callback", async () => {
    const onDismissMock = vi.fn();
    const screen = await render(<Toast durationMs={1000} message="TESTING" onDismiss={onDismissMock} />);

    const toast = screen.getByRole("status", { name: "Toast" });

    await expect.element(toast).toBeInTheDocument();
    await vi.waitFor(
      () => {
        expect(onDismissMock).toHaveBeenCalledOnce();
      },
      { timeout: 2000 },
    );
  });

  test("invokes onDismiss callback when dismiss button is clicked", async () => {
    const onDismissMock = vi.fn();
    const screen = await render(<Toast durationMs={1000} message="TESTING" onDismiss={onDismissMock} />);

    const toast = screen.getByRole("status", { name: "Toast" });
    await expect.element(toast).toBeInTheDocument();

    await clickButton(screen, "Dismiss");

    expect(onDismissMock).toHaveBeenCalledOnce();
  });
});
