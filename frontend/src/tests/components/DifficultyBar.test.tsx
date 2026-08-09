import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import DifficultyBar from "../../components/DifficultyBar";

/**
 * Checks the number of filled, half, and empty bars matches the provided counts.
 */
function checkBars(screen: RenderResult, filledCount: number, halfCount: number, emptyCount: number) {
  const filledBars = screen.getByRole("img", { exact: true, name: "Filled bar" });
  expect(filledBars.length).toBe(filledCount);

  const halfBars = screen.getByRole("img", { exact: true, name: "Half bar" });
  expect(halfBars.length).toBe(halfCount);

  const emptyBars = screen.getByRole("img", { exact: true, name: "Empty bar" });
  expect(emptyBars.length).toBe(emptyCount);
}

describe("DifficultyBar", () => {
  test("displays 10 empty bars when difficulty is 0", async () => {
    const screen = await render(<DifficultyBar difficulty={0} />);
    checkBars(screen, 0, 0, 10);
  });

  test("displays 10 filled bars when difficulty is 10", async () => {
    const screen = await render(<DifficultyBar difficulty={10} />);
    checkBars(screen, 10, 0, 0);
  });

  test("displays 5 filled bars, 1 half bar, 4 empty bars when difficulty is 5.5", async () => {
    const screen = await render(<DifficultyBar difficulty={5.5} />);
    checkBars(screen, 5, 1, 4);
  });

  test("clamps difficulty below 0 to 0", async () => {
    const screen = await render(<DifficultyBar difficulty={-5} />);
    checkBars(screen, 0, 0, 10);
  });

  test("clamps difficulty above 10 to 10", async () => {
    const screen = await render(<DifficultyBar difficulty={10.5} />);
    checkBars(screen, 10, 0, 0);
  });

  test("rounds a difficulty that isn't a multiple of 0.5 to the nearest 0.5 step", async () => {
    const screen = await render(<DifficultyBar difficulty={7.9} />);
    checkBars(screen, 8, 0, 2);
  });
});
