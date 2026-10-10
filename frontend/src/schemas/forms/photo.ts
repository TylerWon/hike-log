import * as z from "zod";

import { CaptionSchema, DisplayOrderSchema } from "../models/photo";

const VALID_FILE_TYPES = ["image/jpeg", "image/png", "image/webp", "image/heic", "image/heif"];
export const MAX_PHOTO_SIZE = 10485760; // 10 MB in bytes

const PhotoFormDataSchema = z.object({
  caption: CaptionSchema,
  displayOrder: DisplayOrderSchema,
  file: z.file().mime(VALID_FILE_TYPES).min(0).max(MAX_PHOTO_SIZE, "Image must be under 10MB"),
});
export type PhotoFormData = z.infer<typeof PhotoFormDataSchema>;

export const PhotoFormDataListSchema = z.array(PhotoFormDataSchema);
