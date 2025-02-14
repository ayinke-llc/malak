"use client";

import { Card } from "@/components/ui/card";
import { RiBarChart2Line, RiPieChartLine, RiSettings4Line } from "@remixicon/react";
import { useParams } from "next/navigation";
import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis, Tooltip, PieChart, Pie, Cell } from "recharts";
import { ChartContainer } from "@/components/ui/chart";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
  SheetFooter,
} from "@/components/ui/sheet";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

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
    title: "Cost Distribution",
    type: "pie",
    description: "Distribution of operational costs"
  },
  {
    id: "4",
    title: "Conversion Rate",
    type: "bar",
    description: "User conversion metrics"
  },
  {
    id: "5",
    title: "Team Distribution",
    type: "pie",
    description: "Team members by department"
  },
  {
    id: "6",
    title: "Storage Usage",
    type: "pie",
    description: "Storage allocation by type"
  },
  {
    id: "7",
    title: "Support Tickets",
    type: "bar",
    description: "Monthly support ticket volume"
  },
  {
    id: "8",
    title: "API Requests",
    type: "bar",
    description: "Daily API request count"
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

const ticketData = [
  { month: "Jan", tickets: 145 },
  { month: "Feb", tickets: 132 },
  { month: "Mar", tickets: 164 },
  { month: "Apr", tickets: 128 },
  { month: "May", tickets: 155 },
  { month: "Jun", tickets: 147 },
];

const apiRequestData = [
  { month: "Jan", requests: 25000 },
  { month: "Feb", requests: 32000 },
  { month: "Mar", requests: 38000 },
  { month: "Apr", requests: 42000 },
  { month: "May", requests: 45000 },
  { month: "Jun", requests: 51000 },
];

// Mock data for pie charts
const costData = [
  { name: "Infrastructure", value: 400, color: "#0088FE" },
  { name: "Marketing", value: 300, color: "#00C49F" },
  { name: "Development", value: 500, color: "#FFBB28" },
  { name: "Operations", value: 200, color: "#FF8042" },
];

const teamData = [
  { name: "Engineering", value: 40, color: "#0088FE" },
  { name: "Product", value: 15, color: "#00C49F" },
  { name: "Marketing", value: 20, color: "#FFBB28" },
  { name: "Sales", value: 25, color: "#FF8042" },
];

const storageData = [
  { name: "Documents", value: 450, color: "#0088FE" },
  { name: "Media", value: 800, color: "#00C49F" },
  { name: "Backups", value: 300, color: "#FFBB28" },
  { name: "Other", value: 150, color: "#FF8042" },
];

// Rename from availableCharts to chartOptions
const chartOptions = [
  {
    title: "Bar Charts",
    items: [
      {
        id: "revenue",
        title: "Revenue Chart",
        description: "Track revenue over time",
        type: "bar",
        icon: <RiBarChart2Line className="h-4 w-4" />,
      },
      {
        id: "users",
        title: "User Growth",
        description: "Monitor user growth trends",
        type: "bar",
        icon: <RiBarChart2Line className="h-4 w-4" />,
      },
      {
        id: "conversion",
        title: "Conversion Rate",
        description: "Track conversion metrics",
        type: "bar",
        icon: <RiBarChart2Line className="h-4 w-4" />,
      },
    ],
  },
  {
    title: "Pie Charts",
    items: [
      {
        id: "distribution",
        title: "Cost Distribution",
        description: "Analyze cost breakdown",
        type: "pie",
        icon: <RiPieChartLine className="h-4 w-4" />,
      },
      {
        id: "team",
        title: "Team Distribution",
        description: "View team composition",
        type: "pie",
        icon: <RiPieChartLine className="h-4 w-4" />,
      },
    ],
  },
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
      case "7":
        return { data: ticketData, key: "tickets" };
      case "8":
        return { data: apiRequestData, key: "requests" };
      default:
        return { data: revenueData, key: "revenue" };
    }
  };

  const getPieData = (chartId: string) => {
    switch (chartId) {
      case "3":
        return costData;
      case "5":
        return teamData;
      case "6":
        return storageData;
      default:
        return costData;
    }
  };

  const renderChart = (type: Chart["type"]) => {
    switch (type) {
      case "bar":
        const { data, key } = getBarData(chart.id);
        return (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={160}>
              <BarChart data={data} margin={{ top: 5, right: 5, left: -15, bottom: 0 }}>
                <XAxis dataKey="month" stroke="#888888" fontSize={11} />
                <YAxis stroke="#888888" fontSize={11} />
                <Tooltip />
                <Bar dataKey={key} fill="#8884d8" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </ChartContainer>
        );
      case "pie":
        const pieData = getPieData(chart.id);
        return (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={160}>
              <PieChart margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={60}
                  dataKey="value"
                >
                  {pieData.map((entry, index) => (
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
    <Card className="p-3">
      <div className="flex items-center justify-between mb-1">
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
            <button className="p-1.5 hover:bg-muted rounded-md">
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
  const [isOpen, setIsOpen] = useState(false);
  const [selectedChart, setSelectedChart] = useState<string>("");
  const [selectedChartLabel, setSelectedChartLabel] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);
  const [availableCharts, setAvailableCharts] = useState<Array<{
    value: string;
    label: string;
    type: "bar" | "pie";
  }>>([]);

  // Simulating API call to fetch available charts
  useEffect(() => {
    const fetchCharts = async () => {
      setIsLoading(true);
      try {
        // Simulated API response
        const data: Array<{ value: string; label: string; type: "bar" | "pie" }> = [
          { value: "revenue", label: "Revenue Chart", type: "bar" },
          { value: "users", label: "User Growth", type: "bar" },
          { value: "conversion", label: "Conversion Rate", type: "bar" },
          { value: "distribution", label: "Cost Distribution", type: "pie" },
          { value: "team", label: "Team Distribution", type: "pie" },
        ];

        setAvailableCharts(data);
      } catch (error) {
        console.error("Failed to fetch charts:", error);
        // Initialize with empty array on error
        setAvailableCharts([]);
      } finally {
        setIsLoading(false);
      }
    };

    if (isOpen) {
      fetchCharts();
    }
  }, [isOpen]);

  const handleAddChart = () => {
    if (!selectedChart) return;

    // Here you would typically add the chart to your dashboard
    const chartToAdd = availableCharts.find(chart => chart.value === selectedChart);
    console.log("Adding chart:", chartToAdd);

    setSelectedChart("");
    setSelectedChartLabel("");
    setIsOpen(false);
  };

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
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{dashboard.title}</h1>
          <p className="text-muted-foreground">{dashboard.description}</p>
        </div>
        <Sheet open={isOpen} onOpenChange={setIsOpen}>
          <SheetTrigger asChild>
            <Button>Add Chart</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>Add Chart</SheetTitle>
              <SheetDescription>
                Select a chart to add to your dashboard
              </SheetDescription>
            </SheetHeader>
            <div className="mt-6">
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label>Chart Type</Label>
                  <Select
                    value={selectedChart}
                    onValueChange={(value) => {
                      setSelectedChart(value);
                      const chart = availableCharts.find(c => c.value === value);
                      if (chart) {
                        setSelectedChartLabel(chart.label);
                      }
                    }}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a chart" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectLabel>Bar Charts</SelectLabel>
                        {availableCharts
                          .filter(chart => chart.type === "bar")
                          .map(chart => (
                            <SelectItem
                              key={chart.value}
                              value={chart.value}
                              className="flex items-center gap-2"
                            >
                              <div className="flex items-center gap-2">
                                <RiBarChart2Line className="h-4 w-4" />
                                <span>{chart.label}</span>
                              </div>
                            </SelectItem>
                          ))}
                      </SelectGroup>
                      <SelectGroup>
                        <SelectLabel>Pie Charts</SelectLabel>
                        {availableCharts
                          .filter(chart => chart.type === "pie")
                          .map(chart => (
                            <SelectItem
                              key={chart.value}
                              value={chart.value}
                              className="flex items-center gap-2"
                            >
                              <div className="flex items-center gap-2">
                                <RiPieChartLine className="h-4 w-4" />
                                <span>{chart.label}</span>
                              </div>
                            </SelectItem>
                          ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                  {selectedChart && (
                    <p className="text-sm text-muted-foreground">
                      Selected: {selectedChartLabel}
                    </p>
                  )}
                </div>
              </div>
            </div>
            <SheetFooter className="mt-4">
              <Button
                onClick={handleAddChart}
                disabled={!selectedChart || isLoading}
              >
                Add to Dashboard
              </Button>
            </SheetFooter>
          </SheetContent>
        </Sheet>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {mockCharts.map((chart) => (
          <ChartCard key={chart.id} chart={chart} />
        ))}
      </div>
    </div>
  );
}
