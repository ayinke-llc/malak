import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { RiLoader4Line, RiBarChartBoxLine, RiAddLine } from "@remixicon/react";
import { Card } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ChartTypeIcon } from "./ChartTypeIcon";
import type { MalakIntegrationChart, MalakIntegrationType, MalakWorkspaceIntegration } from "@/client/Api";
import client from "@/lib/client";
import { FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import { formatChartData, formatTooltipValue } from "@/lib/chart-utils";
import { ColumnDef, flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { AddDataPointDialog } from "./AddDataPointDialog";
import { Button } from "@/components/ui/button";

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

interface ChartDataViewProps {
  chart: MalakIntegrationChart;
  isSystemIntegration?: boolean;
  workspaceIntegration: MalakWorkspaceIntegration;
}

export function ChartDataView({ chart, isSystemIntegration, workspaceIntegration }: ChartDataViewProps) {
  const [open, setOpen] = useState(false);
  const { data: chartData, isLoading } = useQuery({
    queryKey: [FETCH_CHART_DATA_POINTS, chart?.reference],
    queryFn: async () => {
      if (!chart?.reference) return null;
      const response = await client.dashboards.chartsDetail(chart.reference as string);
      return response.data;
    },
    enabled: !!chart?.reference,
  });

  // Memoize the data transformation
  const sortedData = useMemo(() => {
    const formattedData = formatChartData(chartData?.data_points, chart.data_point_type);
    if (!formattedData) return [];
    
    return [...formattedData]
      .sort((a, b) => {
        const aDate = chartData?.data_points?.find(dp => dp.point_name === a.name)?.created_at;
        const bDate = chartData?.data_points?.find(dp => dp.point_name === b.name)?.created_at;
        if (!aDate || !bDate) return 0;
        return new Date(bDate).getTime() - new Date(aDate).getTime();
      })
      .map(point => ({
        name: point.name,
        value: formatTooltipValue(point.value, chart.data_point_type)[0],
      }));
  }, [chartData?.data_points, chart.data_point_type]);

  // Memoize the table instance
  const table = useReactTable({
    data: sortedData,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Card className="h-[calc(100vh-200px)] overflow-hidden">
      <div className="p-4 border-b">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-md bg-primary/10">
              <ChartTypeIcon type={chart.chart_type} />
            </div>
            <div>
              <h4 className="font-medium">{chart.user_facing_name}</h4>
              <p className="text-sm text-muted-foreground">{chart.internal_name || "No description available"}</p>
            </div>
          </div>
          {isSystemIntegration && (
            <>
              <Button variant="outline" size="sm" onClick={() => setOpen(true)} className="gap-2">
                <RiAddLine className="h-4 w-4" />
                Add Data Point
              </Button>
              <AddDataPointDialog 
                chart={chart} 
                open={open} 
                onOpenChange={setOpen}
                workspaceIntegration={workspaceIntegration}
              />
            </>
          )}
        </div>
      </div>
      <ScrollArea className="h-[calc(100vh-280px)]">
        <div className="p-4">
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <RiLoader4Line className="h-6 w-6 animate-spin" />
            </div>
          ) : sortedData.length > 0 ? (
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