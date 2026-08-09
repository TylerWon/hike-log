import { describe, expect, test } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";

import type { Photo } from "../../schemas/photo";

import PhotoGallery from "../../components/PhotoGallery";

const PHOTOS: Photo[] = [
  {
    caption: "Upper Joffre Lake with the glacier above",
    hikeId: BigInt(1),
    id: BigInt(1),
    srcUrl: "https://images.unsplash.com/photo-1780945008220-7edb56992450?w=800&h=600&fit=crop&auto=format",
  },
  {
    caption: "Middle lake — flat calm in the morning",
    hikeId: BigInt(1),
    id: BigInt(2),
    srcUrl: "https://images.unsplash.com/photo-1763593862654-52f923fa7d85?w=1200&h=800&fit=crop&auto=format",
  },
];

/**
 * Checks that the gallery displays `photos`.
 */
async function checkGallery(screen: RenderResult, photos: Photo[]) {
  const gallery = screen.getByRole("img", { name: /Photo \d+ gallery/ });
  expect(gallery.length).toEqual(photos.length);

  for (let i: number = 0; i < gallery.length; i++) {
    const photo = gallery.nth(i);
    await expect.element(photo).toHaveAttribute("src", photos[i].srcUrl);
  }
}

/**
 * Checks that the photo at position `index` of `photos` is displayed in the lightbox. Also checks if the close, next,
 * and prev buttons are shown. Note that the latter two controls are only shown when there is more than one photo in
 * `photos`.
 */
async function checkLightbox(screen: RenderResult, photos: Photo[], index: number) {
  const photo = screen.getByRole("img", { exact: true, name: `Photo ${index + 1} lightbox` });
  await expect
    .element(photo)
    .toHaveAttribute("src", photos[index].srcUrl.replace(/w=\d+/, "w=1200").replace(/h=\d+/, "h=800"));

  if (photos[index].caption) {
    const caption = screen.getByText(photos[index].caption);
    await expect.element(caption).toBeInTheDocument();
  }

  const count = screen.getByText(`${index + 1} / ${photos.length}`);
  await expect.element(count).toBeInTheDocument();

  const closeButton = screen.getByRole("button", { exact: true, name: "Close button" });
  await expect.element(closeButton).toBeInTheDocument();

  if (photos.length <= 1) {
    const prevButton = screen.getByRole("button", { exact: true, name: "Previous photo button" });
    await expect.element(prevButton).not.toBeInTheDocument();

    const nextButton = screen.getByRole("button", { exact: true, name: "Next photo button" });
    await expect.element(nextButton).not.toBeInTheDocument();
  } else {
    const prevButton = screen.getByRole("button", { exact: true, name: "Previous photo button" });
    await expect.element(prevButton).toBeInTheDocument();

    const nextButton = screen.getByRole("button", { exact: true, name: "Next photo button" });
    await expect.element(nextButton).toBeInTheDocument();
  }
}

async function clickCloseButtonInLightbox(screen: RenderResult) {
  const closeButton = screen.getByRole("button", { exact: true, name: "Close button" });
  await closeButton.click();
}

async function clickFirstPhotoInGallery(screen: RenderResult) {
  const firstPhotoButton = screen.getByRole("button", { exact: true, name: "Photo 1 button" });
  await firstPhotoButton.click();
}

async function clickNextButtonInLightbox(screen: RenderResult) {
  const nextButton = screen.getByRole("button", { exact: true, name: "Next photo button" });
  await nextButton.click();
}

async function clickPrevButtonInLightbox(screen: RenderResult) {
  const prevButton = screen.getByRole("button", { exact: true, name: "Previous photo button" });
  await prevButton.click();
}

describe("PhotoGallery", () => {
  test("displays nothing when no photos are provided", async () => {
    const screen = await render(<PhotoGallery photos={[]} />);
    expect(screen.container.hasChildNodes()).toBeFalsy();
  });

  test("displays gallery of photos when photos are provided", async () => {
    const screen = await render(<PhotoGallery photos={PHOTOS} />);
    await checkGallery(screen, PHOTOS);
  });

  test("displays photo in a lightbox when it is clicked in the gallery", async () => {
    const screen = await render(<PhotoGallery photos={PHOTOS} />);

    await checkGallery(screen, PHOTOS);
    await clickFirstPhotoInGallery(screen);
    await checkLightbox(screen, PHOTOS, 0);
  });

  test("switches between photos when next and prev buttons are clicked in the lightbox", async () => {
    const screen = await render(<PhotoGallery photos={PHOTOS} />);

    await checkGallery(screen, PHOTOS);
    await clickFirstPhotoInGallery(screen);
    await checkLightbox(screen, PHOTOS, 0);
    await clickNextButtonInLightbox(screen);
    await checkLightbox(screen, PHOTOS, 1);
    await clickNextButtonInLightbox(screen);
    await checkLightbox(screen, PHOTOS, 0);
    await clickPrevButtonInLightbox(screen);
    await checkLightbox(screen, PHOTOS, 1);
  });

  test("does not display next and prev buttons in the lightbox when only one photo is provided", async () => {
    const photos: Photo[] = [structuredClone(PHOTOS[0])];

    const screen = await render(<PhotoGallery photos={photos} />);

    await checkGallery(screen, photos);
    await clickFirstPhotoInGallery(screen);
    await checkLightbox(screen, photos, 0);
  });

  test("does not display caption in the lightbox when the photo has no caption", async () => {
    const photos: Photo[] = [structuredClone(PHOTOS[0])];
    photos[0].caption = "";

    const screen = await render(<PhotoGallery photos={photos} />);

    await checkGallery(screen, photos);
    await clickFirstPhotoInGallery(screen);
    await checkLightbox(screen, photos, 0);
  });

  test("closes lightbox when close button is clicked", async () => {
    const screen = await render(<PhotoGallery photos={PHOTOS} />);

    await checkGallery(screen, PHOTOS);
    await clickFirstPhotoInGallery(screen);
    await checkLightbox(screen, PHOTOS, 0);
    await clickCloseButtonInLightbox(screen);

    const photo = screen.getByRole("img", { exact: true, name: `Photo 1 lightbox` });
    await expect.element(photo).not.toBeInTheDocument();
  });
});
