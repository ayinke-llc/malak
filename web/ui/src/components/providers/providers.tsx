"use client";

import { GOOGLE_CLIENT_ID } from "@/lib/config";
import { GoogleOAuthProvider } from "@react-oauth/google";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Toaster } from "sonner";
import { CSPostHogProvider } from "./posthog";
import UserProvider from "./user";

export default function Providers({ children }: { children: React.ReactNode }) {
  const queryClient = new QueryClient();

  return (
    <CSPostHogProvider>
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={false} />
        <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
          <Toaster position="top-center" richColors />
          <UserProvider>{children}</UserProvider>
        </GoogleOAuthProvider>
      </QueryClientProvider>
    </CSPostHogProvider>
  );
}
