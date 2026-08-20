import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import StatsOverview from "../../components/StatsOverview/StatsOverview";

describe("StatsOverview", () => {
  test("displays label and value for each statistic provided", async () => {
    const stats = [
      { label: "Hikes", value: "2" },
      { label: "Distance", value: "25.6 km" },
      { label: "Elevation", value: <p>1,223 m</p> }, // test non-string value
      { label: "Time", value: "4h 58m" },
    ];

    const screen = await render(<StatsOverview stats={stats} />);

    const hikesLabel = screen.getByText(stats[0].label);
    await expect.element(hikesLabel).toBeInTheDocument();

    const hikesValue = screen.getByText("2");
    await expect.element(hikesValue).toBeInTheDocument();

    const distanceLabel = screen.getByText(stats[1].label);
    await expect.element(distanceLabel).toBeInTheDocument();

    const distanceValue = screen.getByText("25.6 km");
    await expect.element(distanceValue).toBeInTheDocument();

    const elevationLabel = screen.getByText(stats[2].label);
    await expect.element(elevationLabel).toBeInTheDocument();

    const elevationValue = screen.getByText("1,223 m");
    await expect.element(elevationValue).toBeInTheDocument();

    const timeLabel = screen.getByText(stats[3].label);
    await expect.element(timeLabel).toBeInTheDocument();

    const timeValue = screen.getByText("4h 58m");
    await expect.element(timeValue).toBeInTheDocument();
  });
});
