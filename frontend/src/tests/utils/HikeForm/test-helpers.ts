import type { RenderResult } from "vitest-browser-react";

import { expect } from "vitest";
import { userEvent } from "vitest/browser";

/**
 * Uploads photos to a PhotoField.
 *
 * We cannot test photo upload through the actual user path (i.e. using the upload button) because that opens an OS
 * dialog that Vitest cannot access. Instead, we must directly access the hidden file input and upload photos with it.
 */
export async function uploadPhotos(screen: RenderResult, photos: File[]) {
  const fileInput = screen.getByLabelText("Hidden file input"); // file inputs don't have an ARIA role
  await userEvent.upload(fileInput, photos);

  // userEvent.upload() resolves as soon as the hidden input fires a change event so photos might not actually be
  // finished uploading, especially if they are large. As a result, we verify that photos appear in the DOM before
  // proceeding.
  const photoField = screen.getByRole("group", { name: "Photo field" });
  for (let i = 0; i < photos.length; i++) {
    const caption = photoField.getByRole("group", { name: `Caption ${i + 1} field` });
    await expect.element(caption).toBeInTheDocument();
  }
}

/**
 * Fills a Field input with the provided value. The input is looked up by its label.
 */
export async function fillFieldInput(screen: RenderResult, label: string, value: string) {
  const input = screen.getByLabelText(label);
  await input.fill(value);
}

/**
 * Checks that the provided error exists for a Field.
 */
export async function checkFieldError(screen: RenderResult, fieldName: string, error: string) {
  const field = screen.getByRole("group", { name: `${fieldName} field` });
  const errorElement = field.getByText(error);
  await expect.element(errorElement).toBeInTheDocument();
}

/**
 * Checks no error exists for a Field.
 */
export async function checkNoFieldError(screen: RenderResult, fieldName: string) {
  const field = screen.getByRole("group", { name: `${fieldName} field` });
  const errorElement = field.getByRole("paragraph");
  await expect.element(errorElement).not.toBeInTheDocument();
}

/**
 * Clicks the HikeForm submit button.
 */
export async function clickSubmitButton(screen: RenderResult) {
  const submitButton = screen.getByRole("button", { name: "Submit" });
  await submitButton.click();
}

/**
 * Fills every field in a HikeForm with a valid value. Photos can optionally be added.
 */
export async function fillHikeForm(screen: RenderResult, photos: File[] = []) {
  await fillFieldInput(screen, "Trail Name", "Joffre Lakes");
  await fillFieldInput(screen, "Date", "2026-01-01");
  await fillFieldInput(screen, "Duration (minutes)", "240");
  await fillFieldInput(screen, "Rating (0-5)", "4.5");
  await fillFieldInput(screen, "Difficulty (0-10)", "7.5");
  await fillFieldInput(screen, "Distance (km)", "11.8");
  await fillFieldInput(screen, "Elevation gain (m)", "370");
  await fillFieldInput(screen, "AllTrails URL", "https://www.alltrails.com/trail/canada/british-columbia/joffre-lakes");
  await fillFieldInput(screen, "Notes", "Great views at the lake.");

  if (photos.length > 0) {
    await uploadPhotos(screen, photos);
  }
}
