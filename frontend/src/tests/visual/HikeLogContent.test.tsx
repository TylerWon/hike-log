import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeCard from "../../components/HikeCard";
import HikeLogContent from "../../components/HikeLogContent";
import { formatDistance, formatDuration, formatElevation } from "../../utils/formatters";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../fixtures/hike";

describe("HikeLogContent", () => {
  test("displays component", async () => {
    const hikes = [HIKE_FIXTURE_1, HIKE_FIXTURE_2];
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

    const component = screen.getByRole("region", { name: "Hike log content" });
    await expect(component).toMatchScreenshot();
  });
});
