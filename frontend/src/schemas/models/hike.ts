import * as z from "zod";

import { PhotoListSchema } from "./photo";

export const AllTrailsUrlSchema = z.url();
export const DateSchema = z.iso.date(); // i.e. YYYY-MM-DD
export const DifficultySchema = z.number().min(0).max(10).multipleOf(0.5);
export const DistanceSchema = z.number().nonnegative();
export const DurationSchema = z.int().nonnegative();
export const ElevationGainSchema = z.int().nonnegative();
export const NotesSchema = z.string();
export const RatingSchema = z.number().min(0).max(5).multipleOf(0.5);
export const TrailNameSchema = z.string();

export const HikeSchema = z.object({
  allTrailsUrl: AllTrailsUrlSchema,
  date: DateSchema,
  difficulty: DifficultySchema,
  distance: DistanceSchema,
  duration: DurationSchema,
  elevationGain: ElevationGainSchema,
  id: z.int(),
  notes: NotesSchema,
  photos: PhotoListSchema,
  rating: RatingSchema,
  trailName: TrailNameSchema,
});
export type Hike = z.infer<typeof HikeSchema>;

export const HikeListSchema = z.array(HikeSchema);
