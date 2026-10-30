import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import { clickButton } from "../../utils/button";
import { ConfirmationDialogWrapper } from "../../utils/ConfirmationDialog/test-components";

describe("ConfirmationDialog", () => {
  let mockOnAccept: Mock;
  let mockOnCancel: Mock;

  beforeEach(() => {
    mockOnAccept = vi.fn();
    mockOnCancel = vi.fn();
  });

  test("displays component", async () => {
    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);
    const component = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect(component).toMatchScreenshot();
  });

  test("displays spinner and disabled cancel and accept buttons while accept is in progress", async () => {
    mockOnAccept.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);

    await clickButton(screen, "Add");

    const component = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect(component).toMatchScreenshot();
  });

  test("displays error when accept fails", async () => {
    mockOnAccept.mockRejectedValue(new Error("Failed to submit form"));

    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);

    await clickButton(screen, "Add");

    const component = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect(component).toMatchScreenshot();
  });

  test("styles accept button differently when actionType === 'destructive'", async () => {
    const screen = await render(
      <ConfirmationDialogWrapper actionType="destructive" onAccept={mockOnAccept} onCancel={mockOnCancel} />,
    );
    const component = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect(component).toMatchScreenshot();
  });
});
