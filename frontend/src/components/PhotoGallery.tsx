import { useState } from "react";

import { type Photo } from "../data/hikes";

interface PhotoGalleryProps {
  photos: Photo[];
}

export default function PhotoGallery({ photos }: PhotoGalleryProps) {
  const [lightboxIndex, setLightboxIndex] = useState<null | number>(null);

  if (photos.length === 0) return null;

  const closeLightbox = () => {
    setLightboxIndex(null);
  };

  const showPrevPhoto = () => {
    setLightboxIndex((currIdx) => (currIdx !== null ? (currIdx - 1 + photos.length) % photos.length : 0));
  };

  const showNextPhoto = () => {
    setLightboxIndex((currIdx) => (currIdx !== null ? (currIdx + 1) % photos.length : 0));
  };

  return (
    <>
      {/* Thumbnail strip */}
      <div className="flex gap-2 overflow-x-auto scrollbar-thin">
        {photos.map((photo, i) => (
          <button
            aria-label={photo.caption ?? `Photo ${i + 1}`}
            className="w-[120px] h-[80px] shrink-0 rounded overflow-hidden bg-forest-800 cursor-pointer group/photo"
            key={i}
            onClick={() => setLightboxIndex(i)}
          >
            <img
              alt={photo.caption ?? `Photo ${i + 1}`}
              className="w-full h-full object-cover transition-all duration-200 group-hover/photo:scale-105 group-hover/photo:brightness-110"
              loading="lazy"
              src={photo.url}
            />
          </button>
        ))}
      </div>

      {/* Lightbox */}
      {lightboxIndex !== null && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-[#050805]/95" onClick={closeLightbox}>
          {/* Close */}
          <button
            aria-label="Close"
            className="absolute top-4 right-4 text-forest-600 hover:text-cream-100 transition-colors p-2 focus:outline-none"
            onClick={closeLightbox}
          >
            <svg
              fill="none"
              height="20"
              stroke="currentColor"
              strokeLinecap="round"
              strokeWidth="1.5"
              viewBox="0 0 20 20"
              width="20"
            >
              <path d="M4 4l12 12M16 4L4 16" />
            </svg>
          </button>

          {/* Image */}
          <div
            className="relative flex flex-col items-center max-w-4xl w-full px-16"
            onClick={(e) => e.stopPropagation()} // Prevents lightbox from closing when image is clicked
          >
            <img
              alt={photos[lightboxIndex].caption ?? `Photo ${lightboxIndex + 1}`}
              className="max-h-[75vh] max-w-full w-auto object-contain rounded"
              src={photos[lightboxIndex].url.replace(/w=\d+/, "w=1200").replace(/h=\d+/, "h=800")} // Request desired image width and height to save bandwidth
            />
            {photos[lightboxIndex].caption && (
              <p className="mt-3 text-center text-sm font-mono text-forest-600 max-w-md">
                {photos[lightboxIndex].caption}
              </p>
            )}
            <p className="mt-1 text-xs font-mono text-forest-700">
              {lightboxIndex + 1} / {photos.length}
            </p>
          </div>

          {/* Prev/Next */}
          {photos.length > 1 && (
            <>
              <button
                aria-label="Previous photo"
                className="absolute left-4 top-1/2 -translate-y-1/2 text-forest-600 hover:text-cream-100 transition-colors p-3 focus:outline-none"
                onClick={(e) => {
                  e.stopPropagation();
                  showPrevPhoto();
                }}
              >
                <svg
                  fill="none"
                  height="20"
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="1.5"
                  viewBox="0 0 20 20"
                  width="20"
                >
                  <path d="M13 4l-6 6 6 6" />
                </svg>
              </button>
              <button
                aria-label="Next photo"
                className="absolute right-4 top-1/2 -translate-y-1/2 text-forest-600 hover:text-cream-100 transition-colors p-3 focus:outline-none"
                onClick={(e) => {
                  e.stopPropagation();
                  showNextPhoto();
                }}
              >
                <svg
                  fill="none"
                  height="20"
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="1.5"
                  viewBox="0 0 20 20"
                  width="20"
                >
                  <path d="M7 4l6 6-6 6" />
                </svg>
              </button>
            </>
          )}
        </div>
      )}
    </>
  );
}
