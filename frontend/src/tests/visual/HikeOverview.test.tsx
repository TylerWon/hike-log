import { beforeEach, describe, expect, type Mock, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { userEvent } from "vitest/browser";

import { HIKE_MOCK_1 } from "../mocks/models/hike";
import { HikeOverviewWrapper } from "../utils/HikeOverview/test-components";

describe("HikeOverview", () => {
  let mockOnClick: Mock;
  let mockOnDelete: Mock;

  beforeEach(() => {
    mockOnClick = vi.fn();
    mockOnDelete = vi.fn();
  });

  test("displays overview", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
    await expect(component).toMatchScreenshot();
  });

  test("inverts chevron icon when expanded", async () => {
    const screen = await render(
      <HikeOverviewWrapper isExpanded={true} onClick={mockOnClick} onDelete={mockOnDelete} />,
    );
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
    await expect(component).toMatchScreenshot();
  });

  test("shows delete button and enlarges thumbnail when user hovers over", async () => {
    const screen = await render(<HikeOverviewWrapper onClick={mockOnClick} onDelete={mockOnDelete} />);
    const component = screen.getByRole("region", { name: `${HIKE_MOCK_1.trailName} overview` });
    await userEvent.hover(component);
    await expect(component).toMatchScreenshot();
  });
});
