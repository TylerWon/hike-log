import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import StatsOverview from "../../components/StatsOverview/StatsOverview";

describe("StatsOverview", () => {
  test("displays label and value for each statistic provided", async () => {
    const stats = [
      { label: "Distance", value: "22.1 km" },
      { label: "Time", value: <p>2h 30m</p> }, // non-string value
    ];

    const screen = await render(<StatsOverview stats={stats} />);

    const distanceLabel = screen.getByText(stats[0].label);
    await expect.element(distanceLabel).toBeInTheDocument();

    const distanceValue = screen.getByText(stats[0].label);
    await expect.element(distanceValue).toBeInTheDocument();

    const timeLabel = screen.getByText(stats[1].label);
    await expect.element(timeLabel).toBeInTheDocument();

    const timeValue = screen.getByText("2h 30m");
    await expect.element(timeValue).toBeInTheDocument();
  });
});
