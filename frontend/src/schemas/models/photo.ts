import * as z from "zod";

export const PhotoSchema = z.object({
  caption: z.string(),
  displayOrder: z.int().gte(1),
  hikeId: z.int(),
  id: z.int(),
  srcUrl: z.url(),
});
export type Photo = z.infer<typeof PhotoSchema>;

export const PhotoListSchema = z.array(PhotoSchema);
