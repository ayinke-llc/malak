"use client"

import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Toaster } from '../ui/toaster';
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from '../ui/tooltip';
import { ThemeProvider } from "@/components/providers/theme";
import { GoogleOAuthProvider } from '@react-oauth/google';
import { GOOGLE_CLIENT_ID } from "@/lib/config";

export default function Providers({ children }: { children: React.ReactNode }) {

  const queryClient = new QueryClient()

  return (

    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={false} />
        <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
          <TooltipProvider>
            <Toaster />
            <Sonner position="top-center" />
            {children}
          </TooltipProvider>
        </GoogleOAuthProvider>
      </QueryClientProvider >
    </ThemeProvider>
  )
}
