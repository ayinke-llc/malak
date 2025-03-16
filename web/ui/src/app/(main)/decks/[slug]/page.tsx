import DeckDetailsPage from '@/components/pages/deck-details'

export default async function Page(
  {
    params,
  }: {
    params: Promise<{ slug: string }>
  }
) {

  const { slug } = await params;

  return <DeckDetailsPage reference={slug} />;
}
