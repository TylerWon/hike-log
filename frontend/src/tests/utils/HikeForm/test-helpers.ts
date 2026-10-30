import type { RenderResult } from "vitest-browser-react";

import { type Locator } from "vitest/browser";

import { fillField } from "../Field/test-helpers";
import { uploadPhotos } from "../PhotoField/test-helpers";

/**
 * Clicks the HikeForm submit button.
 */
export async function clickSubmitButton(screen: Locator | RenderResult) {
  const submitButton = screen.getByRole("button", { name: "Submit" });
  await submitButton.click();
}

/**
 * Fills every field in a HikeForm with a valid value. Photos can optionally be added.
 */
export async function fillHikeForm(screen: Locator | RenderResult, photos: File[] = []) {
  await fillField(screen, "Trail Name", "Joffre Lakes");
  await fillField(screen, "Date", "2026-01-01");
  await fillField(screen, "Duration (minutes)", "240");
  await fillField(screen, "Rating (0-5)", "4.5");
  await fillField(screen, "Difficulty (0-10)", "7.5");
  await fillField(screen, "Distance (km)", "11.8");
  await fillField(screen, "Elevation gain (m)", "370");
  await fillField(screen, "AllTrails URL", "https://www.alltrails.com/trail/canada/british-columbia/joffre-lakes");
  await fillField(screen, "Notes", "Great views at the lake.");

  if (photos.length > 0) {
    await uploadPhotos(screen, photos);
  }
}
