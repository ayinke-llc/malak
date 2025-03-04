"use client";

import dynamic from 'next/dynamic';

const ContactDetailsPage = dynamic(
  () => import('@/components/pages/contact-details'),
  { ssr: !!false }
)

export default async function Page(
  {
    params,
  }: {
    params: { slug: string }
  }
) {
  return <ContactDetailsPage reference={params.slug} />;
}
