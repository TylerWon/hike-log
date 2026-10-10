import type { RenderResult } from "vitest-browser-react";

import { expect } from "vitest";
import { userEvent } from "vitest/browser";

// Uploads photos to the PhotoField.
// We cannot test photo upload through the actual user path (i.e. clicking the upload button) because that opens an OS
// dialog that Vitest cannot access. Instead, we directly access the hidden file input and upload photos with it.
export async function uploadPhotos(screen: RenderResult, photos: File[]) {
  const fileInput = screen.getByLabelText("Hidden file input"); // file inputs don't have an ARIA role
  await userEvent.upload(fileInput, photos);
}

// Fills a Field input with the provided value. The input is looked up by its associated label.
export async function fillFieldInput(screen: RenderResult, inputLabel: string, value: string) {
  const input = screen.getByLabelText(inputLabel);
  await input.fill(value);
}

// Checks that the provided error exists for a Field.
export async function checkForFieldError(screen: RenderResult, fieldName: string, error: string) {
  const field = screen.getByRole("group", { name: `${fieldName} field` });
  const errorElement = field.getByText(error);
  await expect.element(errorElement).toBeInTheDocument();
}

// Checks no error exists for a Field.
export async function checkNoFieldError(screen: RenderResult, fieldName: string) {
  const field = screen.getByRole("group", { name: `${fieldName} field` });
  const errorElement = field.getByRole("paragraph");
  await expect.element(errorElement).not.toBeInTheDocument();
}
