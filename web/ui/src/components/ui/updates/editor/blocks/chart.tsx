import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import { RiBarChartBoxLine, RiPieChartLine } from "@remixicon/react";
import { cn } from "@/lib/utils";
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip, XAxis, YAxis } from "recharts";

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

interface BarChartConfig {
  id: number;
  name: string;
  type: "bar";
  dataPrefix: string;
  dataCount: number;
}

interface PieChartConfig {
  id: number;
  name: string;
  type: "pie";
  categories: readonly string[];
}

type ChartConfig = BarChartConfig | PieChartConfig;

// Available charts configuration
const availableCharts: readonly ChartConfig[] = [
  { id: 1, name: "Monthly Revenue", type: "bar", dataPrefix: "Month", dataCount: 6 },
  { id: 2, name: "Revenue Distribution", type: "pie", categories: ["Product A", "Product B", "Product C", "Product D"] },
  { id: 3, name: "User Growth", type: "bar", dataPrefix: "Week", dataCount: 4 },
  { id: 4, name: "User Types", type: "pie", categories: ["New", "Active", "Inactive"] },
  { id: 5, name: "Sales Performance", type: "bar", dataPrefix: "Quarter", dataCount: 4 },
  { id: 6, name: "Geographic Split", type: "pie", categories: ["North", "South", "East", "West"] },
] as const;

interface ChartDisplayProps {
  chart: ChartConfig;
}

function ChartDisplay({ chart }: ChartDisplayProps) {
  const chartData = chart.type === "bar" 
    ? generateBarData(chart.dataPrefix, chart.dataCount)
    : generatePieData(chart.categories);

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
      
      {chart.type === "bar" ? (
        <ChartContainer className="w-full h-full" config={{}}>
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
        </ChartContainer>
      ) : (
        <ChartContainer className="w-full h-full" config={{}}>
          <PieChart
            width={500}
            height={300}
            margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
          >
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              labelLine={false}
              outerRadius={100}
              dataKey="value"
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={`hsl(${index * 45}, 70%, 50%)`} />
              ))}
            </Pie>
            <Tooltip />
          </PieChart>
        </ChartContainer>
      )}
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
      const selectedChart = availableCharts.find(
        (chart) => chart.id === props.block.props.selectedChart
      );

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
            <PopoverContent className="p-0" align="start">
              <Command>
                <CommandInput placeholder="Search charts..." />
                <CommandList>
                  <CommandEmpty>No chart found.</CommandEmpty>
                  <CommandGroup>
                    {availableCharts.map((chart) => (
                      <CommandItem
                        key={chart.id}
                        value={chart.name}
                        onSelect={() => {
                          props.editor.updateBlock(props.block, {
                            type: "chart",
                            props: { selectedChart: chart.id },
                          });
                        }}
                      >
                        <div className="flex items-center justify-between w-full">
                          <div className="flex items-center">
                            {chart.type === "bar" ? (
                              <RiBarChartBoxLine className="mr-2 h-4 w-4" />
                            ) : (
                              <RiPieChartLine className="mr-2 h-4 w-4" />
                            )}
                            <span>{chart.name}</span>
                          </div>
                          <Badge variant="secondary">
                            {chart.type === "bar" ? "Bar" : "Pie"}
                          </Badge>
                        </div>
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
