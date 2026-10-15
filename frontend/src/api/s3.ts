import type { UploadPhotosRequest } from "../schemas/requests/s3";
import type { UploadPhotosResponse } from "../schemas/responses/s3";

/**
 * Uploads a file to S3 using the given presigned URL.
 *
 * Throws an error if the response is not 200.
 */
export async function uploadFile(presignedUrl: string, file: File): Promise<void> {
  const response = await fetch(presignedUrl, {
    body: file,
    headers: {
      "Content-Length": String(file.size),
      "Content-Type": file.type,
    },
    method: "PUT",
  });

  if (response.status !== 200) {
    throw new Error(`Failed to upload file to S3: ${response.status} - ${response.statusText}`);
  }
}

/**
 * Uploads a collection of photos to S3 using the given presigned URLs.
 */
export async function uploadPhotos(reqBody: UploadPhotosRequest): Promise<UploadPhotosResponse> {
  const promises = [];
  for (const item of reqBody) {
    // S3 doesn't support batch uploads when using presigned URLs so upload each photo individually with uploadFile()
    promises.push(
      uploadFile(item.presignedUrl, item.photo).then(
        // on resolved
        () => ({ index: item.index, success: true }),
        // on rejected
        () => ({ index: item.index, success: false }),
      ),
    );
  }

  return await Promise.all(promises);
}
