import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import { clickButton } from "../../utils/button";
import { ConfirmationDialogWrapper } from "../../utils/ConfirmationDialog/test-components";

async function clickBackdrop(screen: RenderResult) {
  const backdrop = screen.getByRole("region", { name: "Confirmation dialog backdrop" });

  // click() clicks the center of the locator by default which is the dialog in this case so explicitly set the click
  // coordinates to hit the backdrop
  await backdrop.click({ position: { x: 0, y: 0 } });
}

describe("ConfirmationDialog", () => {
  let mockOnAccept: Mock;
  let mockOnCancel: Mock;

  beforeEach(() => {
    mockOnAccept = vi.fn();
    mockOnCancel = vi.fn();
  });

  test("displays dialog with title, description, and cancel and accept buttons", async () => {
    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);

    const title = screen.getByText("Add item");
    await expect.element(title).toBeInTheDocument();

    const description = screen.getByText("Do you want to add this item?");
    await expect.element(description).toBeInTheDocument();

    const cancelButton = screen.getByRole("button", { name: "Cancel" });
    await expect.element(cancelButton).toBeInTheDocument();

    const acceptButton = screen.getByRole("button", { name: "Accept" });
    await expect.element(acceptButton).toBeInTheDocument();
  });

  test("invokes onAccept callback when accept button is clicked", async () => {
    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);
    await clickButton(screen, "Accept");
    expect(mockOnAccept).toHaveBeenCalledOnce();
  });

  test("disables cancel and accept buttons and cancel on backdrop click when accept is in progress", async () => {
    // Have onAccept sleep for a few seconds so we can check buttons are disabled
    mockOnAccept.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);

    await clickButton(screen, "Accept");

    const cancelButton = screen.getByRole("button", { name: "Cancel" });
    await expect.element(cancelButton).toBeDisabled();

    const acceptButton = screen.getByRole("button", { name: "Accept" });
    await expect.element(acceptButton).toBeDisabled();

    await clickBackdrop(screen);

    expect(mockOnCancel).not.toHaveBeenCalled();
  });

  test("displays error when accept fails", async () => {
    mockOnAccept.mockRejectedValue(new Error("Failed to submit form"));

    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);

    await clickButton(screen, "Accept");

    const error = screen.getByText("Something went wrong. Please try again.");
    await expect.element(error).toBeInTheDocument();
  });

  test("invokes onCancel callback when cancel button is clicked", async () => {
    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);

    const cancelButton = screen.getByRole("button", { name: "Cancel" });
    await cancelButton.click();

    expect(mockOnCancel).toHaveBeenCalledOnce();
  });

  test("invokes onCancel callback when backdrop is clicked", async () => {
    const screen = await render(<ConfirmationDialogWrapper onAccept={mockOnAccept} onCancel={mockOnCancel} />);
    await clickBackdrop(screen);
    expect(mockOnCancel).toHaveBeenCalledOnce();
  });
});
