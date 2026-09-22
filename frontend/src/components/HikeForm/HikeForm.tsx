import type { $ZodErrorTree } from "zod/v4/core";

import { useState } from "react";
import * as z from "zod";

import "./hike-form.css";
import type { PhotoData } from "./types";

import { type HikeFormData, HikeFormDataSchema } from "../../schemas/forms/hike";
import Field from "./Field";
import PhotoField from "./PhotoField";

interface HikeFormProps {
  onCancel: () => void;
  onSubmit: (hike: HikeFormData) => void;
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
  const [errors, setErrors] = useState<$ZodErrorTree<HikeFormData>>();

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    // target = element that was actually clicked
    // currentTarget = element with the handler (i.e. backdrop)
    if (e.target === e.currentTarget) {
      onCancel();
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
      setErrors(z.treeifyError(result.error));
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

    onSubmit(formData);
  };

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
            onClick={onCancel}
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
            <Field error={errors?.properties?.trailName?.errors?.[0]} label="Trail Name" required>
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
              <Field error={errors?.properties?.date?.errors?.[0]} label="Date" required>
                <input
                  className="hike-form-field scheme-dark"
                  name="Date"
                  onChange={(e) => setDate(e.target.value)}
                  type="date"
                  value={date}
                />
              </Field>
              <Field error={errors?.properties?.duration?.errors?.[0]} label="Duration (minutes)" required>
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
              <Field error={errors?.properties?.rating?.errors?.[0]} label="Rating (0-5)" required>
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
              <Field error={errors?.properties?.difficulty?.errors?.[0]} label="Difficulty (0-10)" required>
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
              <Field error={errors?.properties?.distance?.errors?.[0]} label="Distance (km)" required>
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
              <Field error={errors?.properties?.elevationGain?.errors?.[0]} label="Elevation gain (m)" required>
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

            <Field error={errors?.properties?.allTrailsUrl?.errors?.[0]} label="AllTrails URL" required>
              <input
                className="hike-form-field"
                name="AllTrails URL"
                onChange={(e) => setAllTrailsUrl(e.target.value)}
                placeholder="https://www.alltrails.com/trail/…"
                type="url"
                value={allTrailsUrl}
              />
            </Field>

            <Field error={errors?.properties?.notes?.errors?.[0]} label="Notes" required>
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
            <PhotoField errors={errors?.properties?.photos} photos={photos} setPhotos={setPhotos} />
          </div>

          {/* Footer */}
          <div className="flex items-center justify-end gap-3 px-5 py-4 border-t border-forest-800 sticky bottom-0 bg-forest-900">
            <button
              className="font-mono px-4 py-2 text-xs text-forest-600 hover:text-cream-100 transition-colors focus:outline-none"
              onClick={onCancel}
              type="button"
            >
              Cancel
            </button>
            <button className="primary-button px-5 py-2" type="submit">
              Submit
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
