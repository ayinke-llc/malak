import type { MalakIntegrationChart, ServerListDashboardResponse } from "@/client/Api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import client from "@/lib/client";
import { LIST_DASHBOARDS } from "@/lib/query-constants";
import { cn } from "@/lib/utils";
import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import {
  RiArrowDownSLine,
  RiArrowUpSLine,
  RiBarChartBoxLine,
  RiDashboard2Line,
  RiPieChartLine,
  RiSearchLine
} from "@remixicon/react";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip as RechartsTooltip, XAxis, YAxis } from "recharts";
import "./styles.css";

// Function to generate random bar chart data
const generateBarData = (prefix: string, count: number = 5) => {
  return Array.from({ length: count }, (_, i) => ({
    name: `${prefix} ${i + 1}`,
    value: Math.floor(Math.random() * 10000),
  }));
};

// Function to generate random pie chart data
const generatePieData = (categories: string[], total: number = 1000) => {
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

function getChartData(chartName: string, chartType: string) {
  if (chartType === "bar") {
    switch (chartName) {
      case "Monthly Revenue":
      case "Revenue Growth":
        return generateBarData("Month", 6);
      case "Weekly Active Users":
      case "User Growth":
        return generateBarData("Week", 4);
      case "Engagement Metrics":
        return generateBarData("Metric", 5);
      case "Quarterly Sales":
      case "Sales Growth":
        return generateBarData("Q", 4);
      case "Top Products":
        return generateBarData("Product", 5);
      case "KPI Metrics":
      case "Efficiency Metrics":
        return generateBarData("KPI", 4);
      case "Satisfaction Score":
        return generateBarData("Score", 5);
      case "Stock Levels":
      case "Inventory Turnover":
        return generateBarData("Category", 6);
      case "Campaign Results":
      case "ROI by Channel":
        return generateBarData("Campaign", 4);
      case "Team Performance":
        return generateBarData("Team", 5);
      case "Project Status":
      case "Timeline Progress":
      case "Resource Allocation":
        return generateBarData("Project", 4);
      case "Financial Overview":
      case "Profit Margins":
      case "Cash Flow":
      case "Investment Returns":
        return generateBarData("Period", 6);
      default:
        return generateBarData("Item", 5);
    }
  } else {
    switch (chartName) {
      case "Revenue Distribution":
      case "Product Revenue Split":
        return generatePieData(["Product A", "Product B", "Product C", "Product D"]);
      case "User Types":
      case "User Retention":
        return generatePieData(["New", "Active", "Inactive"]);
      case "Geographic Distribution":
        return generatePieData(["North", "South", "East", "West"]);
      case "Sales Channels":
      case "Sales by Region":
        return generatePieData(["Direct", "Partners", "Online", "Retail"]);
      case "Performance Split":
        return generatePieData(["Excellent", "Good", "Average", "Poor"]);
      case "Feedback Categories":
        return generatePieData(["Positive", "Neutral", "Negative"]);
      case "Category Split":
      case "Stock Distribution":
        return generatePieData(["Category A", "Category B", "Category C"]);
      case "Channel Mix":
        return generatePieData(["Social", "Email", "Search", "Direct"]);
      case "Task Distribution":
        return generatePieData(["Completed", "In Progress", "Pending"]);
      case "Project Types":
      case "Project Success Rate":
        return generatePieData(["Development", "Design", "Marketing", "Research"]);
      case "Cost Breakdown":
      case "Expense Categories":
      case "Budget Allocation":
        return generatePieData(["Operations", "Marketing", "R&D", "Admin"]);
      default:
        return generatePieData(["A", "B", "C"]);
    }
  }
}

interface ChartDataPoint {
  name: string;
  value: number;
}

interface DashboardItem {
  id: string;
  title: string;
  value: string;
  description: string;
  charts: number;
  icon: typeof RiBarChartBoxLine;
}

function MiniChartCard({ chart, dashboardType }: { chart: MalakIntegrationChart; dashboardType: string }) {
  // Generate mock data based on chart type
  const chartData = chart.chart_type === "pie"
    ? generatePieData(["Category A", "Category B", "Category C", "Category D"])
    : generateBarData(chart.user_facing_name || "Data", 5);

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
              data={chartData}
              cx="50%"
              cy="50%"
              labelLine={false}
              outerRadius={35}
              dataKey="value"
            >
              {chartData.map((entry: ChartDataPoint, index: number) => (
                <Cell key={`cell-${index}`} fill={`hsl(${index * 45}, 70%, 50%)`} />
              ))}
            </Pie>
            <RechartsTooltip />
          </PieChart>
        </ChartContainer>
      ) : (
        <ChartContainer className="w-full h-full" config={{}}>
          <BarChart
            width={200}
            height={100}
            data={chartData}
            margin={{ top: 5, right: 5, left: -15, bottom: 0 }}
          >
            <XAxis dataKey="name" stroke="#888888" fontSize={8} />
            <YAxis stroke="#888888" fontSize={8} />
            <RechartsTooltip />
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

      // Generate mock charts based on the dashboard's chart count
      const mockCharts = selectedItem ? Array.from({ length: selectedItem.charts }, (_, i) => ({
        reference: `mock-chart-${i}`,
        user_facing_name: `Chart ${i + 1}`,
        chart_type: i % 2 === 0 ? "bar" : "pie",
        description: `Mock chart ${i + 1}`,
        metadata: {},
      } as MalakIntegrationChart)) : [];

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
              {mockCharts.map((chart, index) => (
                <MiniChartCard
                  key={chart.reference}
                  chart={chart}
                  dashboardType={selectedItem.value}
                />
              ))}
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
