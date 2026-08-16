import { http, HttpResponse } from "msw";
import { setupServer, type SetupServer } from "msw/node";
import { afterEach, beforeEach, describe, expect, test } from "vitest";

import { fetchHikes } from "../../api/hikes";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../fixtures/hike";
import { HIKES_API_RESPONSE_FIXTURE } from "../fixtures/hikes-api";

const API_URL = `${import.meta.env.VITE_API_URL}/hikes`;

describe("hikes", () => {
  describe("fetchHikes", () => {
    let server: SetupServer;

    beforeEach(() => {
      server = setupServer();
      server.listen();
    });

    afterEach(() => {
      server.resetHandlers();
      server.close();
    });

    test("returns an empty list when there are no hikes", async () => {
      const handler = http.get(API_URL, () => {
        return HttpResponse.json([]);
      });
      server.use(handler);

      const hikes = await fetchHikes();

      expect(hikes.length).toEqual(0);
    });

    test("returns hikes when there are hikes", async () => {
      const handler = http.get(API_URL, () => {
        return HttpResponse.json(HIKES_API_RESPONSE_FIXTURE);
      });
      server.use(handler);

      const hikes = await fetchHikes();

      expect(hikes.length).toEqual(2);
      expect(hikes).toEqual([HIKE_FIXTURE_1, HIKE_FIXTURE_2]);
    });

    test("throws an error when the response is not 200", async () => {
      const handler = http.get(API_URL, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler);

      await expect(fetchHikes()).rejects.toThrow();
    });

    test("throws an error when the response cannot be parsed", async () => {
      const handler = http.get(API_URL, () => {
        return HttpResponse.json([{ field: "value" }]);
      });
      server.use(handler);

      await expect(fetchHikes()).rejects.toThrow();
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.get(API_URL, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      expect(fetchHikes()).rejects.toThrow();
    });
  });
});
