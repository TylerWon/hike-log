import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import { checkAuth } from "../../../../network/api/auth";
import { fetchHikes } from "../../../../network/api/hikes";
import { CHECK_AUTH_RESPONSE_MOCK_1, CHECK_AUTH_RESPONSE_MOCK_2 } from "../../../mocks/network/api/auth";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../../../mocks/models/hike";
import { clickButton } from "../../../utils/button";
import { HikeLogWrapper } from "../../../utils/HikeLog/test-components";

vi.mock(import("../../../../network/api/hikes"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    fetchHikes: vi.fn(),
  };
});

vi.mock(import("../../../../network/api/auth"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    checkAuth: vi.fn(),
  };
});

describe("HikeLog", () => {
  const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];

  describe("Query error", () => {
    test("displays an error when there is an issue getting the hikes", async () => {
      vi.mocked(fetchHikes).mockRejectedValue(new Error("Failed to retrieve hikes"));
      vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

      const screen = await render(<HikeLogWrapper />);

      await expect(screen.container).toMatchScreenshot();
    });

    test("displays an error when there is an issue checking auth", async () => {
      vi.mocked(fetchHikes).mockResolvedValue(hikes);
      vi.mocked(checkAuth).mockRejectedValue(new Error("Failed to perform auth check"));

      const screen = await render(<HikeLogWrapper />);

      await expect(screen.container).toMatchScreenshot();
    });
  });

  describe("Query pending", () => {
    test("displays skeleton view when hikes are pending", async () => {
      vi.mocked(fetchHikes).mockImplementation(() => new Promise(() => {}));
      vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

      const screen = await render(<HikeLogWrapper />);

      await expect(screen.container).toMatchScreenshot();
    });

    test("displays skeleton view when auth is pending", async () => {
      vi.mocked(fetchHikes).mockResolvedValue(hikes);
      vi.mocked(checkAuth).mockImplementation(() => new Promise(() => {}));

      const screen = await render(<HikeLogWrapper />);

      await expect(screen.container).toMatchScreenshot();
    });
  });

  test("displays hikes when hikes and auth are loaded", async () => {
    vi.mocked(fetchHikes).mockResolvedValue(hikes);
    vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

    const screen = await render(<HikeLogWrapper />);

    await expect(screen.container).toMatchScreenshot();
  });

  test("displays add hike form when add hike button is clicked", async () => {
    vi.mocked(fetchHikes).mockResolvedValue(hikes);
    vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

    const screen = await render(<HikeLogWrapper />);
    await clickButton(screen, "Add hike");

    await expect(screen.container).toMatchScreenshot();
  });

  test("hides add hike and delete hike buttons when user is unauthenticated", async () => {
    vi.mocked(fetchHikes).mockResolvedValue(hikes);
    vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_2); // mock unauthenticated response

    const screen = await render(<HikeLogWrapper />);

    await expect(screen.container).toMatchScreenshot();
  });
});
