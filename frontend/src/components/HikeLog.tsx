import { useState } from "react";
import { hikes } from "../data/hikes";
import HikeCard from "./HikeCard";

function formatTotalDuration(minutes: number): string {
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
}

export default function HikeLog() {
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const sorted = [...hikes].sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
  );

  const totalDistanceKm = sorted.reduce((sum, h) => sum + h.distanceKm, 0);
  const totalElevationM = sorted.reduce((sum, h) => sum + h.elevationGainM, 0);
  const totalMinutes = sorted.reduce((sum, h) => sum + (h.durationMinutes ?? 0), 0);

  const toggle = (id: string) =>
    setExpandedId((prev) => (prev === id ? null : id));

  return (
    <div className="min-h-screen bg-forest-900 text-cream-100">
      <header className="max-w-3xl mx-auto px-6 pt-16 pb-10">
        <p
          className="text-xs font-mono tracking-[0.2em] uppercase text-sage-500 mb-3"
        >
          Field Notes
        </p>
        <h1
          className="text-5xl font-bold font-serif leading-tight mb-6"
        >
          Hiking Log
        </h1>

        {/* Overall stats */}
        <div className="flex flex-wrap gap-6">
          <div>
            <p
              className="text-[10px] font-mono uppercase tracking-widest text-forest-700 mb-1"
            >
              Hikes
            </p>
            <p
              className="text-2xl font-semibold font-serif text-cream-100"
            >
              {sorted.length}
            </p>
          </div>
          <div className="w-px bg-forest-800 self-stretch" />
          <div>
            <p
              className="text-[10px] font-mono uppercase tracking-widest text-forest-700 mb-1"
            >
              Distance
            </p>
            <p
              className="text-2xl font-semibold font-serif text-cream-100"
            >
              {totalDistanceKm.toFixed(1)}{" "}
              <span className="text-sm font-normal text-forest-600">km</span>
            </p>
          </div>
          <div className="w-px bg-forest-800 self-stretch" />
          <div>
            <p
              className="text-[10px] font-mono uppercase tracking-widest text-forest-700 mb-1"
            >
              Elevation
            </p>
            <p
              className="text-2xl font-semibold font-serif text-cream-100"
            >
              {totalElevationM.toLocaleString()}{" "}
              <span className="text-sm font-normal text-forest-600">m</span>
            </p>
          </div>
          {totalMinutes > 0 && (
            <>
              <div className="w-px bg-forest-800 self-stretch" />
              <div>
                <p
                  className="text-[10px] font-mono uppercase tracking-widest text-forest-700 mb-1"
                >
                  Time
                </p>
                <p
                  className="text-2xl font-semibold font-serif text-cream-100"
                >
                  {formatTotalDuration(totalMinutes)}
                </p>
              </div>
            </>
          )}
        </div>

        <div className="mt-8 h-px bg-forest-800" />
      </header>

      <main className="max-w-3xl mx-auto px-6 pb-24">
        <ol className="space-y-3">
          {sorted.map((hike, index) => (
            <li key={hike.id}>
              <HikeCard
                hike={hike}
                index={index + 1}
                isExpanded={expandedId === hike.id}
                onToggle={() => toggle(hike.id)}
              />
            </li>
          ))}
        </ol>
      </main>
    </div>
  );
}
