import * as z from "zod";

import { CaptionSchema, DisplayOrderSchema } from "../models/photo";

const VALID_FILE_TYPES = ["image/jpeg", "image/png", "image/webp", "image/heic", "image/heif"];
const TEN_MB_IN_BYTES = 10485760;

const PhotoFormDataSchema = z.object({
  caption: CaptionSchema,
  displayOrder: DisplayOrderSchema,
  file: z.file().mime(VALID_FILE_TYPES).min(0).max(TEN_MB_IN_BYTES, "Image must be under 10MB"),
});
export type PhotoFormData = z.infer<typeof PhotoFormDataSchema>;

export const PhotoFormDataListSchema = z.array(PhotoFormDataSchema);
