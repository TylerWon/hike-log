import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import HikeOverview from "../../components/HikeOverview/HikeOverview";
import { formatDate, formatDuration } from "../../utils/formatters";
import { HIKE_FIXTURE } from "../fixtures/hikes";

describe("HikeOverview", () => {
  test("displays hike information", async () => {
    const screen = await render(<HikeOverview hike={HIKE_FIXTURE} index={1} isExpanded={false} />);

    const image = screen.getByRole("img", { exact: true, name: "Joffre Lakes thumbnail" });
    await expect.element(image).toBeInTheDocument();

    const hikeNumber = screen.getByText("#01", { exact: true });
    await expect.element(hikeNumber).toBeInTheDocument();

    const trailName = screen.getByText(HIKE_FIXTURE.trailName, { exact: true });
    await expect.element(trailName).toBeInTheDocument();

    const date = screen.getByText(formatDate(HIKE_FIXTURE.date), { exact: true });
    await expect.element(date).toBeInTheDocument();

    const ratingNum = screen.getByText(`${HIKE_FIXTURE.rating}/5`, { exact: true });
    await expect.element(ratingNum).toBeInTheDocument();

    const ratingStars = screen.getByRole("img", { exact: true, name: `${HIKE_FIXTURE.rating} out of 5 stars` });
    await expect.element(ratingStars).toBeInTheDocument();

    const distance = screen.getByText(`${HIKE_FIXTURE.distance} km`, { exact: true });
    await expect.element(distance).toBeInTheDocument();

    const elevationGain = screen.getByText(`${HIKE_FIXTURE.elevationGain} m`, { exact: true });
    await expect.element(elevationGain).toBeInTheDocument();

    const duration = screen.getByText(formatDuration(HIKE_FIXTURE.duration), { exact: true });
    await expect.element(duration).toBeInTheDocument();
  });
});
