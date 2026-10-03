import * as z from "zod";

import { HikeSchema } from "../models/hike";
import { PhotoSchema } from "../models/photo";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const CreateHikeRequestSchema = HikeSchema.omit({ id: true, photos: true });
export type CreateHikeRequest = z.infer<typeof CreateHikeRequestSchema>;

const CreatePhotosRequestItemSchema = PhotoSchema.omit({ hikeId: true, id: true, srcUrl: true }).extend({
  objectKey: z.string(),
});
const CreatePhotosRequestSchema = z.array(CreatePhotosRequestItemSchema);
export type CreatePhotosRequest = z.infer<typeof CreatePhotosRequestSchema>;

interface CreatePresignedUrlsRequestItem {
  contentLength: number;
  contentType: string;
}
export type CreatePresignedUrlsRequest = CreatePresignedUrlsRequestItem[];
