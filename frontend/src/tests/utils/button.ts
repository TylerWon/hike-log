import type { RenderResult } from "vitest-browser-react";
import type { Locator } from "vitest/browser";

/**
 * Clicks the button with the given aria label.
 */
export async function clickButton(screen: Locator | RenderResult, ariaLabel: string) {
  const button = screen.getByRole("button", { name: ariaLabel });
  await button.click();
}
