"use client"

import dynamic from 'next/dynamic';

const DashboardDetailsPage = dynamic(
  () => import('@/components/pages/dashboard-details'),
  { ssr: false }
)

export default async function Page(
  {
    params,
  }: {
    params: Promise<{ slug: string }>
  }
) {

  const {slug} = await params ;

  return <DashboardDetailsPage reference={slug} />;
}
