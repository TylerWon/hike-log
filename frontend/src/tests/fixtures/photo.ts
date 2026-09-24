import type { Photo } from "../../schemas/models/photo";

import photo from "../assets/images/grouse_grind_1.avif";

export const PHOTO_FIXTURE: Photo = {
  caption: "Dense forest corridor on the Grind",
  displayOrder: 1,
  hikeId: 1,
  id: 1,
  srcUrl: new URL(photo, import.meta.url).href,
};
