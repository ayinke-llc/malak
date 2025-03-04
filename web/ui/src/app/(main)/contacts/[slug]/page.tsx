"use client";

import { MalakContact, MalakContactShareItem } from "@/client/Api";
import dynamic from 'next/dynamic'
import { useRouter } from "next/navigation";
import { toast } from "sonner";

const ContactDetailsPage = dynamic(
  () => import('@/components/pages/contact-details'),
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
  return <ContactDetailsPage reference={slug} />;
}
