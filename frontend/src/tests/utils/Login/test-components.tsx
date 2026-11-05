import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { MemoryRouter } from "react-router";

import Login, { type LoginProps } from "../../../pages/Login/Login";

/**
 * Wraps a Login with a Tanstack Query client and react-router MemoryRouter for testing.
 */
export function LoginWrapper({ navigateProp, searchParamsProp }: LoginProps) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <Login navigateProp={navigateProp} searchParamsProp={searchParamsProp} />
      </MemoryRouter>
    </QueryClientProvider>
  );
}
