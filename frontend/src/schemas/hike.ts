import * as z from "zod";

import { PhotoListSchema } from "./photo";

const HikeSchema = z.object({
  allTrailsUrl: z.url(),
  date: z.iso.date(), // i.e. YYYY-MM-DD
  difficulty: z.number().min(0).max(10).multipleOf(0.5),
  distance: z.number().nonnegative(),
  // "coerce" tries to convert the input to a bigint when the schema is parsed. This prevents an error that occurs when
  // parsing a Response that was deserialized with json() as json() deserializes all number-like values to the Number type.
  duration: z.coerce.bigint<bigint>(),
  elevationGain: z.coerce.bigint<bigint>(),
  id: z.coerce.bigint<bigint>(),
  notes: z.string(),
  photos: PhotoListSchema,
  rating: z.number().min(0).max(5).multipleOf(0.5),
  trailName: z.string(),
});

export const HikeListSchema = z.array(HikeSchema);

export type Hike = z.infer<typeof HikeSchema>;
