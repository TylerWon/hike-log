import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import StatsOverview from "../../components/StatsOverview/StatsOverview";

describe("StatsOverview", () => {
  test("displays component", async () => {
    const stats = [
      { label: "Hikes", value: "2" },
      { label: "Distance", value: "25.6 km" },
      { label: "Elevation", value: <p>1,223 m</p> }, // test non-string value
      { label: "Time", value: "4h 58m" },
    ];

    const screen = await render(<StatsOverview stats={stats} />);
    const component = screen.getByRole("region", { name: "Statistic overview" });
    await expect(component).toMatchScreenshot();
  });
});
