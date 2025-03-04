import dynamic from 'next/dynamic'

const CapTablePage = dynamic(
  () => import('@/components/pages/captable'),
  { ssr: !!false }
)

export default function Page() {
  return <CapTablePage />
}
