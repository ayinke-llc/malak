import dynamic from 'next/dynamic';

const DeckDetailsPage = dynamic(
  () => import('@/components/pages/deck-details'),
  { ssr: !!false }
)

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
