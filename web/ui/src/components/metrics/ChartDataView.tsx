import { useQuery } from "@tanstack/react-query";
import { RiLoader4Line, RiBarChartBoxLine } from "@remixicon/react";
import { Card } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ChartTypeIcon } from "./ChartTypeIcon";
import type { MalakIntegrationChart } from "@/client/Api";
import client from "@/lib/client";
import { FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import { formatChartData, formatTooltipValue } from "@/lib/chart-utils";
import { ColumnDef, flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";

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

export function ChartDataView({ chart }: { chart: MalakIntegrationChart }) {
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