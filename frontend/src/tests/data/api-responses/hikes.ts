import type { Hike } from "../../../schemas/models/hike";
import type { CreatePhotosResponse, CreatePresignedUrlsResponse } from "../../../schemas/responses/hikes";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../fixtures/hike";
import { PHOTO_FIXTURE_1, PHOTO_FIXTURE_2 } from "../fixtures/photo";

export const LIST_HIKES_API_RESPONSE: Hike[] = [
  HIKE_FIXTURE_1,
  HIKE_FIXTURE_2
];

export const CREATE_HIKE_API_RESPONSE: Hike = HIKE_FIXTURE_1;

export const CREATE_PHOTOS_API_RESPONSE: CreatePhotosResponse = [
  {
    index: 0,
    success: true,
    result: PHOTO_FIXTURE_1,
  },
  {
    index: 1,
    success: true,
    result: PHOTO_FIXTURE_2,
  }
]

export const CREATE_PRESIGNED_URLS_API_RESPONSE: CreatePresignedUrlsResponse = [
  {
    index: 0,
    success: true,
    result: {
      // invalid values but OK for testing
      objectKey: "1",
      presignedUrl: "https://www.test-url.com/1"
    },
  },
  {
    index: 1,
    success: true,
    result: {
      // invalid values but OK for testing
      objectKey: "2",
      presignedUrl: "https://www.test-url.com/2"
    },
  }
]