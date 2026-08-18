import type { Photo } from "../../schemas/photo";

import photo1 from "../assets/photo1.avif";

export const PHOTO_FIXTURE: Photo = {
  caption: "Dense forest corridor on the Grind",
  hikeId: BigInt(1),
  id: BigInt(1),
  srcUrl: photo1,
};
