// Uploads a file to S3 using the given presigned URL.
// Throws an error if the response is not 200.
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
