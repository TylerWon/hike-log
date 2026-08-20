import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeCardSkeleton from "../../components/HikeCardSkeleton";

describe("HikeCardSkeleton", () => {
  test("displays component", async () => {
    const screen = await render(<HikeCardSkeleton />);
    const component = screen.getByRole("region", { name: "Hike card skeleton" });
    await expect(component).toMatchScreenshot();
  });
});
