/* eslint-disable react-hooks/rules-of-hooks */
import type { MalakIntegrationChart, ServerListDashboardResponse, MalakIntegrationDataPointType } from "@/client/Api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { TooltipContent, TooltipTrigger, Tooltip as UITooltip } from "@/components/ui/tooltip";
import client from "@/lib/client";
import { FETCH_CHART_DATA_POINTS, LIST_CHARTS } from "@/lib/query-constants";
import { cn } from "@/lib/utils";
import { formatChartData, formatTooltipValue, getChartColors } from "@/lib/chart-utils";
import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import { RiBarChartBoxLine, RiPieChartLine, RiSearchLine, RiLoader4Line } from "@remixicon/react";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip, XAxis, YAxis } from "recharts";

interface ChartConfig {
  id: string;
  name: string;
  type: "bar" | "pie";
  description: string;
  reference: string;
  data_point_type?: MalakIntegrationDataPointType;
}

// Convert API chart to internal chart config
const toChartConfig = (chart: MalakIntegrationChart): ChartConfig => {
  return {
    id: chart.id || "",
    name: chart.user_facing_name || "",
    type: chart.chart_type === "pie" ? "pie" : "bar",
    description: `${chart.user_facing_name || "Chart"} visualization`,
    reference: chart.reference || "",
    data_point_type: chart.data_point_type || undefined,
  };
};

interface ChartDisplayProps {
  chart: ChartConfig;
}

function ChartDisplay({ chart }: ChartDisplayProps) {
  const { data: chartData, isLoading: isLoadingChartData, error } = useQuery({
    queryKey: [FETCH_CHART_DATA_POINTS, chart.reference],
    queryFn: async () => {
      const response = await client.dashboards.chartsDetail(chart.reference);
      return response.data;
    },
    enabled: !!chart.reference,
  });

  const formattedData = formatChartData(chartData?.data_points, chart.data_point_type);

  if (isLoadingChartData) {
    return (
      <Card className="p-4">
        <div className="flex items-center justify-center h-[300px]">
          <RiLoader4Line className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="p-4">
        <div className="flex flex-col items-center justify-center h-[300px] text-center">
          <RiBarChartBoxLine className="h-8 w-8 text-muted-foreground mb-2" />
          <p className="text-sm text-muted-foreground">Failed to load chart data</p>
          <p className="text-xs text-muted-foreground mt-1">Please try again later</p>
        </div>
      </Card>
    );
  }

  const hasNoData = !formattedData || formattedData.length === 0;

  if (hasNoData) {
    return (
      <Card className="p-4">
        <div className="flex flex-col items-center justify-center h-[300px] text-center">
          <RiBarChartBoxLine className="h-8 w-8 text-muted-foreground mb-2" />
          <p className="text-sm text-muted-foreground">No data available</p>
          <p className="text-xs text-muted-foreground mt-1">Check back later for updates</p>
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-4">
      <div className="flex items-center gap-2 mb-4">
        {chart.type === "bar" ? (
          <RiBarChartBoxLine className="h-5 w-5 text-muted-foreground" />
        ) : (
          <RiPieChartLine className="h-5 w-5 text-muted-foreground" />
        )}
        <span className="text-sm font-medium">{chart.name}</span>
      </div>

      <div className="w-full aspect-[16/9] min-h-[300px] relative">
        <ChartContainer className="absolute inset-0" config={{}}>
          {chart.type === "bar" ? (
            <BarChart
              width={500}
              height={300}
              data={formattedData}
              margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
            >
              <XAxis dataKey="name" stroke="#888888" />
              <YAxis stroke="#888888" />
              <Tooltip
                formatter={(value: number) =>
                  formatTooltipValue(value, chart.data_point_type)
                }
              />
              <Bar dataKey="value" fill="#3B82F6" radius={[4, 4, 0, 0]} />
            </BarChart>
          ) : (
            <PieChart
              width={500}
              height={300}
              margin={{ top: 5, right: 5, left: 5, bottom: 5 }}
            >
              <Pie
                data={formattedData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={120}
                innerRadius={60}
                paddingAngle={2}
                dataKey="value"
              >
                {formattedData.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={getChartColors(index)}
                    strokeWidth={1}
                  />
                ))}
              </Pie>
              <Tooltip
                formatter={(value: number) =>
                  formatTooltipValue(value, chart.data_point_type)
                }
              />
            </PieChart>
          )}
        </ChartContainer>
      </div>
    </Card>
  );
}

