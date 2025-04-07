import FundraisingPage from '@/components/pages/fundraising'

export default async function Page(
  {
    params,
  }: {
    params: Promise<{ slug: string }>
  }
) {
  const { slug } = await params;

  return <FundraisingPage slug={slug} />
}
