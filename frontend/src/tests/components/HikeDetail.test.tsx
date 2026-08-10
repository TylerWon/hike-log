import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Hike } from "../../schemas/hike";

import HikeDetail from "../../components/HikeDetail";
import { HIKE_FIXTURE } from "../fixtures/hike";

/**
 * Checks that the difficulty, AllTrails link, notes, and photos of `hike` are displayed in the HikeDetails component.
 * Note that the photo gallery is not shown when the Hike has no photos.
 */
async function checkDetails(screen: RenderResult, hike: Hike) {
  const difficultyNum = screen.getByText(`${hike.difficulty}/10`, { exact: true });
  await expect.element(difficultyNum).toBeInTheDocument();

  const difficultyLabel = screen.getByText("Difficulty", { exact: true });
  await expect.element(difficultyLabel).toBeInTheDocument();

  const difficultyBar = screen.getByRole("img", {
    exact: true,
    name: `${hike.difficulty} out of 10 difficulty`,
  });
  await expect.element(difficultyBar).toBeInTheDocument();

  const allTrailsLink = screen.getByRole("link", { exact: true, name: "View on AllTrails" });
  await expect.element(allTrailsLink).toHaveAttribute("href", hike.allTrailsUrl);
  await expect.element(allTrailsLink).toHaveAttribute("target", "_blank");

  const notesLabel = screen.getByText("Notes", { exact: true });
  await expect.element(notesLabel).toBeInTheDocument();

  const notes = screen.getByText(hike.notes, { exact: true });
  await expect.element(notes).toBeInTheDocument();

  if (hike.photos.length > 0) {
    const photosLabel = screen.getByText(`Photos (${hike.photos.length})`, { exact: true });
    await expect.element(photosLabel).toBeInTheDocument();
  }

  const photos = screen.getByRole("img", { name: /Photo \d+ gallery/ });
  expect(photos.length).toEqual(hike.photos.length);
}

describe("HikeDetail", () => {
  test("displays hike details", async () => {
    const screen = await render(<HikeDetail hike={HIKE_FIXTURE} />);
    await checkDetails(screen, HIKE_FIXTURE);
  });

  test("does not display photo gallery when hike has no photos", async () => {
    const hikes = structuredClone(HIKE_FIXTURE);
    hikes.photos = [];
    const screen = await render(<HikeDetail hike={hikes} />);
    await checkDetails(screen, hikes);
  });
});
