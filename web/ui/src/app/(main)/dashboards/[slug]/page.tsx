"use client";

import { Card } from "@/components/ui/card";
import {
  RiBarChart2Line, RiPieChartLine, RiSettings4Line,
  RiArrowDownSLine, RiLoader4Line
} from "@remixicon/react";
import { useParams } from "next/navigation";
import {
  Bar, BarChart, ResponsiveContainer,
  XAxis, YAxis, Tooltip, PieChart,
  Pie, Cell
} from "recharts";
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
import { useQuery, useMutation } from "@tanstack/react-query";
import client from "@/lib/client";
import { LIST_CHARTS, DASHBOARD_DETAIL } from "@/lib/query-constants";
import type {
  ServerAPIStatus, ServerListIntegrationChartsResponse,
  ServerListDashboardChartsResponse, MalakDashboardChart
} from "@/client/Api";
import { toast } from "sonner";
import { AxiosError } from "axios";

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

// Mock data for pie charts
const costData = [
  { name: "Infrastructure", value: 400, color: "#0088FE" },
  { name: "Marketing", value: 300, color: "#00C49F" },
  { name: "Development", value: 500, color: "#FFBB28" },
  { name: "Operations", value: 200, color: "#FF8042" },
];

function ChartCard({ chart }: { chart: MalakDashboardChart }) {
  const getChartIcon = (type: string) => {
    switch (type) {
      case "bar":
        return <RiBarChart2Line className="h-4 w-4" />;
      case "pie":
        return <RiPieChartLine className="h-4 w-4" />;
      default:
        return <RiBarChart2Line className="h-4 w-4" />;
    }
  };

  const getChartData = (chart: MalakDashboardChart) => {
    // TODO: Replace with real data from the chart's data source
    return chart.chart?.chart_type === "bar" ? revenueData : costData;
  };

  if (!chart.chart) {
    return null;
  }

  return (
    <Card className="p-3">
      <div className="flex items-center justify-between mb-1">
        <div className="flex items-center gap-2">
          <div className="text-muted-foreground">
            {getChartIcon(chart.chart.chart_type || "bar")}
          </div>
          <div>
            <h3 className="text-sm font-medium">{chart.chart.user_facing_name}</h3>
            <p className="text-xs text-muted-foreground">{chart.chart.internal_name}</p>
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
        {chart.chart.chart_type === "bar" ? (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={160}>
              <BarChart data={getChartData(chart)} margin={{ top: 5, right: 5, left: -15, bottom: 0 }}>
                <XAxis dataKey="month" stroke="#888888" fontSize={11} />
                <YAxis stroke="#888888" fontSize={11} />
                <Tooltip />
                <Bar dataKey="revenue" fill="#8884d8" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </ChartContainer>
        ) : (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={160}>
              <PieChart margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                <Pie
                  data={getChartData(chart)}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={60}
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
        )}
      </div>
    </Card>
  );
}

export default function DashboardPage() {
  const params = useParams();
  const dashboardId = params.slug as string;

  const [isOpen, setIsOpen] = useState(false);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [selectedChart, setSelectedChart] = useState<string>("");
  const [selectedChartLabel, setSelectedChartLabel] = useState<string>("");

  const { data: dashboardData, isLoading: isLoadingDashboard } = useQuery<ServerListDashboardChartsResponse>({
    queryKey: [DASHBOARD_DETAIL, dashboardId],
    queryFn: async () => {
      const response = await client.dashboards.dashboardsDetail(dashboardId);
      return response.data;
    },
  });

  const { data: chartsData, isLoading: isLoadingCharts } = useQuery<ServerListIntegrationChartsResponse>({
    queryKey: [LIST_CHARTS],
    queryFn: async () => {
      const response = await client.dashboards.chartsList();
      return response.data;
    },
    enabled: isPopoverOpen,
  });

  const addChartMutation = useMutation({
    mutationFn: async (chartReference: string) => {
      const response = await client.dashboards.chartsUpdate(dashboardId, {
        chart_reference: chartReference
      });
      return response.data;
    },
    onSuccess: (data) => {
      setSelectedChart("");
      setSelectedChartLabel("");
      setIsOpen(false);
      toast.success(data.message);
    },
    onError: (err: AxiosError<ServerAPIStatus>): void => {
      toast.error(err?.response?.data?.message || "Failed to add chart to dashboard");
    }
  });

  const barCharts = chartsData?.charts?.filter(chart => chart.chart_type === "bar") ?? [];
  const pieCharts = chartsData?.charts?.filter(chart => chart.chart_type === "pie") ?? [];

  const handleAddChart = () => {
    if (!selectedChart) {
      toast.warning("Select a chart before adding to dashboard")
      return
    };

    addChartMutation.mutate(selectedChart);
  };

  if (isLoadingDashboard) {
    return (
      <div className="flex items-center justify-center h-[50vh]">
        <div className="text-center">
          <RiLoader4Line className="h-8 w-8 animate-spin mx-auto text-muted-foreground" />
          <h1 className="text-2xl font-bold text-muted-foreground mt-4">Loading dashboard...</h1>
        </div>
      </div>
    );
  }

  if (!dashboardData?.dashboard) {
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
          <h1 className="text-2xl font-bold">{dashboardData.dashboard.title}</h1>
          <p className="text-muted-foreground">{dashboardData.dashboard.description}</p>
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
                                      key={chart.reference}
                                      value={`${chart.user_facing_name} ${chart.internal_name}`}
                                      onSelect={() => {
                                        setSelectedChart(chart.reference || "");
                                        setSelectedChartLabel(chart.user_facing_name || "");
                                        setIsPopoverOpen(false);
                                      }}
                                      className="flex items-center gap-2"
                                    >
                                      <RiBarChart2Line className="h-4 w-4" />
                                      <span>{chart.user_facing_name}</span>
                                    </CommandItem>
                                  ))}
                                </CommandGroup>
                              )}
                              {pieCharts.length > 0 && (
                                <CommandGroup heading="Pie Charts">
                                  {pieCharts.map(chart => (
                                    <CommandItem
                                      key={chart.reference}
                                      value={`${chart.user_facing_name} ${chart.internal_name}`}
                                      onSelect={() => {
                                        setSelectedChart(chart.reference || "");
                                        setSelectedChartLabel(chart.user_facing_name || "");
                                        setIsPopoverOpen(false);
                                      }}
                                      className="flex items-center gap-2"
                                    >
                                      <RiPieChartLine className="h-4 w-4" />
                                      <span>{chart.user_facing_name}</span>
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
                disabled={!selectedChart || addChartMutation.isPending}
              >
                {addChartMutation.isPending ? (
                  <div className="flex items-center gap-2">
                    <RiLoader4Line className="h-4 w-4 animate-spin" />
                    Adding Chart...
                  </div>
                ) : (
                  "Add to Dashboard"
                )}
              </Button>
            </SheetFooter>
          </SheetContent>
        </Sheet>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {!dashboardData?.charts || dashboardData.charts.length === 0 ? (
          <div className="col-span-full flex flex-col items-center justify-center py-12 text-center">
            <RiBarChart2Line className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-medium">No charts yet</h3>
            <p className="text-sm text-muted-foreground mt-1 mb-4">Get started by adding your first chart to this dashboard.</p>
            <Button onClick={() => setIsOpen(true)}>Add Your First Chart</Button>
          </div>
        ) : (
          dashboardData.charts.map((chart) => (
            <ChartCard key={chart.reference} chart={chart} />
          ))
        )}
      </div>
    </div>
  );
}
