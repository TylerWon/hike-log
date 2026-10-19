import type { Mock } from "vitest";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import ConfirmationDialog from "../../../components/ConfirmationDialog";

/**
 * Wraps a ConfirmationDialog with a TanStack Query client for testing.
 */
export function ConfirmationDialogWrapper({
  actionType = "default",
  onAccept,
  onCancel,
}: {
  actionType?: "default" | "destructive";
  onAccept: Mock;
  onCancel: Mock;
}) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: { queries: { retry: false } }, // disable retry to prevent error test from timing out due to retries by Tanstack Query
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <ConfirmationDialog
        acceptButtonPendingText="Adding..."
        acceptButtonText="Add"
        actionType={actionType}
        description={<span>Do you want to add this item?</span>}
        onAccept={onAccept}
        onCancel={onCancel}
        title="Add item"
      />
    </QueryClientProvider>
  );
}
