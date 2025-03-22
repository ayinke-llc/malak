"use client"

import { notFound, useParams } from 'next/navigation';
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip, XAxis, YAxis } from 'recharts';
import { useQuery } from '@tanstack/react-query';
import { PUBLIC_SHARED_DASHBOARD, PUBLIC_SHARED_DASHBOARD_CHARTING_DATA } from '@/lib/query-constants';
import type { MalakDashboardChart, ServerListDashboardChartsResponse, MalakDashboardChartPosition } from '@/client/Api';
import client from '@/lib/client';
import { Card } from '@/components/ui/card';
import { ChartContainer } from '@/components/ui/chart';
import { RiLoader4Line, RiBarChart2Line, RiPieChartLine } from '@remixicon/react';
import { formatChartData, formatTooltipValue, getChartColors } from '@/lib/chart-utils';
import { formatDistanceToNow } from 'date-fns';
import { useMemo } from 'react';

function ChartCard({ chart, dashboard_reference }: { chart: MalakDashboardChart, dashboard_reference: string }) {
  const { data: chartData, isLoading: isLoadingChartData, error } = useQuery({
    queryKey: [PUBLIC_SHARED_DASHBOARD_CHARTING_DATA, chart.chart?.reference],
    queryFn: async () => {
      if (!chart.chart?.reference) return null;
      const response = await client.public.dashboardsChartsDetail(dashboard_reference, chart.chart?.reference);
      return response.data;
    },
    enabled: !!chart.chart?.reference,
  });

  const formattedData = formatChartData(chartData?.data_points, chart.chart?.data_point_type);

  if (isLoadingChartData) {
    return (
      <Card className="p-3">
        <div className="flex items-center justify-center h-[160px]">
          <RiLoader4Line className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      </Card>
    );
  }

  if (error || !chart.chart) {
    return (
      <Card className="p-3">
        <div className="flex flex-col items-center justify-center h-[160px] text-center p-4">
          <RiBarChart2Line className="h-8 w-8 text-muted-foreground mb-2" />
          <p className="text-sm text-muted-foreground">Failed to load chart data</p>
          <p className="text-xs text-muted-foreground mt-1">Please try again later</p>
        </div>
      </Card>
    );
  }

  const hasNoData = !formattedData || formattedData.length === 0;

  return (
    <Card className="p-3 transition-colors duration-200 hover:bg-accent/5">
      <div className="flex items-center justify-between mb-1">
        <div className="flex items-center gap-2">
          <div className="text-muted-foreground">
            {chart.chart.chart_type === "pie" ? (
              <RiPieChartLine className="h-4 w-4" />
            ) : (
              <RiBarChart2Line className="h-4 w-4" />
            )}
          </div>
          <div>
            <h3 className="text-sm font-bold">{chart.chart.user_facing_name}</h3>
          </div>
        </div>
      </div>
      <div className="w-full">
        {hasNoData ? (
          <div className="flex flex-col items-center justify-center h-[160px] text-center p-4">
            <RiBarChart2Line className="h-8 w-8 text-muted-foreground mb-2" />
            <p className="text-sm text-muted-foreground">No data available</p>
            <p className="text-xs text-muted-foreground mt-1">Check back later for updates</p>
          </div>
        ) : chart.chart.chart_type === "bar" ? (
          <ChartContainer className="w-full h-full" config={{}}>
            <BarChart
              width={390}
              height={160}
              data={formattedData}
              margin={{ top: 5, right: 5, left: -15, bottom: 0 }}
            >
              <XAxis dataKey="name" stroke="#888888" fontSize={11} />
              <YAxis stroke="#888888" fontSize={11} />
              <Tooltip
                formatter={(value: number) =>
                  formatTooltipValue(value, chart.chart?.data_point_type)
                }
              />
              <Bar dataKey="value" fill="#3B82F6" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ChartContainer>
        ) : (
          <ChartContainer className="w-full h-full" config={{}}>
            <PieChart
              width={390}
              height={160}
              margin={{ top: 5, right: 5, left: 5, bottom: 5 }}
            >
              <Pie
                data={formattedData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={60}
                dataKey="value"
              >
                {formattedData.map((_entry, index) => (
                  <Cell key={`cell-${index}`} fill={getChartColors(index)} />
                ))}
              </Pie>
              <Tooltip
                formatter={(value: number) =>
                  formatTooltipValue(value, chart.chart?.data_point_type)
                }
              />
            </PieChart>
          </ChartContainer>
        )}
      </div>
    </Card>
  );
}

export default function SharedDashboardPage() {

  const params = useParams()

  const slug = params.slug as string

  const { data, isLoading, error } = useQuery<ServerListDashboardChartsResponse>({
    queryKey: [PUBLIC_SHARED_DASHBOARD, slug],
    queryFn: async () => {
      const response = await client.public.dashboardsDetail(slug);
      return response.data;
    },
  });

  // Sort charts based on their positions - moved before conditional returns
  const sortedCharts = useMemo(() => {
    if (!data?.charts || !data?.positions) return [];

    // Create a map of chart_id to position for faster lookup
    const positionMap = new Map<string, number>();
    data.positions.forEach((pos: MalakDashboardChartPosition) => {
      if (pos.chart_id) {
        positionMap.set(pos.chart_id, pos.order_index || 0);
      }
    });

    // Sort charts based on their positions
    return [...data.charts].sort((a, b) => {
      const posA = a.id ? positionMap.get(a.id) || 0 : 0;
      const posB = b.id ? positionMap.get(b.id) || 0 : 0;
      return posA - posB;
    });
  }, [data?.charts, data?.positions]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  if (error || !data?.dashboard) {
    notFound();
  }

  const dashboard = data.dashboard;

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <div className="border-b pb-4 mb-4">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">{dashboard.title}</h2>
          <p className="text-gray-600">{dashboard.description}</p>
        </div>

        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">Reference:</span>
            <span className="ml-2 text-gray-900">{dashboard.reference}</span>
          </div>
          <div>
            <span className="text-gray-500">Charts:</span>
            <span className="ml-2 text-gray-900">{dashboard.chart_count}</span>
          </div>
          <div>
            <span className="text-gray-500">Last Updated:</span>
            <span className="ml-2 text-gray-900">
              {dashboard.updated_at
                ? `${formatDistanceToNow(new Date(dashboard.updated_at))} ago`
                : 'Never'}
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {sortedCharts.map((dashboardChart: MalakDashboardChart) => (
          <ChartCard key={dashboardChart.id} chart={dashboardChart} dashboard_reference={slug} />
        ))}
      </div>

      <div className="bg-white shadow rounded-lg p-6">
        <p className="text-sm text-gray-500 text-center">
          This is a public view of the dashboard. For more detailed information, please contact the dashboard owner.
        </p>
      </div>
    </div>
  );
} 
