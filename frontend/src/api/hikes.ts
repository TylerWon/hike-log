import type { CreateHikeRequest, CreatePhotoRequest, CreatePresignedUrlRequest } from "../schemas/requests/hikes";

import { type Hike, HikeListSchema, HikeSchema } from "../schemas/models/hike";
import { type Photo, PhotoSchema } from "../schemas/models/photo";
import { type CreatePresignedUrlResponse, CreatePresignedUrlResponseSchema } from "../schemas/responses/hikes";

const API_URL = `${import.meta.env.VITE_API_URL}/hikes`;

// Creates a Hike.
// Throws an error if the response is not 201 or cannot be parsed.
export async function createHike(reqBody: CreateHikeRequest): Promise<Hike> {
  const response = await fetch(API_URL, {
    body: JSON.stringify(reqBody),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  if (response.status !== 201) {
    throw new Error(`Failed to create hike: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = HikeSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse hike: ${result.error}`);
  }

  return result.data;
}

// Creates a Photo for a Hike.
// Throws an error if the response is not 201 or cannot be parsed.
export async function createPhoto(hikeId: number, reqBody: CreatePhotoRequest): Promise<Photo> {
  const response = await fetch(`${API_URL}/${hikeId}/photos`, {
    body: JSON.stringify(reqBody),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  if (response.status !== 201) {
    throw new Error(`Failed to create photo: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = PhotoSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse photo: ${result.error}`);
  }

  return result.data;
}

// Creates a S3 presigned URL to upload a photo for a Hike.
// Throws an error if the response is not 200 or cannot be parsed.
export async function createPresignedUrl(
  hikeId: number,
  reqBody: CreatePresignedUrlRequest,
): Promise<CreatePresignedUrlResponse> {
  const response = await fetch(`${API_URL}/${hikeId}/photos/upload-url`, {
    body: JSON.stringify(reqBody),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  if (response.status !== 200) {
    throw new Error(`Failed to create presigned URL: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = CreatePresignedUrlResponseSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse presigned URL response: ${result.error}`);
  }

  return result.data;
}

// Retrieves all hikes from the API.
// Throws an error if the response is not 200 or cannot be parsed.
export async function fetchHikes(): Promise<Hike[]> {
  const response = await fetch(API_URL);

  if (response.status !== 200) {
    throw new Error(`Failed to retrieve hikes: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = HikeListSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse hikes: ${result.error}`);
  }

  return result.data;
}
