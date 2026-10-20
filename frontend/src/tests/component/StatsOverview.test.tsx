import { describe, test } from "vitest";
import { render } from "vitest-browser-react";

import StatsOverview from "../../components/StatsOverview/StatsOverview";
import { checkStatLabels, checkStatValues } from "../utils/StatsOverview/test-helpers";

describe("StatsOverview", () => {
  test("displays label and value for each statistic provided", async () => {
    const stats = [
      { label: "Hikes", value: "2" },
      { label: "Distance", value: "25.6 km" },
      { label: "Elevation", value: <p>1,223 m</p> }, // test non-string value
      { label: "Time", value: "4h 58m" },
    ];

    const screen = await render(<StatsOverview stats={stats} />);

    await checkStatLabels(screen);
    await checkStatValues(screen, 2, 25.6, 1223, 298);
  });
});
