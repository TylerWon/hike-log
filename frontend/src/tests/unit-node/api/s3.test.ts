import { http, HttpResponse } from "msw";
import { setupServer, type SetupServer } from "msw/node";
import { afterEach, beforeEach, describe, expect, test } from "vitest";

import type { UploadPhotosRequest } from "../../../schemas/requests/s3";

import { uploadFile, uploadPhotos } from "../../../api/s3";

describe("s3", () => {
  const presignedUrl1 =
    "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/dbcf4c30-b23b-4fa4-b350-2620a50f0736?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20261005%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20261005T042400Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost&x-id=PutObject&X-Amz-Signature=<64-char-hex>";
  const presignedUrl2 =
    "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/541f2d00-dd4c-438e-8b1e-4ce5be77d7ed?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20261005%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20261005T042400Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost&x-id=PutObject&X-Amz-Signature=<64-char-hex>";
  // excludes the query parameters
  const presignedUrlPath1 =
    "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/dbcf4c30-b23b-4fa4-b350-2620a50f0736";
  const presignedUrlPath2 =
    "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/541f2d00-dd4c-438e-8b1e-4ce5be77d7ed";
  const photo1 = new File(["photo"], "photo1.png", { type: "image/png" });
  const photo2 = new File(["photo"], "photo2.png", { type: "image/png" });

  let server: SetupServer;

  beforeEach(() => {
    server = setupServer();
    server.listen();
  });

  afterEach(() => {
    server.resetHandlers();
    server.close();
  });

  describe("uploadFile", () => {
    test("returns void when upload is successful", async () => {
      let interceptedRequest: Request | undefined;
      const handler = http.put(presignedUrlPath1, ({ request }) => {
        interceptedRequest = request;
        return HttpResponse.json({}, { status: 200 });
      });
      server.use(handler);

      await expect(uploadFile(presignedUrl1, photo1)).resolves.toBeUndefined();
      expect(interceptedRequest?.headers.get("Content-Type")).toBe(photo1.type);
      expect(interceptedRequest?.headers.get("Content-Length")).toBe(String(photo1.size));
    });

    test("throws an error when the response is not 200", async () => {
      const handler = http.put(presignedUrlPath1, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler);

      await expect(uploadFile(presignedUrl1, photo1)).rejects.toThrow();
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.put(presignedUrlPath1, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      await expect(uploadFile(presignedUrl1, photo1)).rejects.toThrow();
    });
  });

  describe("uploadPhotos", () => {
    test("returns upload status for each photo", async () => {
      const handler1 = http.put(presignedUrlPath1, () => {
        return HttpResponse.json({ status: 200 });
      });
      const handler2 = http.put(presignedUrlPath2, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler1);
      server.use(handler2);

      const reqBody: UploadPhotosRequest = [
        {
          index: 0,
          photo: photo1,
          presignedUrl: presignedUrl1,
        },
        {
          index: 1,
          photo: photo2,
          presignedUrl: presignedUrl2,
        },
      ];
      const response = await uploadPhotos(reqBody);

      expect(response).toEqual([
        { index: 0, success: true },
        { index: 1, success: false },
      ]);
    });
  });
});
