import { useState } from "react";

import "../assets/styles/text.css";
import type { Hike } from "../schemas/models/hike";

import { classnames } from "../utils/classnames";
import { formatDate, formatDistance, formatDuration, formatElevation } from "../utils/formatters";
import ConfirmationDialog from "./ConfirmationDialog";
import StarRating from "./StarRating";
import Thumbnail from "./Thumbnail";

interface HikeOverviewProps {
  hike: Hike;
  index: number;
  isExpanded: boolean;
  onClick: () => void;
  onDelete: () => Promise<void>;
}

export default function HikeOverview({ hike, index, isExpanded, onClick, onDelete }: HikeOverviewProps) {
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState<boolean>(false);

  const handleDeleteButtonClick = () => {
    setShowDeleteConfirmation(true);
  };

  const handleDeleteCancel = () => {
    setShowDeleteConfirmation(false);
  };

  const handleDelete = async () => {
    await onDelete();
    setShowDeleteConfirmation(false);
  };

  return (
    <>
      <div
        aria-label={`${hike.trailName} overview`}
        className="flex items-stretch gap-1 p-4 group/hike-card"
        role="region"
      >
        {/* Content */}
        <button aria-expanded={isExpanded} className="w-full text-left cursor-pointer" onClick={onClick}>
          <div className="flex gap-4">
            <Thumbnail photo={hike.photos.length > 0 ? hike.photos[0] : null} trailName={hike.trailName} />

            <div className="flex-1 min-w-0">
              {/* Trail name, date */}
              <div className="flex items-start justify-between gap-2">
                <div className="min-w-0">
                  <div className="flex items-baseline gap-2">
                    <span className="text-[10px] font-mono text-forest-700 tabular-nums shrink-0">
                      #{String(index).padStart(2, "0")}
                    </span>
                    <span className="block text-base font-semibold font-serif leading-snug text-cream-100 truncate">
                      {hike.trailName}
                    </span>
                  </div>
                  <p className="text-xs font-mono text-forest-600 mt-0.5">{formatDate(hike.date)}</p>
                </div>
              </div>

              {/* Stats */}
              <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-2.5">
                <span className="stat-value inline-flex items-center gap-1.5">
                  <StarRating rating={hike.rating} />
                  {hike.rating}/5
                </span>
                <span className="stat-value inline-flex items-center gap-1">
                  <svg
                    fill="none"
                    height="10"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeWidth="1.5"
                    viewBox="0 0 10 10"
                    width="10"
                  >
                    <circle cx="5" cy="5" r="3.5" />
                    <path d="M5 1.5v1M5 7.5v1M1.5 5h1M7.5 5h1" />
                  </svg>
                  {formatDistance(hike.distance)}
                </span>
                <span className="stat-value inline-flex items-center gap-1">
                  <svg
                    fill="none"
                    height="10"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="1.5"
                    viewBox="0 0 10 10"
                    width="10"
                  >
                    <path d="M2 8L5 2l3 6" />
                  </svg>
                  {formatElevation(hike.elevationGain)}
                </span>
                <span className="stat-value inline-flex items-center gap-1">
                  <svg
                    fill="none"
                    height="10"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeWidth="1.5"
                    viewBox="0 0 10 10"
                    width="10"
                  >
                    <circle cx="5" cy="5" r="3.5" />
                    <path d="M5 3v2l1.5 1.5" />
                  </svg>
                  {formatDuration(hike.duration)}
                </span>
              </div>
            </div>
          </div>
        </button>

        {/* Controls */}
        <div className="flex flex-col items-center justify-start gap-2">
          <button
            aria-label={isExpanded ? "Collapse" : "Expand"}
            className="text-forest-700 focus:outline-none"
            onClick={onClick}
          >
            <svg
              className={classnames("transition-transform duration-200", isExpanded ? "rotate-180" : "rotate-0")}
              fill="none"
              height="16"
              stroke="currentColor"
              strokeLinecap="round"
              strokeWidth="1.5"
              viewBox="0 0 16 16"
              width="16"
            >
              <path d="M4 6l4 4 4-4" />
            </svg>
          </button>
          <button
            aria-label="Delete hike"
            className="text-forest-700 opacity-0 group-hover/hike-card:opacity-100 hover:text-coral-500 transition-all rounded focus:outline-none focus:opacity-100"
            onClick={handleDeleteButtonClick}
          >
            <svg
              fill="none"
              height="13"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="1.4"
              viewBox="0 0 13 13"
              width="13"
            >
              <path d="M2 3.5h9M5 3.5V2.5h3v1M5.5 6v3.5M7.5 6v3.5M3 3.5l.5 7h6l.5-7" />
            </svg>
          </button>
        </div>
      </div>

      {showDeleteConfirmation && (
        <ConfirmationDialog
          acceptButtonPendingText="Deleting..."
          acceptButtonText="Delete"
          actionType="destructive"
          description={
            <span>
              <span className="text-coral-500">{hike.trailName}</span> will permanently be deleted from the log.
            </span>
          }
          onAccept={handleDelete}
          onCancel={handleDeleteCancel}
          title="Delete hike"
        />
      )}
    </>
  );
}
