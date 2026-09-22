import * as z from "zod";

import { PhotoSchema } from "../models/photo";

const VALID_FILE_TYPES = ["image/jpeg", "image/png", "image/webp", "image/heic", "image/heif"];
// const TEN_MB_IN_BYTES = 10485760;
const TEN_MB_IN_BYTES = 50000;

const PhotoFormDataSchema = PhotoSchema.omit({ hikeId: true, id: true, srcUrl: true }).extend({
  file: z.file().mime(VALID_FILE_TYPES).min(0).max(TEN_MB_IN_BYTES, "Image must be under 10MB"),
});

export const PhotoFormDataListSchema = z.array(PhotoFormDataSchema);

export type PhotoFormData = z.infer<typeof PhotoFormDataSchema>;
