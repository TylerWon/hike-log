import { type Hike } from "../data/hikes";
import StarRating from "./StarRating";
import HikeDetail from "./HikeDetail";

interface HikeCardProps {
  hike: Hike;
  index: number;
  isExpanded: boolean;
  onToggle: () => void;
}

function formatDate(iso: string): string {
  const d = new Date(iso + "T12:00:00");
  return d.toLocaleDateString("en-CA", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function formatDuration(minutes: number): string {
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
}

export default function HikeCard({ hike, index, isExpanded, onToggle }: HikeCardProps) {
  const allPhotos = [
    ...(hike.coverPhoto ? [hike.coverPhoto] : []),
    ...(hike.photos ?? []),
  ];

  return (
    <article
      className="rounded-lg overflow-hidden border transition-all duration-200"
      style={{
        borderColor: isExpanded ? "var(--color-forest-700)" : "var(--color-forest-800)",
        backgroundColor: isExpanded ? "var(--color-forest-850)" : "var(--color-forest-900)",
      }}
    >
      {/* Collapsed header — always visible */}
      <button
        onClick={onToggle}
        className="w-full text-left cursor-pointer group"
        aria-expanded={isExpanded}
      >
        <div className="flex gap-4 p-4">
          {/* Cover photo thumbnail */}
          {hike.coverPhoto ? (
            <div
              className="flex-shrink-0 rounded overflow-hidden bg-forest-800"
              style={{ width: 96, height: 72 }}
            >
              <img
                src={hike.coverPhoto.url}
                alt={hike.coverPhoto.caption ?? hike.trail}
                className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
                loading="lazy"
              />
            </div>
          ) : (
            <div
              className="flex-shrink-0 rounded bg-[#1f2d1a] flex items-center justify-center"
              style={{ width: 96, height: 72 }}
            >
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="var(--color-forest-700)"
                strokeWidth="1.5"
              >
                <path d="M3 17L9 9l4 5 3-3 5 6H3z" />
              </svg>
            </div>
          )}

          {/* Content */}
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0">
                <div className="flex items-baseline gap-2">
                  <span
                    className="text-[10px] font-mono text-forest-700 tabular-nums flex-shrink-0"
                  >
                    #{String(index).padStart(2, "0")}
                  </span>
                  <h2
                    className="text-base font-semibold font-serif leading-snug text-cream-100 truncate"
                  >
                    {hike.trail}
                  </h2>
                </div>
                <p
                  className="text-xs font-mono text-forest-600 mt-0.5"
                >
                  {formatDate(hike.date)}
                </p>
              </div>
              {/* Expand chevron */}
              <svg
                width="16"
                height="16"
                viewBox="0 0 16 16"
                fill="none"
                stroke="var(--color-forest-700)"
                strokeWidth="1.5"
                strokeLinecap="round"
                className="flex-shrink-0 mt-1 transition-transform duration-200"
                style={{ transform: isExpanded ? "rotate(180deg)" : "rotate(0deg)" }}
              >
                <path d="M4 6l4 4 4-4" />
              </svg>
            </div>

            {/* Stats row */}
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-2.5">
              <span className="inline-flex items-center gap-1.5">
                <StarRating rating={hike.rating} />
                <span
                  className="text-[11px] font-mono text-forest-600"
                >
                  {hike.rating}/5
                </span>
              </span>

              <span
                className="text-[11px] font-mono text-forest-600 inline-flex items-center gap-1"
              >
                <svg
                  width="10"
                  height="10"
                  viewBox="0 0 10 10"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                >
                  <circle cx="5" cy="5" r="3.5" />
                  <path d="M5 1.5v1M5 7.5v1M1.5 5h1M7.5 5h1" />
                </svg>
                {hike.distanceKm} km
              </span>

              <span
                className="text-[11px] font-mono text-forest-600 inline-flex items-center gap-1"
              >
                <svg
                  width="10"
                  height="10"
                  viewBox="0 0 10 10"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M2 8L5 2l3 6" />
                </svg>
                {hike.elevationGainM} m
              </span>

              {hike.durationMinutes !== undefined && (
                <span
                  className="text-[11px] font-mono text-forest-600 inline-flex items-center gap-1"
                >
                  <svg
                    width="10"
                    height="10"
                    viewBox="0 0 10 10"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                  >
                    <circle cx="5" cy="5" r="3.5" />
                    <path d="M5 3v2l1.5 1.5" />
                  </svg>
                  {formatDuration(hike.durationMinutes)}
                </span>
              )}
            </div>
          </div>
        </div>
      </button>

      {/* Expandable detail panel */}
      <div
        style={{
          display: "grid",
          gridTemplateRows: isExpanded ? "1fr" : "0fr",
          transition: "grid-template-rows 0.3s ease",
        }}
      >
        <div style={{ overflow: "hidden" }}>
          <HikeDetail hike={hike} photos={allPhotos} />
        </div>
      </div>
    </article>
  );
}
