import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { FieldWrapper } from "../../utils/HikeForm/test-components";

describe("Field", () => {
  test("displays label, input, and required symbol", async () => {
    const screen = await render(<FieldWrapper error="Invalid" required />);

    const label = screen.getByText("Name");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Name");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).toBeInTheDocument();

    const error = screen.getByText("Invalid");
    await expect.element(error).toBeInTheDocument();
  });
});
