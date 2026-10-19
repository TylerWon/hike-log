import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import { formatDate, formatDuration } from "../../utils/formatters";
import { HIKE_MOCK_1 } from "../mocks/models/hike";
import { HikeOverviewWrapper } from "../utils/HikeOverview/test-components";

async function clickDeleteHikeButton(screen: RenderResult) {
  const deleteHikeButton = screen.getByRole("button", { name: "Delete hike" });
  await deleteHikeButton.click();
}

describe("HikeOverview", () => {
  let mockOnClick: Mock;
  let mockOnDelete: Mock;

  beforeEach(() => {
    mockOnClick = vi.fn();
    mockOnDelete = vi.fn();
  });

  test("displays hike information", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);

    const image = screen.getByRole("img", { name: "Joffre Lakes thumbnail" });
    await expect.element(image).toBeInTheDocument();

    const hikeNumber = screen.getByText("#01", { exact: true });
    await expect.element(hikeNumber).toBeInTheDocument();

    const trailName = screen.getByText(HIKE_MOCK_1.trailName, { exact: true });
    await expect.element(trailName).toBeInTheDocument();

    const date = screen.getByText(formatDate(HIKE_MOCK_1.date), { exact: true });
    await expect.element(date).toBeInTheDocument();

    const ratingNum = screen.getByText(`${HIKE_MOCK_1.rating}/5`, { exact: true });
    await expect.element(ratingNum).toBeInTheDocument();

    const ratingStars = screen.getByRole("img", { name: `${HIKE_MOCK_1.rating} out of 5 stars` });
    await expect.element(ratingStars).toBeInTheDocument();

    const distance = screen.getByText(`${HIKE_MOCK_1.distance} km`, { exact: true });
    await expect.element(distance).toBeInTheDocument();

    const elevationGain = screen.getByText(`${HIKE_MOCK_1.elevationGain} m`, { exact: true });
    await expect.element(elevationGain).toBeInTheDocument();

    const duration = screen.getByText(formatDuration(HIKE_MOCK_1.duration), { exact: true });
    await expect.element(duration).toBeInTheDocument();
  });

  test("invokes onClick callback when clicked", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);

    const overview = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
    await overview.click();

    expect(mockOnClick).toHaveBeenCalledOnce();
  });

  test("opens delete confirmation dialog when delete hike button is clicked", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);

    await clickDeleteHikeButton(screen);

    const deleteDialog = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect.element(deleteDialog).toBeInTheDocument();
  });

  test("closes delete confirmation dialog when dialog cancel button is clicked", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);

    await clickDeleteHikeButton(screen);

    const deleteDialog = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect.element(deleteDialog).toBeInTheDocument();

    const cancelDeleteButton = deleteDialog.getByRole("button", { name: "Cancel" });
    await cancelDeleteButton.click();

    await expect.element(deleteDialog).not.toBeInTheDocument();
  });

  test("invokes onDelete callback and closes confirmation dialog when dialog delete button is clicked", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);

    await clickDeleteHikeButton(screen);

    const deleteDialog = screen.getByRole("region", { name: "Confirmation dialog backdrop" });
    await expect.element(deleteDialog).toBeInTheDocument();

    const acceptDeleteButton = deleteDialog.getByRole("button", { name: "Accept" });
    await acceptDeleteButton.click();

    await expect.element(deleteDialog).not.toBeInTheDocument();
    expect(mockOnDelete).toHaveBeenCalledOnce();
  });
});
