import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import { fetchHikes } from "../../../../api/hikes";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../../../mocks/models/hike";
import { clickButton } from "../../../utils/button";
import { HikeLogWrapper } from "../../../utils/HikeLog/test-components";

vi.mock(import("../../../../api/hikes"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    fetchHikes: vi.fn(),
  };
});

describe("HikeLog", () => {
  test("displays an error when there was an issue getting the hikes", async () => {
    vi.mocked(fetchHikes).mockRejectedValue(new Error("Failed to retrieve hikes"));
    const screen = await render(<HikeLogWrapper />);
    await expect(screen.container).toMatchScreenshot();
  });

  test("displays skeleton view when hikes are loading", async () => {
    vi.mocked(fetchHikes).mockImplementation(() => new Promise(() => {}));
    const screen = await render(<HikeLogWrapper />);
    await expect(screen.container).toMatchScreenshot();
  });

  test("displays hikes when hikes are loaded", async () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);

    const screen = await render(<HikeLogWrapper />);

    await expect(screen.container).toMatchScreenshot();
  });

  test("displays add hike form when add hike button is clicked", async () => {
    const screen = await render(<HikeLogWrapper />);
    await clickButton(screen, "Add hike");
    await expect(screen.container).toMatchScreenshot();
  });
});
