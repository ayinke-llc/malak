"use client";
import { PostHogProvider } from "posthog-js/react";
import type React from "react";

// Singleton flag to track initialization
let isInitialized = false;

export function CSPostHogProvider({ children }: { children: React.ReactNode }) {
  // If PostHog is disabled or already initialized, just render children
  if (!process.env.NEXT_PUBLIC_MALAK_ENABLE_POSTHOG || isInitialized) {
    return <>{children}</>;
  }

  // Mark as initialized before rendering the provider
  isInitialized = true;

  return (
    <PostHogProvider
      apiKey={process.env.NEXT_PUBLIC_MALAK_POSTHOG_KEY}
      options={{
        api_host: process.env.NEXT_PUBLIC_MALAK_POSTHOG_HOST,
      }}
    >
      {children}
    </PostHogProvider>
  );
}
