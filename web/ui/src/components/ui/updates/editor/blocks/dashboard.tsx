import type { MalakIntegrationChart, ServerListDashboardResponse } from "@/client/Api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import client from "@/lib/client";
import { FETCH_CHART_DATA_POINTS, LIST_DASHBOARDS } from "@/lib/query-constants";
import { cn } from "@/lib/utils";
import { formatChartData, formatTooltipValue, getChartColors } from "@/lib/chart-utils";
import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import {
  RiArrowDownSLine,
  RiArrowUpSLine,
  RiBarChartBoxLine,
  RiDashboard2Line,
  RiPieChartLine,
  RiSearchLine,
  RiLoader4Line
} from "@remixicon/react";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip as RechartsTooltip, XAxis, YAxis } from "recharts";
import "./styles.css";

interface DashboardItem {
  id: string;
  title: string;
  value: string;
  description: string;
  charts: number;
  icon: typeof RiBarChartBoxLine;
}

function MiniChartCard({ chart, dashboardType }: { chart: MalakIntegrationChart; dashboardType: string }) {
  const { data: chartData, isLoading: isLoadingChartData, error } = useQuery({
    queryKey: [FETCH_CHART_DATA_POINTS, chart.reference],
    queryFn: async () => {
      if (!chart.reference) return null;
      const response = await client.dashboards.chartsDetail(chart.reference);
      return response.data;
    },
    enabled: !!chart.reference,
  });

  const formattedData = formatChartData(chartData?.data_points);

  if (isLoadingChartData) {
    return (
      <Card className="p-2">
        <div className="flex items-center justify-center h-[100px]">
          <RiLoader4Line className="h-4 w-4 animate-spin text-muted-foreground" />
        </div>
      </Card>
    );
  }

  if (error || !chart.reference) {
    return (
      <Card className="p-2">
        <div className="flex flex-col items-center justify-center h-[100px] text-center">
          <RiBarChartBoxLine className="h-4 w-4 text-muted-foreground mb-1" />
          <p className="text-xs text-muted-foreground">Failed to load chart</p>
        </div>
      </Card>
    );
  }

  const hasNoData = !formattedData || formattedData.length === 0;

  if (hasNoData) {
    return (
      <Card className="p-2">
        <div className="flex flex-col items-center justify-center h-[100px] text-center">
          <RiBarChartBoxLine className="h-4 w-4 text-muted-foreground mb-1" />
          <p className="text-xs text-muted-foreground">No data available</p>
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-2">
      <div className="flex items-center gap-2 mb-1">
        {chart.chart_type === "pie" ? (
          <RiPieChartLine className="h-4 w-4 text-muted-foreground" />
        ) : (
          <RiBarChartBoxLine className="h-4 w-4 text-muted-foreground" />
        )}
        <span className="text-xs font-medium truncate">{chart.user_facing_name}</span>
      </div>
      {chart.chart_type === "pie" ? (
        <ChartContainer className="w-full h-full" config={{}}>
          <PieChart
            width={200}
            height={100}
            margin={{ top: 5, right: 5, left: 5, bottom: 5 }}
          >
            <Pie
              data={formattedData}
              cx="50%"
              cy="50%"
              labelLine={false}
              outerRadius={35}
              dataKey="value"
            >
              {formattedData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={getChartColors(index)} />
              ))}
            </Pie>
            <RechartsTooltip 
              formatter={(value: number) => 
                formatTooltipValue(value, chartData?.data_points?.[0]?.data_point_type)
              } 
            />
          </PieChart>
        </ChartContainer>
      ) : (
        <ChartContainer className="w-full h-full" config={{}}>
          <BarChart
            width={200}
            height={100}
            data={formattedData}
            margin={{ top: 5, right: 5, left: -15, bottom: 0 }}
          >
            <XAxis dataKey="name" stroke="#888888" fontSize={8} />
            <YAxis stroke="#888888" fontSize={8} />
            <RechartsTooltip 
              formatter={(value: number) => 
                formatTooltipValue(value, chartData?.data_points?.[0]?.data_point_type)
              } 
            />
            <Bar dataKey="value" fill="#3B82F6" radius={[2, 2, 0, 0]} />
          </BarChart>
        </ChartContainer>
      )}
    </Card>
  );
}

