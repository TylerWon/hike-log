import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import { MAX_PHOTO_SIZE } from "../../../../schemas/forms/photo";
import { HikeFormWrapper } from "../../../utils/HikeForm/test-components";
import {
  checkFieldError,
  checkNoFieldError,
  clickSubmitButton,
  fillFieldInput,
  fillHikeForm,
  uploadPhotos,
} from "../../../utils/HikeForm/test-helpers";

async function clickBackdrop(screen: RenderResult) {
  const backdrop = screen.getByRole("region", { name: "Hike form backdrop" });

  // click() clicks the center of the locator by default which is the form in this case so explicitly set the click
  // coordinates to hit the backdrop
  await backdrop.click({ position: { x: 0, y: 0 } });
}

describe("HikeForm", () => {
  let mockOnCancel: Mock;
  let mockOnSubmit: Mock;

  beforeEach(() => {
    mockOnCancel = vi.fn();
    mockOnSubmit = vi.fn();
  });

  test("invokes onCancel callback when backdrop is clicked", async () => {
    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
    await clickBackdrop(screen);
    expect(mockOnCancel).toHaveBeenCalledOnce();
  });

  test("invokes onCancel callback when close button is clicked", async () => {
    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    const closeButton = screen.getByRole("button", { name: "Close" });
    await closeButton.click();

    expect(mockOnCancel).toHaveBeenCalledOnce();
  });

  test("invokes onCancel callback when cancel button is clicked", async () => {
    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    const closeButton = screen.getByRole("button", { name: "Cancel" });
    await closeButton.click();

    expect(mockOnCancel).toHaveBeenCalledOnce();
  });

  describe("Field validation", () => {
    describe("Trail Name", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Trail Name", "Required");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Trail Name", "Joffre Lakes");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Trail Name");
      });
    });

    describe("Date", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Date", "Invalid ISO date");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Date", "2026-01-01");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Date");
      });
    });

    describe("Duration", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Duration (minutes)", "Required");
      });

      test("displays error on submit when value is a decimal", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Duration (minutes)", "1.5");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Duration (minutes)", "Must be a positive integer");
      });

      test("displays error on submit when value is negative", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Duration (minutes)", "-1");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Duration (minutes)", "Must be a positive integer");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Duration (minutes)", "240");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Duration (minutes)");
      });
    });

    describe("Rating", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Rating (0-5)", "Required");
      });

      test("displays error on submit when value is below 0", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Rating (0-5)", "-0.5");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Rating (0-5)", "Too small: expected number to be >=0");
      });

      test("displays error on submit when value is above 5", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Rating (0-5)", "5.5");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Rating (0-5)", "Too big: expected number to be <=5");
      });

      test("displays error on submit when value is not a multiple of 0.5", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Rating (0-5)", "1.2");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Rating (0-5)", "Invalid number: must be a multiple of 0.5");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Rating (0-5)", "4.5");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Rating (0-5)");
      });
    });

    describe("Difficulty", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Difficulty (0-10)", "Required");
      });

      test("displays error on submit when value is below 0", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Difficulty (0-10)", "-0.5");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Difficulty (0-10)", "Too small: expected number to be >=0");
      });

      test("displays error on submit when value is above 10", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Difficulty (0-10)", "10.5");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Difficulty (0-10)", "Too big: expected number to be <=10");
      });

      test("displays error on submit when value is not a multiple of 0.5", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Difficulty (0-10)", "1.2");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Difficulty (0-10)", "Invalid number: must be a multiple of 0.5");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Difficulty (0-10)", "7.5");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Difficulty (0-10)");
      });
    });

    describe("Distance", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Distance (km)", "Required");
      });

      test("displays error on submit when value is negative", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Distance (km)", "-1");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Distance (km)", "Too small: expected number to be >=0");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Distance (km)", "11.8");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Distance (km)");
      });
    });

    describe("Elevation Gain", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Elevation gain (m)", "Required");
      });

      test("displays error on submit when value is a decimal", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Elevation gain (m)", "1.5");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Elevation gain (m)", "Must be a positive integer");
      });

      test("displays error on submit when value is negative", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Elevation gain (m)", "-1");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Elevation gain (m)", "Must be a positive integer");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Elevation gain (m)", "370");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Elevation gain (m)");
      });
    });

    describe("AllTrails URL", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "AllTrails URL", "Required");
      });

      test("displays error on submit when value is not a valid URL", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "AllTrails URL", "not a url");
        await clickSubmitButton(screen);
        await checkFieldError(screen, "AllTrails URL", "Invalid URL");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(
          screen,
          "AllTrails URL",
          "https://www.alltrails.com/trail/canada/british-columbia/joffre-lakes",
        );
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "AllTrails URL");
      });
    });

    describe("Notes", () => {
      test("displays error on submit when value is empty", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await clickSubmitButton(screen);
        await checkFieldError(screen, "Notes", "Required");
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await fillFieldInput(screen, "Notes", "Great views at the lake.");
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Notes");
      });
    });

    describe("Photos", () => {
      const oversizedPhoto = new File([new Uint8Array(MAX_PHOTO_SIZE + 1)], "oversized.png", { type: "image/png" });
      const invalidTypePhoto = new File(["notes"], "notes.txt", { type: "text/plain" });
      const validPhoto = new File(["photo"], "summit.png", { type: "image/png" });

      test("displays error on submit when photo exceeds 10MB", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await uploadPhotos(screen, [oversizedPhoto]);
        await clickSubmitButton(screen);
        // An error for a photo is displayed on its Caption field
        await checkFieldError(screen, "Caption 1", "Image must be under 10MB");
      });

      test("displays error on submit when file is not a valid image type", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await uploadPhotos(screen, [invalidTypePhoto]);
        await clickSubmitButton(screen);
        await checkFieldError(
          screen,
          "Caption 1",
          'Invalid option: expected one of "image/jpeg"|"image/png"|"image/webp"|"image/heic"|"image/heif"',
        );
      });

      test("does not display error on submit when value is valid", async () => {
        const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);
        await uploadPhotos(screen, [validPhoto]);
        await clickSubmitButton(screen);
        await checkNoFieldError(screen, "Caption 1");
      });
    });
  });

  test("invokes onSubmit callback when valid form is submitted and disables cancel and submit buttons while submission is in progress", async () => {
    // Have onSubmit sleep for a few seconds so we can check buttons are disabled
    mockOnSubmit.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    expect(mockOnSubmit).toHaveBeenCalledOnce();

    const cancelButton = screen.getByRole("button", { name: "Cancel" });
    await expect.element(cancelButton).toBeDisabled();

    const submitButton = screen.getByRole("button", { name: "Submit" });
    await expect.element(submitButton).toBeDisabled();
  });

  test("does not invoke onCancel callback when backdrop is clicked while submission is in progress", async () => {
    // Have onSubmit sleep for a few seconds so we can check buttons are disabled
    mockOnSubmit.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    await fillHikeForm(screen);
    await clickSubmitButton(screen);
    await clickBackdrop(screen);

    expect(mockOnSubmit).toHaveBeenCalledOnce();
    expect(mockOnCancel).not.toHaveBeenCalled();
  });

  test("displays error message in footer when submission of valid form fails", async () => {
    mockOnSubmit.mockRejectedValue(new Error("Failed to submit form"));

    const screen = await render(<HikeFormWrapper onCancel={mockOnCancel} onSubmit={mockOnSubmit} />);

    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    expect(mockOnSubmit).toHaveBeenCalledOnce();

    const footer = screen.getByRole("region", { name: "Footer" });
    const errorMessage = footer.getByText("Something went wrong. Please try again.");
    await expect.element(errorMessage).toBeInTheDocument();
  });
});
