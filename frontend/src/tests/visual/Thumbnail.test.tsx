import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import Thumbnail from "../../components/Thumbnail";
import { PHOTO_MOCK_1 } from "../mocks/models/photo";

describe("Thumbnail", () => {
  test("displays placeholder image when no photo is provided", async () => {
    const screen = await render(<Thumbnail photo={null} trailName="Grouse Grind" />);
    const component = screen.getByRole("region", { name: "Grouse Grind thumbnail" });
    await expect(component).toMatchScreenshot();
  });

  test("displays image when photo is provided", async () => {
    const screen = await render(<Thumbnail photo={PHOTO_MOCK_1} trailName="Grouse Grind" />);
    const component = screen.getByRole("region", { name: "Grouse Grind thumbnail" });
    await expect(component).toMatchScreenshot();
  });
});
