import type { RenderResult } from "vitest-browser-react";

import { userEvent } from "vitest/browser";

// Uploads photos to the PhotoField.
// We cannot test photo upload through the actual user path (i.e. clicking the upload button) because that opens an OS
// dialog that Vitest cannot access. Instead, we directly access the hidden file input and upload photos with it.
export async function uploadPhotos(screen: RenderResult, photos: File[]) {
  const fileInput = screen.getByLabelText("Hidden file input"); // file inputs don't have an ARIA role
  await userEvent.upload(fileInput, photos);
}
