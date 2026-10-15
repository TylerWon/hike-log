/**
 * Vite imports static assets as URLs. Fetch the bytes so they can be uploaded as real image files.
 */
export async function createFileFromUrl(assetUrl: string, filename: string) {
  const response = await fetch(new URL(assetUrl, import.meta.url));
  const blob = await response.blob();
  return new File([blob], filename, { type: blob.type });
}
