import { useState } from "react";
import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import HikeCard from "../../components/HikeCard";
import { HIKE_FIXTURE_1 } from "../fixtures/hike";

/**
 * Checks that the card displays an overview of the hike. If the card is expanded, also verifies that more details about
 * the hike are shown.
 */
async function checkCard(screen: RenderResult, isExpanded: boolean) {
  const overview = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} overview` });
  await expect.element(overview).toBeInTheDocument();

  // The details section appears in the DOM regardless of whether it should be visible to the user or not. We check
  // whether it is visible below
  const details = screen.getByRole("region", {
    includeHidden: true,
    name: `${HIKE_FIXTURE_1.trailName} details`,
  });
  await expect.element(details).toBeInTheDocument();

  const detailsWrapper = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} details wrapper` });

  if (!isExpanded) {
    // When the card is not expanded, the details section be hidden from view. The wrapper of the detail section handles
    // this by collapsing it to zero height
    await expect.poll(() => detailsWrapper.element().getBoundingClientRect().height).toEqual(0);
  } else {
    await expect.poll(() => detailsWrapper.element().getBoundingClientRect().height).toBeGreaterThan(0);
  }
}

async function clickDetails(screen: RenderResult) {
  const details = screen.getByRole("region", {
    includeHidden: true,
    name: `${HIKE_FIXTURE_1.trailName} details`,
  });
  await details.click();
}

async function clickOverview(screen: RenderResult) {
  const overview = screen.getByRole("region", { name: `${HIKE_FIXTURE_1.trailName} overview` });
  await overview.click();
}

// Wrapper around a HikeCard which allows it to be expanded/unexpanded
function StatefulHikeCard() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false);

  const handleClick = () => {
    setIsExpanded(!isExpanded);
  };

  return <HikeCard hike={HIKE_FIXTURE_1} index={1} isExpanded={isExpanded} onClick={handleClick} />;
}

describe("HikeCard", () => {
  test("displays overview of the hike but no additional details by default", async () => {
    const screen = await render(<HikeCard hike={HIKE_FIXTURE_1} index={1} isExpanded={false} onClick={() => {}} />);
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
