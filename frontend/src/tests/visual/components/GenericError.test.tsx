import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import GenericError from "../../../components/GenericError";

describe("GenericError", () => {
  test("displays component", async () => {
    const screen = await render(<GenericError />);
    const component = screen.getByRole("region", { name: "Generic error" });
    await expect(component).toMatchScreenshot();
  });
});
