import type { $ZodErrorTree } from "zod/v4/core";

import { useRef } from "react";

import type { PhotoFormData } from "../../schemas/forms/photo";
import type { PhotoData } from "./types";

import { readFileAsDataUrl } from "../../utils/file";
import Field from "./Field";

interface PhotoErrors {
  errors: string[];
  items?: $ZodErrorTree<PhotoFormData>[];
}

interface PhotoFieldProps {
  errors?: PhotoErrors; // Errors should be in the same order as Photos
  photos: PhotoData[];
  setPhotos: React.Dispatch<React.SetStateAction<PhotoData[]>>;
}

export default function PhotoField({ errors, photos, setPhotos }: PhotoFieldProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const openFileSelectDialog = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? []);
    if (files.length < 1) {
      return;
    }

    try {
      const newPhotos = await Promise.all(
        files.map(async (file) => ({
          caption: "",
          file: file,
          previewUrl: await readFileAsDataUrl(file),
        })),
      );

      setPhotos((prev) => [...prev, ...newPhotos]);
    } catch (e) {
      console.error("Failed to handle selected photos", e);
      return;
    } finally {
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const removePhoto = (index: number) => {
    setPhotos((prev) => prev.filter((_, idx) => idx !== index));
  };

  const movePhoto = (index: number, direction: -1 | 1) => {
    setPhotos((prev) => {
      const newPhotos = [...prev];
      const newPos = index + direction;
      if (newPos < 0 || newPos >= newPhotos.length) return prev;
      [newPhotos[index], newPhotos[newPos]] = [newPhotos[newPos], newPhotos[index]];
      return newPhotos;
    });
  };

  const updateCaption = (index: number, caption: string) => {
    setPhotos((prev) => prev.map((photo, idx) => (idx === index ? { ...photo, caption } : photo)));
  };

  return (
    <div className="border-t border-forest-800 pt-4">
      <div className="mb-3">
        <p className="font-mono text-[10px] uppercase tracking-widest text-forest-700">Photos</p>
        <input
          accept="image/*"
          className="hidden"
          multiple
          name="Photos"
          onChange={handleFileChange} // Called after files are selected by user
          ref={fileInputRef}
          type="file"
        />
      </div>
      {photos.length === 0 ? (
        <button
          className="w-full border border-dashed border-forest-800 rounded-lg py-8 flex flex-col items-center gap-2 text-forest-700 hover:border-forest-700 hover:text-forest-600 transition-colors focus:outline-none"
          onClick={openFileSelectDialog}
          type="button"
        >
          <svg
            fill="none"
            height="24"
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="1.5"
            viewBox="0 0 24 24"
            width="24"
          >
            <path d="M4 16.5V19a1 1 0 001 1h14a1 1 0 001-1v-2.5" />
            <path d="M12 15V4M8.5 7.5L12 4l3.5 3.5" />
          </svg>
          <span className="font-mono text-xs">Upload photos</span>
        </button>
      ) : (
        <div className="flex flex-col gap-3">
          {photos.map((photo, idx) => (
            <div className="flex gap-3 items-center p-2 rounded-lg border border-forest-800 bg-forest-950" key={idx}>
              <div className="shrink-0 rounded overflow-hidden w-[96px] h-[72px] bg-forest-800">
                <img
                  alt={photo.caption || `Photo ${idx + 1}`}
                  className="w-full h-full object-cover"
                  src={photo.previewUrl}
                />
              </div>

              <Field error={errors?.items?.[idx]?.properties?.file?.errors?.[0]}>
                <input
                  className="hike-form-field py-1.5 text-xs"
                  name="Caption"
                  onChange={(e) => updateCaption(idx, e.target.value)}
                  placeholder="Caption (optional)"
                  type="text"
                  value={photo.caption}
                />
              </Field>

              <div className="flex flex-col items-center gap-1 shrink-0">
                <button
                  aria-label="Remove photo"
                  className="p-1 text-forest-700 hover:text-coral-500 transition-colors focus:outline-none"
                  onClick={() => removePhoto(idx)}
                  type="button"
                >
                  <svg
                    fill="none"
                    height="12"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeWidth="1.5"
                    viewBox="0 0 12 12"
                    width="12"
                  >
                    <path d="M2 2l8 8M10 2L2 10" />
                  </svg>
                </button>
                <button
                  aria-label="Move up"
                  className="p-1 text-forest-700 hover:text-forest-600 transition-colors disabled:opacity-25 focus:outline-none"
                  disabled={idx === 0}
                  onClick={() => movePhoto(idx, -1)}
                  type="button"
                >
                  <svg
                    fill="none"
                    height="12"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="1.5"
                    viewBox="0 0 12 12"
                    width="12"
                  >
                    <path d="M2 8l4-4 4 4" />
                  </svg>
                </button>
                <button
                  aria-label="Move down"
                  className="p-1 text-forest-700 hover:text-forest-600 transition-colors disabled:opacity-25 focus:outline-none"
                  disabled={idx === photos.length - 1}
                  onClick={() => movePhoto(idx, 1)}
                  type="button"
                >
                  <svg
                    fill="none"
                    height="12"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="1.5"
                    viewBox="0 0 12 12"
                    width="12"
                  >
                    <path d="M2 4l4 4 4-4" />
                  </svg>
                </button>
              </div>
            </div>
          ))}

          <button
            className="font-mono w-full border border-dashed border-forest-800 rounded py-2.5 text-xs text-forest-700 hover:border-forest-700 hover:text-forest-600 transition-colors focus:outline-none flex items-center justify-center gap-1.5"
            onClick={openFileSelectDialog}
            type="button"
          >
            <svg
              fill="none"
              height="11"
              stroke="currentColor"
              strokeLinecap="round"
              strokeWidth="1.5"
              viewBox="0 0 11 11"
              width="11"
            >
              <path d="M5.5 1v9M1 5.5h9" />
            </svg>
            Add more photos
          </button>
        </div>
      )}
    </div>
  );
}
