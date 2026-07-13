import { useState } from "react";
import { mockHikeData } from "../data/hikes";
import HikeCard from "./HikeCard";
import Header from "./Header/Header";

export default function HikeLog() {
  const [expandedCardId, setExpandedCardId] = useState<string | null>(null);

  const hikes = [...mockHikeData].sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());

  const handleCardClick = (id: string) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

  return (
    <div className="min-h-screen bg-forest-900 text-cream-100">
      <Header hikes={hikes} />

      <main className="max-w-3xl mx-auto px-6 pb-24">
        <ol className="flex flex-col gap-3">
          {hikes.map((hike, index) => (
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
