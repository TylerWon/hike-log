import type { RenderResult } from "vitest-browser-react";
import type { Locator } from "vitest/browser";

import { expect } from "vitest";

import { formatDistance, formatDuration, formatElevation } from "../../../utils/formatters";

export async function checkStatLabels(screen: Locator | RenderResult) {
  const hikesStatLabel = screen.getByText("Hikes");
  await expect.element(hikesStatLabel).toBeInTheDocument();

  const distanceStatLabel = screen.getByText("Distance");
  await expect.element(distanceStatLabel).toBeInTheDocument();

  const elevationStatLabel = screen.getByText("Elevation");
  await expect.element(elevationStatLabel).toBeInTheDocument();

  const timeStatLabel = screen.getByText("Time");
  await expect.element(timeStatLabel).toBeInTheDocument();
}

export async function checkStatValues(
  screen: Locator | RenderResult,
  hikes: number,
  distance: number,
  elevationGain: number,
  duration: number,
) {
  const hikesStatValue = screen.getByText(`${hikes}`);
  await expect.element(hikesStatValue).toBeInTheDocument();

  const distanceStatValue = screen.getByText(formatDistance(distance));
  await expect.element(distanceStatValue).toBeInTheDocument();

  const elevationStatValue = screen.getByText(formatElevation(elevationGain));
  await expect.element(elevationStatValue).toBeInTheDocument();

  const timeStatValue = screen.getByText(formatDuration(duration));
  await expect.element(timeStatValue).toBeInTheDocument();
}
