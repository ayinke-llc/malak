"use client";

import type {
  MalakDashboardChart, ServerAPIStatus, ServerListDashboardChartsResponse,
  ServerListIntegrationChartsResponse
} from "@/client/Api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { ShareDialog } from "@/components/ui/dashboard/share-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { formatChartData, formatTooltipValue, getChartColors } from "@/lib/chart-utils";
import client from "@/lib/client";
import {
  ADD_CHART_DASHBOARD,
  DASHBOARD_DETAIL,
  FETCH_CHART_DATA_POINTS,
  LIST_CHARTS,
  REMOVE_CHART_DASHBOARD
} from "@/lib/query-constants";
import {
  closestCenter,
  DndContext,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  arrayMove,
  rectSortingStrategy,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import {
  RiArrowDownSLine,
  RiBarChart2Line, RiLoader4Line, RiPieChartLine,
  RiSettings4Line
} from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { useCallback, useEffect, useRef, useState } from "react";
import {
  Bar, BarChart,
  Cell,
  Pie,
  PieChart,
  Tooltip,
  XAxis,
  YAxis
} from "recharts";
import { toast } from "sonner";
import styles from "./styles.module.css";

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
                  formatTooltipValue(value, chartData?.data_points?.[0]?.data_point_type)
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
                  formatTooltipValue(value, chartData?.data_points?.[0]?.data_point_type)
                }
              />
            </PieChart>
          </ChartContainer>
        )}
      </div>
    </Card>
  );
}

function SortableChartCard({ chart, onRemove }: { chart: MalakDashboardChart; onRemove: (chartRef: string) => void }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: chart.reference || '' });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    zIndex: isDragging ? 50 : undefined,
    position: 'relative' as const,
    opacity: isDragging ? 0.5 : undefined,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={`touch-none ${styles.sortableChart} ${isDragging ? styles.sortableChartDragging : ''} cursor-grab active:cursor-grabbing`}
    >
      <div className="absolute top-2 right-2 z-50">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="p-1.5 hover:bg-muted rounded-md" onClick={(e) => e.stopPropagation()}>
              <RiSettings4Line className="h-4 w-4 text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem
              className="text-destructive cursor-pointer"
              onClick={() => onRemove(chart.chart?.reference || '')}
            >
              Remove from dashboard
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <ChartCard chart={chart} />
    </div>
  );
}

