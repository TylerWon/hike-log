import * as z from "zod";

export const CaptionSchema = z.string();
export const DisplayOrderSchema = z.int().gte(1);

export const PhotoSchema = z.object({
  caption: CaptionSchema,
  displayOrder: DisplayOrderSchema,
  hikeId: z.int(),
  id: z.int(),
  srcUrl: z.url(),
});
export type Photo = z.infer<typeof PhotoSchema>;

export const PhotoListSchema = z.array(PhotoSchema);
