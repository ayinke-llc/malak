import dynamic from 'next/dynamic'

const DashboardsPage = dynamic(
  () => import('@/components/pages/dashboards'),
  { ssr: false }
)

export default function Page() {
  return <DashboardsPage />
}
