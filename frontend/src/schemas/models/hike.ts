import * as z from "zod";

import { PhotoListSchema } from "./photo";

export const DifficultySchema = z.number().min(0).max(10).multipleOf(0.5);
export const DistanceSchema = z.number().nonnegative();
export const DurationSchema = z.bigint();
export const ElevationGainSchema = z.bigint();
export const NotesSchema = z.string();
export const RatingSchema = z.number().min(0).max(5).multipleOf(0.5);
export const TrailNameSchema = z.string();

// "coerce" tries to convert the input to the respective type before the schema is parsed. This is useful for number-
// like fields that may have a string number as their value in the raw input as coerce will convert the string to a
// number first so that parsing does not fail.
export const HikeSchema = z.object({
  allTrailsUrl: z.url(),
  date: z.iso.date(), // i.e. YYYY-MM-DD
  difficulty: z.coerce.number().pipe(DifficultySchema),
  distance: z.coerce.number().pipe(DistanceSchema),
  duration: z.coerce.bigint().pipe(DurationSchema),
  elevationGain: z.coerce.bigint().pipe(ElevationGainSchema),
  id: z.coerce.bigint(),
  notes: NotesSchema,
  photos: PhotoListSchema,
  rating: z.coerce.number().pipe(RatingSchema),
  trailName: TrailNameSchema,
});

export const HikeListSchema = z.array(HikeSchema);

export type Hike = z.infer<typeof HikeSchema>;
