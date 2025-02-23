import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { ChartContainer } from "@/components/ui/chart";
import { RiDashboard2Line, RiBarChartBoxLine, RiPieChartLine } from "@remixicon/react";
import { cn } from "@/lib/utils";
import { Bar, BarChart, Cell, Pie, PieChart, Tooltip, XAxis, YAxis } from "recharts";
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

function MiniChartCard({ chart, dashboardType }: { chart: any; dashboardType: string }) {
  const chartData = getChartData(chart.user_facing_name, chart.chart_type);

  return (
    <Card className="p-2">
      <div className="flex items-center gap-2 mb-1">
        {chart.chart_type === "bar" ? (
          <RiBarChartBoxLine className="h-4 w-4 text-muted-foreground" />
        ) : (
          <RiPieChartLine className="h-4 w-4 text-muted-foreground" />
        )}
        <span className="text-xs font-medium truncate">{chart.user_facing_name}</span>
      </div>
      {chart.chart_type === "bar" ? (
        <ChartContainer className="w-full h-full" config={{}}>
          <BarChart
            width={200}
            height={100}
            data={chartData}
            margin={{ top: 5, right: 5, left: -15, bottom: 0 }}
          >
            <XAxis dataKey="name" stroke="#888888" fontSize={8} />
            <YAxis stroke="#888888" fontSize={8} />
            <Tooltip />
            <Bar dataKey="value" fill="#3B82F6" radius={[2, 2, 0, 0]} />
          </BarChart>
        </ChartContainer>
      ) : (
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
              {chartData.map((entry: any, index: number) => (
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

export const dashboardItems = [
  { id: 1, title: "Revenue Overview", value: "revenue", charts: 4 },
  { id: 2, title: "User Analytics", value: "users", charts: 6 },
  { id: 3, title: "Sales Report", value: "sales", charts: 5 },
  { id: 4, title: "Performance Metrics", value: "performance", charts: 3 },
  { id: 5, title: "Customer Feedback", value: "feedback", charts: 2 },
  { id: 6, title: "Inventory Status", value: "inventory", charts: 4 },
  { id: 7, title: "Marketing ROI", value: "marketing", charts: 3 },
  { id: 8, title: "Team Progress", value: "team", charts: 2 },
  { id: 9, title: "Project Timeline", value: "projects", charts: 5 },
  { id: 10, title: "Financial Summary", value: "finance", charts: 7 },
] as const;

// Mock chart configurations for each dashboard
const mockDashboardCharts = {
  revenue: [
    { user_facing_name: "Monthly Revenue", chart_type: "bar" },
    { user_facing_name: "Revenue Distribution", chart_type: "pie" },
    { user_facing_name: "Revenue Growth", chart_type: "bar" },
    { user_facing_name: "Product Revenue Split", chart_type: "pie" },
  ],
  users: [
    { user_facing_name: "Weekly Active Users", chart_type: "bar" },
    { user_facing_name: "User Types", chart_type: "pie" },
    { user_facing_name: "User Growth", chart_type: "bar" },
    { user_facing_name: "Engagement Metrics", chart_type: "bar" },
    { user_facing_name: "User Retention", chart_type: "pie" },
    { user_facing_name: "Geographic Distribution", chart_type: "pie" },
  ],
  sales: [
    { user_facing_name: "Quarterly Sales", chart_type: "bar" },
    { user_facing_name: "Sales Channels", chart_type: "pie" },
    { user_facing_name: "Sales Growth", chart_type: "bar" },
    { user_facing_name: "Top Products", chart_type: "bar" },
    { user_facing_name: "Sales by Region", chart_type: "pie" },
  ],
  performance: [
    { user_facing_name: "KPI Metrics", chart_type: "bar" },
    { user_facing_name: "Performance Split", chart_type: "pie" },
    { user_facing_name: "Efficiency Metrics", chart_type: "bar" },
  ],
  feedback: [
    { user_facing_name: "Satisfaction Score", chart_type: "bar" },
    { user_facing_name: "Feedback Categories", chart_type: "pie" },
  ],
  inventory: [
    { user_facing_name: "Stock Levels", chart_type: "bar" },
    { user_facing_name: "Category Split", chart_type: "pie" },
    { user_facing_name: "Inventory Turnover", chart_type: "bar" },
    { user_facing_name: "Stock Distribution", chart_type: "pie" },
  ],
  marketing: [
    { user_facing_name: "Campaign Results", chart_type: "bar" },
    { user_facing_name: "Channel Mix", chart_type: "pie" },
    { user_facing_name: "ROI by Channel", chart_type: "bar" },
  ],
  team: [
    { user_facing_name: "Team Performance", chart_type: "bar" },
    { user_facing_name: "Task Distribution", chart_type: "pie" },
  ],
  projects: [
    { user_facing_name: "Project Status", chart_type: "bar" },
    { user_facing_name: "Project Types", chart_type: "pie" },
    { user_facing_name: "Timeline Progress", chart_type: "bar" },
    { user_facing_name: "Resource Allocation", chart_type: "bar" },
    { user_facing_name: "Project Success Rate", chart_type: "pie" },
  ],
  finance: [
    { user_facing_name: "Financial Overview", chart_type: "bar" },
    { user_facing_name: "Cost Breakdown", chart_type: "pie" },
    { user_facing_name: "Profit Margins", chart_type: "bar" },
    { user_facing_name: "Expense Categories", chart_type: "pie" },
    { user_facing_name: "Cash Flow", chart_type: "bar" },
    { user_facing_name: "Investment Returns", chart_type: "bar" },
    { user_facing_name: "Budget Allocation", chart_type: "pie" },
  ],
};

export const Dashboard = createReactBlockSpec(
  {
    type: "dashboard",
    propSchema: {
      textAlignment: defaultProps.textAlignment,
      textColor: defaultProps.textColor,
      selectedItem: {
        default: "",
        values: ["", ...dashboardItems.map(item => item.value)],
      },
    },
    content: "inline",
  },
  {
    render: (props) => {
      const selectedItem = dashboardItems.find(
        (item) => item.value === props.block.props.selectedItem
      );

      const dashboardCharts = selectedItem 
        ? mockDashboardCharts[selectedItem.value as keyof typeof mockDashboardCharts] || []
        : [];

      return (
        <Card className="dashboard">
          <div className="flex items-center gap-4">
            <Popover>
              <PopoverTrigger asChild>
                <Button 
                  variant="outline" 
                  role="combobox" 
                  className={cn(
                    "dashboard-button",
                    !selectedItem && "text-muted-foreground"
                  )}
                >
                  <RiDashboard2Line className="mr-2 h-5 w-5" />
                  {selectedItem ? selectedItem.title : "Select Dashboard"}
                  {selectedItem && (
                    <Badge variant="secondary" className="ml-2">
                      <RiBarChartBoxLine className="mr-1 h-4 w-4" />
                      {selectedItem.charts} charts
                    </Badge>
                  )}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="p-0" align="start">
                <Command>
                  <CommandInput placeholder="Search dashboards..." />
                  <CommandList>
                    <CommandEmpty>No dashboard found.</CommandEmpty>
                    <CommandGroup>
                      {dashboardItems.map((item) => (
                        <CommandItem
                          key={item.value}
                          value={item.value}
                          onSelect={() =>
                            props.editor.updateBlock(props.block, {
                              type: "dashboard",
                              props: { selectedItem: item.value },
                            })
                          }
                        >
                          <div className="flex items-center justify-between w-full">
                            <span>{item.title}</span>
                            <Badge variant="secondary" className="ml-2">
                              <RiBarChartBoxLine className="mr-1 h-3 w-3" />
                              {item.charts}
                            </Badge>
                          </div>
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
          </div>

          {selectedItem ? (
            <div className="grid grid-cols-2 gap-4 mt-4">
              {dashboardCharts.map((chart, index) => (
                <MiniChartCard 
                  key={index} 
                  chart={chart} 
                  dashboardType={selectedItem.value}
                />
              ))}
            </div>
          ) : (
            <div className="flex items-center justify-center p-6 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed mt-4">
              Select a dashboard to view charts
            </div>
          )}
        </Card>
      );
    },
  }
);
