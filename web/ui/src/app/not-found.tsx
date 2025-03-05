'use client'

import dynamic from 'next/dynamic'

const NotFoundPage = dynamic(
  () => import('@/components/pages/not-found'),
  {
    ssr: false,
    loading: () => (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-6">
          <h1 className="text-6xl font-bold text-primary">404</h1>
          <h2 className="text-2xl font-semibold">Loading...</h2>
        </div>
      </div>
    )
  }
)

export default function NotFound() {
  return <NotFoundPage />
}
