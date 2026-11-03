import type { LoginRequest } from "../../../schemas/requests/auth";

import { type CheckAuthResponse, CheckAuthResponseSchema } from "../../../schemas/responses/auth";
import { apiFetch } from "../fetch";
import { AuthError } from "./errors";

/**
 * Checks if the current user is authenticated.
 *
 * Throws an error if the response is not 200 or cannot be parsed.
 */
export async function checkAuth(): Promise<CheckAuthResponse> {
  const response = await apiFetch("/auth");

  if (response.status !== 200) {
    throw new Error(`Failed to perform auth check: ${response.status} - ${response.statusText}`);
  }

  const responseJson = await response.json();
  const result = CheckAuthResponseSchema.safeParse(responseJson);
  if (!result.success) {
    throw new Error(`Failed to parse auth response: ${result.error}`);
  }

  return result.data;
}

/**
 * Logs a user in to the admin account.
 *
 * If login is successful, a session cookie will be attached to the response which should be used to authenticate
 * requests to restricted API endpoints.
 *
 * Throws an error if the response is not 200. An AuthError is thrown for a 401 response, otherwise a normal Error is
 * thrown.
 */
export async function login(reqBody: LoginRequest) {
  const response = await apiFetch("/auth", {
    body: JSON.stringify(reqBody),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  if (response.status === 401) {
    throw new AuthError("Incorrect credentials");
  } else if (response.status !== 200) {
    throw new Error(`Failed to login: ${response.status} - ${response.statusText}`);
  }
}
