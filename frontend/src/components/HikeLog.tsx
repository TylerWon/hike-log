import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

import { fetchHikes } from "../api/hikes";
import { formatDistance, formatDuration, formatElevation } from "../utils/formatters";
import HikeCard from "./HikeCard";
import HikeCardSkeleton from "./HikeCardSkeleton";
import HikeLogContent from "./HikeLogContent";
import HikeLogError from "./HikeLogError";
import "../assets/styles/animation.css";
import StatValueSkeleton from "./StatValueSkeleton";

export default function HikeLog() {
  const [expandedCardId, setExpandedCardId] = useState<bigint | null>(null);

  const hikes = useQuery({
    queryFn: fetchHikes,
    queryKey: ["hikes"],
    refetchOnWindowFocus: false,
  });

  if (hikes.isError) {
    return <HikeLogError />;
  }

  if (hikes.isPending) {
    const overallStats = [
      { label: "Hikes", value: <StatValueSkeleton widthClass="w-8" /> },
      { label: "Distance", value: <StatValueSkeleton widthClass="w-20" /> },
      { label: "Elevation", value: <StatValueSkeleton widthClass="w-24" /> },
      { label: "Time", value: <StatValueSkeleton widthClass="w-16" /> },
    ];

    return (
      <HikeLogContent overallStats={overallStats}>
        <ol className="flex flex-col gap-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <li key={i}>
              <HikeCardSkeleton />
            </li>
          ))}
        </ol>
      </HikeLogContent>
    );
  }

  const handleCardClick = (id: bigint) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

  const totalDistance = hikes.data.reduce<number>((sum, h) => sum + h.distance, 0);
  const totalElevation = hikes.data.reduce<bigint>((sum, h) => sum + h.elevationGain, BigInt(0));
  const totalMinutes = hikes.data.reduce<bigint>((sum, h) => sum + h.duration, BigInt(0));

  const overallStats = [
    { label: "Hikes", value: hikes.data.length },
    { label: "Distance", value: formatDistance(totalDistance) },
    { label: "Elevation", value: formatElevation(totalElevation) },
    { label: "Time", value: formatDuration(totalMinutes) },
  ];

  return (
    <HikeLogContent overallStats={overallStats}>
      <ol className="flex flex-col gap-3">
        {hikes.data.map((hike, i) => (
          <li key={i}>
            <HikeCard
              hike={hike}
              index={i + 1}
              isExpanded={expandedCardId === hike.id}
              onClick={() => handleCardClick(hike.id)}
            />
          </li>
        ))}
      </ol>
    </HikeLogContent>
  );
}
