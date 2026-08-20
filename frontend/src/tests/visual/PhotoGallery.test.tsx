import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Photo } from "../../schemas/photo";

import PhotoGallery from "../../components/PhotoGallery";
import { HIKE_FIXTURE_1 } from "../fixtures/hike";

async function clickFirstPhotoInGallery(screen: RenderResult) {
  const firstPhotoButton = screen.getByRole("button", { name: "Photo 1 button" });
  await firstPhotoButton.click();
}

describe("PhotoGallery", () => {
  test("displays gallery of photos when photos are provided", async () => {
    const screen = await render(<PhotoGallery photos={HIKE_FIXTURE_1.photos} />);
    const component = screen.getByRole("region", { name: "Photo gallery" });
    await expect(component).toMatchScreenshot();
  });

  test("displays photo in a lightbox when it is clicked in the gallery", async () => {
    const screen = await render(<PhotoGallery photos={HIKE_FIXTURE_1.photos} />);

    await clickFirstPhotoInGallery(screen);

    const component = screen.getByRole("region", { name: "Photo lightbox" });
    await expect(component).toMatchScreenshot();
  });

  test("does not display next and prev buttons in the lightbox when only one photo is provided", async () => {
    const photos = [HIKE_FIXTURE_1.photos[0]];

    const screen = await render(<PhotoGallery photos={photos} />);

    await clickFirstPhotoInGallery(screen);

    const component = screen.getByRole("region", { name: "Photo lightbox" });
    await expect(component).toMatchScreenshot();
  });

  test("does not display caption in the lightbox when the photo has no caption", async () => {
    const photos: Photo[] = [structuredClone(HIKE_FIXTURE_1.photos[0])];
    photos[0].caption = "";

    const screen = await render(<PhotoGallery photos={photos} />);

    await clickFirstPhotoInGallery(screen);

    const component = screen.getByRole("region", { name: "Photo lightbox" });
    await expect(component).toMatchScreenshot();
  });
});
