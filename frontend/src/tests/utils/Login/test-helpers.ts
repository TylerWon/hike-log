import type { RenderResult } from "vitest-browser-react";

export async function clickLoginButton(screen: RenderResult) {
  const button = screen.getByRole("button", { name: "Login" });
  await button.click();
}
