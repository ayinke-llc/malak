"use client";

import { Card } from "@/components/ui/card";
import { RiBarChart2Line, RiPieChartLine, RiSettings4Line, RiArrowDownSLine, RiLoader4Line } from "@remixicon/react";
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
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { useQuery } from "@tanstack/react-query";
import client from "@/lib/client";
import { LIST_CHARTS } from "@/lib/query-constants";
import type { ServerListIntegrationChartsResponse } from "@/client/Api";

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
  "dashboard_lLDM9lfnEk": {
    id: "dashboard_C5B_LEc_D7",
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
  { month: "Day 1", revenue: 2400 },
  { month: "Day 2", revenue: 1398 },
  { month: "Day 3", revenue: 9800 },
  { month: "Day 4", revenue: 3908 },
  { month: "Day 5", revenue: 4800 },
  { month: "Day 6", revenue: 3800 },
  { month: "Day 7", revenue: 5200 },
  { month: "Day 8", revenue: 4100 },
  { month: "Day 9", revenue: 6300 },
  { month: "Day 10", revenue: 5400 },
  { month: "Day 11", revenue: 4700 },
  { month: "Day 12", revenue: 3900 },
  { month: "Day 13", revenue: 5600 },
  { month: "Day 14", revenue: 4800 },
  { month: "Day 15", revenue: 6100 },
  { month: "Day 16", revenue: 5300 },
  { month: "Day 17", revenue: 4500 },
  { month: "Day 18", revenue: 3700 },
  { month: "Day 19", revenue: 5900 },
  { month: "Day 20", revenue: 4200 },
  { month: "Day 21", revenue: 6400 },
  { month: "Day 22", revenue: 5500 },
  { month: "Day 23", revenue: 4600 },
  { month: "Day 24", revenue: 3800 },
  { month: "Day 25", revenue: 5700 },
  { month: "Day 26", revenue: 4900 },
  { month: "Day 27", revenue: 6200 },
  { month: "Day 28", revenue: 5100 },
  { month: "Day 29", revenue: 4300 },
  { month: "Day 30", revenue: 3600 }
];

const userGrowthData = [
  { month: "Day 1", users: 1200 },
  { month: "Day 2", users: 1800 },
  { month: "Day 3", users: 2400 },
  { month: "Day 4", users: 3600 },
  { month: "Day 5", users: 4200 },
  { month: "Day 6", users: 5000 },
  { month: "Day 7", users: 5500 },
  { month: "Day 8", users: 4800 },
  { month: "Day 9", users: 6200 },
  { month: "Day 10", users: 5800 },
  { month: "Day 11", users: 4900 },
  { month: "Day 12", users: 4100 },
  { month: "Day 13", users: 5900 },
  { month: "Day 14", users: 5100 },
  { month: "Day 15", users: 6400 },
  { month: "Day 16", users: 5600 },
  { month: "Day 17", users: 4700 },
  { month: "Day 18", users: 3900 },
  { month: "Day 19", users: 6100 },
  { month: "Day 20", users: 4400 },
  { month: "Day 21", users: 6600 },
  { month: "Day 22", users: 5700 },
  { month: "Day 23", users: 4800 },
  { month: "Day 24", users: 4000 },
  { month: "Day 25", users: 5900 },
  { month: "Day 26", users: 5100 },
  { month: "Day 27", users: 6400 },
  { month: "Day 28", users: 5300 },
  { month: "Day 29", users: 4500 },
  { month: "Day 30", users: 3800 },
  { month: "Day 31", users: 6200 },
  { month: "Day 32", users: 5400 }
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

interface ChartOption {
  value: string;
  label: string;
  type: "bar" | "pie";
  description: string;
}

export default function DashboardPage() {
  const params = useParams();
  const dashboardId = params.slug as string;
  const dashboard = mockDashboards[dashboardId];
  const [isOpen, setIsOpen] = useState(false);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [selectedChart, setSelectedChart] = useState<string>("");
  const [selectedChartLabel, setSelectedChartLabel] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);

  const { data: chartsData, isLoading: isLoadingCharts } = useQuery<ServerListIntegrationChartsResponse>({
    queryKey: [LIST_CHARTS],
    queryFn: async () => {
      const response = await client.dashboards.chartsList();
      return response.data;
    },
    enabled: isPopoverOpen, // Only fetch when popover is open
  });

  const availableCharts: ChartOption[] = (chartsData?.charts || []).map(chart => ({
    value: chart.reference || "",
    label: chart.user_facing_name || "",
    type: chart.internal_name?.toLowerCase().includes("transaction") ? "bar" : "pie",
    description: chart.internal_name || "",
  }));

  const barCharts = availableCharts.filter(chart => chart.type === "bar");
  const pieCharts = availableCharts.filter(chart => chart.type === "pie");

  const handleAddChart = () => {
    if (!selectedChart) return;

    const chartToAdd = availableCharts.find(chart => chart.value === selectedChart);

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
          <SheetContent className="sm:max-w-xl">
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
                  <Popover open={isPopoverOpen} onOpenChange={setIsPopoverOpen}>
                    <PopoverTrigger asChild>
                      <Button
                        variant="outline"
                        role="combobox"
                        className="w-full justify-between"
                        aria-expanded={isPopoverOpen}
                      >
                        {selectedChart ? selectedChartLabel : "Select a chart..."}
                        <RiArrowDownSLine className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-[--radix-popover-trigger-width] p-0" align="start" side="bottom">
                      <Command className="w-full">
                        <CommandInput placeholder="Search charts..." />
                        <CommandList>
                          <CommandEmpty>No charts found.</CommandEmpty>
                          {isLoadingCharts ? (
                            <CommandItem disabled className="flex items-center gap-2 opacity-60">
                              <RiLoader4Line className="h-4 w-4 animate-spin" />
                              <span>Loading available charts...</span>
                            </CommandItem>
                          ) : (
                            <>
                              {barCharts.length > 0 && (
                                <CommandGroup heading="Bar Charts">
                                  {barCharts.map(chart => (
                                    <CommandItem
                                      key={chart.value}
                                      value={`${chart.label} ${chart.description}`}
                                      onSelect={(currentValue) => {
                                        setSelectedChart(chart.value);
                                        setSelectedChartLabel(chart.label);
                                        setIsPopoverOpen(false);
                                      }}
                                      className="flex items-center gap-2"
                                    >
                                      <RiBarChart2Line className="h-4 w-4" />
                                      <span>{chart.label}</span>
                                    </CommandItem>
                                  ))}
                                </CommandGroup>
                              )}
                              {pieCharts.length > 0 && (
                                <CommandGroup heading="Pie Charts">
                                  {pieCharts.map(chart => (
                                    <CommandItem
                                      key={chart.value}
                                      value={`${chart.label} ${chart.description}`}
                                      onSelect={(currentValue) => {
                                        setSelectedChart(chart.value);
                                        setSelectedChartLabel(chart.label);
                                        setIsPopoverOpen(false);
                                      }}
                                      className="flex items-center gap-2"
                                    >
                                      <RiPieChartLine className="h-4 w-4" />
                                      <span>{chart.label}</span>
                                    </CommandItem>
                                  ))}
                                </CommandGroup>
                              )}
                            </>
                          )}
                        </CommandList>
                      </Command>
                    </PopoverContent>
                  </Popover>
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
