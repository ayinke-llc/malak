import DashboardDetailsPage from '@/components/pages/dashboard-details'

export default async function Page(
  {
    params,
  }: {
    params: Promise<{ slug: string }>
  }
) {

  const { slug } = await params;

  return <DashboardDetailsPage reference={slug} />;
}
