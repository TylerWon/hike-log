import "../../assets/styles/animation.css";
import "../../assets/styles/button.css";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import type { HikeFormData } from "../../schemas/forms/hike";

import { fetchHikes } from "../../api/hikes";
import { formatDistance, formatDuration, formatElevation } from "../../utils/formatters";
import HikeCard from "../HikeCard";
import HikeCardSkeleton from "../HikeCardSkeleton";
import HikeForm from "../HikeForm/HikeForm";
import HikeLogContent from "../HikeLogContent";
import HikeLogError from "../HikeLogError";
import StatValueSkeleton from "../StatValueSkeleton";
import Toast from "../Toast";
import { PhotoCreationError, submitHikeForm } from "./submit-hike-form";

const HIKES_QUERY_KEY = "hikes";

export default function HikeLog() {
  const queryClient = useQueryClient();

  const [expandedCardId, setExpandedCardId] = useState<null | number>(null);
  const [showHikeForm, setShowHikeForm] = useState<boolean>(false);
  const [toastMessage, setToastMessage] = useState<null | string>(null);

  const hikes = useQuery({
    queryFn: fetchHikes, // called on component load
    queryKey: [HIKES_QUERY_KEY],
    refetchOnWindowFocus: false,
  });

  const handleCardClick = (id: number) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

  const handleHikeFormClose = () => {
    setShowHikeForm(false);
  };

  const handleHikeFormSubmit = async (formData: HikeFormData): Promise<void> => {
    try {
      await submitHikeForm(formData);
    } catch (e) {
      if (e instanceof PhotoCreationError) {
        console.warn(`Failed to create some photos: ${e}`);
        setToastMessage("Some photos could not be uploaded. Please try again.");
        handleHikeFormClose();
        return;
      }

      throw e;
    }

    handleHikeFormClose();
  };

  const handleHikeFormSuccess = async () => {
    // Invalidates the list hikes query so all Hikes get refetched
    await queryClient.invalidateQueries({ queryKey: [HIKES_QUERY_KEY] });
  };

  if (hikes.isError) {
    console.error("Failed to fetch hikes: ", hikes.error);
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
  const totalElevation = hikes.data.reduce<number>((sum, h) => sum + h.elevationGain, 0);
  const totalMinutes = hikes.data.reduce<number>((sum, h) => sum + h.duration, 0);

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
      {showHikeForm && (
        <HikeForm onClose={handleHikeFormClose} onSubmit={handleHikeFormSubmit} onSuccess={handleHikeFormSuccess} />
      )}
      {toastMessage && <Toast message={toastMessage} onDismiss={() => setToastMessage(null)} />}
    </>
  );
}
