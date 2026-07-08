import { useState } from "react";
import { type Photo } from "../data/hikes";

interface PhotoGalleryProps {
  photos: Photo[];
}

export default function PhotoGallery({ photos }: PhotoGalleryProps) {
  const [lightboxIndex, setLightboxIndex] = useState<number | null>(null);

  if (photos.length === 0) return null;

  const closeLightbox = () => setLightboxIndex(null);
  const prev = () =>
    setLightboxIndex((i) => (i !== null ? (i - 1 + photos.length) % photos.length : 0));
  const next = () =>
    setLightboxIndex((i) => (i !== null ? (i + 1) % photos.length : 0));

  return (
    <>
      {/* Thumbnail strip */}
      <div className="flex gap-2 overflow-x-auto pb-1" style={{ scrollbarWidth: "thin" }}>
        {photos.map((photo, i) => (
          <button
            key={i}
            onClick={() => setLightboxIndex(i)}
            className="flex-shrink-0 rounded overflow-hidden bg-[#2e3d2a] cursor-pointer group focus:outline-none focus:ring-2 focus:ring-[#7fa86a]"
            style={{ width: 120, height: 80 }}
            aria-label={photo.caption ?? `Photo ${i + 1}`}
          >
            <img
              src={photo.url}
              alt={photo.caption ?? `Photo ${i + 1}`}
              className="w-full h-full object-cover transition-all duration-200 group-hover:scale-105 group-hover:brightness-110"
              loading="lazy"
            />
          </button>
        ))}
      </div>

      {/* Lightbox */}
      {lightboxIndex !== null && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center"
          style={{ backgroundColor: "rgba(5, 8, 5, 0.95)" }}
          onClick={closeLightbox}
        >
          {/* Close */}
          <button
            className="absolute top-4 right-4 text-[#8a9e82] hover:text-[#e8dfc8] transition-colors p-2 focus:outline-none"
            onClick={closeLightbox}
            aria-label="Close"
          >
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round">
              <path d="M4 4l12 12M16 4L4 16" />
            </svg>
          </button>

          {/* Image */}
          <div
            className="relative flex flex-col items-center max-w-4xl w-full px-16"
            onClick={(e) => e.stopPropagation()}
          >
            <img
              src={photos[lightboxIndex].url.replace(/w=\d+/, "w=1200").replace(/h=\d+/, "h=800")}
              alt={photos[lightboxIndex].caption ?? `Photo ${lightboxIndex + 1}`}
              className="max-h-[75vh] w-auto object-contain rounded"
              style={{ maxWidth: "100%" }}
            />
            {photos[lightboxIndex].caption && (
              <p
                className="mt-3 text-center text-sm text-[#8a9e82] max-w-md"
                style={{ fontFamily: "'DM Mono', monospace" }}
              >
                {photos[lightboxIndex].caption}
              </p>
            )}
            {/* Counter */}
            <p
              className="mt-1 text-xs text-[#4a5c44]"
              style={{ fontFamily: "'DM Mono', monospace" }}
            >
              {lightboxIndex + 1} / {photos.length}
            </p>
          </div>

          {/* Prev/Next */}
          {photos.length > 1 && (
            <>
              <button
                className="absolute left-4 top-1/2 -translate-y-1/2 text-[#8a9e82] hover:text-[#e8dfc8] transition-colors p-3 focus:outline-none"
                onClick={(e) => { e.stopPropagation(); prev(); }}
                aria-label="Previous photo"
              >
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M13 4l-6 6 6 6" />
                </svg>
              </button>
              <button
                className="absolute right-4 top-1/2 -translate-y-1/2 text-[#8a9e82] hover:text-[#e8dfc8] transition-colors p-3 focus:outline-none"
                onClick={(e) => { e.stopPropagation(); next(); }}
                aria-label="Next photo"
              >
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
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
