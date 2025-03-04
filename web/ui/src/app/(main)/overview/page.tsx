import dynamic from 'next/dynamic'

const OverviewPage = dynamic(
  () => import('@/components/pages/overview'),
  { ssr: false }
)

export default function Page() {
  return <OverviewPage />
}
