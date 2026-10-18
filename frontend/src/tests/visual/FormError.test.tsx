import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import FormError from "../../components/FormError";

describe("FormError", () => {
  test("displays component", async () => {
    const screen = await render(<FormError />);
    const component = screen.getByRole("region", { name: "Form error" });
    await expect(component).toMatchScreenshot();
  });
});
