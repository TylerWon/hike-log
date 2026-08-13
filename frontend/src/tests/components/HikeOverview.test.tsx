import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeOverview from "../../components/HikeOverview";
import { formatDate, formatDuration } from "../../utils/formatters";
import { HIKE_FIXTURE_1 } from "../fixtures/hike";

describe("HikeOverview", () => {
  test("displays hike information", async () => {
    const screen = await render(<HikeOverview hike={HIKE_FIXTURE_1} index={1} isExpanded={false} />);

    const image = screen.getByRole("img", { name: "Joffre Lakes thumbnail" });
    await expect.element(image).toBeInTheDocument();

    const hikeNumber = screen.getByText("#01", { exact: true });
    await expect.element(hikeNumber).toBeInTheDocument();

    const trailName = screen.getByText(HIKE_FIXTURE_1.trailName, { exact: true });
    await expect.element(trailName).toBeInTheDocument();

    const date = screen.getByText(formatDate(HIKE_FIXTURE_1.date), { exact: true });
    await expect.element(date).toBeInTheDocument();

    const ratingNum = screen.getByText(`${HIKE_FIXTURE_1.rating}/5`, { exact: true });
    await expect.element(ratingNum).toBeInTheDocument();

    const ratingStars = screen.getByRole("img", { name: `${HIKE_FIXTURE_1.rating} out of 5 stars` });
    await expect.element(ratingStars).toBeInTheDocument();

    const distance = screen.getByText(`${HIKE_FIXTURE_1.distance} km`, { exact: true });
    await expect.element(distance).toBeInTheDocument();

    const elevationGain = screen.getByText(`${HIKE_FIXTURE_1.elevationGain} m`, { exact: true });
    await expect.element(elevationGain).toBeInTheDocument();

    const duration = screen.getByText(formatDuration(HIKE_FIXTURE_1.duration), { exact: true });
    await expect.element(duration).toBeInTheDocument();
  });
});
