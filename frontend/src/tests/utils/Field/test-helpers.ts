import type { RenderResult } from "vitest-browser-react";
import type { Locator } from "vitest/browser";

import { expect } from "vitest";

/**
 * Fills a Field with the provided value.
 */
export async function fillField(screen: Locator | RenderResult, fieldName: string, value: string) {
  const input = screen.getByLabelText(fieldName);
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
