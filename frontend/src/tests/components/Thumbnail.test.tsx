import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import type { Photo } from "../../schemas/photo";

import Thumbnail from "../../components/Thumbnail";

describe("Thumbnail", () => {
  test("displays placeholder svg when photo is not provided", async () => {
    const screen = await render(<Thumbnail photo={null} trailName="Grouse Grind" />);

    const svg = screen.getByRole("img", { exact: true, name: "Grouse Grind thumbnail" });
    await expect.element(svg).not.toHaveAttribute("src");
  });

  test("displays image when photo is provided", async () => {
    const photo: Photo = {
      caption: "Dense forest corridor on the Grind",
      hikeId: BigInt(1),
      id: BigInt(1),
      srcUrl: "https://images.unsplash.com/photo-1562593028-1fe2d15bde36?w=800&h=600&fit=crop&auto=format",
    };
    const screen = await render(<Thumbnail photo={photo} trailName="Grouse Grind" />);

    const image = screen.getByRole("img", { exact: true, name: "Grouse Grind thumbnail" });
    await expect.element(image).toHaveAttribute("src", photo.srcUrl);
  });
});
