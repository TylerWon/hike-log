import * as z from "zod";

const PhotoSchema = z.object({
  caption: z.string(),
  hikeId: z.uint64(),
  id: z.uint64(),
  srcUrl: z.url(),
});

export const PhotoListSchema = z.array(PhotoSchema)

export type Photo = z.infer<typeof PhotoSchema>;
