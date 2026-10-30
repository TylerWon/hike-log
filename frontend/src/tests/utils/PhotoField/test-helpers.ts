import type { RenderResult } from "vitest-browser-react";

import { expect } from "vitest";
import { type Locator, userEvent } from "vitest/browser";

/**
 * Uploads photos to a PhotoField.
 *
 * We cannot test photo upload through the actual user path (i.e. using the upload button) because that opens an OS
 * dialog that Vitest cannot access. Instead, we must directly access the hidden file input and upload photos with it.
 */
export async function uploadPhotos(screen: Locator | RenderResult, photos: File[]) {
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
