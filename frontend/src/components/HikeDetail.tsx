import "../assets/styles/text.css";
import { type Hike, type Photo } from "../data/hikes";
import DifficultyBar from "./DifficultyBar";
import PhotoGallery from "./PhotoGallery";

interface HikeDetailProps {
  hike: Hike;
  photos: Photo[];
}

export default function HikeDetail({ hike, photos }: HikeDetailProps) {
  return (
    <div className="border-t border-forest-800 px-4 py-5 flex flex-col gap-y-5">
      {/* Difficulty and AllTrails link */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <span className="field-label">Difficulty</span>
          <DifficultyBar difficulty={hike.difficulty} />
          <span className="stat-value">{hike.difficulty}/10</span>
        </div>

        <a
          href={hike.allTrailsUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1.5 text-xs font-mono text-sage-500 hover:text-sage-400 transition-colors border border-sage-500 rounded px-2.5 py-1 hover:border-sage-400"
        >
          <svg
            width="11"
            height="11"
            viewBox="0 0 11 11"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M4.5 2H2.5A1 1 0 001.5 3v5.5A1 1 0 002.5 9.5H8a1 1 0 001-1V6.5" />
            <path d="M6.5 1.5H9.5v3M9.5 1.5L5 7" />
          </svg>
          View on AllTrails
        </a>
      </div>

      <div>
        <p className="field-label mb-2">Notes</p>
        <p className="text-sm leading-relaxed text-cream-200">{hike.notes}</p>
      </div>

      {photos.length > 0 && (
        <div>
          <p className="field-label mb-2">Photos ({photos.length})</p>
          <PhotoGallery photos={photos} />
        </div>
      )}
    </div>
  );
}
