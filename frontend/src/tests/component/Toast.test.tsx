import { useState } from "react";
import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import Toast from "../../components/Toast";

/**
 * Wrapper around a Toast which allows it to appear/disappear.
 */
function StatefulToast({ durationMs }: { durationMs: number }) {
  const [showToast, setShowToast] = useState<boolean>(true);

  return <>{showToast && <Toast durationMs={durationMs} message="TESTING" onDismiss={() => setShowToast(false)} />}</>;
}

describe("Toast", () => {
  test("displays toast for 1 second then disappears", async () => {
    const screen = await render(<StatefulToast durationMs={1000} />);

    const toast = screen.getByRole("status", { name: "Toast" });
    await expect.element(toast).toBeInTheDocument();
    await expect.element(toast, { timeout: 2000 }).not.toBeInTheDocument();
  });

  test("hides toast when dismiss button is clicked", async () => {
    const screen = await render(<StatefulToast durationMs={10000} />);

    const toast = screen.getByRole("status", { name: "Toast" });
    await expect.element(toast).toBeInTheDocument();

    const dismissButton = screen.getByRole("button", { name: "Dismiss" });
    await dismissButton.click();

    await expect.element(toast).not.toBeInTheDocument();
  });
});
