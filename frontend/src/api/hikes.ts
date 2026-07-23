import { type Hike, HikeListSchema } from "../schemas/hike";

const API_URL = `${import.meta.env.VITE_API_URL}/hikes`;

// Retrieves all hikes from the API. 
// Throws an error if the response is not 200 or cannot be parsed.
export async function fetchHikes(): Promise<Hike[]> {
  const response = await fetch(API_URL);

  if (!response.ok) {
    throw new Error(`Failed to retrieve hikes: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();

  const result = HikeListSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse hikes: ${result.error}`);
  }

  return result.data;
}
