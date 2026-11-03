import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import { HIKE_MOCK_1 } from "../../../../mocks/models/hike";
import { HikeCardWrapper } from "../../../../utils/HikeCard/test-components";

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

describe("HikeCard", () => {
  test("displays overview of the hike", async () => {
    const screen = await render(<HikeCardWrapper />);
    await checkCard(screen, false);
  });

  test("displays more details about the hike when expanded", async () => {
    const screen = await render(<HikeCardWrapper isExpanded={true} />);
    await checkCard(screen, true);
  });
});
