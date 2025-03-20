"use client"

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Card } from "@/components/ui/card";
import { 
  RiErrorWarningLine, 
  RiBarChartBoxLine, 
  RiPieChartLine, 
  RiLoader4Line,
  RiApps2Line,
  RiArrowRightSLine,
  RiCheckLine,
  RiCloseLine,
  RiArrowLeftLine
} from "@remixicon/react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import client from "@/lib/client";
import { LIST_INTEGRATIONS, LIST_CHARTS, FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import { IntegrationsList } from "@/components/ui/integrations/list";
import type { MalakWorkspaceIntegration, MalakIntegrationChart, MalakIntegrationChartType } from "@/client/Api";
import { cn } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { formatChartData, formatTooltipValue } from "@/lib/chart-utils";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";

interface DataPoint {
  name: string;
  value: number | string;
}

const columns: ColumnDef<DataPoint>[] = [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "value",
    header: "Value",
  },
];

function ChartTypeIcon({ type }: { type?: MalakIntegrationChartType }) {
  if (type === "pie") {
    return <RiPieChartLine className="h-4 w-4" />;
  }
  return <RiBarChartBoxLine className="h-4 w-4" />;
}

function IntegrationCard({ integration, isSelected, onClick }: { 
  integration: MalakWorkspaceIntegration; 
  isSelected: boolean;
  onClick: () => void;
}) {
  const isDisabled = !integration.is_enabled || !integration.integration?.is_enabled;

  return (
    <Button
      key={integration.id}
      variant={isSelected ? "secondary" : "outline"}
      className={cn(
        "w-full justify-between p-4 h-auto",
        isDisabled && "opacity-50 cursor-not-allowed",
        "text-left"
      )}
      onClick={onClick}
      disabled={isDisabled}
    >
      <div className="flex items-center gap-3">
        {integration.integration?.logo_url ? (
          <img 
            src={integration.integration.logo_url} 
            alt={integration.integration.integration_name || "Integration"} 
            className="w-6 h-6 rounded"
          />
        ) : (
          <RiApps2Line className="w-6 h-6" />
        )}
        <div className="text-left">
          <div className="font-medium text-foreground">{integration.integration?.integration_name}</div>
          <div className="text-xs text-muted-foreground truncate max-w-[180px]">
            {integration.integration?.description || "No description available"}
          </div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Badge variant={integration.is_active ? "default" : "destructive"} className="h-5 shrink-0">
          {integration.is_active ? (
            <RiCheckLine className="h-3 w-3 mr-1" />
          ) : (
            <RiCloseLine className="h-3 w-3 mr-1" />
          )}
          {integration.is_active ? "Active" : "Inactive"}
        </Badge>
        <RiArrowRightSLine className={cn(
          "h-4 w-4 transition-transform shrink-0",
          isSelected && "transform rotate-90"
        )} />
      </div>
    </Button>
  );
}

