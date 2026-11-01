import * as z from "zod";

import {
  AllTrailsUrlSchema,
  DateSchema,
  DifficultySchema,
  DistanceSchema,
  DurationSchema,
  ElevationGainSchema,
  NotesSchema,
  RatingSchema,
  TrailNameSchema,
} from "../models/hike";
import { CaptionSchema, DisplayOrderSchema } from "../models/photo";

const VALID_FILE_TYPES = ["image/jpeg", "image/png", "image/webp", "image/heic", "image/heif"];
export const MAX_PHOTO_SIZE = 10485760; // 10 MB in bytes

const IntegerStringValidation = z.string().min(1, "Required").regex(/^\d+$/, "Must be a positive integer");
const NumberStringValidation = z.string().min(1, "Required");

const HikeFormPhotoDataSchema = z.object({
  caption: CaptionSchema,
  displayOrder: DisplayOrderSchema,
  file: z.file().mime(VALID_FILE_TYPES).min(0).max(MAX_PHOTO_SIZE, "Image must be under 10MB"),
});
export type HikeFormPhotoData = z.infer<typeof HikeFormPhotoDataSchema>;

export const HikeFormPhotoDataListSchema = z.array(HikeFormPhotoDataSchema);

// Raw form values are strings so for fields with an actual type that is not a string, this schema will transform
// the value to the correct type after some initial validation.
export const HikeFormDataSchema = z.object({
  // Check if string is empty first otherwise empty value errors with "Invalid URL"
  allTrailsUrl: z.string().min(1, "Required").pipe(AllTrailsUrlSchema),

  // A date input will remain empty when the date entered is invalid so do not add validation for an empty value
  // otherwise an invalid date errors with "Required"
  date: DateSchema,

  difficulty: NumberStringValidation.transform(Number).pipe(DifficultySchema),
  distance: NumberStringValidation.transform(Number).pipe(DistanceSchema),
  duration: IntegerStringValidation.transform(Number).pipe(DurationSchema),
  elevationGain: IntegerStringValidation.transform(Number).pipe(ElevationGainSchema),
  notes: NotesSchema.min(1, "Required"),
  photos: HikeFormPhotoDataListSchema,
  rating: NumberStringValidation.transform(Number).pipe(RatingSchema),
  trailName: TrailNameSchema.min(1, "Required"),
});
export type HikeFormData = z.infer<typeof HikeFormDataSchema>;
