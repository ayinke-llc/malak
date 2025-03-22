import { RiPieChartLine, RiBarChartBoxLine } from "@remixicon/react";
import type { MalakIntegrationChartType } from "@/client/Api";

export function ChartTypeIcon({ type }: { type?: MalakIntegrationChartType }) {
  if (type === "pie") {
    return <RiPieChartLine className="h-4 w-4" />;
  }
  return <RiBarChartBoxLine className="h-4 w-4" />;
} 