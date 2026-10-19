import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import TextWithSpinner from "../../components/TextWithSpinner";

describe("TextWithSpinner", () => {
  test("displays spinner and text", async () => {
    const screen = await render(<TextWithSpinner text="TESTING" />);

    const text = screen.getByText("TESTING");
    await expect.element(text).toBeInTheDocument();

    const spinner = screen.getByRole("img", { name: "Spinner" });
    await expect.element(spinner).toBeInTheDocument();
  });
});
