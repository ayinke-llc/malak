"use client";

import { notFound } from 'next/navigation'
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip, XAxis, YAxis } from 'recharts';

// Mock data to simulate a public dashboard
const mockDashboard = {
  id: 'mock-dashboard-1',
  title: 'Sample Public Dashboard',
  description: 'This is a sample shared dashboard showing public metrics',
  lastUpdated: '2024-03-20T10:00:00Z',
  status: 'Active',
  owner: 'Demo Team',
  charts: [
    {
      id: 'chart-1',
      type: 'bar',
      title: 'Monthly Active Users',
      description: 'Public user activity trends',
      data: [
        { name: 'Jan', value: 4000 },
        { name: 'Feb', value: 3000 },
        { name: 'Mar', value: 2000 },
        { name: 'Apr', value: 2780 },
        { name: 'May', value: 1890 },
        { name: 'Jun', value: 2390 },
      ]
    },
    {
      id: 'chart-2',
      type: 'pie',
      title: 'Usage Distribution',
      description: 'Distribution of platform usage',
      data: [
        { name: 'Mobile', value: 400 },
        { name: 'Desktop', value: 300 },
        { name: 'Tablet', value: 200 },
      ]
    }
  ],
  publicMetrics: [
    {
      id: 'metric-1',
      name: 'Uptime',
      value: '99.9%',
      label: 'Last 30 days'
    },
    {
      id: 'metric-2',
      name: 'Status',
      value: 'Operational',
      label: 'All systems normal'
    }
  ]
}

function ChartCard({ chart }: { chart: typeof mockDashboard.charts[0] }) {
  const getChartColors = (index: number) => {
    const colors = ['#3B82F6', '#10B981', '#6366F1', '#F59E0B', '#EF4444'];
    return colors[index % colors.length];
  };

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <div className="mb-4">
        <h3 className="text-lg font-medium text-gray-900">{chart.title}</h3>
        <p className="text-sm text-gray-500">{chart.description}</p>
      </div>
      
      <div className="h-[200px] w-full">
        {chart.type === 'bar' ? (
          <BarChart
            width={390}
            height={200}
            data={chart.data}
            margin={{ top: 5, right: 5, left: -15, bottom: 5 }}
          >
            <XAxis dataKey="name" stroke="#888888" fontSize={12} />
            <YAxis stroke="#888888" fontSize={12} />
            <Tooltip />
            <Bar dataKey="value" fill="#3B82F6" radius={[4, 4, 0, 0]} />
          </BarChart>
        ) : (
          <PieChart width={390} height={200}>
            <Pie
              data={chart.data}
              cx="50%"
              cy="50%"
              labelLine={false}
              label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
              outerRadius={80}
              dataKey="value"
            >
              {chart.data.map((_entry, index) => (
                <Cell key={`cell-${index}`} fill={getChartColors(index)} />
              ))}
            </Pie>
            <Tooltip />
          </PieChart>
        )}
      </div>
    </div>
  );
}

export default function SharedDashboardPage({
  params,
}: {
  params: { slug: string }
}) {
  // Simulate fetching dashboard data
  const dashboard = params.slug === mockDashboard.id ? mockDashboard : null;

  if (!dashboard) {
    notFound();
  }

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <div className="border-b pb-4 mb-4">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">{dashboard.title}</h2>
          <p className="text-gray-600">{dashboard.description}</p>
        </div>
        
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">Owner:</span>
            <span className="ml-2 text-gray-900">{dashboard.owner}</span>
          </div>
          <div>
            <span className="text-gray-500">Status:</span>
            <span className="ml-2 text-gray-900">{dashboard.status}</span>
          </div>
          <div>
            <span className="text-gray-500">Last Updated:</span>
            <span className="ml-2 text-gray-900">
              {new Date(dashboard.lastUpdated).toLocaleDateString()}
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {dashboard.publicMetrics.map((metric) => (
          <div key={metric.id} className="bg-white shadow rounded-lg p-6">
            <h3 className="text-sm font-medium text-gray-500">{metric.name}</h3>
            <div className="mt-2">
              <p className="text-2xl font-semibold text-gray-900">{metric.value}</p>
              <p className="text-sm text-gray-500 mt-1">{metric.label}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {dashboard.charts.map((chart) => (
          <ChartCard key={chart.id} chart={chart} />
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