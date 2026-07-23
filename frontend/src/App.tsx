import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import HikeLog from "./components/HikeLog/HikeLog";

const queryClient = new QueryClient();

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <HikeLog />;
    </QueryClientProvider>
  );
}
