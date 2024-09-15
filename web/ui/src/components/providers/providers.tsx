"use client"

import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { ThemeProvider } from "@/components/providers/theme";
import { GoogleOAuthProvider } from '@react-oauth/google';
import { Toaster } from "sonner"
import { GOOGLE_CLIENT_ID } from "@/lib/config";
import UserProvider from './user';

export default function Providers({ children }: { children: React.ReactNode }) {

  const queryClient = new QueryClient()

  return (
    <QueryClientProvider client={queryClient}>
      <ReactQueryDevtools initialIsOpen={false} />
      <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
        <Toaster position="top-center" />
        <UserProvider>
          {children}
        </UserProvider>
      </GoogleOAuthProvider>
    </QueryClientProvider >
  )
}