export const Dashboard = createReactBlockSpec(
  {
    type: "dashboard",
    propSchema: {
      textAlignment: defaultProps.textAlignment,
      textColor: defaultProps.textColor,
      selectedItem: {
        default: "",
        values: [] as string[],
      },
    },
    content: "inline",
  },
  {
    render: (props) => {
      const [isExpanded, setIsExpanded] = useState(false);
      const [search, setSearch] = useState("");

      const { data: dashboardsResponse, isLoading, error } = useQuery<ServerListDashboardResponse>({
        queryKey: [LIST_DASHBOARDS],
        queryFn: () => client.dashboards.dashboardsList().then(res => res.data),
      });

      // Convert API dashboards to internal format and filter out dashboards with no charts
      const availableDashboards: DashboardItem[] = dashboardsResponse?.dashboards
        ?.filter(dashboard => (dashboard.chart_count || 0) > 0)
        ?.map(dashboard => ({
          id: dashboard.reference || "",
          title: dashboard.title || "Untitled Dashboard",
          value: dashboard.reference || "",
          description: dashboard.description || "No description available",
          charts: dashboard.chart_count || 0,
          icon: RiBarChartBoxLine
        })) || [];

      const selectedItem = availableDashboards.find(
        (item) => item.value === props.block.props.selectedItem
      );

      const filteredItems = availableDashboards.filter((item) => {
        if (!search) return true;

        const searchLower = search.toLowerCase();
        return item.title.toLowerCase().includes(searchLower);
      });

      // Get actual charts for the selected dashboard
      const { data: selectedDashboardCharts } = useQuery({
        queryKey: [LIST_DASHBOARDS, selectedItem?.id, "charts"],
        queryFn: async () => {
          if (!selectedItem?.id) return null;
          const response = await client.dashboards.dashboardsDetail(selectedItem.id);
          return response.data;
        },
        enabled: !!selectedItem?.id,
      });

      // Sort charts based on their positions if available
      const sortedCharts = [...(selectedDashboardCharts?.charts || [])].sort((a, b) => {
        const posA = selectedDashboardCharts?.positions?.find(p => p.chart_id === a.id)?.order_index ?? 0;
        const posB = selectedDashboardCharts?.positions?.find(p => p.chart_id === b.id)?.order_index ?? 0;
        return posA - posB;
      });

      if (isLoading) {
        return (
          <div className="flex items-center justify-center p-6 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed">
            Loading available dashboards...
          </div>
        );
      }

      if (error) {
        return (
          <div className="flex items-center justify-center p-6 text-sm text-destructive bg-destructive/10 rounded-md border border-dashed border-destructive">
            Error loading dashboards. Please try again.
          </div>
        );
      }

      return (
        <Card className="dashboard">
          <div className="flex items-center justify-between gap-4">
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  role="combobox"
                  className={cn(
                    "w-[300px] justify-between",
                    !selectedItem && "text-muted-foreground"
                  )}
                >
                  <div className="flex items-center gap-2">
                    <RiDashboard2Line className="h-4 w-4" />
                    {selectedItem ? selectedItem.title : "Select Dashboard"}
                  </div>
                  {selectedItem && (
                    <Badge variant="secondary" className="ml-2">
                      {selectedItem.charts} charts
                    </Badge>
                  )}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-[300px] p-0">
                <Command shouldFilter={false}>
                  <div className="flex items-center border-b px-3">
                    <RiSearchLine className="mr-2 h-4 w-4 shrink-0 opacity-50" />
                    <CommandInput
                      placeholder="Search dashboards..."
                      className="h-9 w-full border-0 bg-transparent p-0 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-0"
                      value={search}
                      onValueChange={setSearch}
                    />
                  </div>
                  <CommandList>
                    <CommandEmpty>No dashboards found.</CommandEmpty>
                    <CommandGroup>
                      {filteredItems.map((item) => (
                        <CommandItem
                          key={item.value}
                          value={item.value}
                          onSelect={() => {
                            props.editor.updateBlock(props.block, {
                              type: "dashboard",
                              props: { selectedItem: item.value },
                            });
                            setIsExpanded(true);
                            setSearch("");
                          }}
                          className="flex items-start justify-between py-3"
                        >
                          <div className="flex flex-col gap-1 flex-1 min-w-0">
                            <div className="flex items-center">
                              <item.icon className="mr-2 h-4 w-4 shrink-0" />
                              <span className="font-medium truncate">{item.title}</span>
                            </div>
                            <Tooltip delayDuration={200}>
                              <TooltipTrigger asChild>
                                <p className="text-xs text-muted-foreground truncate">
                                  {item.description}
                                </p>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p className="text-xs">{item.description}</p>
                              </TooltipContent>
                            </Tooltip>
                          </div>
                          <Badge variant="secondary" className="ml-4 shrink-0 self-center whitespace-nowrap">
                            {item.charts} charts
                          </Badge>
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
            {selectedItem && (
              <Button
                variant="ghost"
                size="sm"
                className="text-muted-foreground hover:text-foreground"
                onClick={() => setIsExpanded(!isExpanded)}
              >
                {isExpanded ? (
                  <>
                    <RiArrowUpSLine className="mr-1 h-4 w-4" />
                    Hide Charts
                  </>
                ) : (
                  <>
                    <RiArrowDownSLine className="mr-1 h-4 w-4" />
                    Show Charts
                  </>
                )}
              </Button>
            )}
          </div>

          {selectedItem && !isExpanded && (
            <div className="flex items-center justify-center p-4 mt-4 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed cursor-pointer hover:bg-muted/70 transition-colors"
              onClick={() => setIsExpanded(true)}
            >
              <RiBarChartBoxLine className="mr-2 h-4 w-4" />
              Click to view {selectedItem.charts} charts
            </div>
          )}

          {selectedItem && isExpanded && (
            <div className="grid grid-cols-2 gap-4 mt-4">
              {sortedCharts?.map((chart) => (
                <MiniChartCard
                  key={chart.reference}
                  chart={chart.chart || chart}
                  dashboardType={selectedItem.value}
                />
              )) || (
                <div className="col-span-2 flex items-center justify-center p-4 text-sm text-muted-foreground">
                  <RiLoader4Line className="mr-2 h-4 w-4 animate-spin" />
                  Loading charts...
                </div>
              )}
            </div>
          )}

          {!selectedItem && (
            <div className="flex items-center justify-center p-6 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed mt-4">
              Select a dashboard to view charts
            </div>
          )}
        </Card>
      );
    },
  }
);
