import dynamic from 'next/dynamic'

const BillingPage = dynamic(
  () => import('@/components/pages/billing'),
  { ssr: false }
)

export default function Page() {
  return <BillingPage />
}
