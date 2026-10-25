import type { Hike } from "../../schemas/models/hike";

import { classnames } from "../../utils/classnames";
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
      <HikeOverview hike={hike} index={index} isExpanded={isExpanded} onClick={onClick} onDelete={onDelete} />
      {isExpanded && <HikeDetail hike={hike} />}
    </div>
  );
}
