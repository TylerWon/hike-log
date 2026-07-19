import { type Hike } from "../data/hikes";
import { classnames } from "../utils/classnames";
import HikeDetail from "./HikeDetail";
import HikeOverview from "./HikeOverview/HikeOverview";

interface HikeCardProps {
  hike: Hike;
  index: number;
  isExpanded: boolean;
  onClick: () => void;
}

export default function HikeCard({ hike, index, isExpanded, onClick }: HikeCardProps) {
  return (
    <article
      className={classnames(
        "border rounded-lg overflow-hidden transition-all duration-200",
        isExpanded ? "bg-forest-850 border-forest-700" : "bg-forest-900 border-forest-800",
      )}
    >
      {/* Hike overview (always shown) */}
      <button onClick={onClick} className="w-full text-left cursor-pointer group/hike-card" aria-expanded={isExpanded}>
        <HikeOverview hike={hike} index={index} isExpanded={isExpanded} />
      </button>

      {/* Hike details (shown when overview is clicked) */}
      <div
        style={{
          display: "grid",
          gridTemplateRows: isExpanded ? "1fr" : "0fr",
          transition: "grid-template-rows 0.3s ease",
        }}
      >
        <div className="overflow-hidden">
          <HikeDetail hike={hike} photos={hike.photos} />
        </div>
      </div>
    </article>
  );
}
