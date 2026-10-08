import type { $ZodErrorTree } from "zod/v4/core";

import type { PhotoFormData } from "../../schemas/forms/photo";

export interface PhotoData {
  caption: string;
  file: File;
  previewUrl: string;
}

export interface PhotoErrors {
  errors: string[];
  items?: $ZodErrorTree<PhotoFormData>[];
}
