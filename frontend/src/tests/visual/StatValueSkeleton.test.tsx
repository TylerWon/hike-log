import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import StatValueSkeleton from "../../components/StatValueSkeleton";

describe("StatValueSkeleton", () => {
  test("displays component", async () => {
    const screen = await render(<StatValueSkeleton widthClass="w-16" />);
    const component = screen.getByRole("region", { name: "Statistic skeleton" });
    await expect(component).toMatchScreenshot();
  });
});
