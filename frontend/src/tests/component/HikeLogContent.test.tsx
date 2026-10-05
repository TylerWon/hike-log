import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeCard from "../../components/HikeCard";
import HikeLogContent from "../../components/HikeLogContent";
import { formatDistance, formatDuration, formatElevation } from "../../utils/formatters";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../mocks/models/hike";

describe("HikeLogContent", () => {
  test("displays title, stats overview, and body", async () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
    const overallStats = [
      { label: "Hikes", value: hikes.length },
      { label: "Distance", value: formatDistance(hikes[0].distance + hikes[1].distance) },
      { label: "Elevation", value: formatElevation(hikes[0].elevationGain + hikes[1].elevationGain) },
      { label: "Time", value: formatDuration(hikes[0].duration + hikes[1].duration) },
    ];

    const screen = await render(
      <HikeLogContent overallStats={overallStats}>
        {hikes.map((hike, i) => (
          <li key={i}>
            <HikeCard hike={hike} index={i + 1} isExpanded={false} onClick={() => {}} />
          </li>
        ))}
      </HikeLogContent>,
    );

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
