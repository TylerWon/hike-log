import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeLogContent from "../../components/HikeLogContent";

describe("HikeLogContent", () => {
  test("displays title, stats overview, and body", async () => {
    const overallStats = [
      { label: "Hikes", value: "2" },
      { label: "Distance", value: "25.6 km" },
    ];

    const screen = await render(
      <HikeLogContent overallStats={overallStats}>
        <ol aria-label="body">
          <li>Hike 1</li>
          <li>Hike 2</li>
        </ol>
      </HikeLogContent>,
    );

    const hikeStatLabel = screen.getByText(overallStats[0].label);
    await expect.element(hikeStatLabel).toBeInTheDocument();

    const hikeStatValue = screen.getByText(overallStats[0].value);
    await expect.element(hikeStatValue).toBeInTheDocument();

    const distanceStatLabel = screen.getByText(overallStats[1].label);
    await expect.element(distanceStatLabel).toBeInTheDocument();

    const distanceStatValue = screen.getByText(overallStats[1].value);
    await expect.element(distanceStatValue).toBeInTheDocument();

    const body = screen.getByRole("list", { name: "body" });
    await expect.element(body).toBeInTheDocument();
  });
});
