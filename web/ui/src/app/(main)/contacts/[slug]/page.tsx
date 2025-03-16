import ContactDetailsPage from "@/components/pages/contact-details"

export default async function Page(
  {
    params,
  }: {
    params: Promise<{ slug: string }>
  }
) {
  const { slug } = await params

  return <ContactDetailsPage reference={slug} />;
}