export default function DashboardDetailsPage({ reference }: { reference: string }) {
  const queryClient = useQueryClient();
  const dashboardID = reference;

  const [isOpen, setIsOpen] = useState(false);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [selectedChart, setSelectedChart] = useState<string>("");
  const [selectedChartLabel, setSelectedChartLabel] = useState<string>("");
  const [charts, setCharts] = useState<MalakDashboardChart[]>([]);
  const isUpdatingRef = useRef(false);

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

  useEffect(() => {
    if (dashboardData?.charts) {
      // Sort charts based on their positions if available
      const sortedCharts = [...dashboardData.charts].sort((a, b) => {
        const posA = dashboardData.positions?.find(p => p.chart_id === a.id)?.order_index ?? 0;
        const posB = dashboardData.positions?.find(p => p.chart_id === b.id)?.order_index ?? 0;
        return posA - posB;
      });
      setCharts(sortedCharts);
    }
  }, [dashboardData?.charts, dashboardData?.positions]);

  const updatePositionsMutation = useMutation({
    mutationFn: async (positions: { chart_id: string; index: number }[]) => {
      const response = await client.dashboards.positionsCreate(dashboardID, { positions });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DASHBOARD_DETAIL, dashboardID] });
      toast.success("Chart positions updated");
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err?.response?.data?.message || "Failed to update chart positions");
      if (dashboardData?.charts) {
        setCharts(dashboardData.charts);
      }
    }
  });

  const updatePositionsDebounced = useCallback(
    (positions: { chart_id: string; index: number }[]) => {
      if (isUpdatingRef.current) return;
      isUpdatingRef.current = true;

      setTimeout(() => {
        updatePositionsMutation.mutate(positions);
        isUpdatingRef.current = false;
      }, 100);
    },
    [updatePositionsMutation]
  );

  const addChartMutation = useMutation({
    mutationKey: [ADD_CHART_DASHBOARD],
    mutationFn: async (chartReference: string) => {
      const response = await client.dashboards.chartsUpdate(dashboardID, {
        chart_reference: chartReference
      });
      return { data: response.data, chartReference };
    },
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: [DASHBOARD_DETAIL, dashboardID] });
      setSelectedChart("");
      setSelectedChartLabel("");
      setIsOpen(false);
      toast.success(result.data.message);
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err?.response?.data?.message || "Failed to add chart to dashboard");
    }
  });

  const deleteChartMutation = useMutation({
    mutationKey: [REMOVE_CHART_DASHBOARD],
    mutationFn: async (chartReference: string) => {
      const response = await client.dashboards.chartsDelete(dashboardID, {
        chart_reference: chartReference
      });
      return { data: response.data, chartReference };
    },
    onSuccess: (result) => {
      setCharts(prevCharts => prevCharts.filter(chart => chart.chart?.reference !== result.chartReference));
      queryClient.invalidateQueries({ queryKey: [DASHBOARD_DETAIL, dashboardID] });
      toast.success(result.data.message);
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err?.response?.data?.message || "Failed to remove chart from dashboard");
    }
  });

  const handleRemoveChart = (chartReference: string) => {
    if (!chartReference) return;
    deleteChartMutation.mutate(chartReference);
  };

  const handleAddChart = () => {
    if (!selectedChart) {
      toast.warning("Select a chart before adding to dashboard");
      return;
    }
    addChartMutation.mutate(selectedChart);
  };

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  function handleDragEnd(event: any) {
    const { active, over } = event;

    if (!active || !over || active.id === over.id) {
      return;
    }

    setCharts((items) => {
      const oldIndex = items.findIndex((item) => item.reference === active.id);
      const newIndex = items.findIndex((item) => item.reference === over.id);
      const newItems = arrayMove(items, oldIndex, newIndex);

      const positions = newItems.map((item, index) => ({
        chart_id: item.id || '',
        index
      }));

      updatePositionsDebounced(positions);
      return newItems;
    });
  }

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
          <p className="text-muted-foreground mt-2">The dashboard you&apos;re looking for doesn&apos;t exist.</p>
        </div>
      </div>
    );
  }

  const dashboard = dashboardData.dashboard;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{dashboard.title}</h1>
          <p className="text-muted-foreground">{dashboard.description}</p>
        </div>
        <div className="flex items-center gap-2">
          <ShareDialog
            title={dashboard.title || "Dashboard"}
            reference={dashboardID}
            token={dashboardData?.link?.token as string}
          />
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
                        <Command>
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
                                {(chartsData?.charts || []).filter(chart =>
                                  chart.chart_type === "bar" &&
                                  !charts.some(dashboardChart => dashboardChart.chart?.reference === chart.reference)
                                ).length > 0 && (
                                    <CommandGroup heading="Bar Charts">
                                      {(chartsData?.charts || [])
                                        .filter(chart =>
                                          chart.chart_type === "bar" &&
                                          !charts.some(dashboardChart => dashboardChart.chart?.reference === chart.reference)
                                        )
                                        .map(chart => (
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
                                {(chartsData?.charts || []).filter(chart =>
                                  chart.chart_type === "pie" &&
                                  !charts.some(dashboardChart => dashboardChart.chart?.reference === chart.reference)
                                ).length > 0 && (
                                    <CommandGroup heading="Pie Charts">
                                      {(chartsData?.charts || [])
                                        .filter(chart =>
                                          chart.chart_type === "pie" &&
                                          !charts.some(dashboardChart => dashboardChart.chart?.reference === chart.reference)
                                        )
                                        .map(chart => (
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
      </div>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={handleDragEnd}
        onDragStart={() => {
          if (window.navigator.vibrate) {
            window.navigator.vibrate(100);
          }
        }}
      >
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {!charts || charts.length === 0 ? (
            <div className="col-span-full flex flex-col items-center justify-center py-12 text-center">
              <RiBarChart2Line className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium">No charts yet</h3>
              <p className="text-sm text-muted-foreground mt-1 mb-4">Get started by adding your first chart to this dashboard.</p>
              <Button onClick={() => setIsOpen(true)}>Add Your First Chart</Button>
            </div>
          ) : (
            <SortableContext
              items={charts.map(chart => chart.reference || '')}
              strategy={rectSortingStrategy}
            >
              {charts.map((chart) => (
                <SortableChartCard
                  key={chart.reference}
                  chart={chart}
                  onRemove={handleRemoveChart}
                />
              ))}
            </SortableContext>
          )}
        </div>
      </DndContext>
    </div>
  );
} 
