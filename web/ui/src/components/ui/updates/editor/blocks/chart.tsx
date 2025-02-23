import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import { RiBarChartBoxLine, RiPieChartLine, RiSearchLine } from "@remixicon/react";
import { cn } from "@/lib/utils";
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip, XAxis, YAxis } from "recharts";
import { useState } from "react";
import { Tooltip as UITooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

interface ChartDataPoint {
  name: string;
  value: number;
}

// Function to generate random bar chart data
const generateBarData = (prefix: string, count: number = 5): ChartDataPoint[] => {
  return Array.from({ length: count }, (_, i) => ({
    name: `${prefix} ${i + 1}`,
    value: Math.floor(Math.random() * 10000),
  }));
};

// Function to generate random pie chart data
const generatePieData = (categories: readonly string[], total: number = 1000): ChartDataPoint[] => {
  const data = categories.map(name => ({
    name,
    value: Math.floor(Math.random() * (total / categories.length)),
  }));
  
  // Ensure values sum up to total
  const currentSum = data.reduce((sum, item) => sum + item.value, 0);
  const factor = total / currentSum;
  return data.map(item => ({
    ...item,
    value: Math.floor(item.value * factor),
  }));
};

interface ChartConfig {
  id: number;
  name: string;
  type: "bar" | "pie";
  description: string;
  dataPrefix?: string;
  dataCount?: number;
  categories?: readonly string[];
}

// Available charts configuration
const availableCharts: readonly ChartConfig[] = [
  { 
    id: 1, 
    name: "Monthly Revenue", 
    type: "bar", 
    description: "Track monthly revenue performance",
    dataPrefix: "Month", 
    dataCount: 6 
  },
  { 
    id: 2, 
    name: "Revenue Distribution", 
    type: "pie", 
    description: "View revenue split across products",
    categories: ["Product A", "Product B", "Product C", "Product D"] 
  },
  { 
    id: 3, 
    name: "User Growth", 
    type: "bar", 
    description: "Monitor weekly user growth trends",
    dataPrefix: "Week", 
    dataCount: 4 
  },
  { 
    id: 4, 
    name: "User Types", 
    type: "pie", 
    description: "Analyze user activity distribution",
    categories: ["New", "Active", "Inactive"] 
  },
  { 
    id: 5, 
    name: "Sales Performance", 
    type: "bar", 
    description: "Track quarterly sales metrics",
    dataPrefix: "Quarter", 
    dataCount: 4 
  },
  { 
    id: 6, 
    name: "Geographic Split", 
    type: "pie", 
    description: "View regional distribution data",
    categories: ["North", "South", "East", "West"] 
  },
] as const;

interface ChartDisplayProps {
  chart: ChartConfig;
}

function ChartDisplay({ chart }: ChartDisplayProps) {
  const chartData = chart.type === "bar" 
    ? generateBarData(chart.dataPrefix || "Item", chart.dataCount || 5)
    : generatePieData(chart.categories || ["A", "B", "C"]);

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
              data={chartData}
              margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
            >
              <XAxis dataKey="name" stroke="#888888" />
              <YAxis stroke="#888888" />
              <Tooltip />
              <Bar dataKey="value" fill="#3B82F6" radius={[4, 4, 0, 0]} />
            </BarChart>
          ) : (
            <PieChart
              width={500}
              height={300}
              margin={{ top: 5, right: 5, left: 5, bottom: 5 }}
            >
              <Pie
                data={chartData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, value }) => `${name}: ${value}`}
                outerRadius={120}
                innerRadius={60}
                paddingAngle={2}
                dataKey="value"
              >
                {chartData.map((entry, index) => (
                  <Cell 
                    key={`cell-${index}`} 
                    fill={`hsl(${index * 45}, 70%, 50%)`}
                    strokeWidth={1}
                  />
                ))}
              </Pie>
              <Tooltip />
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
        default: 0,
        values: [0, ...availableCharts.map(chart => chart.id)],
      },
    },
    content: "inline",
  },
  {
    render: (props) => {
      const [search, setSearch] = useState("");
      
      const selectedChart = availableCharts.find(
        (chart) => chart.id === props.block.props.selectedChart
      );

      const filteredCharts = availableCharts.filter((chart) => {
        if (!search) return true;
        const searchLower = search.toLowerCase();
        return (
          chart.name.toLowerCase().includes(searchLower) ||
          chart.description.toLowerCase().includes(searchLower)
        );
      });

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
                            props: { selectedChart: chart.id },
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
