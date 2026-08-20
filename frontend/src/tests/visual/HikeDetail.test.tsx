import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeDetail from "../../components/HikeDetail";
import { HIKE_FIXTURE_1 } from "../fixtures/hike";

describe("HikeDetail", () => {
  test("displays hike details without photos when hike has no photos", async () => {
    const hike = structuredClone(HIKE_FIXTURE_1);
    hike.photos = [];

    const screen = await render(<HikeDetail hike={hike} />);
    const component = screen.getByRole("region", { name: `${hike.trailName} details` });
    await expect(component).toMatchScreenshot();
  });

  test("displays hike details with photos when hike has photos", async () => {
    const screen = await render(<HikeDetail hike={HIKE_FIXTURE_1} />);
    const component = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} details` });
    await expect(component).toMatchScreenshot();
  });
});
