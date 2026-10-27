import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import GenericError from "../../../../components/GenericError";

describe("FormError", () => {
  test("displays error message", async () => {
    const screen = await render(<GenericError />);
    const message = screen.getByText("Something went wrong. Please try again.");
    await expect.element(message).toBeInTheDocument();
  });
});
