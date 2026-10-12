import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import joffreLakes1 from "../../assets/images/joffre_lakes_1.avif";
import joffreLakes2 from "../../assets/images/joffre_lakes_2.avif";
import { HikeFormWrapper } from "../../utils/HikeForm/test-components";
import { clickSubmitButton, createPhotoFile, fillFieldInput, fillHikeForm } from "../../utils/HikeForm/test-helpers";

describe("HikeForm", () => {
  let mockOnClose: Mock;
  let mockOnSubmit: Mock;
  let mockOnSuccess: Mock;

  beforeEach(() => {
    mockOnClose = vi.fn();
    mockOnSubmit = vi.fn();
    mockOnSuccess = vi.fn();
  });

  test("displays empty form", async () => {
    const screen = await render(
      <HikeFormWrapper onClose={mockOnClose} onSubmit={mockOnSubmit} onSuccess={mockOnSuccess} />,
    );

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays completed form", async () => {
    const screen = await render(
      <HikeFormWrapper onClose={mockOnClose} onSubmit={mockOnSubmit} onSuccess={mockOnSuccess} />,
    );

    const photos = await Promise.all([
      createPhotoFile(joffreLakes1, "joffre_lakes_1.avif"),
      createPhotoFile(joffreLakes2, "joffre_lakes_2.avif"),
    ]);
    await fillHikeForm(screen, photos);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays field errors when form is submitted with invalid values", async () => {
    const screen = await render(
      <HikeFormWrapper onClose={mockOnClose} onSubmit={mockOnSubmit} onSuccess={mockOnSuccess} />,
    );

    // Fill in some fields correctly, some incorrectly, and leave some empty
    await fillFieldInput(screen, "Duration (minutes)", "240");
    await fillFieldInput(screen, "Difficulty (0-10)", "11.5");
    await fillFieldInput(screen, "AllTrails URL", "abc");

    await clickSubmitButton(screen);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays spinner and disabled cancel and submit buttons when submission is in progress", async () => {
    // Have onSubmit sleep for a few seconds so we can check buttons are disabled
    mockOnSubmit.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(
      <HikeFormWrapper onClose={mockOnClose} onSubmit={mockOnSubmit} onSuccess={mockOnSuccess} />,
    );

    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays error when submission fails", async () => {
    mockOnSubmit.mockRejectedValue(new Error("Failed to submit form"));

    const screen = await render(
      <HikeFormWrapper onClose={mockOnClose} onSubmit={mockOnSubmit} onSuccess={mockOnSuccess} />,
    );

    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });
});
