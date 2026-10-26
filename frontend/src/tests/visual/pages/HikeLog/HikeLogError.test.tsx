import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeLogError from "../../../../pages/HikeLog/HikeLogError";

describe("HikeLogError", () => {
  test("displays component", async () => {
    const screen = await render(<HikeLogError />);
    const component = screen.getByRole("region", { name: "Hike log error" });
    await expect(component).toMatchScreenshot();
  });
});
