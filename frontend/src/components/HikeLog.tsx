import "../assets/styles/animation.css";
import "../assets/styles/button.css";

import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

import type { HikeFormData } from "../schemas/forms/hike";

import { fetchHikes } from "../api/hikes";
import { formatDistance, formatDuration, formatElevation } from "../utils/formatters";
import HikeCard from "./HikeCard";
import HikeCardSkeleton from "./HikeCardSkeleton";
import HikeForm from "./HikeForm/HikeForm";
import HikeLogContent from "./HikeLogContent";
import HikeLogError from "./HikeLogError";
import StatValueSkeleton from "./StatValueSkeleton";

const HIKES_QUERY_KEY = "hikes";

export default function HikeLog() {
  const [expandedCardId, setExpandedCardId] = useState<bigint | null>(null);
  const [showHikeForm, setShowHikeForm] = useState<boolean>(false);

  const hikes = useQuery({
    queryFn: fetchHikes,
    queryKey: [HIKES_QUERY_KEY],
    refetchOnWindowFocus: false,
  });

  const handleCardClick = (id: bigint) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

  const handleHikeFormSubmit = (hike: HikeFormData) => {
    console.log(hike);
    setShowHikeForm(false);
    return;
  };

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
        {Array.from({ length: 5 }).map((_, i) => (
          <li key={i}>
            <HikeCardSkeleton />
          </li>
        ))}
      </HikeLogContent>
    );
  }

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
    <>
      <HikeLogContent overallStats={overallStats}>
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
        <li className="mt-4 flex justify-center">
          <button className="primary-button px-3 py-1.5" onClick={() => setShowHikeForm(true)}>
            <svg
              fill="none"
              height="11"
              stroke="currentColor"
              strokeLinecap="round"
              strokeWidth="1.5"
              viewBox="0 0 11 11"
              width="11"
            >
              <path d="M5.5 1v9M1 5.5h9" />
            </svg>
          </button>
        </li>
      </HikeLogContent>
      {showHikeForm && <HikeForm onCancel={() => setShowHikeForm(false)} onSubmit={handleHikeFormSubmit} />}
    </>
  );
}
