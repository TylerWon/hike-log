import * as z from "zod";

const PhotoSchema = z.object({
  caption: z.string(),
  // "coerce" tries to convert the input to a bigint when the schema is parsed. This prevents an error that occurs when
  // parsing a Response that was deserialized with json() as this json() deserializes all numbers to the Number type.
  hikeId: z.coerce.bigint(),
  id: z.coerce.bigint(),
  srcUrl: z.url(),
});

export const PhotoListSchema = z.array(PhotoSchema);

export type Photo = z.infer<typeof PhotoSchema>;
