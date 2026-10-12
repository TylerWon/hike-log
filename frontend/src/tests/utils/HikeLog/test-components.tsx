import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import HikeLog from "../../../components/HikeLog/HikeLog";

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
