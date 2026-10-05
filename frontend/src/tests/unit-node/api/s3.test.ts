import { http, HttpResponse } from "msw";
import { setupServer, type SetupServer } from "msw/node";
import { afterEach, beforeEach, describe, expect, test } from "vitest";

import { uploadFile } from "../../../api/s3";

describe("s3", () => {
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
    const presignedUrl = "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/dbcf4c30-b23b-4fa4-b350-2620a50f0736?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20261005%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20261005T042400Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost&x-id=PutObject&X-Amz-Signature=<64-char-hex>";
    const presignedUrlPath = "https://hike-log.s3.us-east-1.amazonaws.com/hikes/1/photos/dbcf4c30-b23b-4fa4-b350-2620a50f0736" // excludes the query parameters
    const file = new File(["photo"], "photo.png", { type: "image/png" });

    test("returns void when upload is successful", async () => {
      let interceptedRequest: Request | undefined;
      const handler = http.put(presignedUrlPath, ({ request }) => {
        interceptedRequest = request;
        return HttpResponse.json({}, { status: 200 });
      });
      server.use(handler);

      await expect(uploadFile(presignedUrl, file)).resolves.toBeUndefined();
      expect(interceptedRequest?.headers.get("Content-Type")).toBe(file.type)
      expect(interceptedRequest?.headers.get("Content-Length")).toBe(String(file.size))
    });

    test("throws an error when the response is not 200", async () => {
      const handler = http.put(presignedUrlPath, () => {
        return HttpResponse.json("Internal Service Error", { status: 500 });
      });
      server.use(handler);

      await expect(uploadFile(presignedUrl, file)).rejects.toThrow();
    });

    test("throws an error when there is a network error", async () => {
      const handler = http.put(presignedUrlPath, () => {
        return HttpResponse.error();
      });
      server.use(handler);

      await expect(uploadFile(presignedUrl, file)).rejects.toThrow();
    });
  });
});
