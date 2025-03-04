import dynamic from 'next/dynamic'

const NotFoundPage = dynamic(
  () => import('@/components/pages/not-found'),
  { ssr: !!false }
)

export default function NotFound() {
  return <NotFoundPage />
}
