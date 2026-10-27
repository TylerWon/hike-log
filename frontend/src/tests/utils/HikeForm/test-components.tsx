import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import HikeForm, { type HikeFormProps } from "../../../pages/HikeLog/HikeForm/HikeForm";

/**
 * Wraps a HikeForm with a Tanstack Query client for testing.
 */
export function HikeFormWrapper({ onCancel, onSubmit }: HikeFormProps) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <HikeForm onCancel={onCancel} onSubmit={onSubmit} />
    </QueryClientProvider>
  );
}
