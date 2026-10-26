import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../../../mocks/models/hike";
import { HikeLogContentWrapper } from "../../../utils/HikeLog/test-components";

describe("HikeLogContent", () => {
  test("displays title, stats overview, and body", async () => {
    const screen = await render(<HikeLogContentWrapper />);

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
    await expect.element(statsOverview).toBeInTheDocument();

    const card1 = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} card` });
    await expect.element(card1).toBeInTheDocument();

    const card2 = screen.getByRole("region", { name: `${HIKE_MOCK_2.trailName} card` });
    await expect.element(card2).toBeInTheDocument();
  });
});
