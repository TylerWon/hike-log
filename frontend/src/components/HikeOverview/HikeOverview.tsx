import { type Hike } from "../../data/hikes";
import StarRating from "../StarRating";
import Thumbnail from "../Thumbnail";
import { classnames } from "../../utils/classnames";
import "./hike-overview.css";

interface HikeOverviewProps {
  hike: Hike;
  index: number;
  isExpanded: boolean;
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

export default function HikeOverview({ hike, index, isExpanded }: HikeOverviewProps) {
  return (
    <div className="flex gap-4 p-4">
      <Thumbnail photo={hike.coverPhoto} trail={hike.trail} />

      <div className="flex-1 min-w-0">
        {/* Trail name, date, chevron */}
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <div className="flex items-baseline gap-2">
              <span className="text-[10px] font-mono text-forest-700 tabular-nums shrink-0">
                #{String(index).padStart(2, "0")}
              </span>
              <span className="block text-base font-semibold font-serif leading-snug text-cream-100 truncate">
                {hike.trail}
              </span>
            </div>
            <p className="text-xs font-mono text-forest-600 mt-0.5">{formatDate(hike.date)}</p>
          </div>
          <svg
            className={classnames("shrink-0 mt-1 transition-transform duration-200 stroke-forest-700", {
              "rotate-180": isExpanded,
              "rotate-0": !isExpanded,
            })}
            width="16"
            height="16"
            strokeWidth="1.5"
            strokeLinecap="round"
            viewBox="0 0 16 16"
            fill="none"
          >
            <path d="M4 6l4 4 4-4" />
          </svg>
        </div>

        {/* Stats */}
        <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-2.5">
          <span className="card-stat gap-1.5">
            <StarRating rating={hike.rating} />
            {hike.rating}/5
          </span>

          <span className="card-stat gap-1">
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

          <span className="card-stat gap-1">
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

          <span className="card-stat gap-1">
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
        </div>
      </div>
    </div>
  );
}
