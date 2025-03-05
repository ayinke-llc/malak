"use client"

import dynamic from 'next/dynamic'

const ContactsPage = dynamic(
  () => import('@/components/pages/contacts'),
  { ssr: false }
)

export default function Page() {
  return <ContactsPage />
}
