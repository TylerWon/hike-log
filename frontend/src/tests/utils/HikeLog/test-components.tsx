import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import HikeCard from "../../../pages/HikeLog/HikeCard";
import HikeLog from "../../../pages/HikeLog/HikeLog/HikeLog";
import HikeLogContent from "../../../pages/HikeLog/HikeLog/HikeLogContent";
import { formatDistance, formatDuration, formatElevation } from "../../../utils/formatters";
import { HIKE_MOCK_1, HIKE_MOCK_2 } from "../../mocks/models/hike";

/**
 * Wraps a HikeLog with a TanStack Query client for testing.
 */
export function HikeLogWrapper() {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: { queries: { retry: false } }, // disable retry to prevent error test from timing out due to retries by Tanstack Query
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <HikeLog />
    </QueryClientProvider>
  );
}

/**
 * Wraps a HikeLogContent to provide default prop values.
 */
export function HikeLogContentWrapper() {
  const hikes = [HIKE_MOCK_1, HIKE_MOCK_2];
  const overallStats = [
    { label: "Hikes", value: hikes.length },
    { label: "Distance", value: formatDistance(hikes[0].distance + hikes[1].distance) },
    { label: "Elevation", value: formatElevation(hikes[0].elevationGain + hikes[1].elevationGain) },
    { label: "Time", value: formatDuration(hikes[0].duration + hikes[1].duration) },
  ];

  return (
    <HikeLogContent overallStats={overallStats}>
      {hikes.map((hike, i) => (
        <li key={i}>
          <HikeCard hike={hike} index={i + 1} isExpanded={false} onClick={() => {}} onDelete={async () => {}} />
        </li>
      ))}
    </HikeLogContent>
  );
}
