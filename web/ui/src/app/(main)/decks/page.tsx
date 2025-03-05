"use client"

import dynamic from 'next/dynamic'

const DecksPage = dynamic(
  () => import('@/components/pages/decks'),
  { ssr: false }
)

export default function Page() {
  return <DecksPage />
}
