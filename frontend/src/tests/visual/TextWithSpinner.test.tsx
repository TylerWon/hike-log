import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import TextWithSpinner from "../../components/TextWithSpinner";

describe("TextWithSpinner", () => {
  test("displays component", async () => {
    const screen = await render(<TextWithSpinner text="TESTING" />);
    const component = screen.getByRole("group", { name: "Text with spinner" });
    await expect(component).toMatchScreenshot();
  });
});
