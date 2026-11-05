import type { HikeFormData } from "../../../schemas/forms/hike";
import type { HikeFormPhotoData } from "../../../schemas/forms/hike";
import type { CreatePhotosResponse, CreatePresignedUrlsResponse } from "../../../schemas/responses/hikes";
import type { UploadPhotosResponse } from "../../../schemas/responses/s3";

import { AuthError } from "../../../network/api/errors";
import { createHike, createPhotos, createPresignedUrls } from "../../../network/api/hikes";
import { uploadPhotos } from "../../../network/s3/s3";
import { HikeCreationError, PhotoCreationError } from "./errors";

/**
 * Handles submission of data collected by a HikeForm.
 *
 * Creates a hike and any photos that were selected by the user.
 */
export async function submitHikeForm(formData: HikeFormData): Promise<void> {
  const { photos, ...hikeData } = formData;

  let hike;
  try {
    hike = await createHike(hikeData);
  } catch (e) {
    if (e instanceof AuthError) {
      // Propogate AuthError up
      throw e;
    }

    // Consolidate other errors into HikeCreationError then throw
    throw new HikeCreationError(`Failed to create hike: ${e}`);
  }

  if (photos.length < 1) {
    return;
  }

  try {
    const partialErrors = await uploadAndCreatePhotos(photos, hike.id);
    if (partialErrors.size > 0) {
      throw new PhotoCreationError(`Failed to create some photos: ${Array.from(partialErrors.values())}`);
    }
  } catch (e) {
    if (e instanceof AuthError) {
      // Propogate AuthError up
      throw e;
    }

    // Consolidate other errors into PhotoCreationError then throw
    throw new PhotoCreationError(`Failed to create photos: ${e}`);
  }
}

/**
 * Uploads and creates photos from data submitted in a HikeForm.
 *
 * Photo creation is a multi-step process that involves obtaining a S3 presigned PUT URL, using those URLs to upload the
 * photos to S3, and finally creating the Photo models in the backend.
 *
 * Returns errors when there are partial failures: creation steps complete, but some individual photos fail to be
 * uploaded or created. If one of the creation steps outright fail, throws an error instead.
 */
export async function uploadAndCreatePhotos(photos: HikeFormPhotoData[], hikeId: number): Promise<Map<number, string>> {
  const errors: Map<number, string> = new Map();

  const createPresignedUrlsResponse = await createPresignedUrlsHelper(photos, hikeId);
  for (const response of createPresignedUrlsResponse) {
    if (!errors.has(response.index) && !response.success) {
      errors.set(response.index, `Failed to create presigned URL for photo ${response.index}`);
    }
  }

  const uploadPhotosResponse = await uploadPhotosHelper(photos, createPresignedUrlsResponse);
  for (const response of uploadPhotosResponse) {
    if (!errors.has(response.index) && !response.success) {
      errors.set(response.index, `Failed to upload photo ${response.index} to S3`);
    }
  }

  const createPhotosResponse = await createPhotosHelper(
    photos,
    hikeId,
    createPresignedUrlsResponse,
    uploadPhotosResponse,
  );
  for (const response of createPhotosResponse) {
    if (!errors.has(response.index) && !response.success) {
      errors.set(response.index, `Failed to create model for photo ${response.index}`);
    }
  }

  return errors;
}

/**
 * Creates presigned PUT URLs for each provided photo. The URL can be used to upload the photo to S3.
 */
export async function createPresignedUrlsHelper(
  photos: HikeFormPhotoData[],
  hikeId: number,
): Promise<CreatePresignedUrlsResponse> {
  const reqBody = [];
  for (let i = 0; i < photos.length; i++) {
    const photo = photos[i];
    const reqItem = {
      contentLength: photo.file.size,
      contentType: photo.file.type,
      index: i,
    };
    reqBody.push(reqItem);
  }

  return await createPresignedUrls(hikeId, reqBody);
}

/**
 * Uploads photos to S3 using the provided presigned PUT URLs. Skips photos that failed to have a presigned URL created.
 */
export async function uploadPhotosHelper(
  photos: HikeFormPhotoData[],
  createPresignedUrlsResponse: CreatePresignedUrlsResponse,
): Promise<UploadPhotosResponse> {
  const reqBody = [];
  for (const response of createPresignedUrlsResponse) {
    if (!response.success) {
      continue;
    }

    const photoIndex = response.index;
    const reqItem = {
      index: photoIndex,
      photo: photos[photoIndex].file,
      presignedUrl: response.result.presignedUrl,
    };
    reqBody.push(reqItem);
  }

  return await uploadPhotos(reqBody);
}

/**
 * Creates Photo models for each provided photo. Skips photos that failed to have a presigned URL created or could not
 * be uploaded to S3.
 */
export async function createPhotosHelper(
  photos: HikeFormPhotoData[],
  hikeId: number,
  createPresignedUrlsResponse: CreatePresignedUrlsResponse,
  uploadPhotosResponse: UploadPhotosResponse,
): Promise<CreatePhotosResponse> {
  const reqBody = [];
  for (const uploadPhotoResponse of uploadPhotosResponse) {
    const photoIndex = uploadPhotoResponse.index;
    const createPresignedUrlResponse = createPresignedUrlsResponse.find((response) => response.index === photoIndex);

    if (!createPresignedUrlResponse?.success || !uploadPhotoResponse.success) {
      continue;
    }

    const reqItem = {
      caption: photos[photoIndex].caption,
      displayOrder: photos[photoIndex].displayOrder,
      index: photoIndex,
      objectKey: createPresignedUrlResponse.result.objectKey,
    };
    reqBody.push(reqItem);
  }

  return await createPhotos(hikeId, reqBody);
}
