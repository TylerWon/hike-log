import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeCard from "../../components/HikeCard";
import { HIKE_FIXTURE_1 } from "../fixtures/hike";

describe("HikeCard", () => {
  test("displays collapsed card", async () => {
    const screen = await render(<HikeCard hike={HIKE_FIXTURE_1} index={1} isExpanded={false} onClick={() => {}} />);
    const component = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} card` });
    await expect(component).toMatchScreenshot();
  });

  test("displays expanded card", async () => {
    const screen = await render(<HikeCard hike={HIKE_FIXTURE_1} index={1} isExpanded={true} onClick={() => {}} />);
    const component = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} card` });
    await expect(component).toMatchScreenshot();
  });
});
