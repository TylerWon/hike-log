import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeOverview from "../../components/HikeOverview";
import { HIKE_MOCK_1 } from "../mocks/models/hike";

describe("HikeOverview", () => {
  test("displays component", async () => {
    const screen = await render(<HikeOverview hike={HIKE_MOCK_1} index={0} isExpanded={false} />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
    await expect(component).toMatchScreenshot();
  });
});
