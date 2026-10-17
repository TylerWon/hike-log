import type { Hike } from "../schemas/models/hike";

import { classnames } from "../utils/classnames";
import HikeDetail from "./HikeDetail";
import HikeOverview from "./HikeOverview";

interface HikeCardProps {
  hike: Hike;
  index: number;
  isExpanded: boolean;
  onClick: () => void;
  onDelete: () => Promise<void>;
}

export default function HikeCard({ hike, index, isExpanded, onClick, onDelete }: HikeCardProps) {
  return (
    <div
      aria-label={`${hike.trailName} card`}
      className={classnames(
        "border rounded-lg overflow-hidden transition-all duration-200",
        isExpanded ? "bg-forest-850 border-forest-700" : "bg-forest-900 border-forest-800",
      )}
      role="region"
    >
      {/* Hike overview (always shown) */}
      <HikeOverview hike={hike} index={index} isExpanded={isExpanded} onClick={onClick} onDelete={onDelete} />

      {/* Hike details (shown when overview is clicked) */}
      <div
        aria-label={`${hike.trailName} details wrapper`}
        role="region"
        // Easier to use inline CSS here than Tailwind
        style={{
          display: "grid",
          gridTemplateRows: isExpanded ? "1fr" : "0fr",
          transition: "grid-template-rows 0.3s ease",
        }}
      >
        <div className="overflow-hidden" {...(!isExpanded && { inert: true })}>
          <HikeDetail hike={hike} />
        </div>
      </div>
    </div>
  );
}
