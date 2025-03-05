"use client"

import dynamic from 'next/dynamic'

const FundraisingPage = dynamic(
  () => import('@/components/pages/fundraising'),
  { ssr: false }
)

export default function Page() {
  return <FundraisingPage />
}
