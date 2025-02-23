import { defaultProps } from "@blocknote/core";
import { createReactBlockSpec } from "@blocknote/react";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { RiDashboard2Line, RiBarChartBoxLine } from "@remixicon/react";
import { cn } from "@/lib/utils";
import "./styles.css";

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

      return (
        <Card className="dashboard">
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
          {selectedItem ? (
            <div className="inline-content" ref={props.contentRef} />
          ) : (
            <div className="flex items-center justify-center p-6 text-sm text-muted-foreground bg-muted/50 rounded-md border border-dashed">
              Click to select a dashboard
            </div>
          )}
        </Card>
      );
    },
  }
);
