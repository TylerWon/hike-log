import { describe, expect, test, vi } from "vitest";

import type { HikeFormData } from "../../../../schemas/forms/hike";
import type { PhotoFormData } from "../../../../schemas/forms/photo";

import { createHike, createPhotos, createPresignedUrls } from "../../../../network/api/hikes";
import { uploadPhotos } from "../../../../network/s3/s3";
import { HikeCreationError, PhotoCreationError } from "../../../../pages/HikeLog/HikeLog/errors";
import {
  createPhotosHelper,
  createPresignedUrlsHelper,
  submitHikeForm,
  uploadAndCreatePhotos,
  uploadPhotosHelper,
} from "../../../../pages/HikeLog/HikeLog/submit-hike-form";
import {
  CREATE_HIKE_RESPONSE_MOCK,
  CREATE_PHOTOS_RESPONSE_MOCK_1,
  CREATE_PHOTOS_RESPONSE_MOCK_2,
  CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1,
  CREATE_PRESIGNED_URLS_RESPONSE_MOCK_2,
} from "../../../mocks/api-responses/hikes";
import { UPLOAD_PHOTOS_RESPONSE_MOCK_1, UPLOAD_PHOTOS_RESPONSE_MOCK_2 } from "../../../mocks/api-responses/s3";

vi.mock(import("../../../../network/api/hikes"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    createHike: vi.fn(),
    createPhotos: vi.fn(),
    createPresignedUrls: vi.fn(),
  };
});

vi.mock(import("../../../../network/s3/s3"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    uploadPhotos: vi.fn(),
  };
});

