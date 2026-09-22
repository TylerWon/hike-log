/**
 * Returns a data URL for the given file.
 *
 * A data URL embeds a file's contents directly in a URL string, so the browser can use it without a network request.
 * We use this to display the a preview of the image in the form.
 */
export function readFileAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}
