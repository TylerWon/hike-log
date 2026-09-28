import * as z from "zod";

export const CreatePresignedUrlResponseSchema = z.object({
  objectKey: z.string(),
  uploadUrl: z.url(),
});
export type CreatePresignedUrlResponse = z.infer<typeof CreatePresignedUrlResponseSchema>;
