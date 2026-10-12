import type { RenderResult } from "vitest-browser-react";

/**
 * Clicks the add hike button.
 */
export async function clickAddHikeButton(screen: RenderResult) {
  const button = screen.getByRole("button", { name: "Add hike" });
  await button.click();
}
