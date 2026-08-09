import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

import { fetchHikes } from "../api/hikes";
import { formatDuration } from "../utils/formatters";
import HikeCard from "./HikeCard";
import HikeCardSkeleton from "./HikeCardSkeleton";
import HikeLogContent from "./HikeLogContent/HikeLogContent";
import HikeLogError from "./HikeLogError";

export default function HikeLog() {
  const [expandedCardId, setExpandedCardId] = useState<bigint | null>(null);

  const hikes = useQuery({
    queryFn: fetchHikes,
    queryKey: ["hikes"],
    refetchOnWindowFocus: false,
  });

  if (hikes.isPending) {
    return (
      <HikeLogContent
        hikeCards={Array.from({ length: 5 }).map((_, i) => (
          <HikeCardSkeleton key={i} />
        ))}
        overallStats={[
          { label: "Hikes", value: <div className="shimmer rounded h-7 w-8" /> },
          { label: "Distance", value: <div className="shimmer rounded h-7 w-20" /> },
          { label: "Elevation", value: <div className="shimmer rounded h-7 w-24" /> },
          { label: "Time", value: <div className="shimmer rounded h-7 w-16" /> },
        ]}
      />
    );
  }

  if (hikes.isError) {
    return <HikeLogError />;
  }

  const totalDistanceKm = hikes.data.reduce<number>((sum, h) => sum + h.distance, 0);
  const totalElevationM = hikes.data.reduce<bigint>((sum, h) => sum + h.elevationGain, BigInt(0));
  const totalMinutes = hikes.data.reduce<bigint>((sum, h) => sum + h.duration, BigInt(0));

  const handleCardClick = (id: bigint) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

  return (
    <HikeLogContent
      hikeCards={hikes.data.map((hike, i) => (
        <HikeCard
          hike={hike}
          index={i + 1}
          isExpanded={expandedCardId === hike.id}
          key={i}
          onClick={() => handleCardClick(hike.id)}
        />
      ))}
      overallStats={[
        { label: "Hikes", value: <p className="overall-stat-value">{hikes.data.length}</p> },
        { label: "Distance", value: <p className="overall-stat-value">{totalDistanceKm.toFixed(1)} km</p> },
        { label: "Elevation", value: <p className="overall-stat-value">{totalElevationM.toLocaleString()} m</p> },
        { label: "Time", value: <p className="overall-stat-value">{formatDuration(totalMinutes)}</p> },
      ]}
    />
  );
}
