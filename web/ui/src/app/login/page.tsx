"use client";

import dynamic from 'next/dynamic'

const LoginPage = dynamic(
  () => import('@/components/pages/login'),
  { ssr: false }
)

export default function Page() {
  return <LoginPage />
}
