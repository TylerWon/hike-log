import { describe, expect, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Hike } from "../../schemas/models/hike";

import { fetchHikes } from "../../api/hikes";
import { HikeCreationError, PhotoCreationError } from "../../components/HikeLog/errors";
import { submitHikeForm } from "../../components/HikeLog/submit-hike-form";
import { formatDistance, formatDuration, formatElevation } from "../../utils/formatters";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../mocks/models/hike";
import { clickSubmitButton, fillHikeForm } from "../utils/HikeForm/test-helpers";
import { HikeLogWrapper } from "../utils/HikeLog/test-components";
import { clickAddHikeButton } from "../utils/HikeLog/test-helpers";

vi.mock(import("../../components/HikeLog/submit-hike-form"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    submitHikeForm: vi.fn(),
  };
});

vi.mock(import("../../api/hikes"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    fetchHikes: vi.fn(),
  };
});

/**
 * Checks that cards are displayed/not displayed for each of the given `hikes`.
 */
async function checkCards(screen: RenderResult, hikes: Hike[], shouldExist: boolean = true) {
  for (const hike of hikes) {
    const card = screen.getByRole("region", { name: `${hike.trailName} card` });
    if (shouldExist) {
      await expect.element(card).toBeInTheDocument();
    } else {
      await expect.element(card).not.toBeInTheDocument();
    }
  }
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

  test("displays title, overall stats, hike cards, and button to add hike when hikes are loaded", async () => {
    const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
    vi.mocked(fetchHikes).mockResolvedValue(hikes);

    const screen = await render(<HikeLogWrapper />);

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

    const addHikeButton = screen.getByRole("button", { name: "Add hike" });
    await expect.element(addHikeButton).toBeInTheDocument();
  });

  test("displays add hike form when add hike button is clicked", async () => {
    vi.mocked(fetchHikes).mockResolvedValueOnce([]);

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
    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("dialog", { name: "Hike form" });
    const error = form.getByText("Something went wrong. Please try again.");
    await expect.element(error).toBeInTheDocument();
  });

  test("closes add hike form, displays hike card for new hike, and displays warning when form submission partially fails", async () => {
    const hikes = [HIKE_MOCK_1];
    // Mock fetchHikes twice since it is called on component load and again when submission succeeds
    vi.mocked(fetchHikes).mockResolvedValueOnce([]).mockResolvedValueOnce(hikes);
    vi.mocked(submitHikeForm).mockRejectedValue(new PhotoCreationError("Failed to create some photos"));

    const screen = await render(<HikeLogWrapper />);

    await checkCards(screen, hikes, false);

    await clickAddHikeButton(screen);
    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("dialog", { name: "Hike form" });
    await expect.element(form).not.toBeInTheDocument();

    await checkCards(screen, hikes);

    const warning = screen.getByRole("status", { name: "Toast" });
    await expect.element(warning).toBeInTheDocument();
  });

  test("closes add hike form and displays hike card for new hike when form submission succeeds", async () => {
    const hikes = [HIKE_MOCK_1];
    // Mock fetchHikes twice since it is called on component load and again when submission succeeds
    vi.mocked(fetchHikes).mockResolvedValueOnce([]).mockResolvedValueOnce(hikes);

    const screen = await render(<HikeLogWrapper />);

    await checkCards(screen, hikes, false);

    await clickAddHikeButton(screen);
    await fillHikeForm(screen);
    await clickSubmitButton(screen);

    const form = screen.getByRole("dialog", { name: "Hike form" });
    await expect.element(form).not.toBeInTheDocument();

    await checkCards(screen, hikes);
  });
});
