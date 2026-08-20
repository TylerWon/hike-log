import type { Photo } from "../../schemas/photo";

import photo from "../assets/images/grouse_grind_1.avif";

export const PHOTO_FIXTURE: Photo = {
  caption: "Dense forest corridor on the Grind",
  hikeId: BigInt(1),
  id: BigInt(1),
  srcUrl: new URL(photo, import.meta.url).href,
};
