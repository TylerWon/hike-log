import type { LoginRequest } from "../../schemas/requests/auth";

import { AuthError } from "./errors";

const API_URL = `${import.meta.env.VITE_API_URL}/auth`;

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
  const response = await fetch(API_URL, {
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
