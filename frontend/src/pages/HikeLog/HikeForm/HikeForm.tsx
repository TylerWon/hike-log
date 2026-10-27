import "../../../assets/styles/dialog.css";

import type { $ZodErrorTree } from "zod/v4/core";

import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import * as z from "zod";

import type { PhotoData } from "./types";

import PhotoField from "../../../components/form/PhotoField";
import TextAreaField from "../../../components/form/TextAreaField";
import TextField from "../../../components/form/TextField";
import GenericError from "../../../components/GenericError";
import TextWithSpinner from "../../../components/TextWithSpinner";
import { type HikeFormData, HikeFormDataSchema } from "../../../schemas/forms/hike";

export interface HikeFormProps {
  onCancel: () => void;
  onSubmit: (formData: HikeFormData) => Promise<void>;
}

export default function HikeForm({ onCancel, onSubmit }: HikeFormProps) {
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

  const addHikeMutation = useMutation({
    mutationFn: onSubmit, // called when mutate() is invoked for this mutation
    onError: (error) => console.error("Failed to submit HikeForm: ", error),
  });

  const submissionPending = addHikeMutation.isPending;
  const submissionError = addHikeMutation.isError;

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
    // Don't allow form to be cancelled while submission is in progress
    if (!submissionPending) {
      onCancel();
    }
  };

  return (
    <div aria-label="Hike form backdrop" className="dialog-backdrop" onClick={handleCancel} role="region">
      <div
        aria-label="Hike form"
        aria-modal="true"
        className="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-lg border border-forest-800 bg-forest-900 scrollbar-thin"
        onClick={(e) => e.stopPropagation()} // Prevents form from closing when it is clicked
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
            <TextField
              error={fieldErrors?.properties?.trailName?.errors?.[0]}
              id="trail-name"
              inputProps={{
                autoFocus: true,
                onChange: (e) => setTrailName(e.target.value),
                placeholder: "e.g. Joffre Lakes",
                type: "text",
                value: trailName,
              }}
              label="Trail Name"
              required
            />

            <div className="grid grid-cols-2 gap-3">
              <TextField
                error={fieldErrors?.properties?.date?.errors?.[0]}
                id="date"
                inputProps={{
                  className: "scheme-dark",
                  onChange: (e) => setDate(e.target.value),
                  type: "date",
                  value: date,
                }}
                label="Date"
                required
              />
              <TextField
                error={fieldErrors?.properties?.duration?.errors?.[0]}
                id="duration"
                inputProps={{
                  min: 1,
                  onChange: (e) => setDuration(e.target.value),
                  placeholder: "e.g. 240",
                  step: 1,
                  type: "number",
                  value: duration,
                }}
                label="Duration (minutes)"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <TextField
                error={fieldErrors?.properties?.rating?.errors?.[0]}
                id="rating"
                inputProps={{
                  max: 5,
                  min: 0,
                  onChange: (e) => setRating(e.target.value),
                  placeholder: "e.g. 4.5",
                  step: 0.5,
                  type: "number",
                  value: rating,
                }}
                label="Rating (0-5)"
                required
              />
              <TextField
                error={fieldErrors?.properties?.difficulty?.errors?.[0]}
                id="difficulty"
                inputProps={{
                  max: 10,
                  min: 0,
                  onChange: (e) => setDifficulty(e.target.value),
                  placeholder: "e.g. 7.5",
                  step: 0.5,
                  type: "number",
                  value: difficulty,
                }}
                label="Difficulty (0-10)"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <TextField
                error={fieldErrors?.properties?.distance?.errors?.[0]}
                id="distance"
                inputProps={{
                  min: 0,
                  onChange: (e) => setDistance(e.target.value),
                  placeholder: "e.g. 11.8",
                  step: 0.1,
                  type: "number",
                  value: distance,
                }}
                label="Distance (km)"
                required
              />
              <TextField
                error={fieldErrors?.properties?.elevationGain?.errors?.[0]}
                id="elevation-gain"
                inputProps={{
                  min: 0,
                  onChange: (e) => setElevationGain(e.target.value),
                  placeholder: "e.g. 370",
                  step: 1,
                  type: "number",
                  value: elevationGain,
                }}
                label="Elevation gain (m)"
                required
              />
            </div>

            <TextField
              error={fieldErrors?.properties?.allTrailsUrl?.errors?.[0]}
              id="all-trails-url"
              inputProps={{
                onChange: (e) => setAllTrailsUrl(e.target.value),
                placeholder: "https://www.alltrails.com/trail/…",
                type: "url",
                value: allTrailsUrl,
              }}
              label="AllTrails URL"
              required
            />

            <TextAreaField
              error={fieldErrors?.properties?.notes?.errors?.[0]}
              id="notes"
              label="Notes"
              required
              textAreaProps={{
                onChange: (e) => setNotes(e.target.value),
                placeholder: "How was the hike?",
                rows: 3,
                value: notes,
              }}
            />

            {/* Photo field */}
            <PhotoField errors={fieldErrors?.properties?.photos} photos={photos} setPhotos={setPhotos} />
          </div>

          {/* Footer */}
          <div aria-label="Footer" className="border-t border-forest-800 sticky bottom-0 bg-forest-900" role="region">
            {submissionError && <GenericError />}

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
