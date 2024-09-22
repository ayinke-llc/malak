'use client'
import { PostHogProvider } from 'posthog-js/react'
import React from 'react'

export function CSPostHogProvider({ children }: { children: React.ReactNode }) {

  if (!process.env.NEXT_PUBLIC_MALAK_ENABLE_POSTHOG) {
    return (
      <>
        {children}
      </>
    )
  }

  return (
    <PostHogProvider
      apiKey={process.env.NEXT_PUBLIC_MALAK_POSTHOG_KEY}
      options={
        {
          api_host: process.env.NEXT_PUBLIC_MALAK_POSTHOG_HOST
        }
      }>
      {children}
    </PostHogProvider>
  )
}
