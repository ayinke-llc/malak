import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { MalakIntegrationChart } from "@/client/Api";
import { ChartTypeIcon } from "./ChartTypeIcon";

export function ChartCard({ chart, onClick, isSelected }: {
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
