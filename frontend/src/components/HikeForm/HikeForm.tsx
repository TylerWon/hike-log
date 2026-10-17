import "./hike-form.css";
import "../../assets/styles/dialog.css";

import type { $ZodErrorTree } from "zod/v4/core";

import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import * as z from "zod";

import type { PhotoData } from "./types";

import { type HikeFormData, HikeFormDataSchema } from "../../schemas/forms/hike";
import FormError from "../FormError";
import TextWithSpinner from "../TextWithSpinner";
import Field from "./Field";
import PhotoField from "./PhotoField";

export interface HikeFormProps {
  onClose: () => void;
  onSubmit: (formData: HikeFormData) => Promise<void>;
  onSuccess: () => Promise<void>;
}

export default function HikeForm({ onClose, onSubmit, onSuccess }: HikeFormProps) {
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
    // currentTarget = element that is assigned this handler
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  const validate = (): HikeFormData | null => {
    const hikePhotos = photos.map((photo: PhotoData, idx: number) => ({
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

  const addHikeMutation = useMutation({
    mutationFn: onSubmit, // called when mutate() is invoked for this mutation
    onError: (error) => console.error("Failed to submit HikeForm: ", error),
    onSuccess: onSuccess,
  });

  const submissionPending = addHikeMutation.isPending;
  const submissionError = addHikeMutation.isError;

  return (
    <div
      aria-label="Hike form backdrop"
      className="dialog-backdrop"
      onClick={(e) => handleBackdropClick(e)}
      role="region"
    >
      <div
        aria-label="Hike form"
        aria-modal="true"
        className="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-lg border border-forest-800 bg-forest-900 scrollbar-thin"
        role="dialog"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-forest-800 bg-forest-900 sticky top-0 z-10">
          <h2 className="font-serif text-lg font-semibold text-cream-100">Add hike</h2>
          <button
            aria-label="Close"
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
            <Field error={fieldErrors?.properties?.trailName?.errors?.[0]} id="trail-name" label="Trail Name" required>
              <input
                autoFocus
                className="hike-form-field"
                id="trail-name"
                onChange={(e) => setTrailName(e.target.value)}
                placeholder="e.g. Joffre Lakes"
                type="text"
                value={trailName}
              />
            </Field>

            <div className="grid grid-cols-2 gap-3">
              <Field error={fieldErrors?.properties?.date?.errors?.[0]} id="date" label="Date" required>
                <input
                  className="hike-form-field scheme-dark"
                  id="date"
                  onChange={(e) => setDate(e.target.value)}
                  type="date"
                  value={date}
                />
              </Field>
              <Field
                error={fieldErrors?.properties?.duration?.errors?.[0]}
                id="duration"
                label="Duration (minutes)"
                required
              >
                <input
                  className="hike-form-field"
                  id="duration"
                  min={1}
                  onChange={(e) => setDuration(e.target.value)}
                  placeholder="e.g. 240"
                  step={1}
                  type="number"
                  value={duration}
                />
              </Field>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <Field error={fieldErrors?.properties?.rating?.errors?.[0]} id="rating" label="Rating (0-5)" required>
                <input
                  className="hike-form-field"
                  id="rating"
                  max={5}
                  min={0}
                  onChange={(e) => setRating(e.target.value)}
                  placeholder="e.g. 4.5"
                  step={0.5}
                  type="number"
                  value={rating}
                />
              </Field>
              <Field
                error={fieldErrors?.properties?.difficulty?.errors?.[0]}
                id="difficulty"
                label="Difficulty (0-10)"
                required
              >
                <input
                  className="hike-form-field"
                  id="difficulty"
                  max={10}
                  min={0}
                  onChange={(e) => setDifficulty(e.target.value)}
                  placeholder="e.g. 7.5"
                  step={0.5}
                  type="number"
                  value={difficulty}
                />
              </Field>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <Field
                error={fieldErrors?.properties?.distance?.errors?.[0]}
                id="distance"
                label="Distance (km)"
                required
              >
                <input
                  className="hike-form-field"
                  id="distance"
                  min={0}
                  onChange={(e) => setDistance(e.target.value)}
                  placeholder="e.g. 11.8"
                  step={0.1}
                  type="number"
                  value={distance}
                />
              </Field>
              <Field
                error={fieldErrors?.properties?.elevationGain?.errors?.[0]}
                id="elevation-gain"
                label="Elevation gain (m)"
                required
              >
                <input
                  className="hike-form-field"
                  id="elevation-gain"
                  min={0}
                  onChange={(e) => setElevationGain(e.target.value)}
                  placeholder="e.g. 370"
                  step={1}
                  type="number"
                  value={elevationGain}
                />
              </Field>
            </div>

            <Field
              error={fieldErrors?.properties?.allTrailsUrl?.errors?.[0]}
              id="all-trails-url"
              label="AllTrails URL"
              required
            >
              <input
                className="hike-form-field"
                id="all-trails-url"
                onChange={(e) => setAllTrailsUrl(e.target.value)}
                placeholder="https://www.alltrails.com/trail/…"
                type="url"
                value={allTrailsUrl}
              />
            </Field>

            <Field error={fieldErrors?.properties?.notes?.errors?.[0]} id="notes" label="Notes" required>
              <textarea
                className="hike-form-field"
                id="notes"
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
          <div aria-label="Footer" className="border-t border-forest-800 sticky bottom-0 bg-forest-900" role="region">
            {submissionError && <FormError />}

            <div className="flex items-center justify-end gap-3 px-5 py-4 ">
              <button
                aria-label="Cancel"
                className="secondary-button px-4 py-2"
                disabled={submissionPending}
                onClick={handleCancel}
                type="button"
              >
                Cancel
              </button>
              <button
                aria-label="Submit"
                className="primary-button px-5 py-2"
                disabled={submissionPending}
                type="submit"
              >
                {submissionPending ? <TextWithSpinner text="Adding..." /> : "Submit"}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
