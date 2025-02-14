"use client";

import { Card } from "@/components/ui/card";
import { RiBarChart2Line, RiPieChartLine, RiSettings4Line } from "@remixicon/react";
import { useParams } from "next/navigation";
import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis, Tooltip, PieChart, Pie, Cell } from "recharts";
import { ChartContainer, ChartTooltip } from "@/components/ui/chart";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

// Mock type for demonstration
type Chart = {
  id: string;
  title: string;
  type: "bar" | "line" | "pie";
  description: string;
};

type Dashboard = {
  id: string;
  title: string;
  description: string;
};

type DashboardsMap = {
  [key: string]: Dashboard;
};

// Mock dashboard data
const mockDashboards: DashboardsMap = {
  "1": {
    id: "1",
    title: "Revenue Overview",
    description: "Monthly revenue trends and projections",
  },
  "2": {
    id: "2",
    title: "User Analytics",
    description: "User engagement and activity metrics",
  },
  "3": {
    id: "3",
    title: "Cost Distribution",
    description: "Breakdown of operational costs",
  }
};

// Mock data for demonstration
const mockCharts: Chart[] = [
  {
    id: "1",
    title: "Monthly Revenue",
    type: "bar",
    description: "Revenue trends over the past 12 months"
  },
  {
    id: "2",
    title: "User Growth",
    type: "bar",
    description: "Daily active users growth"
  },
  {
    id: "3",
    title: "Cost Breakdown",
    type: "pie",
    description: "Distribution of operational costs"
  },
  {
    id: "4",
    title: "Conversion Rate",
    type: "bar",
    description: "User conversion metrics"
  }
];

// Mock data for bar charts
const revenueData = [
  { month: "Jan", revenue: 2400 },
  { month: "Feb", revenue: 1398 },
  { month: "Mar", revenue: 9800 },
  { month: "Apr", revenue: 3908 },
  { month: "May", revenue: 4800 },
  { month: "Jun", revenue: 3800 },
];

const userGrowthData = [
  { month: "Jan", users: 1200 },
  { month: "Feb", users: 1800 },
  { month: "Mar", users: 2400 },
  { month: "Apr", users: 3600 },
  { month: "May", users: 4200 },
  { month: "Jun", users: 5000 },
];

const conversionData = [
  { month: "Jan", rate: 45 },
  { month: "Feb", rate: 52 },
  { month: "Mar", rate: 48 },
  { month: "Apr", rate: 61 },
  { month: "May", rate: 55 },
  { month: "Jun", rate: 67 },
];

// Mock data for pie chart
const costData = [
  { name: "Infrastructure", value: 400, color: "#0088FE" },
  { name: "Marketing", value: 300, color: "#00C49F" },
  { name: "Development", value: 500, color: "#FFBB28" },
  { name: "Operations", value: 200, color: "#FF8042" },
];

function ChartCard({ chart }: { chart: Chart }) {
  const getChartIcon = (type: Chart["type"]) => {
    switch (type) {
      case "bar":
        return <RiBarChart2Line className="h-4 w-4" />;
      case "pie":
        return <RiPieChartLine className="h-4 w-4" />;
    }
  };

  const getBarData = (chartId: string) => {
    switch (chartId) {
      case "1":
        return { data: revenueData, key: "revenue" };
      case "2":
        return { data: userGrowthData, key: "users" };
      case "4":
        return { data: conversionData, key: "rate" };
      default:
        return { data: revenueData, key: "revenue" };
    }
  };

  const renderChart = (type: Chart["type"]) => {
    switch (type) {
      case "bar":
        const { data, key } = getBarData(chart.id);
        return (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={200}>
              <BarChart data={data} margin={{ top: 10, right: 10, left: -15, bottom: 0 }}>
                <XAxis dataKey="month" stroke="#888888" fontSize={12} />
                <YAxis stroke="#888888" fontSize={12} />
                <Tooltip />
                <Bar dataKey={key} fill="#8884d8" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </ChartContainer>
        );
      case "pie":
        return (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={200}>
              <PieChart margin={{ top: 10, right: 10, left: 10, bottom: 10 }}>
                <Pie
                  data={costData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  dataKey="value"
                >
                  {costData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </ChartContainer>
        );
    }
  };

  return (
    <Card className="p-4">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <div className="text-muted-foreground">
            {getChartIcon(chart.type)}
          </div>
          <div>
            <h3 className="text-sm font-medium">{chart.title}</h3>
            <p className="text-xs text-muted-foreground">{chart.description}</p>
          </div>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="p-2 hover:bg-muted rounded-md">
              <RiSettings4Line className="h-4 w-4 text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem>Edit Chart</DropdownMenuItem>
            <DropdownMenuItem>Duplicate</DropdownMenuItem>
            <DropdownMenuItem className="text-destructive">Delete</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="w-full">
        {renderChart(chart.type)}
      </div>
    </Card>
  );
}

export default function DashboardPage() {
  const params = useParams();
  const dashboardId = params.slug as string;
  const dashboard = mockDashboards[dashboardId];

  if (!dashboard) {
    return (
      <div className="flex items-center justify-center h-[50vh]">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-muted-foreground">Dashboard not found</h1>
          <p className="text-muted-foreground mt-2">The dashboard you're looking for doesn't exist.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{dashboard.title}</h1>
          <p className="text-muted-foreground">{dashboard.description}</p>
        </div>
        <div>
          <button className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90">
            Add Chart
          </button>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {mockCharts.map((chart) => (
          <ChartCard key={chart.id} chart={chart} />
        ))}
      </div>
    </div>
  );
}
