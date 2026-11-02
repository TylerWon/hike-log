import { describe, expect, test, vi } from "vitest";

import { apiFetch } from "../../../../network/api/fetch";

const API_URL = import.meta.env.VITE_API_URL;

describe("fetch", () => {
  describe("apiFetch", () => {
    test("sends request to API with credentials", async () => {
      const fetchSpy = vi.spyOn(globalThis, "fetch");
      fetchSpy.mockResolvedValue(new Response(null, { status: 200 }));

      await apiFetch("/hikes");

      expect(fetchSpy).toHaveBeenCalledWith(`${API_URL}/hikes`, expect.objectContaining({ credentials: "include" }));
    });
  });
});
