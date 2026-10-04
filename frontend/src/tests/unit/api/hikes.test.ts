import { http, HttpResponse } from "msw";
import { setupServer, type SetupServer } from "msw/node";
import { afterEach, beforeEach, describe, expect, test } from "vitest";

import { createHike, createPhotos, createPresignedUrls, fetchHikes } from "../../../api/hikes";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "../../data/fixtures/hike";
import { CREATE_HIKE_API_RESPONSE, CREATE_PHOTOS_API_RESPONSE, CREATE_PRESIGNED_URLS_API_RESPONSE, LIST_HIKES_API_RESPONSE } from "../../data/api-responses/hikes";
import type { CreateHikeRequest, CreatePhotosRequest, CreatePresignedUrlsRequest } from "../../../schemas/requests/hikes";

const API_URL = `${import.meta.env.VITE_API_URL}/hikes`;

describe("hikes", () => {
  let server: SetupServer;

  beforeEach(() => {
    server = setupServer();
    server.listen();
  });

  afterEach(() => {
    server.resetHandlers();
    server.close();
  });

  describe("createHike", () => {
    const reqBody: CreateHikeRequest = {
      difficulty: 7.5,
      distance: 25,
      duration: 360,
      elevationGain: 1500,
      notes: "Difficult hike",
      rating: 3.5,
      trailName: "Trail #1",
      allTrailsUrl: "https://alltrails.com",
      date: "2025-01-01"
    };

    test("returns the newly created hike when successful", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.json(CREATE_HIKE_API_RESPONSE, { status: 201 });
      });
      server.use(handler);

      const hike = await createHike(reqBody);
      expect(hike).toEqual(CREATE_HIKE_API_RESPONSE);
    });

    test("throws an error when the response is not 201", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler);

      await expect(createHike(reqBody)).rejects.toThrow();
    });

    test("throws an error when the response cannot be parsed", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.json({ field: "value" }, { status: 201 });
      });
      server.use(handler);

      await expect(createHike(reqBody)).rejects.toThrow();
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.post(API_URL, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      await expect(createHike(reqBody)).rejects.toThrow();
    });
  });

  describe("createPhotos", () => {
    const hikeId = 1;
    const reqBody: CreatePhotosRequest = [
      {
        caption: "",
        displayOrder: 1,
        index: 0,
        objectKey: "1" // invalid objectKey but OK for testing
      },
      {
        caption: "Viewpoint photo",
        displayOrder: 2,
        index: 1,
        objectKey: "2" // invalid objectKey but OK for testing
      },
    ]

    test("returns newly created photos when successful", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos`, () => {
        return HttpResponse.json(CREATE_PHOTOS_API_RESPONSE, { status: 201 });
      });
      server.use(handler);

      const photos = await createPhotos(hikeId, reqBody);
      expect(photos).toEqual(CREATE_PHOTOS_API_RESPONSE);
    });

    test("throws an error when the response is not 201", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos`, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler);

      await expect(createPhotos(hikeId, reqBody)).rejects.toThrow();
    });

    test("throws an error when the response cannot be parsed", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos`, () => {
        return HttpResponse.json({ field: "value" }, { status: 201 });
      });
      server.use(handler);

      await expect(createPhotos(hikeId, reqBody)).rejects.toThrow();
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos`, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      await expect(createPhotos(hikeId, reqBody)).rejects.toThrow();
    });
  })

  describe("createPresignedUrls", () => {
    const hikeId = 1;
    const reqBody: CreatePresignedUrlsRequest = [
      {
        contentLength: 100,
        contentType: "image/png",
        index: 0
      },
      {
        contentLength: 200,
        contentType: "image/jpeg",
        index: 1
      },
    ]

    test("returns newly created presigned URLs when successful", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos/presigned-urls`, () => {
        return HttpResponse.json(CREATE_PRESIGNED_URLS_API_RESPONSE, { status: 200 });
      });
      server.use(handler);

      const presignedUrls = await createPresignedUrls(hikeId, reqBody);
      expect(presignedUrls).toEqual(CREATE_PRESIGNED_URLS_API_RESPONSE);
    });

    test("throws an error when the response is not 200", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos/presigned-urls`, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler);

      await expect(createPresignedUrls(hikeId, reqBody)).rejects.toThrow();
    });

    test("throws an error when the response cannot be parsed", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos/presigned-urls`, () => {
        return HttpResponse.json({ field: "value" }, { status: 200 });
      });
      server.use(handler);

      await expect(createPresignedUrls(hikeId, reqBody)).rejects.toThrow();
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.post(`${API_URL}/${hikeId}/photos/presigned-urls`, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      await expect(createPresignedUrls(hikeId, reqBody)).rejects.toThrow();
    });
  })

  describe("fetchHikes", () => {
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
        return HttpResponse.json(LIST_HIKES_API_RESPONSE);
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

      await expect(fetchHikes()).rejects.toThrow();
    });
  });
});
