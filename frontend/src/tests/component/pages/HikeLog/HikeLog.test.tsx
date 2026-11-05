import { beforeAll, describe, expect, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Hike } from "../../../../schemas/models/hike";

import { checkAuth } from "../../../../network/api/auth";
import { AuthError } from "../../../../network/api/errors";
import { deleteHike, fetchHikes } from "../../../../network/api/hikes";
import { HikeCreationError, PhotoCreationError } from "../../../../pages/HikeLog/HikeLog/errors";
import { submitHikeForm } from "../../../../pages/HikeLog/HikeLog/submit-hike-form";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../../../mocks/models/hike";
import { CHECK_AUTH_RESPONSE_MOCK_1, CHECK_AUTH_RESPONSE_MOCK_2 } from "../../../mocks/network/api/auth";
import { clickButton } from "../../../utils/button";
import { fillHikeForm } from "../../../utils/HikeForm/test-helpers";
import { HikeLogWrapper } from "../../../utils/HikeLog/test-components";
import { checkStatsOverviewLabels, checkStatsOverviewValues } from "../../../utils/StatsOverview/test-helpers";

vi.mock(import("../../../../pages/HikeLog/HikeLog/submit-hike-form"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    submitHikeForm: vi.fn(),
  };
});

vi.mock(import("../../../../network/api/hikes"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    deleteHike: vi.fn(),
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

/**
 * Checks that cards are displayed/not displayed for each of the given hikes.
 */
async function checkCardsExist(screen: RenderResult, hikes: Hike[], shouldExist: boolean) {
  for (const hike of hikes) {
    const card = screen.getByRole("region", { name: `${hike.trailName} card` });
    if (shouldExist) {
      await expect.element(card).toBeInTheDocument();
    } else {
      await expect.element(card).not.toBeInTheDocument();
    }
  }
}

/**
 * Checks that the card with the given trail name is expanded/unexpanded.
 */
async function checkCardExpanded(screen: RenderResult, trailName: string, isExpanded: boolean) {
  const card = screen.getByRole("region", { name: `${trailName} details` });
  if (isExpanded) {
    await expect.element(card).toBeInTheDocument();
  } else {
    await expect.element(card).not.toBeInTheDocument();
  }
}

/**
 * Checks that the HikeLog displays the correct components when a query fails.
 */
async function checkError(screen: RenderResult) {
  const error = screen.getByRole("region", { name: "Hike log error" });
  await expect.element(error).toBeInTheDocument();
}

/**
 * Checks that the HikeLog displays the correct components while waiting for a query to finish.
 */
async function checkPending(screen: RenderResult) {
  const title = screen.getByText("Hike Log");
  await expect.element(title).toBeInTheDocument();

  const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
  await checkStatsOverviewLabels(statsOverview);

  const statValueSkeletons = screen.getByRole("region", { name: "Statistic skeleton" });
  expect(statValueSkeletons.length).toEqual(4);

  const cardSkeletons = screen.getByRole("region", { name: "Hike card skeleton" });
  expect(cardSkeletons.length).toEqual(5);
}

/**
 * Clicks the card with the given trail name.
 */
async function clickCard(screen: RenderResult, trailName: string) {
  // Only the overview portion of card reacts to clicks
  const card = screen.getByRole("region", { name: `${trailName} overview` });
  await card.click();
}

/**
 * Deletes the card with the given trail name.
 */
async function deleteCard(screen: RenderResult, trailName: string) {
  const card = screen.getByRole("region", { name: `${trailName} card` });

  await clickButton(card, "Delete hike");

  const deleteDialog = card.getByRole("alertdialog", { name: "Confirmation dialog" });
  await clickButton(deleteDialog, "Delete");
}

describe("HikeLog", () => {
  test("displays title, overall stats, hike cards, and button to add hike when hikes and auth are loaded", async () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);
    vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

    const screen = await render(<HikeLogWrapper />);

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
    await checkStatsOverviewLabels(statsOverview);
    await checkStatsOverviewValues(
      statsOverview,
      hikes.length,
      hikes[0].distance + hikes[1].distance,
      hikes[0].elevationGain + hikes[1].elevationGain,
      hikes[0].duration + hikes[1].duration,
    );

    await checkCardsExist(screen, hikes, true);

    const addHikeButton = screen.getByRole("button", { name: "Add hike" });
    await expect.element(addHikeButton).toBeInTheDocument();
  });

  test("hides add hike and delete hike buttons when user is unauthenticated", async () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);
    vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_2); // mock unauthenticated response

    const screen = await render(<HikeLogWrapper />);

    const addHikeButton = screen.getByRole("button", { name: "Add hike" });
    await expect.element(addHikeButton).not.toBeInTheDocument();

    const card1 = screen.getByRole("region", { name: `${hikes[0].trailName} card` });
    const deleteHikeButton = card1.getByRole("button", { name: "Delete hike" });
    await expect.element(deleteHikeButton).not.toBeInTheDocument();
  });

  describe("Query error", () => {
    test("displays an error when there is an issue getting the hikes", async () => {
      vi.mocked(fetchHikes).mockRejectedValue(new Error("Failed to retrieve hikes"));
      vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

      const screen = await render(<HikeLogWrapper />);

      await checkError(screen);
    });

    test("displays an error when there is an issue checking auth", async () => {
      const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
      vi.mocked(fetchHikes).mockResolvedValue(hikes);
      vi.mocked(checkAuth).mockRejectedValue(new Error("Failed to perform auth check"));

      const screen = await render(<HikeLogWrapper />);

      await checkError(screen);
    });
  });

  describe("Query pending", () => {
    test("displays title and skeletons of the overall stats and hike cards when hikes are pending", async () => {
      vi.mocked(fetchHikes).mockImplementation(() => new Promise(() => {})); // never resolves
      vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);

      const screen = await render(<HikeLogWrapper />);

      await checkPending(screen);
    });

    test("displays title and skeletons of the overall stats and hike cards when auth is pending", async () => {
      const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
      vi.mocked(fetchHikes).mockResolvedValue(hikes);
      vi.mocked(checkAuth).mockImplementation(() => new Promise(() => {})); // never resolves

      const screen = await render(<HikeLogWrapper />);

      await checkPending(screen);
    });
  });

  describe("Add hike", () => {
    beforeAll(() => {
      vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);
    });

    test("displays add hike form when add hike button is clicked", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);

      const screen = await render(<HikeLogWrapper />);

      await clickButton(screen, "Add hike");

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await expect.element(form).toBeInTheDocument();
    });

    test("hides add hike form when it is closed", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);

      const screen = await render(<HikeLogWrapper />);

      await clickButton(screen, "Add hike");

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await expect.element(form).toBeInTheDocument();

      await clickButton(screen, "Close");

      await expect.element(form).not.toBeInTheDocument();
    });

    test("displays error in add hike form when submission fails because hike cannot be created", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);
      vi.mocked(submitHikeForm).mockRejectedValue(new HikeCreationError("Failed to create hike"));

      const screen = await render(<HikeLogWrapper />);

      await clickButton(screen, "Add hike");

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickButton(form, "Submit");

      const error = form.getByText("Something went wrong. Please try again.");
      await expect.element(error).toBeInTheDocument();
    });

    test("closes add hike form, displays hike card for new hike, updates overall stats, and displays warning when submission succeeds but some photos cannot be created", async () => {
      const hikes = [HIKE_MOCK_1];
      // Mock fetchHikes twice since it is called on component load and again when submission succeeds
      vi.mocked(fetchHikes).mockResolvedValueOnce([]).mockResolvedValueOnce(hikes);
      vi.mocked(submitHikeForm).mockRejectedValue(new PhotoCreationError("Failed to create some photos"));

      const screen = await render(<HikeLogWrapper />);

      await checkCardsExist(screen, hikes, false);

      await clickButton(screen, "Add hike");

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickButton(form, "Submit");

      await expect.element(form).not.toBeInTheDocument();

      await checkCardsExist(screen, hikes, true);

      const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
      await checkStatsOverviewValues(statsOverview, 1, hikes[0].distance, hikes[0].elevationGain, hikes[0].duration);

      const warning = screen.getByRole("status", { name: "Toast" });
      await expect.element(warning).toBeInTheDocument();
    });

    test("redirects to login page when submission fails because the user's session expired", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);
      vi.mocked(submitHikeForm).mockRejectedValue(new AuthError("Unauthorized"));

      const navigateMock = vi.fn();

      const screen = await render(<HikeLogWrapper navigateProp={navigateMock} />);

      await clickButton(screen, "Add hike");

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickButton(form, "Submit");

      expect(navigateMock).toHaveBeenCalledWith("/login?reason=expired");
    });

    test("closes add hike form, displays hike card for new hike, and updates overall stats when form submission succeeds", async () => {
      const hikes = [HIKE_MOCK_1];
      // Mock fetchHikes twice since it is called on component load and again when submission succeeds
      vi.mocked(fetchHikes).mockResolvedValueOnce([]).mockResolvedValueOnce(hikes);
      vi.mocked(submitHikeForm).mockResolvedValue();

      const screen = await render(<HikeLogWrapper />);

      await checkCardsExist(screen, hikes, false);

      await clickButton(screen, "Add hike");

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickButton(form, "Submit");

      await expect.element(form).not.toBeInTheDocument();

      await checkCardsExist(screen, hikes, true);

      const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
      await checkStatsOverviewValues(statsOverview, 1, hikes[0].distance, hikes[0].elevationGain, hikes[0].duration);
    });
  });

  describe("Hike card", () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];

    describe("Expand hike card", () => {
      beforeAll(() => {
        vi.mocked(fetchHikes).mockResolvedValue(hikes);
        vi.mocked(checkAuth).mockResolvedValue(CHECK_AUTH_RESPONSE_MOCK_1);
      });

      test("expands hike card when clicked", async () => {
        const screen = await render(<HikeLogWrapper />);
        await clickCard(screen, hikes[0].trailName);
        await checkCardExpanded(screen, hikes[0].trailName, true);
      });

      test("unexpands expanded hike card when clicked again", async () => {
        const screen = await render(<HikeLogWrapper />);
        await clickCard(screen, hikes[0].trailName);
        await checkCardExpanded(screen, hikes[0].trailName, true);
        await clickCard(screen, hikes[0].trailName);
        await checkCardExpanded(screen, hikes[0].trailName, false);
      });

      test("unexpands expanded hike card when another card is clicked", async () => {
        const screen = await render(<HikeLogWrapper />);
        await clickCard(screen, hikes[0].trailName);
        await checkCardExpanded(screen, hikes[0].trailName, true);
        await clickCard(screen, hikes[1].trailName);
        await checkCardExpanded(screen, hikes[0].trailName, false);
        await checkCardExpanded(screen, hikes[1].trailName, true);
      });
    });

    describe("Delete hike card", () => {
      test("displays error when delete fails", async () => {
        vi.mocked(fetchHikes).mockResolvedValue(hikes);
        vi.mocked(deleteHike).mockRejectedValue(new Error("Failed to delete hike"));

        const screen = await render(<HikeLogWrapper />);

        await deleteCard(screen, hikes[0].trailName);

        const error = screen.getByText("Something went wrong. Please try again.");
        await expect.element(error).toBeInTheDocument();

        const card1 = screen.getByRole("region", { name: `${hikes[0].trailName} card` });
        await expect.element(card1).toBeInTheDocument();

        const card2 = screen.getByRole("region", { name: `${hikes[1].trailName} card` });
        await expect.element(card2).toBeInTheDocument();
      });

      test("redirects to login page when delete fails because the user's session expired", async () => {
        vi.mocked(fetchHikes).mockResolvedValue(hikes);
        vi.mocked(deleteHike).mockRejectedValue(new AuthError("Unauthorized"));

        const navigateMock = vi.fn();

        const screen = await render(<HikeLogWrapper navigateProp={navigateMock} />);

        await deleteCard(screen, hikes[0].trailName);

        expect(navigateMock).toHaveBeenCalledWith("/login?reason=expired");
      });

      test("removes hike card and updates overall stats when delete succeeds", async () => {
        // Mock fetchHikes twice since it is called on component load and again when delete succeeds
        vi.mocked(fetchHikes).mockResolvedValueOnce(hikes).mockResolvedValueOnce([hikes[1]]);
        vi.mocked(deleteHike).mockResolvedValue();

        const screen = await render(<HikeLogWrapper />);

        await deleteCard(screen, hikes[0].trailName);

        const card1 = screen.getByRole("region", { name: `${hikes[0].trailName} card` });
        await expect.element(card1).not.toBeInTheDocument();

        const card2 = screen.getByRole("region", { name: `${hikes[1].trailName} card` });
        await expect.element(card2).toBeInTheDocument();

        const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
        await checkStatsOverviewValues(statsOverview, 1, hikes[1].distance, hikes[1].elevationGain, hikes[1].duration);
      });

      test("removes expanded hike card and leaves remaining cards unexpanded when delete succeeds", async () => {
        // Mock fetchHikes twice since it is called on component load and again when delete succeeds
        vi.mocked(fetchHikes).mockResolvedValueOnce(hikes).mockResolvedValueOnce([hikes[1]]);
        vi.mocked(deleteHike).mockResolvedValue();

        const screen = await render(<HikeLogWrapper />);

        await clickCard(screen, hikes[0].trailName);
        await checkCardExpanded(screen, hikes[0].trailName, true);

        await deleteCard(screen, hikes[0].trailName);

        const card1 = screen.getByRole("region", { name: `${hikes[0].trailName} card` });
        await expect.element(card1).not.toBeInTheDocument();

        const card2 = screen.getByRole("region", { name: `${hikes[1].trailName} card` });
        await expect.element(card2).toBeInTheDocument();
        await checkCardExpanded(screen, hikes[1].trailName, false);
      });

      test("removes hike card and leaves expanded card untouched when delete succeeds", async () => {
        // Mock fetchHikes twice since it is called on component load and again when delete succeeds
        vi.mocked(fetchHikes).mockResolvedValueOnce(hikes).mockResolvedValueOnce([hikes[1]]);
        vi.mocked(deleteHike).mockResolvedValue();

        const screen = await render(<HikeLogWrapper />);

        await clickCard(screen, hikes[1].trailName);
        await checkCardExpanded(screen, hikes[1].trailName, true);

        await deleteCard(screen, hikes[0].trailName);

        const card1 = screen.getByRole("region", { name: `${hikes[0].trailName} card` });
        await expect.element(card1).not.toBeInTheDocument();

        const card2 = screen.getByRole("region", { name: `${hikes[1].trailName} card` });
        await expect.element(card2).toBeInTheDocument();
        await checkCardExpanded(screen, hikes[1].trailName, true);
      });
    });
  });
});
