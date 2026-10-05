import type { Photo } from "../../../schemas/models/photo";

import photo1 from "../../assets/images/joffre_lakes_1.avif";
import photo2 from "../../assets/images/joffre_lakes_2.avif";

export const PHOTO_MOCK_1: Photo = {
  caption: "Joffre Lakes 1",
  displayOrder: 1,
  hikeId: 1,
  id: 1,
  srcUrl: new URL(photo1, import.meta.url).href,
};

export const PHOTO_MOCK_2: Photo = {
  caption: "Joffre Lakes 2",
  displayOrder: 2,
  hikeId: 1,
  id: 2,
  srcUrl: new URL(photo2, import.meta.url).href,
};
