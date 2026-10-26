import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import Toast from "../../components/Toast";

describe("Toast", () => {
  test("displays component", async () => {
    const screen = await render(<Toast message="TESTING" onDismiss={() => {}} />);
    const component = screen.getByRole("status", { name: "Toast" });
    await expect(component).toMatchScreenshot();
  });
});
