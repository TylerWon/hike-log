import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeOverview from "../../components/HikeOverview";
import { HIKE_FIXTURE_1 } from "../fixtures/hike";

describe("HikeOverview", () => {
  test("displays component", async () => {
    const screen = await render(<HikeOverview hike={HIKE_FIXTURE_1} index={0} isExpanded={false} />);
    const component = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} overview` });
    await expect(component).toMatchScreenshot();
  });
});
