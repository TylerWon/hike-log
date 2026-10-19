import type { RenderResult } from "vitest-browser-react";

export async function clickAcceptButton(screen: RenderResult) {
  const acceptButton = screen.getByRole("button", { name: "Accept" });
  await acceptButton.click();
}
