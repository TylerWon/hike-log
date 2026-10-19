import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { HIKE_MOCK_1 } from "../mocks/models/hike";
import { HikeCardWrapper } from "../utils/HikeCard/test-components";

describe("HikeCard", () => {
  test("displays collapsed card", async () => {
    const screen = await render(<HikeCardWrapper />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} card` });
    await expect(component).toMatchScreenshot();
  });

  test("displays expanded card", async () => {
    const screen = await render(<HikeCardWrapper isExpanded={true} />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} card` });
    await expect(component).toMatchScreenshot();
  });
});
