import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import StarRating from "../../components/StarRating";

describe("StarRating", () => {
  test("displays 5 filled stars when rating is 5", async () => {
    const screen = await render(<StarRating rating={5} />);

    const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
    expect(filledStars.length).toBe(5);

    const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
    expect(halfStars.length).toBe(0);

    const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
    expect(emptyStars.length).toBe(0);
  });

  test("displays 5 empty stars when rating is 0", async () => {
    const screen = await render(<StarRating rating={0} />);

    const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
    expect(filledStars.length).toBe(0);

    const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
    expect(halfStars.length).toBe(0);

    const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
    expect(emptyStars.length).toBe(5);
  });

  test("displays 3 filled stars, 1 half star, 1 empty star when rating is 3.5", async () => {
    const screen = await render(<StarRating rating={3.5} />);

    const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
    expect(filledStars.length).toBe(3);

    const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
    expect(halfStars.length).toBe(1);

    const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
    expect(emptyStars.length).toBe(1);
  });

  test("clamps ratings below 0 to 0", async () => {
    const screen = await render(<StarRating rating={-2} />);

    const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
    expect(filledStars.length).toBe(0);

    const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
    expect(halfStars.length).toBe(0);

    const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
    expect(emptyStars.length).toBe(5);
  });

  test("clamps ratings above 5 to 5", async () => {
    const screen = await render(<StarRating rating={7} />);

    const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
    expect(filledStars.length).toBe(5);

    const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
    expect(halfStars.length).toBe(0);

    const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
    expect(emptyStars.length).toBe(0);
  });

  test("rounds a rating that isn't a multiple of 0.5 to the nearest 0.5 step", async () => {
    const screen = await render(<StarRating rating={2.2} />);

    const filledStars = screen.getByRole("img", { exact: true, name: "Filled star" });
    expect(filledStars.length).toBe(2);

    const halfStars = screen.getByRole("img", { exact: true, name: "Half star" });
    expect(halfStars.length).toBe(0);

    const emptyStars = screen.getByRole("img", { exact: true, name: "Empty star" });
    expect(emptyStars.length).toBe(3);
  });
});
