import dynamic from 'next/dynamic'

const GeneralSettingsPage = dynamic(
  () => import('@/components/pages/general-settings'),
  { ssr: false }
)

export default function Page() {
  return <GeneralSettingsPage />
}
