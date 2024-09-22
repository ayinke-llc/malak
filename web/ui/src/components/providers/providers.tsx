"use client"

import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { GoogleOAuthProvider } from '@react-oauth/google';
import { Toaster } from "sonner"
import { GOOGLE_CLIENT_ID } from "@/lib/config";
import UserProvider from './user';
import { CSPostHogProvider } from './posthog';

export default function Providers({ children }: { children: React.ReactNode }) {

  const queryClient = new QueryClient()

  return (
    <CSPostHogProvider>
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={false} />
        <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
          <Toaster position="top-center" richColors />
          <UserProvider>
            {children}
          </UserProvider>
        </GoogleOAuthProvider>
      </QueryClientProvider >
    </CSPostHogProvider>
  )
}


