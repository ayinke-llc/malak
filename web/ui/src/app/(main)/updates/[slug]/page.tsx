"use client"

import dynamic from 'next/dynamic';

const UpdateDetailsPage = dynamic(
  () => import('@/components/pages/update-details'),
  { ssr: false }
)

export default async function Page(
  {
    params,
  }: {
    params: Promise<{ slug: string }>
  }
) {
  const { slug } = await params;

  return <UpdateDetailsPage reference={slug} />;
}
