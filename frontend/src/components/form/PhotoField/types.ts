import type { $ZodErrorTree } from "zod/v4/core";

import type { HikeFormPhotoData } from "../../../schemas/forms/hike";

export interface PhotoData {
  caption: string;
  file: File;
  previewUrl: string;
}

export interface PhotoErrors {
  errors: string[];
  items?: $ZodErrorTree<HikeFormPhotoData>[];
}
