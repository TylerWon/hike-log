import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import HikeOverview from "../../../components/HikeOverview";
import { HIKE_MOCK_1 } from "../../mocks/models/hike";

/**
 * Wraps a HikeOverview with a TanStack Query client for testing.
 */
export function HikeOverviewWrapper({
  isExpanded = false,
  onClick,
  onDelete,
}: {
  isExpanded?: boolean;
  onClick: () => void;
  onDelete: () => Promise<void>;
}) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: { queries: { retry: false } }, // disable retry to prevent error test from timing out due to retries by Tanstack Query
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <HikeOverview hike={HIKE_MOCK_1} index={1} isExpanded={isExpanded} onClick={onClick} onDelete={onDelete} />
    </QueryClientProvider>
  );
}
