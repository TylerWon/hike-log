import type { HikeFormData } from "../../schemas/forms/hike";

import { createHike, createPhotos, createPresignedUrls } from "../../api/hikes";
import { uploadFile } from "../../api/s3";

class HikeCreationError extends Error {
  constructor(message: string, options?: ErrorOptions) {
    super(message, options);
    this.name = this.constructor.name;
  }
}

export class PhotoCreationError extends Error {
  constructor(message: string, options?: ErrorOptions) {
    super(message, options);
    this.name = this.constructor.name;
  }
}

// Handles submission of HikeForm data.
// Creates the Hike and any Photos that were uploaded. Photo creation is a multi-step process that involves obtaining
// a S3 presigned PUT URL, using those URLs to upload the photos to S3, and finally creating the Photo models in the
// backend.
export async function submitHikeForm(formData: HikeFormData): Promise<void> {
  const { photos, ...hikeData } = formData;

  let hike;
  try {
    hike = await createHike(hikeData);
  } catch (e) {
    throw new HikeCreationError(`Failed to create hike: ${e}`, { cause: e });
  }

  if (photos.length < 1) {
    return;
  }

  const photoErrors = [];
  let createPresignedUrlsResponse;
  try {
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

    createPresignedUrlsResponse = await createPresignedUrls(hike.id, reqBody);
  } catch (e) {
    throw new PhotoCreationError(`Failed to create presigned URLs: ${e}`, { cause: e });
  }

  let s3UploadResponses;
  try {
    const promises = [];
    for (let i = 0; i < photos.length; i++) {
      const createPresignedUrlResponse = createPresignedUrlsResponse[i];
      if (createPresignedUrlResponse.success) {
        promises.push(uploadFile(createPresignedUrlResponse.result.presignedUrl, photos[i].file));
      } else {
        photoErrors.push(`Failed to create presigned URL for photo ${i}`);
        promises.push(Promise.reject());
      }
    }

    s3UploadResponses = await Promise.allSettled(promises);
  } catch (e) {
    throw new PhotoCreationError(`Failed to upload to S3: ${e}`, { cause: e });
  }

  let createPhotosResponse;
  try {
    const reqBody = [];
    for (let i = 0; i < photos.length; i++) {
      const createPresignedUrlResponse = createPresignedUrlsResponse[i];
      if (!createPresignedUrlResponse.success) {
        continue;
      }

      const s3UploadResponse = s3UploadResponses[i];
      if (s3UploadResponse.status == "fulfilled") {
        const reqItem = {
          caption: photos[i].caption,
          displayOrder: photos[i].displayOrder,
          index: i,
          objectKey: createPresignedUrlResponse.result.objectKey,
        };
        reqBody.push(reqItem);
      } else {
        photoErrors.push(`Failed to upload photo ${i} to S3`);
      }
    }

    createPhotosResponse = await createPhotos(hike.id, reqBody);
  } catch (e) {
    throw new PhotoCreationError(`Failed to create models: ${e}`, { cause: e });
  }

  for (const createPhotoResponse of createPhotosResponse) {
    if (!createPhotoResponse.success) {
      photoErrors.push(`Failed to create Photo model for photo ${createPhotoResponse.index}`);
    }
  }

  if (photoErrors.length > 0) {
    throw new PhotoCreationError(`Some photos could not be created: ${photoErrors}`);
  }
}
