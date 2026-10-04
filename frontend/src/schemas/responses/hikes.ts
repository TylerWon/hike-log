import * as z from "zod";

import { PhotoSchema } from "../models/photo";

const CreatePhotosResponseItemSchema = z.discriminatedUnion("success", [
  // Schema when Success is true
  z.object({
    index: z.int().gte(0),
    result: PhotoSchema,
    success: z.literal(true),
  }),
  // Schema when Success is false
  z.object({
    error: z.string(),
    index: z.int().gte(0),
    success: z.literal(false),
  }),
]);
export const CreatePhotosResponseSchema = z.array(CreatePhotosResponseItemSchema);
export type CreatePhotosResponse = z.infer<typeof CreatePhotosResponseSchema>;

const CreatePresignedUrlsResponseItemResultSchema = z.object({
  objectKey: z.string(),
  presignedUrl: z.url(),
});
const CreatePresignedUrlsResponseItemSchema = z.discriminatedUnion("success", [
  // Schema when Success is true
  z.object({
    index: z.int().gte(0),
    result: CreatePresignedUrlsResponseItemResultSchema,
    success: z.literal(true),
  }),
  // Schema when Success is false
  z.object({
    error: z.string(),
    index: z.int().gte(0),
    success: z.literal(false),
  }),
]);
export const CreatePresignedUrlsResponseSchema = z.array(CreatePresignedUrlsResponseItemSchema);
export type CreatePresignedUrlsResponse = z.infer<typeof CreatePresignedUrlsResponseSchema>;
