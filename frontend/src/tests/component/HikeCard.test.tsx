import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import { HIKE_MOCK_1 } from "../mocks/models/hike";
import { HikeCardWrapper, StatefulHikeCard } from "../utils/HikeCard/test-components";

/**
 * Checks that the card displays an overview of the hike. If the card is expanded, also verifies that more details about
 * the hike are shown.
 */
async function checkCard(screen: RenderResult, isExpanded: boolean) {
  const overview = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
  await expect.element(overview).toBeInTheDocument();

  const details = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} details` });
  if (!isExpanded) {
    await expect.element(details).not.toBeInTheDocument();
  } else {
    await expect.element(details).toBeInTheDocument();
  }
}

async function clickDetails(screen: RenderResult) {
  const details = screen.getByRole("region", {
    includeHidden: true,
    name: `${HIKE_MOCK_1.trailName} details`,
  });
  await details.click();
}

async function clickOverview(screen: RenderResult) {
  const overview = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
  await overview.click();
}

describe("HikeCard", () => {
  test("displays overview of the hike but no additional details by default", async () => {
    const screen = await render(<HikeCardWrapper />);
    await checkCard(screen, false);
  });

  test("displays more details about the hike when overview is clicked", async () => {
    const screen = await render(<StatefulHikeCard />);
    await checkCard(screen, false);
    await clickOverview(screen);
    await checkCard(screen, true);
  });

  test("hides details about the hike when overview is clicked again", async () => {
    const screen = await render(<StatefulHikeCard />);
    await checkCard(screen, false);
    await clickOverview(screen);
    await checkCard(screen, true);
    await clickDetails(screen);
    await checkCard(screen, true); // details should not be hidden when details section is clicked
    await clickOverview(screen);
    await checkCard(screen, false);
  });
});
