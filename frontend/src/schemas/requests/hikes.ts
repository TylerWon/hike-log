import * as z from "zod";

import { HikeSchema } from "../models/hike";
import { PhotoSchema } from "../models/photo";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const CreateHikeRequestSchema = HikeSchema.omit({ id: true, photos: true });
export type CreateHikeRequest = z.infer<typeof CreateHikeRequestSchema>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const CreatePhotoRequestSchema = PhotoSchema.omit({ hikeId: true, id: true, srcUrl: true }).extend({
  objectKey: z.string(),
});
export type CreatePhotoRequest = z.infer<typeof CreatePhotoRequestSchema>;

interface CreatePresignedUrlsRequestItem {
  contentLength: number;
  contentType: string;
}
export type CreatePresignedUrlsRequest = CreatePresignedUrlsRequestItem[];
