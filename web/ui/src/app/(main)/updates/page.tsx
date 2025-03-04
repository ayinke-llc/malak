import dynamic from 'next/dynamic'

const UpdatesPage = dynamic(
  () => import('@/components/pages/updates'),
  { ssr: false }
)

export default function Page() {
  return <UpdatesPage />
}
