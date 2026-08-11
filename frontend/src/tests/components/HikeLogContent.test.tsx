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

    const title = screen.getByText("Hiking Log");
    await expect.element(title).toBeInTheDocument();

    const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
    await expect.element(statsOverview).toBeInTheDocument();

    const body = screen.getByRole("list", { name: "body" });
    await expect.element(body).toBeInTheDocument();
  });
});
