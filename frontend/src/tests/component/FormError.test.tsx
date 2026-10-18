import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import FormError from "../../components/FormError";

describe("FormError", () => {
  test("displays error message", async () => {
    const screen = await render(<FormError />);
    const message = screen.getByText("Something went wrong. Please try again.");
    await expect.element(message).toBeInTheDocument();
  });
});
