import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import DifficultyBar from "../../components/DifficultyBar";

describe("DifficultyBar", () => {
  test("displays 10 empty bars when difficulty is 0", async () => {
    const screen = await render(<DifficultyBar difficulty={0} />);
    const component = screen.getByRole("img", { name: "0 out of 10 difficulty" });
    await expect(component).toMatchScreenshot();
  });

  test("displays 10 filled bars when difficulty is 10", async () => {
    const screen = await render(<DifficultyBar difficulty={10} />);
    const component = screen.getByRole("img", { name: "10 out of 10 difficulty" });
    await expect(component).toMatchScreenshot();
  });

  test("displays 5 filled, 1 half, 4 empty bars when difficulty is 5.5", async () => {
    const screen = await render(<DifficultyBar difficulty={5.5} />);
    const component = screen.getByRole("img", { name: "5.5 out of 10 difficulty" });
    await expect(component).toMatchScreenshot();
  });
});
