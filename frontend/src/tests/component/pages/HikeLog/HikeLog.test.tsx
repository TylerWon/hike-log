import { beforeAll, describe, expect, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Hike } from "../../../../schemas/models/hike";

import { deleteHike, fetchHikes } from "../../../../api/hikes";
import { HikeCreationError, PhotoCreationError } from "../../../../pages/HikeLog/HikeLog/errors";
import { submitHikeForm } from "../../../../pages/HikeLog/HikeLog/submit-hike-form";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../../../mocks/models/hike";
import { clickSubmitButton, fillHikeForm } from "../../../utils/HikeForm/test-helpers";
import { HikeLogWrapper } from "../../../utils/HikeLog/test-components";
import { clickAddHikeButton } from "../../../utils/HikeLog/test-helpers";
import { checkStatLabels, checkStatValues } from "../../../utils/StatsOverview/test-helpers";

vi.mock(import("../../../../pages/HikeLog/HikeLog/submit-hike-form"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    submitHikeForm: vi.fn(),
  };
});

vi.mock(import("../../../../api/hikes"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    deleteHike: vi.fn(),
    fetchHikes: vi.fn(),
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
 * Deletes the card with the given trail name.
 */
async function deleteCard(screen: RenderResult, trailName: string) {
  const card = screen.getByRole("region", { name: `${trailName} card` });

  const deleteButton = card.getByRole("button", { name: "Delete hike" });
  await deleteButton.click();

  const deleteDialog = card.getByRole("alertdialog", { name: "Confirmation dialog" });
  const acceptDeleteButton = deleteDialog.getByRole("button", { name: "Accept" });
  await acceptDeleteButton.click();
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
 * Clicks the card with the given trail name.
 */
async function clickCard(screen: RenderResult, trailName: string) {
  // Only the overview portion of card reacts to clicks
  const card = screen.getByRole("region", { name: `${trailName} overview` });
  await card.click();
}

describe("HikeLog", () => {
  test("displays an error when there was an issue getting the hikes", async () => {
    vi.mocked(fetchHikes).mockRejectedValue(new Error("Failed to retrieve hikes"));

    const screen = await render(<HikeLogWrapper />);

    const error = screen.getByRole("region", { name: "Hike log error" });
    await expect.element(error).toBeInTheDocument();
  });

  test("displays title and skeletons of the overall stats and hike cards when hikes are loading", async () => {
    vi.mocked(fetchHikes).mockImplementation(() => new Promise(() => {})); // never resolves

    const screen = await render(<HikeLogWrapper />);

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
    await checkStatLabels(statsOverview);

    const statValueSkeletons = screen.getByRole("region", { name: "Statistic skeleton" });
    expect(statValueSkeletons.length).toEqual(4);

    const cardSkeletons = screen.getByRole("region", { name: "Hike card skeleton" });
    expect(cardSkeletons.length).toEqual(5);
  });

  test("displays title, overall stats, hike cards, and button to add hike when hikes are loaded", async () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);

    const screen = await render(<HikeLogWrapper />);

    const title = screen.getByText("Hike Log");
    await expect.element(title).toBeInTheDocument();

    const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
    await checkStatLabels(statsOverview);
    await checkStatValues(
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

  describe("Add hike", () => {
    test("displays add hike form when add hike button is clicked", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);

      const screen = await render(<HikeLogWrapper />);

      await clickAddHikeButton(screen);

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await expect.element(form).toBeInTheDocument();
    });

    test("hides add hike form when it is closed", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);

      const screen = await render(<HikeLogWrapper />);

      await clickAddHikeButton(screen);

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await expect.element(form).toBeInTheDocument();

      const closeButton = form.getByRole("button", { name: "Close" });
      await closeButton.click();

      await expect.element(form).not.toBeInTheDocument();
    });

    test("displays error in add hike form when form submission fails", async () => {
      vi.mocked(fetchHikes).mockResolvedValue([]);
      vi.mocked(submitHikeForm).mockRejectedValue(new HikeCreationError("Failed to create hike"));

      const screen = await render(<HikeLogWrapper />);

      await clickAddHikeButton(screen);

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickSubmitButton(form);

      const error = form.getByText("Something went wrong. Please try again.");
      await expect.element(error).toBeInTheDocument();
    });

    test("closes add hike form, displays hike card for new hike, updates overall stats, and displays warning when form submission partially fails", async () => {
      const hikes = [HIKE_MOCK_1];
      // Mock fetchHikes twice since it is called on component load and again when submission succeeds
      vi.mocked(fetchHikes).mockResolvedValueOnce([]).mockResolvedValueOnce(hikes);
      vi.mocked(submitHikeForm).mockRejectedValue(new PhotoCreationError("Failed to create some photos"));

      const screen = await render(<HikeLogWrapper />);

      await checkCardsExist(screen, hikes, false);

      await clickAddHikeButton(screen);

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickSubmitButton(form);

      await expect.element(form).not.toBeInTheDocument();

      await checkCardsExist(screen, hikes, true);

      const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
      await checkStatValues(statsOverview, 1, hikes[0].distance, hikes[0].elevationGain, hikes[0].duration);

      const warning = screen.getByRole("status", { name: "Toast" });
      await expect.element(warning).toBeInTheDocument();
    });

    test("closes add hike form, displays hike card for new hike, and updates overall stats when form submission succeeds", async () => {
      const hikes = [HIKE_MOCK_1];
      // Mock fetchHikes twice since it is called on component load and again when submission succeeds
      vi.mocked(fetchHikes).mockResolvedValueOnce([]).mockResolvedValueOnce(hikes);

      const screen = await render(<HikeLogWrapper />);

      await checkCardsExist(screen, hikes, false);

      await clickAddHikeButton(screen);

      const form = screen.getByRole("dialog", { name: "Hike form" });
      await fillHikeForm(form);
      await clickSubmitButton(form);

      await expect.element(form).not.toBeInTheDocument();

      await checkCardsExist(screen, hikes, true);

      const statsOverview = screen.getByRole("region", { name: "Statistic overview" });
      await checkStatValues(statsOverview, 1, hikes[0].distance, hikes[0].elevationGain, hikes[0].duration);
    });
  });

  describe("Hike cards", () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];

    describe("Expand hike card", () => {
      beforeAll(() => {
        vi.mocked(fetchHikes).mockResolvedValue(hikes);
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
        await checkStatValues(statsOverview, 1, hikes[1].distance, hikes[1].elevationGain, hikes[1].duration);
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
