import * as z from "zod";

export const PhotoSchema = z.object({
  caption: z.string(),
  displayOrder: z.coerce.bigint<bigint>().gte(BigInt(1)),
  hikeId: z.coerce.bigint<bigint>(),
  id: z.coerce.bigint<bigint>(),
  srcUrl: z.url(),
});
export type Photo = z.infer<typeof PhotoSchema>;

export const PhotoListSchema = z.array(PhotoSchema);
