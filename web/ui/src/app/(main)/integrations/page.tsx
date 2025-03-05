"use client"

import dynamic from 'next/dynamic'

const IntegrationsPage = dynamic(
  () => import('@/components/pages/integrations'),
  { ssr: false }
)

export default function Page() {
  return <IntegrationsPage />
}