describe("submit-hike-form", () => {
  const photoFormData: PhotoFormData[] = [
    {
      caption: "Caption 1",
      displayOrder: 1,
      file: new File(["photo"], "photo1.png", { type: "image/png" }),
    },
    {
      caption: "Caption 2",
      displayOrder: 2,
      file: new File(["photo"], "photo2.png", { type: "image/jpeg" }),
    },
    {
      caption: "Caption 3",
      displayOrder: 3,
      file: new File(["photo"], "photo3.png", { type: "image/heic" }),
    },
  ];
  const formData: HikeFormData = {
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/joffre-lakes",
    date: "2026-05-17",
    difficulty: 5.5,
    distance: 11.8,
    duration: 240,
    elevationGain: 370,
    notes:
      "Arrived at the trailhead at 7am and got one of the last parking spots. All three lakes are stunning but Upper Joffre is the showstopper — glacier calving ice directly into the teal water. Light snow made the upper section slippery.",
    photos: photoFormData,
    rating: 4.5,
    trailName: "Joffre Lakes",
  };

  describe("submitHikeForm", () => {
    test("throws HikeCreationError when hike fails to be created", async () => {
      vi.mocked(createHike).mockRejectedValue(new Error("Failed to create hike"));
      await expect(submitHikeForm(formData)).rejects.toThrow(HikeCreationError);
    });

    test("throws PhotoCreationError when photos fail to be created", async () => {
      vi.mocked(createHike).mockResolvedValue(CREATE_HIKE_RESPONSE_MOCK);
      vi.mocked(createPresignedUrls).mockRejectedValue(new Error("Failed to create presigned URLs"));

      await expect(submitHikeForm(formData)).rejects.toThrow(PhotoCreationError);
    });

    test("returns successfully when hike and photo are created", async () => {
      vi.mocked(createHike).mockResolvedValue(CREATE_HIKE_RESPONSE_MOCK);
      vi.mocked(createPresignedUrls).mockResolvedValue(CREATE_PRESIGNED_URLS_RESPONSE_MOCK_2);
      vi.mocked(uploadPhotos).mockResolvedValue(UPLOAD_PHOTOS_RESPONSE_MOCK_2);
      vi.mocked(createPhotos).mockResolvedValue(CREATE_PHOTOS_RESPONSE_MOCK_2);

      await expect(submitHikeForm(formData)).resolves.toBeUndefined();
    });
  });

  describe("uploadAndCreatePhotos", () => {
    test("returns errors when presigned URLs fail to be created, photos fail to be uploaded to S3, and photos fail to be created", async () => {
      vi.mocked(createPresignedUrls).mockResolvedValue(CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1); // third photo fails
      vi.mocked(uploadPhotos).mockResolvedValue(UPLOAD_PHOTOS_RESPONSE_MOCK_1); // second photo fails
      vi.mocked(createPhotos).mockResolvedValue(CREATE_PHOTOS_RESPONSE_MOCK_1); // first photo fails

      const errors = await uploadAndCreatePhotos(photoFormData, 1);
      expect(errors).toEqual(
        new Map([
          [0, "Failed to create model for photo 0"],
          [1, "Failed to upload photo 1 to S3"],
          [2, "Failed to create presigned URL for photo 2"],
        ]),
      );
    });

    test("returns error when unexpected exception is thrown", async () => {
      vi.mocked(createPresignedUrls).mockRejectedValue(new Error("Failed to create presigned URLs"));

      const errors = await uploadAndCreatePhotos(photoFormData, 1);
      expect(errors).toEqual(new Map([[0, "Failed to create photos: Error: Failed to create presigned URLs"]]));
    });
  });

  describe("createPresignedUrlsHelper", () => {
    test("creates presigned URLs", async () => {
      vi.mocked(createPresignedUrls).mockResolvedValue(CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1);

      const response = await createPresignedUrlsHelper(photoFormData, 1);
      expect(createPresignedUrls).toHaveBeenCalledWith(1, [
        { contentLength: photoFormData[0].file.size, contentType: photoFormData[0].file.type, index: 0 },
        { contentLength: photoFormData[1].file.size, contentType: photoFormData[1].file.type, index: 1 },
        { contentLength: photoFormData[2].file.size, contentType: photoFormData[2].file.type, index: 2 },
      ]);
      expect(response).toEqual(CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1);
    });

    test("propogates error when presigned URLs cannot be created", async () => {
      vi.mocked(createPresignedUrls).mockRejectedValue(new Error("Failed to create presigned URLs"));
      await expect(createPresignedUrlsHelper(photoFormData, 1)).rejects.toThrow();
    });
  });

  describe("uploadPhotosHelper", () => {
    test("uploads photos to S3, skipping ones without presigned URLs", async () => {
      vi.mocked(uploadPhotos).mockResolvedValue(UPLOAD_PHOTOS_RESPONSE_MOCK_1);

      const response = await uploadPhotosHelper(photoFormData, CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1);

      const [presignedUrlResponse1, presignedUrlResponse2] = CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1;
      if (!presignedUrlResponse1.success || !presignedUrlResponse2.success) {
        throw new Error("Expected successful presigned URLs");
      }
      // Skips third photo because presigned URL failed to be created
      expect(uploadPhotos).toHaveBeenCalledWith([
        {
          index: 0,
          photo: photoFormData[0].file,
          presignedUrl: presignedUrlResponse1.result.presignedUrl,
        },
        {
          index: 1,
          photo: photoFormData[1].file,
          presignedUrl: presignedUrlResponse2.result.presignedUrl,
        },
      ]);
      expect(response).toEqual(UPLOAD_PHOTOS_RESPONSE_MOCK_1);
    });
  });

  describe("createPhotosHelper", () => {
    test("creates photos, skipping ones without presigned URLs or that failed to be uploaded to S3", async () => {
      vi.mocked(createPhotos).mockResolvedValue(CREATE_PHOTOS_RESPONSE_MOCK_1);

      const response = await createPhotosHelper(
        photoFormData,
        1,
        CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1,
        UPLOAD_PHOTOS_RESPONSE_MOCK_1,
      );

      const presignedUrlResponse1 = CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1[0];
      if (!presignedUrlResponse1.success) {
        throw new Error("Expected successful presigned URLs");
      }
      // Skips second photo because it failed to be uploaded to S3, third photo because presigned URL failed to be created
      expect(createPhotos).toHaveBeenCalledWith(1, [
        {
          caption: photoFormData[0].caption,
          displayOrder: photoFormData[0].displayOrder,
          index: 0,
          objectKey: presignedUrlResponse1.result.objectKey,
        },
      ]);
      expect(response).toEqual(CREATE_PHOTOS_RESPONSE_MOCK_1);
    });

    test("propogates error when photos cannot be created", async () => {
      vi.mocked(createPhotos).mockRejectedValue(new Error("Failed to create photos"));
      await expect(
        createPhotosHelper(photoFormData, 1, CREATE_PRESIGNED_URLS_RESPONSE_MOCK_1, UPLOAD_PHOTOS_RESPONSE_MOCK_1),
      ).rejects.toThrow();
    });
  });
});
