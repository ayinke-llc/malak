import UpdateDetailsPage from "@/components/pages/update-details"

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
