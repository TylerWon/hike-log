import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { TextAreaFieldWrapper } from "../../../utils/TextAreaField/test-components";

describe("TextAreaField", () => {
  test("displays label and input", async () => {
    const screen = await render(<TextAreaFieldWrapper />);
    const component = screen.getByRole("group", { name: "Notes field" });
    await expect(component).toMatchScreenshot();
  });

  test("displays error when provided", async () => {
    const screen = await render(<TextAreaFieldWrapper error="Invalid" />);
    const component = screen.getByRole("group", { name: "Notes field" });
    await expect(component).toMatchScreenshot();
  });

  test("displays required symbol when specified", async () => {
    const screen = await render(<TextAreaFieldWrapper required />);
    const component = screen.getByRole("group", { name: "Notes field" });
    await expect(component).toMatchScreenshot();
  });
});
