import "../assets/styles/animation.css";
import "../assets/styles/button.css";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import type { HikeFormData } from "../schemas/forms/hike";

import { createHike, createPhoto, createPresignedUrl, fetchHikes } from "../api/hikes";
import { uploadFile } from "../api/s3";
import { formatDistance, formatDuration, formatElevation } from "../utils/formatters";
import HikeCard from "./HikeCard";
import HikeCardSkeleton from "./HikeCardSkeleton";
import HikeForm from "./HikeForm/HikeForm";
import HikeLogContent from "./HikeLogContent";
import HikeLogError from "./HikeLogError";
import StatValueSkeleton from "./StatValueSkeleton";

const HIKES_QUERY_KEY = "hikes";

export default function HikeLog() {
  const queryClient = useQueryClient();

  const [expandedCardId, setExpandedCardId] = useState<null | number>(null);
  const [showHikeForm, setShowHikeForm] = useState<boolean>(false);

  const hikes = useQuery({
    queryFn: fetchHikes, // called on component load
    queryKey: [HIKES_QUERY_KEY],
    refetchOnWindowFocus: false,
  });

  const refetchHikes = async () => {
    // Invalidates "hikes" query so all Hikes get refetched
    await queryClient.invalidateQueries({ queryKey: [HIKES_QUERY_KEY] });
  };

  const addHike = async (formData: HikeFormData) => {
    const { photos, ...hikeData } = formData;

    // Hike must be created so fail-close
    let hike;
    try {
      hike = await createHike(hikeData);
    } catch (e) {
      throw new Error(`Failed to create hike: ${e}`, { cause: e });
    }

    // Photo upload is best effort so fail-open
    let promises = [];
    let presignedUrls;
    try {
      for (const photo of photos) {
        const reqBody = {
          contentLength: photo.file.size,
          contentType: photo.file.type,
        };
        promises.push(createPresignedUrl(hike.id, reqBody));
      }
      presignedUrls = await Promise.all(promises);
    } catch (e) {
      console.warn("Failed to create presigned URL for photo: ", e);
      return;
    }

    try {
      promises = [];
      for (let i = 0; i < photos.length; i++) {
        promises.push(uploadFile(presignedUrls[i].uploadUrl, photos[i].file));
      }
      await Promise.all(promises);
    } catch (e) {
      console.warn("Failed to upload photo to S3: ", e);
      return;
    }

    try {
      promises = [];
      for (let i = 0; i < photos.length; i++) {
        const reqBody = {
          caption: photos[i].caption,
          displayOrder: photos[i].displayOrder,
          objectKey: presignedUrls[i].objectKey,
        };
        promises.push(createPhoto(hike.id, reqBody));
      }
      await Promise.all(promises);
    } catch (e) {
      console.warn("Failed to create photo: ", e);
      return;
    }

    setShowHikeForm(false);
  };

  const addHikeMutation = useMutation({
    mutationFn: addHike, // called when mutate() is invoked for this mutation
    onSuccess: refetchHikes,
  });

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

  const handleCardClick = (id: number) => {
    setExpandedCardId((prev) => (prev === id ? null : id));
  };

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
        <HikeForm onCancel={() => setShowHikeForm(false)} onSubmit={(hike) => addHikeMutation.mutate(hike)} />
      )}
    </>
  );
}
