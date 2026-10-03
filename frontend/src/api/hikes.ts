import type { CreateHikeRequest, CreatePhotosRequest as CreatePhotosRequest, CreatePresignedUrlsRequest } from "../schemas/requests/hikes";

import { type Hike, HikeListSchema, HikeSchema } from "../schemas/models/hike";
import { type CreatePhotosResponse, CreatePhotosResponseSchema, type CreatePresignedUrlsResponse, CreatePresignedUrlsResponseSchema } from "../schemas/responses/hikes";

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

// Creates Photos for a Hike.
// Throws an error if the response is not 201 or cannot be parsed.
export async function createPhotos(hikeId: number, reqBody: CreatePhotosRequest): Promise<CreatePhotosResponse> {
  const response = await fetch(`${API_URL}/${hikeId}/photos`, {
    body: JSON.stringify(reqBody),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  if (response.status !== 201) {
    throw new Error(`Failed to create photos: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = CreatePhotosResponseSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse photo response: ${result.error}`);
  }

  return result.data;
}

// Creates S3 presigned URLs to upload a photos for a Hike.
// Throws an error if the response is not 200 or cannot be parsed.
export async function createPresignedUrls(
  hikeId: number,
  reqBody: CreatePresignedUrlsRequest,
): Promise<CreatePresignedUrlsResponse> {
  const response = await fetch(`${API_URL}/${hikeId}/photos/presigned-urls`, {
    body: JSON.stringify(reqBody),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  if (response.status !== 200) {
    throw new Error(`Failed to create presigned URLs: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = CreatePresignedUrlsResponseSchema.safeParse(responseJson);
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