function ChartDataView({ chart }: { chart: MalakIntegrationChart }) {
  const { data: chartData, isLoading } = useQuery({
    queryKey: [FETCH_CHART_DATA_POINTS, chart?.reference],
    queryFn: async () => {
      if (!chart?.reference) return null;
      const response = await client.dashboards.chartsDetail(chart.reference);
      return response.data;
    },
    enabled: !!chart?.reference,
  });

  const formattedData = formatChartData(chartData?.data_points);
  
  // Sort data points by created_at in descending order
  const sortedData = [...(formattedData || [])].sort((a, b) => {
    const aDate = chartData?.data_points?.find(dp => dp.point_name === a.name)?.created_at;
    const bDate = chartData?.data_points?.find(dp => dp.point_name === b.name)?.created_at;
    if (!aDate || !bDate) return 0;
    return new Date(bDate).getTime() - new Date(aDate).getTime();
  }).map(point => ({
    name: point.name,
    value: formatTooltipValue(point.value, chartData?.data_points?.[0]?.data_point_type)[0],
  }));

  const table = useReactTable({
    data: sortedData,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Card className="h-[calc(100vh-200px)] overflow-hidden">
      <div className="p-4 border-b">
        <div className="flex items-center gap-2">
          <div className="p-2 rounded-md bg-primary/10">
            <ChartTypeIcon type={chart.chart_type} />
          </div>
          <div>
            <h4 className="font-medium">{chart.user_facing_name}</h4>
            <p className="text-sm text-muted-foreground">{chart.internal_name || "No description available"}</p>
          </div>
        </div>
      </div>
      <ScrollArea className="h-[calc(100vh-280px)]">
        <div className="p-4">
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <RiLoader4Line className="h-6 w-6 animate-spin" />
            </div>
          ) : sortedData && sortedData.length > 0 ? (
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  {table.getHeaderGroups().map((headerGroup) => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map((header) => (
                        <TableHead key={header.id}>
                          {header.isPlaceholder
                            ? null
                            : flexRender(
                                header.column.columnDef.header,
                                header.getContext()
                              )}
                        </TableHead>
                      ))}
                    </TableRow>
                  ))}
                </TableHeader>
                <TableBody>
                  {table.getRowModel().rows?.length ? (
                    table.getRowModel().rows.map((row) => (
                      <TableRow
                        key={row.id}
                        data-state={row.getIsSelected() && "selected"}
                      >
                        {row.getVisibleCells().map((cell) => (
                          <TableCell key={cell.id}>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        ))}
                      </TableRow>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell colSpan={columns.length} className="h-24 text-center">
                        No results.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <RiBarChartBoxLine className="h-8 w-8 text-muted-foreground mb-2" />
              <p className="text-sm text-muted-foreground">No data available</p>
            </div>
          )}
        </div>
      </ScrollArea>
    </Card>
  );
}

function ChartCard({ chart, onClick, isSelected }: { 
  chart: MalakIntegrationChart; 
  onClick: () => void;
  isSelected: boolean;
}) {
  return (
    <Card 
      className={cn(
        "p-4 hover:shadow-md transition-shadow duration-200 cursor-pointer",
        isSelected && "border-primary"
      )}
      onClick={onClick}
    >
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <div className="p-2 rounded-md bg-primary/10">
            <ChartTypeIcon type={chart.chart_type} />
          </div>
          <h4 className="font-medium truncate">{chart.user_facing_name}</h4>
        </div>
        <Badge variant="outline" className="capitalize">
          {chart.chart_type?.replace("IntegrationChartType", "")}
        </Badge>
      </div>
      <p className="text-sm text-muted-foreground line-clamp-2">
        {chart.internal_name || "No description available"}
      </p>
    </Card>
  );
}

export default function Metrics() {
  const [selectedIntegration, setSelectedIntegration] = useState<MalakWorkspaceIntegration | null>(null);
  const [selectedChart, setSelectedChart] = useState<MalakIntegrationChart | null>(null);

  const { data: integrationsData, isLoading: isLoadingIntegrations, error: integrationsError, refetch: refetchIntegrations } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  const { data: chartsData, isLoading: isLoadingCharts } = useQuery({
    queryKey: [LIST_CHARTS],
    queryFn: () => client.dashboards.chartsList(),
    enabled: !!selectedIntegration,
  });

  // Filter charts for the selected integration
  const integrationCharts = chartsData?.data?.charts?.filter(
    (chart: MalakIntegrationChart) => chart.workspace_integration_id === selectedIntegration?.id
  ) || [];

  if (integrationsError) {
    return (
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3 className="text-lg font-medium">Integration Metrics</h3>
              <p className="text-sm text-muted-foreground">View metrics from your connected integrations</p>
            </div>
          </div>
        </section>

        <section className="mt-10">
          <Card className="flex flex-col items-center justify-center py-16 px-4 bg-background">
            <div className="flex flex-col items-center justify-center text-center max-w-sm">
              <div className="rounded-full bg-destructive/10 p-4">
                <RiErrorWarningLine className="h-8 w-8 text-destructive" />
              </div>
              <h3 className="mt-6 text-lg font-medium text-foreground">
                Error loading integrations
              </h3>
              <p className="mt-2 text-sm text-muted-foreground">
                We could not load your integrations. Please try again.
              </p>
              <Button
                onClick={() => refetchIntegrations()}
                className="mt-6"
                variant="outline"
              >
                Try Again
              </Button>
            </div>
          </Card>
        </section>
      </div>
    );
  }

  return (
    <div className="pt-6 bg-background min-h-screen">
      <section>
        <div className="sm:flex sm:items-center sm:justify-between mb-8">
          <div>
            <h3 className="text-2xl font-semibold">Integration Metrics</h3>
            <p className="text-sm text-muted-foreground mt-1">View and analyze metrics from your connected integrations</p>
          </div>
        </div>
      </section>

      <section>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-4">
          {/* Integrations List */}
          <div className="lg:col-span-1">
            <Card className="overflow-hidden">
              <div className="p-4 border-b">
                <h4 className="font-medium">Available Integrations</h4>
                <p className="text-sm text-muted-foreground mt-1">Select an integration to view its metrics</p>
              </div>
              <ScrollArea className="h-[calc(100vh-300px)]">
                <div className="p-4 space-y-2">
                  {isLoadingIntegrations ? (
                    <>
                      <Skeleton className="h-20 w-full" />
                      <Skeleton className="h-20 w-full" />
                      <Skeleton className="h-20 w-full" />
                    </>
                  ) : (
                    (integrationsData?.data?.integrations || []).map((integration) => (
                      <IntegrationCard
                        key={integration.id}
                        integration={integration}
                        isSelected={selectedIntegration?.id === integration.id}
                        onClick={() => {
                          setSelectedIntegration(integration);
                          setSelectedChart(null);
                        }}
                      />
                    ))
                  )}
                </div>
              </ScrollArea>
            </Card>
          </div>

          {/* Charts and Data View */}
          <div className="lg:col-span-3">
            {selectedChart ? (
              <div>
                <div className="mb-4">
                  <Button
                    variant="ghost"
                    className="gap-2"
                    onClick={() => setSelectedChart(null)}
                  >
                    <RiArrowLeftLine className="h-4 w-4" />
                    Back to Charts
                  </Button>
                </div>
                <ChartDataView chart={selectedChart} />
              </div>
            ) : (
              <Card className="h-[calc(100vh-200px)] overflow-hidden">
                <div className="p-4 border-b">
                  <h4 className="font-medium">
                    {selectedIntegration ? (
                      <>Charts for {selectedIntegration.integration?.integration_name}</>
                    ) : (
                      <>Available Charts</>
                    )}
                  </h4>
                  <p className="text-sm text-muted-foreground mt-1">
                    {selectedIntegration ? (
                      <>Select a chart to view its data points</>
                    ) : (
                      <>Select an integration to view its charts</>
                    )}
                  </p>
                </div>
                <ScrollArea className="h-[calc(100vh-280px)]">
                  <div className="p-4">
                    {selectedIntegration ? (
                      isLoadingCharts ? (
                        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                          <Skeleton className="h-32 w-full" />
                          <Skeleton className="h-32 w-full" />
                          <Skeleton className="h-32 w-full" />
                        </div>
                      ) : integrationCharts.length > 0 ? (
                        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                          {integrationCharts.map((chart: MalakIntegrationChart) => (
                            <ChartCard 
                              key={chart.id} 
                              chart={chart}
                              isSelected={false}
                              onClick={() => setSelectedChart(chart)}
                            />
                          ))}
                        </div>
                      ) : (
                        <div className="flex flex-col items-center justify-center py-16 text-center">
                          <div className="rounded-full bg-muted p-4 mb-4">
                            <RiBarChartBoxLine className="h-6 w-6 text-muted-foreground" />
                          </div>
                          <h4 className="text-lg font-medium">No charts available</h4>
                          <p className="text-sm text-muted-foreground mt-1">
                            This integration doesn't have any charts configured yet
                          </p>
                        </div>
                      )
                    ) : (
                      <div className="flex flex-col items-center justify-center py-16 text-center">
                        <div className="rounded-full bg-muted p-4 mb-4">
                          <RiApps2Line className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <h4 className="text-lg font-medium">Select an integration</h4>
                        <p className="text-sm text-muted-foreground mt-1">
                          Choose an integration from the sidebar to view its charts
                        </p>
                      </div>
                    )}
                  </div>
                </ScrollArea>
              </Card>
            )}
          </div>
        </div>
      </section>
    </div>
  );
} 