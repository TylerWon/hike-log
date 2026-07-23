import * as z from "zod";

import { PhotoListSchema } from "./photo";

const HikeSchema = z.object({
  allTrailsUrl: z.url(),
  date: z.iso.date(), // i.e. YYYY-MM-DD
  difficulty: z.number().min(0).max(10).multipleOf(0.5),
  distance: z.number().nonnegative(),
  duration: z.uint64(),
  elevationGain: z.uint64(),
  id: z.uint64(),
  notes: z.string(),
  photos: PhotoListSchema,
  rating: z.number().min(0).max(5).multipleOf(0.5),
  trailName: z.string()
});

export const HikeListSchema = z.array(HikeSchema)

export type Hike = z.infer<typeof HikeSchema>
