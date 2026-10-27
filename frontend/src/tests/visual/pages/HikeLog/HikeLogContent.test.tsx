import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { HikeLogContentWrapper } from "../../../utils/HikeLogContent/test-components";

describe("HikeLogContent", () => {
  test("displays component", async () => {
    const screen = await render(<HikeLogContentWrapper />);
    const component = screen.getByRole("region", { name: "Hike log content" });
    await expect(component).toMatchScreenshot();
  });
});
