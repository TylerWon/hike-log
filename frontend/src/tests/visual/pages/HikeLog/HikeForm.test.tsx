import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import joffreLakes1 from "../../../assets/images/joffre_lakes_1.avif";
import joffreLakes2 from "../../../assets/images/joffre_lakes_2.avif";
import { fillField } from "../../../utils/Field/test-helpers";
import { createFileFromUrl } from "../../../utils/files";
import { HikeFormWrapper } from "../../../utils/HikeForm/test-components";
import { clickSubmitButton, fillHikeForm } from "../../../utils/HikeForm/test-helpers";

describe("HikeForm", () => {
  let mockOnCancel: Mock;
  let mockOnSubmit: Mock;

  beforeEach(() => {
    mockOnCancel = vi.fn();
    mockOnSubmit = vi.fn();
  });

  test("displays empty form", async () => {
    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays completed form", async () => {
    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    const photos = await Promise.all([
      createFileFromUrl(joffreLakes1, "joffre_lakes_1.avif"),
      createFileFromUrl(joffreLakes2, "joffre_lakes_2.avif"),
    ]);
    await fillHikeForm(screen, photos);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays field errors when form is submitted with invalid values", async () => {
    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    // Fill in some fields correctly, some incorrectly, and leave some empty
    await fillField(screen, "Duration (minutes)", "240");
    await fillField(screen, "Difficulty (0-10)", "11.5");
    await fillField(screen, "AllTrails URL", "abc");

    await clickSubmitButton(screen);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays spinner and disabled cancel and submit buttons when submission is in progress", async () => {
    // Have onSubmit sleep for a few seconds so we can check buttons are disabled
    mockOnSubmit.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });

  test("displays error when submission fails", async () => {
    mockOnSubmit.mockRejectedValue(new Error("Failed to submit form"));

    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("region", { name: "Hike form backdrop" });
    await expect(form).toMatchScreenshot();
  });
});
