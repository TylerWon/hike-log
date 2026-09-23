import * as z from "zod";

import {
  DifficultySchema,
  DistanceSchema,
  DurationSchema,
  ElevationGainSchema,
  HikeSchema,
  NotesSchema,
  RatingSchema,
  TrailNameSchema,
} from "../models/hike";
import { PhotoFormDataListSchema } from "./photo";

const IntegerStringSchema = z.string().min(1, "Required").regex(/^\d+$/, "Must be a whole number").transform(BigInt);

// HikeFormDataSchema Has the same fields as HikeSchema (excluding id) but with extra validation (empty values not
// allowed)
export const HikeFormDataSchema = HikeSchema.omit({ id: true }).extend({
  difficulty: z.string().min(1, "Required").transform(Number).pipe(DifficultySchema),
  distance: z.string().min(1, "Required").transform(Number).pipe(DistanceSchema),
  duration: IntegerStringSchema.pipe(DurationSchema),
  elevationGain: IntegerStringSchema.pipe(ElevationGainSchema),
  notes: NotesSchema.min(1, "Required"),
  photos: PhotoFormDataListSchema,
  rating: z.string().min(1, "Required").transform(Number).pipe(RatingSchema),
  trailName: TrailNameSchema.min(1, "Required"),
});
export type HikeFormData = z.infer<typeof HikeFormDataSchema>;
