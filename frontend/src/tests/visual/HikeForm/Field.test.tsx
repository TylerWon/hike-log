import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { FieldWrapper } from "../../utils/HikeForm/test-components";

describe("Field", () => {
  test("displays component", async () => {
    const screen = await render(<FieldWrapper error="Invalid" required />);
    const field = screen.getByRole("group", { name: "Name field" });
    await expect(field).toMatchScreenshot();
  });
});
