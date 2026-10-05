import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeCard from "../../components/HikeCard";
import { HIKE_MOCK_1 } from "../mocks/models/hike";

describe("HikeCard", () => {
  test("displays collapsed card", async () => {
    const screen = await render(<HikeCard hike={HIKE_MOCK_1} index={1} isExpanded={false} onClick={() => {}} />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} card` });
    await expect(component).toMatchScreenshot();
  });

  test("displays expanded card", async () => {
    const screen = await render(<HikeCard hike={HIKE_MOCK_1} index={1} isExpanded={true} onClick={() => {}} />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} card` });
    await expect(component).toMatchScreenshot();
  });
});
