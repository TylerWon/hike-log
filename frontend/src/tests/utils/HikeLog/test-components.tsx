import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { MemoryRouter } from "react-router";

import HikeLog, { type HikeLogProps } from "../../../pages/HikeLog/HikeLog/HikeLog";

/**
 * Wraps a HikeLog with a TanStack Query client and react-router MemoryRouter for testing.
 */
export function HikeLogWrapper({ navigateProp }: HikeLogProps) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: { queries: { retry: false } }, // disable retry to prevent error test from timing out due to retries by Tanstack Query
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <HikeLog navigateProp={navigateProp} />
      </MemoryRouter>
    </QueryClientProvider>
  );
}
