import type { $ZodErrorTree } from "zod/v4/core";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import "./hike-form.css";
import * as z from "zod";

import type { PhotoData } from "./types";

import { createHike, createPhotos, createPresignedUrls } from "../../api/hikes";
import { uploadFile } from "../../api/s3";
import { type HikeFormData, HikeFormDataSchema } from "../../schemas/forms/hike";
import { HIKES_QUERY_KEY } from "../HikeLog";
import Field from "./Field";
import PhotoField from "./PhotoField";

interface HikeFormProps {
  onClose: () => void;
  setToastMessage: (message: string) => void;
}

export default function HikeForm({ onClose, setToastMessage }: HikeFormProps) {
  const queryClient = useQueryClient();

  const [trailName, setTrailName] = useState<string>("");
  const [date, setDate] = useState<string>("");
  const [rating, setRating] = useState<string>("");
  const [difficulty, setDifficulty] = useState<string>("");
  const [distance, setDistance] = useState<string>("");
  const [elevationGain, setElevationGain] = useState<string>("");
  const [duration, setDuration] = useState<string>("");
  const [allTrailsUrl, setAllTrailsUrl] = useState<string>("");
  const [notes, setNotes] = useState<string>("");
  const [photos, setPhotos] = useState<PhotoData[]>([]);
  const [fieldErrors, setFieldErrors] = useState<$ZodErrorTree<HikeFormData>>();

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    // target = element that was actually clicked
    // currentTarget = element with the handler (i.e. backdrop)
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  const validate = (): HikeFormData | null => {
    const hikePhotos = photos.map((photo, idx) => ({
      caption: photo.caption.trim(),
      displayOrder: idx + 1,
      file: photo.file,
    }));

    const hike = {
      allTrailsUrl: allTrailsUrl.trim(),
      date,
      difficulty: difficulty,
      distance: distance,
      duration: duration,
      elevationGain: elevationGain,
      notes: notes.trim(),
      photos: hikePhotos,
      rating: rating,
      trailName: trailName.trim(),
    };

    const result = HikeFormDataSchema.safeParse(hike);
    if (!result.success) {
      setFieldErrors(z.treeifyError(result.error));
      return null;
    }

    return result.data;
  };

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault();

    const formData = validate();
    if (!formData) {
      return;
    }

    addHikeMutation.mutate(formData);
  };

  const handleCancel = () => {
    // Clears any errors for the mutation
    addHikeMutation.reset();
    onClose();
  };

  const handlePhotoError = (message: string) => {
    setToastMessage(message);
    onClose();
  };

  const refetchHikes = async () => {
    // Invalidates the "hikes" query so all Hikes get refetched
    await queryClient.invalidateQueries({ queryKey: [HIKES_QUERY_KEY] });
  };

  const addHike = async (formData: HikeFormData) => {
    const { photos, ...hikeData } = formData;

    let hike;
    try {
      hike = await createHike(hikeData);
    } catch (e) {
      // Hike must be created so fail-close
      throw new Error(`Failed to create hike: ${e}`, { cause: e });
    }

    if (photos.length < 1) {
      onClose();
      return;
    }

    const photoErrors = [];
    let createPresignedUrlsResponse;
    try {
      const reqBody = [];
      for (let i = 0; i < photos.length; i++) {
        const photo = photos[i];
        const reqItem = {
          contentLength: photo.file.size,
          contentType: photo.file.type,
          index: i,
        };
        reqBody.push(reqItem);
      }
      createPresignedUrlsResponse = await createPresignedUrls(hike.id, reqBody);
    } catch (e) {
      console.warn("Failed to create presigned URLs: ", e);
      handlePhotoError("Photos could not be uploaded. Please try again.");
      return;
    }

    let s3UploadResponses;
    try {
      const promises = [];
      for (let i = 0; i < photos.length; i++) {
        const createPresignedUrlResponse = createPresignedUrlsResponse[i];
        if (createPresignedUrlResponse.success) {
          promises.push(uploadFile(createPresignedUrlResponse.result.presignedUrl, photos[i].file));
        } else {
          photoErrors.push(`Failed to create presigned URL for photo ${i}`);
          promises.push(Promise.reject());
        }
      }
      s3UploadResponses = await Promise.allSettled(promises);
    } catch (e) {
      console.warn("Failed to upload photos to S3: ", e);
      handlePhotoError("Photos could not be uploaded. Please try again.");
      return;
    }

    let createPhotosResponse;
    try {
      const reqBody = [];
      for (let i = 0; i < photos.length; i++) {
        const createPresignedUrlResponse = createPresignedUrlsResponse[i];
        if (!createPresignedUrlResponse.success) {
          continue;
        }

        const s3UploadResponse = s3UploadResponses[i];
        if (s3UploadResponse.status == "fulfilled") {
          const reqItem = {
            caption: photos[i].caption,
            displayOrder: photos[i].displayOrder,
            index: i,
            objectKey: createPresignedUrlResponse.result.objectKey,
          };
          reqBody.push(reqItem);
        } else {
          photoErrors.push(`Failed to upload photo ${i} to S3`);
        }
      }
      createPhotosResponse = await createPhotos(hike.id, reqBody);
    } catch (e) {
      console.warn("Failed to create photos: ", e);
      handlePhotoError("Photos could not be uploaded. Please try again.");
      return;
    }

    for (const createPhotoResponse of createPhotosResponse) {
      if (!createPhotoResponse.success) {
        photoErrors.push(`Failed to create Photo model for photo ${createPhotoResponse.index}`);
      }
    }

    if (photoErrors.length > 0) {
      console.warn(photoErrors);
      handlePhotoError("Some photos could not be uploaded. Please try again.");
      return;
    }

    onClose();
  };

  const addHikeMutation = useMutation({
    mutationFn: addHike, // called when mutate() is invoked for this mutation
    onError: (error) => console.error("addHikeMutation failed: ", error),
    onSuccess: refetchHikes,
  });

  const submissionPending = addHikeMutation.isPending;
  const submissionError = addHikeMutation.isError;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-forest-950"
      onClick={(e) => handleBackdropClick(e)}
    >
      <div
        aria-label="Add hike"
        aria-modal="true"
        className="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-lg border border-forest-800 bg-forest-900 scrollbar-thin"
        role="dialog"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-forest-800 bg-forest-900 sticky top-0 z-10">
          <h2 className="font-serif text-lg font-semibold text-cream-100">Add hike</h2>
          <button
            aria-label="Cancel"
            className="text-forest-700 hover:text-forest-600 transition-colors p-1 focus:outline-none"
            onClick={handleCancel}
          >
            <svg
              fill="none"
              height="18"
              stroke="currentColor"
              strokeLinecap="round"
              strokeWidth="1.5"
              viewBox="0 0 18 18"
              width="18"
            >
              <path d="M4 4l10 10M14 4L4 14" />
            </svg>
          </button>
        </div>

        {/* Hike fields */}
        <form noValidate onSubmit={handleSubmit}>
          <div className="flex flex-col gap-4 px-5 py-5">
            <Field error={fieldErrors?.properties?.trailName?.errors?.[0]} label="Trail Name" required>
              <input
                autoFocus
                className="hike-form-field"
                name="Trail Name"
                onChange={(e) => setTrailName(e.target.value)}
                placeholder="e.g. Joffre Lakes"
                type="text"
                value={trailName}
              />
            </Field>

            <div className="grid grid-cols-2 gap-3">
              <Field error={fieldErrors?.properties?.date?.errors?.[0]} label="Date" required>
                <input
                  className="hike-form-field scheme-dark"
                  name="Date"
                  onChange={(e) => setDate(e.target.value)}
                  type="date"
                  value={date}
                />
              </Field>
              <Field error={fieldErrors?.properties?.duration?.errors?.[0]} label="Duration (minutes)" required>
                <input
                  className="hike-form-field"
                  min={1}
                  name="Duration"
                  onChange={(e) => setDuration(e.target.value)}
                  placeholder="e.g. 240"
                  step={1}
                  type="number"
                  value={duration}
                />
              </Field>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <Field error={fieldErrors?.properties?.rating?.errors?.[0]} label="Rating (0-5)" required>
                <input
                  className="hike-form-field"
                  max={5}
                  min={0}
                  name="Rating"
                  onChange={(e) => setRating(e.target.value)}
                  placeholder="e.g. 4.5"
                  step={0.5}
                  type="number"
                  value={rating}
                />
              </Field>
              <Field error={fieldErrors?.properties?.difficulty?.errors?.[0]} label="Difficulty (0-10)" required>
                <input
                  className="hike-form-field"
                  max={10}
                  min={0}
                  name="Difficulty"
                  onChange={(e) => setDifficulty(e.target.value)}
                  placeholder="e.g. 7.5"
                  step={0.5}
                  type="number"
                  value={difficulty}
                />
              </Field>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <Field error={fieldErrors?.properties?.distance?.errors?.[0]} label="Distance (km)" required>
                <input
                  className="hike-form-field"
                  min={0}
                  name="Distance"
                  onChange={(e) => setDistance(e.target.value)}
                  placeholder="e.g. 11.8"
                  step={0.1}
                  type="number"
                  value={distance}
                />
              </Field>
              <Field error={fieldErrors?.properties?.elevationGain?.errors?.[0]} label="Elevation gain (m)" required>
                <input
                  className="hike-form-field"
                  min={0}
                  name="Elevation Gain"
                  onChange={(e) => setElevationGain(e.target.value)}
                  placeholder="e.g. 370"
                  step={1}
                  type="number"
                  value={elevationGain}
                />
              </Field>
            </div>

            <Field error={fieldErrors?.properties?.allTrailsUrl?.errors?.[0]} label="AllTrails URL" required>
              <input
                className="hike-form-field"
                name="AllTrails URL"
                onChange={(e) => setAllTrailsUrl(e.target.value)}
                placeholder="https://www.alltrails.com/trail/…"
                type="url"
                value={allTrailsUrl}
              />
            </Field>

            <Field error={fieldErrors?.properties?.notes?.errors?.[0]} label="Notes" required>
              <textarea
                className="hike-form-field"
                name="Notes"
                onChange={(e) => setNotes(e.target.value)}
                placeholder="How was the hike?"
                rows={3}
                value={notes}
              />
            </Field>

            {/* Photo field */}
            <PhotoField errors={fieldErrors?.properties?.photos} photos={photos} setPhotos={setPhotos} />
          </div>

          {/* Footer */}
          <div className="border-t border-forest-800 sticky bottom-0 bg-forest-900">
            {submissionError && (
              <div className="flex items-center gap-2 px-5 py-3 bg-coral-950">
                <svg
                  className="shrink-0"
                  fill="none"
                  height="14"
                  stroke="#c0604a"
                  strokeLinecap="round"
                  strokeWidth="1.5"
                  viewBox="0 0 14 14"
                  width="14"
                >
                  <circle cx="7" cy="7" r="5.5" />
                  <path d="M7 4.5v3M7 9.5v.5" />
                </svg>
                <p className="font-mono flex-1 text-xs text-coral-500">Something went wrong. Please try again.</p>
              </div>
            )}

            <div className="flex items-center justify-end gap-3 px-5 py-4 ">
              <button
                className="font-mono px-4 py-2 text-xs text-forest-600 hover:text-cream-100 transition-colors focus:outline-none"
                disabled={submissionPending}
                onClick={handleCancel}
                type="button"
              >
                Cancel
              </button>

              <button
                className="primary-button px-5 py-2 inline-flex items-center gap-2"
                disabled={submissionPending}
                type="submit"
              >
                {submissionPending ? (
                  <>
                    <svg
                      className="animate-spin"
                      fill="none"
                      height="12"
                      stroke="currentColor"
                      strokeLinecap="round"
                      strokeWidth="1.5"
                      viewBox="0 0 12 12"
                      width="12"
                    >
                      <path d="M6 1v2M6 9v2M1 6h2M9 6h2M2.5 2.5l1.4 1.4M8.1 8.1l1.4 1.4M9.5 2.5L8.1 3.9M3.9 8.1L2.5 9.5" />
                    </svg>
                    Adding...
                  </>
                ) : (
                  "Submit"
                )}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
