import type { HikeFormData } from "../../schemas/forms/hike";
import type { PhotoFormData } from "../../schemas/forms/photo";
import type { CreatePhotosResponse, CreatePresignedUrlsResponse } from "../../schemas/responses/hikes";
import type { UploadPhotosResponse } from "../../schemas/responses/s3";

import { createHike, createPhotos, createPresignedUrls } from "../../api/hikes";
import { uploadPhotos } from "../../api/s3";
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
    throw new HikeCreationError(`Failed to create hike: ${e}`);
  }

  if (photos.length < 1) {
    return;
  }

  const errors = await uploadAndCreatePhotos(photos, hike.id);
  if (errors.size > 0) {
    throw new PhotoCreationError(`${Array.from(errors.values())}`);
  }
}

/**
 * Uploads and creates photos from data submitted in a HikeForm.
 *
 * Photo creation is a multi-step process that involves obtaining a S3 presigned PUT URL, using those URLs to upload the
 * photos to S3, and finally creating the Photo models in the backend.
 *
 * Returns any errors encountered.
 */
export async function uploadAndCreatePhotos(photos: PhotoFormData[], hikeId: number): Promise<Map<number, string>> {
  const errors: Map<number, string> = new Map();

  try {
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
  } catch (e) {
    return new Map([[0, `Failed to create photos: ${e}`]]);
  }

  return errors;
}

/**
 * Creates presigned PUT URLs for each provided photo. The URL can be used to upload the photo to S3.
 */
export async function createPresignedUrlsHelper(
  photos: PhotoFormData[],
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
  photos: PhotoFormData[],
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
  photos: PhotoFormData[],
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
