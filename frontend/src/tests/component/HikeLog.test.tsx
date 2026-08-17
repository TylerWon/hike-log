import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { describe, expect, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Hike } from "../../schemas/hike";

import { fetchHikes } from "../../api/hikes";
import HikeLog from "../../components/HikeLog";
import { formatDistance, formatDuration, formatElevation } from "../../utils/formatters";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../fixtures/hike";

vi.mock(import("../../api/hikes"), () => ({
  fetchHikes: vi.fn(),
}));

/**
 * Checks that cards are displayed for each of the given `hikes`.
 */
async function checkCards(screen: RenderResult, hikes: Hike[]) {
  for (const hike of hikes) {
    const card = screen.getByRole("region", { name: `${hike.trailName} card` });
    await expect.element(card).toBeInTheDocument();
  }
}

/**
 * Wraps HikeLog with TanStack Query client for testing.
 */
function renderHikeLog() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } }, // disable retry to prevent error test from timing out due to retries by Tanstack Query
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <HikeLog />
    </QueryClientProvider>,
  );
}

describe("HikeLog", () => {
  test("displays an error when there was an issue getting the hike data", async () => {
    vi.mocked(fetchHikes).mockRejectedValue(new Error("Failed to retrieve hikes"));

    const screen = await renderHikeLog();

    const error = screen.getByRole("region", { name: "Hike log error" });
    await expect.element(error).toBeInTheDocument();
  });

  test("displays title and skeletons of the overall stats and hike cards when hike data is loading", async () => {
    vi.mocked(fetchHikes).mockImplementation(() => new Promise(() => {})); // never resolves

    const screen = await renderHikeLog();

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    const hikesStatLabel = screen.getByText("Hikes");
    await expect.element(hikesStatLabel).toBeInTheDocument();

    const distanceStatLabel = screen.getByText("Distance");
    await expect.element(distanceStatLabel).toBeInTheDocument();

    const elevationStatLabel = screen.getByText("Elevation");
    await expect.element(elevationStatLabel).toBeInTheDocument();

    const timeStatLabel = screen.getByText("Time");
    await expect.element(timeStatLabel).toBeInTheDocument();

    const statValueSkeletons = screen.getByRole("region", { name: "Statistic skeleton" });
    expect(statValueSkeletons.length).toEqual(4);

    const cardSkeletons = screen.getByRole("region", { name: "Hike card skeleton" });
    expect(cardSkeletons.length).toEqual(5);
  });

  test("displays title, overall stats, and hike cards when hike data exists", async () => {
    const hikes = [HIKE_FIXTURE_1, HIKE_FIXTURE_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);

    const screen = await renderHikeLog();

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    const hikesStatLabel = screen.getByText("Hikes");
    await expect.element(hikesStatLabel).toBeInTheDocument();

    const hikesStatValue = screen.getByText(`${hikes.length}`);
    await expect.element(hikesStatValue).toBeInTheDocument();

    const distanceStatLabel = screen.getByText("Distance");
    await expect.element(distanceStatLabel).toBeInTheDocument();

    const distanceStatValue = screen.getByText(formatDistance(hikes[0].distance + hikes[1].distance));
    await expect.element(distanceStatValue).toBeInTheDocument();

    const elevationStatLabel = screen.getByText("Elevation");
    await expect.element(elevationStatLabel).toBeInTheDocument();

    const elevationStatValue = screen.getByText(formatElevation(hikes[0].elevationGain + hikes[1].elevationGain));
    await expect.element(elevationStatValue).toBeInTheDocument();

    const timeStatLabel = screen.getByText("Time");
    await expect.element(timeStatLabel).toBeInTheDocument();

    const timeStatValue = screen.getByText(formatDuration(hikes[0].duration + hikes[1].duration));
    await expect.element(timeStatValue).toBeInTheDocument();

    await checkCards(screen, hikes);
  });
});
