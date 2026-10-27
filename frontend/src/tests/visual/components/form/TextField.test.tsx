import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { TextFieldWrapper } from "../../../utils/TextField/test-components";

describe("TextField", () => {
  test("displays label and input", async () => {
    const screen = await render(<TextFieldWrapper />);
    const component = screen.getByRole("group", { name: "Name field" });
    await expect(component).toMatchScreenshot();
  });

  test("displays error when provided", async () => {
    const screen = await render(<TextFieldWrapper error="Invalid" />);
    const component = screen.getByRole("group", { name: "Name field" });
    await expect(component).toMatchScreenshot();
  });

  test("displays required symbol when specified", async () => {
    const screen = await render(<TextFieldWrapper required />);
    const component = screen.getByRole("group", { name: "Name field" });
    await expect(component).toMatchScreenshot();
  });
});
