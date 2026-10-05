/**
 * Returns a data URL for the given file.
 *
 * A data URL embeds a file's contents directly in a URL string, so the browser can use it without a network request.
 * This can used to display a preview of an image on the page.
 */
export function readFileAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    // Fires after a successful read
    reader.onload = () => resolve(reader.result as string);

    // Fires on read error
    reader.onerror = reject;

    // Initiates file read
    reader.readAsDataURL(file);
  });
}
