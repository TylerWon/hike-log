import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import StarRating from "../../components/StarRating";

/**
 * Checks the number of filled, half, and empty stars matches the provided counts.
 */
function checkStars(screen: RenderResult, filledCount: number, halfCount: number, emptyCount: number) {
  const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
  expect(filledStars.length).toEqual(filledCount);

  const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
  expect(halfStars.length).toEqual(halfCount);

  const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
  expect(emptyStars.length).toEqual(emptyCount);
}

describe("StarRating", () => {
  test("displays 5 filled stars when rating is 5", async () => {
    const screen = await render(<StarRating rating={5} />);
    checkStars(screen, 5, 0, 0);
  });

  test("displays 5 empty stars when rating is 0", async () => {
    const screen = await render(<StarRating rating={0} />);
    checkStars(screen, 0, 0, 5);
  });

  test("displays 3 filled stars, 1 half star, 1 empty star when rating is 3.5", async () => {
    const screen = await render(<StarRating rating={3.5} />);
    checkStars(screen, 3, 1, 1);
  });

  test("clamps ratings below 0 to 0", async () => {
    const screen = await render(<StarRating rating={-2} />);
    checkStars(screen, 0, 0, 5);
  });

  test("clamps ratings above 5 to 5", async () => {
    const screen = await render(<StarRating rating={7} />);
    checkStars(screen, 5, 0, 0);
  });

  test("rounds a rating that isn't a multiple of 0.5 to the nearest 0.5 step", async () => {
    const screen = await render(<StarRating rating={2.2} />);
    checkStars(screen, 2, 0, 3);
  });
});
