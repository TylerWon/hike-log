import { useQuery } from "@tanstack/react-query";

import "../../assets/styles/text.css";

import { useState } from "react";

import "./hike-log.css";
import { fetchHikes } from "../../api/hikes";
import HikeCard from "../HikeCard";

export default function HikeLog() {
  const [expandedCardId, setExpandedCardId] = useState<bigint | null>(null);

  const hikes = useQuery({
    queryFn: fetchHikes,
    queryKey: ["hikes"],
  });

  if (hikes.isPending) return "Loading...";

  if (hikes.isError) {
    console.error(`Error while loading hikes: ${hikes.error}`);
    return "Something went wrong. Please try again later.";
  }

  const totalDistanceKm = hikes.data.reduce<number>((sum, h) => sum + h.distance, 0);
  const totalElevationM = hikes.data.reduce<bigint>((sum, h) => sum + h.elevationGain, BigInt(0));
  const totalMinutes = hikes.data.reduce<bigint>((sum, h) => sum + h.duration, BigInt(0));

  const handleCardClick = (id: bigint) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

  return (
    <div className="min-h-screen bg-forest-900 text-cream-100">
      <header className="max-w-3xl mx-auto px-6 pt-16 pb-10">
        <p className="text-xs font-mono tracking-[0.2em] uppercase text-sage-500 mb-3">Field Notes</p>
        <h1 className="text-5xl font-bold font-serif leading-tight mb-6">Hiking Log</h1>

        <div className="flex flex-wrap gap-6">
          <div>
            <p className="field-label mb-1">Hikes</p>
            <p className="overall-stat-value">{hikes.data.length}</p>
          </div>
          <div className="overall-stat-divider" />
          <div>
            <p className="field-label mb-1">Distance</p>
            <p className="overall-stat-value">{totalDistanceKm.toFixed(1)} km</p>
          </div>
          <div className="overall-stat-divider" />
          <div>
            <p className="field-label mb-1">Elevation</p>
            <p className="overall-stat-value">{totalElevationM.toLocaleString()} m</p>
          </div>
          <div className="overall-stat-divider" />
          <div>
            <p className="field-label mb-1">Time</p>
            <p className="overall-stat-value">{formatTotalDuration(Number(totalMinutes))}</p>
          </div>
        </div>

        <div className="mt-8 h-px bg-forest-800" />
      </header>

      <main className="max-w-3xl mx-auto px-6 pb-24">
        <ol className="flex flex-col gap-3">
          {hikes.data.map((hike, index) => (
            <li key={hike.id}>
              <HikeCard
                hike={hike}
                index={index + 1}
                isExpanded={expandedCardId === hike.id}
                onClick={() => handleCardClick(hike.id)}
              />
            </li>
          ))}
        </ol>
      </main>
    </div>
  );
}

function formatTotalDuration(minutes: number): string {
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
}
