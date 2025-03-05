"use client"

import dynamic from 'next/dynamic'

const SettingsPage = dynamic(
  () => import('@/components/pages/settings'),
  { ssr: false }
)

export default function Page() {
  return <SettingsPage />
}
