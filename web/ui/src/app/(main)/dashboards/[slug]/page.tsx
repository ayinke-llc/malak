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
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import client from "@/lib/client";
import { LIST_CHARTS, DASHBOARD_DETAIL, FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import type {
  ServerAPIStatus, ServerListIntegrationChartsResponse,
  ServerListDashboardChartsResponse, MalakDashboardChart, MalakIntegrationDataPoint
} from "@/client/Api";
import { toast } from "sonner";
import { AxiosError } from "axios";

function ChartCard({ chart }: { chart: MalakDashboardChart }) {
  const { data: chartData, isLoading: isLoadingChartData, error } = useQuery({
    queryKey: [FETCH_CHART_DATA_POINTS, chart.chart?.reference],
    queryFn: async () => {
      if (!chart.chart?.reference) return null;
      const response = await client.dashboards.chartsDetail(chart.chart.reference);
      return response.data;
    },
    enabled: !!chart.chart?.reference,
  });

  const getChartIcon = (type: string | undefined) => {
    switch (type) {
      case "bar":
        return <RiBarChart2Line className="h-4 w-4" />;
      case "pie":
        return <RiPieChartLine className="h-4 w-4" />;
      default:
        return <RiBarChart2Line className="h-4 w-4" />;
    }
  };

  const formatChartData = (dataPoints: MalakIntegrationDataPoint[] | undefined): Array<{
    name: string;
    value: number;
  }> => {
    if (!dataPoints) return [];

    // Transform data points into the format expected by recharts
    return dataPoints.map(point => {
      const value = point.data_point_type === 'currency'
        ? (point.point_value || 0) / 100
        : point.point_value || 0;

      return {
        name: point.point_name || '',
        value,
      };
    });
  };

  if (!chart.chart) {
    return null;
  }

  const formattedData = formatChartData(chartData?.data_points);

  if (isLoadingChartData) {
    return (
      <Card className="p-3">
        <div className="flex items-center justify-center h-[160px]">
          <RiLoader4Line className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      </Card>
    );
  }

  if (error) {
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
    <Card className="p-3">
      <div className="flex items-center justify-between mb-1">
        <div className="flex items-center gap-2">
          <div className="text-muted-foreground">
            {getChartIcon(chart.chart.chart_type)}
          </div>
          <div>
            <h3 className="text-sm font-bold">{chart.chart.user_facing_name}</h3>
          </div>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="p-1.5 hover:bg-muted rounded-md">
              <RiSettings4Line className="h-4 w-4 text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem className="text-destructive cursor-pointer">Remove from dashboard</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
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
            <ResponsiveContainer width="100%" height={160}>
              <BarChart data={formattedData} margin={{ top: 5, right: 5, left: -15, bottom: 0 }}>
                <XAxis dataKey="name" stroke="#888888" fontSize={11} />
                <YAxis stroke="#888888" fontSize={11} />
                <Tooltip formatter={(value: number) => {
                  if (chartData?.data_points?.[0]?.data_point_type === 'currency') {
                    return [`$${value.toFixed(2)}`, 'Value'];
                  }
                  return [value, 'Value'];
                }} />
                <Bar dataKey="value" fill="#8884d8" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </ChartContainer>
        ) : (
          <ChartContainer className="w-full h-full" config={{}}>
            <ResponsiveContainer width="100%" height={160}>
              <PieChart margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                <Pie
                  data={formattedData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={60}
                  dataKey="value"
                >
                  {formattedData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={`hsl(${index * 45}, 70%, 50%)`} />
                  ))}
                </Pie>
                <Tooltip formatter={(value: number) => {
                  if (chartData?.data_points?.[0]?.data_point_type === 'currency') {
                    return [`$${value.toFixed(2)}`, 'Value'];
                  }
                  return [value, 'Value'];
                }} />
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
  const dashboardID = params.slug as string;

  const queryClient = useQueryClient();

  const [isOpen, setIsOpen] = useState(false);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [selectedChart, setSelectedChart] = useState<string>("");
  const [selectedChartLabel, setSelectedChartLabel] = useState<string>("");

  const { data: dashboardData, isLoading: isLoadingDashboard } = useQuery<ServerListDashboardChartsResponse>({
    queryKey: [DASHBOARD_DETAIL, dashboardID],
    queryFn: async () => {
      const response = await client.dashboards.dashboardsDetail(dashboardID);
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
      const response = await client.dashboards.chartsUpdate(dashboardID, {
        chart_reference: chartReference
      });
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: [DASHBOARD_DETAIL, dashboardID] });
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
