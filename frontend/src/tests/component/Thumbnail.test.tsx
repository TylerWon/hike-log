import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import Thumbnail from "../../components/Thumbnail";
import { PHOTO_FIXTURE } from "../fixtures/photo";

describe("Thumbnail", () => {
  test("displays placeholder svg when photo is not provided", async () => {
    const screen = await render(<Thumbnail photo={null} trailName="Grouse Grind" />);

    const svg = screen.getByRole("img");
    await expect.element(svg).not.toHaveAttribute("src");
  });

  test("displays image when photo is provided", async () => {
    const screen = await render(<Thumbnail photo={PHOTO_FIXTURE} trailName="Grouse Grind" />);

    const image = screen.getByRole("img");
    await expect.element(image).toHaveAttribute("src", PHOTO_FIXTURE.srcUrl);
  });
});
