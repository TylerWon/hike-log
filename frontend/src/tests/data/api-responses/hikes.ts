import type { Hike } from "../../../schemas/models/hike";
import type { CreatePhotosResponse, CreatePresignedUrlsResponse } from "../../../schemas/responses/hikes";

import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../fixtures/hike";
import { PHOTO_FIXTURE_1, PHOTO_FIXTURE_2 } from "../fixtures/photo";

export const LIST_HIKES_API_RESPONSE: Hike[] = [HIKE_FIXTURE_1, HIKE_FIXTURE_2];

export const CREATE_HIKE_API_RESPONSE: Hike = HIKE_FIXTURE_1;

export const CREATE_PHOTOS_API_RESPONSE: CreatePhotosResponse = [
  {
    index: 0,
    result: PHOTO_FIXTURE_1,
    success: true,
  },
  {
    index: 1,
    result: PHOTO_FIXTURE_2,
    success: true,
  },
];

export const CREATE_PRESIGNED_URLS_API_RESPONSE: CreatePresignedUrlsResponse = [
  {
    index: 0,
    result: {
      objectKey: "hikes/1/photos/dbcf4c30-b23b-4fa4-b350-2620a50f0736",
      presignedUrl: "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/dbcf4c30-b23b-4fa4-b350-2620a50f0736?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20261005%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20261005T042400Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost&x-id=PutObject&X-Amz-Signature=<64-char-hex>",
    },
    success: true,
  },
  {
    index: 1,
    result: {
      objectKey: "hikes/1/photos/e78f95c0-0bac-4071-afce-c0305ed2a7e1",
      presignedUrl: "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/e78f95c0-0bac-4071-afce-c0305ed2a7e1?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20261005%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20261005T042400Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost&x-id=PutObject&X-Amz-Signature=<64-char-hex>",
    },
    success: true,
  },
];