export const Chart = createReactBlockSpec(
  {
    type: "chart",
    propSchema: {
      textAlignment: defaultProps.textAlignment,
      textColor: defaultProps.textColor,
      selectedChart: {
        default: "",
        values: [] as string[],
      },
    },
    content: "inline",
  },
  {
    render: (props) => {
      const [search, setSearch] = useState("");

      const { data: chartsResponse, isLoading, error } = useQuery({
        queryKey: [LIST_CHARTS],
        queryFn: () => client.dashboards.chartsList(),
      });

      // Convert API charts to internal chart configs
      const availableCharts = chartsResponse?.data.charts?.map(toChartConfig) || [];

      const selectedChart = availableCharts.find(
        (chart) => chart.reference === props.block.props.selectedChart
      );

      const filteredCharts = availableCharts.filter((chart) => {
        if (!search) return true;
        const searchLower = search.toLowerCase();
        return (
          chart.name.toLowerCase().includes(searchLower) ||
          chart.description.toLowerCase().includes(searchLower)
        );
      });

      if (isLoading) {
        return (
          <div className="flex items-center justify-center p-6 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed">
            Loading available charts...
          </div>
        );
      }

      if (error) {
        return (
          <div className="flex items-center justify-center p-6 text-sm text-destructive bg-destructive/10 rounded-md border border-dashed border-destructive">
            Error loading charts. Please try again.
          </div>
        );
      }

      return (
        <div className="chart-block">
          <Popover>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                className={cn(
                  "w-full justify-between",
                  !selectedChart && "text-muted-foreground"
                )}
              >
                {selectedChart ? (
                  <>
                    {selectedChart.type === "bar" ? (
                      <RiBarChartBoxLine className="mr-2 h-5 w-5" />
                    ) : (
                      <RiPieChartLine className="mr-2 h-5 w-5" />
                    )}
                    {selectedChart.name}
                  </>
                ) : (
                  <>
                    <RiBarChartBoxLine className="mr-2 h-5 w-5" />
                    Select a chart
                  </>
                )}
                {selectedChart && (
                  <Badge variant="secondary" className="ml-2">
                    {selectedChart.type === "bar" ? "Bar Chart" : "Pie Chart"}
                  </Badge>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[300px] p-0" align="start">
              <Command shouldFilter={false}>
                <div className="flex items-center border-b px-3">
                  <RiSearchLine className="mr-2 h-4 w-4 shrink-0 opacity-50" />
                  <CommandInput
                    placeholder="Search charts..."
                    className="h-9 w-full border-0 bg-transparent p-0 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-0"
                    value={search}
                    onValueChange={setSearch}
                  />
                </div>
                <CommandList>
                  <CommandEmpty>No chart found.</CommandEmpty>
                  <CommandGroup>
                    {filteredCharts.map((chart) => (
                      <CommandItem
                        key={chart.id}
                        value={chart.name}
                        onSelect={() => {
                          props.editor.updateBlock(props.block, {
                            type: "chart",
                            props: { selectedChart: chart.reference },
                          });
                          setSearch("");
                        }}
                        className="flex items-start justify-between py-3"
                      >
                        <div className="flex flex-col gap-1 flex-1 min-w-0">
                          <div className="flex items-center">
                            {chart.type === "bar" ? (
                              <RiBarChartBoxLine className="mr-2 h-4 w-4 shrink-0" />
                            ) : (
                              <RiPieChartLine className="mr-2 h-4 w-4 shrink-0" />
                            )}
                            <span className="font-medium truncate">{chart.name}</span>
                          </div>
                          <UITooltip delayDuration={200}>
                            <TooltipTrigger asChild>
                              <p className="text-xs text-muted-foreground truncate">
                                {chart.description}
                              </p>
                            </TooltipTrigger>
                            <TooltipContent>
                              <p className="text-xs">{chart.description}</p>
                            </TooltipContent>
                          </UITooltip>
                        </div>
                        <Badge variant="secondary" className="ml-4 shrink-0 self-center">
                          {chart.type === "bar" ? "Bar" : "Pie"}
                        </Badge>
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>

          {selectedChart ? (
            <div className="mt-4">
              <ChartDisplay chart={selectedChart} />
            </div>
          ) : (
            <div className="flex items-center justify-center p-6 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed mt-4">
              Select a chart to display
            </div>
          )}
        </div>
      );
    },
  }
);
