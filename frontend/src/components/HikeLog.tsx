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
    <div className="min-h-screen bg-[#0f1410] text-[#e8dfc8]">
      <header className="max-w-3xl mx-auto px-6 pt-16 pb-10">
        <p
          className="text-xs tracking-[0.2em] uppercase text-[#7fa86a] mb-3"
          style={{ fontFamily: "'DM Mono', monospace" }}
        >
          Field Notes
        </p>
        <h1
          className="text-5xl font-bold leading-tight mb-6"
          style={{ fontFamily: "'Playfair Display', serif" }}
        >
          Hiking Log
        </h1>

        {/* Overall stats */}
        <div className="flex flex-wrap gap-6">
          <div>
            <p
              className="text-[10px] uppercase tracking-widest text-[#4a5c44] mb-1"
              style={{ fontFamily: "'DM Mono', monospace" }}
            >
              Hikes
            </p>
            <p
              className="text-2xl font-semibold text-[#e8dfc8]"
              style={{ fontFamily: "'Playfair Display', serif" }}
            >
              {sorted.length}
            </p>
          </div>
          <div className="w-px bg-[#2e3d2a] self-stretch" />
          <div>
            <p
              className="text-[10px] uppercase tracking-widest text-[#4a5c44] mb-1"
              style={{ fontFamily: "'DM Mono', monospace" }}
            >
              Distance
            </p>
            <p
              className="text-2xl font-semibold text-[#e8dfc8]"
              style={{ fontFamily: "'Playfair Display', serif" }}
            >
              {totalDistanceKm.toFixed(1)}{" "}
              <span className="text-sm font-normal text-[#8a9e82]">km</span>
            </p>
          </div>
          <div className="w-px bg-[#2e3d2a] self-stretch" />
          <div>
            <p
              className="text-[10px] uppercase tracking-widest text-[#4a5c44] mb-1"
              style={{ fontFamily: "'DM Mono', monospace" }}
            >
              Elevation
            </p>
            <p
              className="text-2xl font-semibold text-[#e8dfc8]"
              style={{ fontFamily: "'Playfair Display', serif" }}
            >
              {totalElevationM.toLocaleString()}{" "}
              <span className="text-sm font-normal text-[#8a9e82]">m</span>
            </p>
          </div>
          {totalMinutes > 0 && (
            <>
              <div className="w-px bg-[#2e3d2a] self-stretch" />
              <div>
                <p
                  className="text-[10px] uppercase tracking-widest text-[#4a5c44] mb-1"
                  style={{ fontFamily: "'DM Mono', monospace" }}
                >
                  Time
                </p>
                <p
                  className="text-2xl font-semibold text-[#e8dfc8]"
                  style={{ fontFamily: "'Playfair Display', serif" }}
                >
                  {formatTotalDuration(totalMinutes)}
                </p>
              </div>
            </>
          )}
        </div>

        <div className="mt-8 h-px bg-[#2e3d2a]" />
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
