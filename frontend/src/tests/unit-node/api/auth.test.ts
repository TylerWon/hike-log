import { http, HttpResponse } from "msw";
import { type SetupServer, setupServer } from "msw/node";
import { afterEach, beforeEach, describe, expect, test } from "vitest";

import type { LoginRequest } from "../../../schemas/requests/auth";

import { login } from "../../../api/auth/auth";
import { AuthError } from "../../../api/auth/errors";

const API_URL = `${import.meta.env.VITE_API_URL}/auth`;

describe("auth", () => {
  let server: SetupServer;

  beforeEach(() => {
    server = setupServer();
    server.listen();
  });

  afterEach(() => {
    server.resetHandlers();
    server.close();
  });

  describe("login", () => {
    const reqBody: LoginRequest = {
      password: "admin",
    };

    test("returns when successful", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.json({}, { status: 200 });
      });
      server.use(handler);

      await expect(login(reqBody)).resolves.toBeUndefined();
    });

    test("throws AuthError when response is 401", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.json({ error: "Unauthorized" }, { status: 401 });
      });
      server.use(handler);

      await expect(login(reqBody)).rejects.toThrow(AuthError);
    });

    test("throws error when response is not 200 or 401", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.json({ error: "Internal service error" }, { status: 500 });
      });
      server.use(handler);

      await expect(login(reqBody)).rejects.toThrow(Error);
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      await expect(login(reqBody)).rejects.toThrow();
    });
  });
});
