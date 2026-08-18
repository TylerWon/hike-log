import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import StarRating from "../../components/StarRating";

describe("StarRating", () => {
  test("displays 5 empty stars when rating is 0", async () => {
    const screen = await render(<StarRating rating={0} />);
    const component = screen.getByRole("img", { name: "0 out of 5 stars" });
    await expect(component).toMatchScreenshot();
  });

  test("displays 5 filled stars when rating is 0", async () => {
    const screen = await render(<StarRating rating={5} />);
    const component = screen.getByRole("img", { name: "5 out of 5 stars" });
    await expect(component).toMatchScreenshot();
  });

  test("displays 2 filled, 1 half, 2 empty stars when rating is 2.5", async () => {
    const screen = await render(<StarRating rating={2.5} />);
    const component = screen.getByRole("img", { name: "2.5 out of 5 stars" });
    await expect(component).toMatchScreenshot();
  });
});
