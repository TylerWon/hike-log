import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import { fetchHikes } from "../../api/hikes";
import HikeLog from "../../components/HikeLog";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../fixtures/hike";

vi.mock(import("../../api/hikes"), () => ({
  fetchHikes: vi.fn(),
}));

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
  test("displays an error when there was an issue getting the hikes", async () => {
    vi.mocked(fetchHikes).mockRejectedValue(new Error("Failed to retrieve hikes"));

    const screen = await renderHikeLog();

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    await expect(screen.container).toMatchScreenshot();
  });

  test("displays skeleton view when hikes are loading", async () => {
    vi.mocked(fetchHikes).mockImplementation(() => new Promise(() => {}));

    const screen = await renderHikeLog();

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    await expect(screen.container).toMatchScreenshot();
  });

  test("displays hikes when hikes are loaded", async () => {
    const hikes = [HIKE_FIXTURE_1, HIKE_FIXTURE_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);

    const screen = await renderHikeLog();

    const card = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} card` });
    await expect.element(card).toBeInTheDocument();

    await expect(screen.container).toMatchScreenshot();
  });
});
