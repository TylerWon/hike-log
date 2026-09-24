import * as z from "zod";

import { PhotoListSchema } from "./photo";

export const DifficultySchema = z.number().min(0).max(10).multipleOf(0.5);
export const DistanceSchema = z.number().nonnegative();
export const DurationSchema = z.int();
export const ElevationGainSchema = z.int();
export const NotesSchema = z.string();
export const RatingSchema = z.number().min(0).max(5).multipleOf(0.5);
export const TrailNameSchema = z.string();

export const HikeSchema = z.object({
  allTrailsUrl: z.url(),
  date: z.iso.date(), // i.e. YYYY-MM-DD
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
