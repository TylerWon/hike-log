const API_URL = import.meta.env.VITE_API_URL;

/**
 * Wrapper around `fetch` that configures it for use with the application's API.
 *
 * The `path` argument should start with a forward slash.
 */
export function apiFetch(path: string, init: RequestInit = {}) {
  return fetch(`${API_URL}${path}`, {
    // Allows cookies to be included in cross-origin requests/responses.
    // Required for session cookie authentication to work.
    credentials: "include",
    ...init,
  });
}
