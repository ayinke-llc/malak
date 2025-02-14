"use client"

import { Card } from "@/components/ui/card";
import { RiDashboardLine, RiBarChart2Line, RiPieChartLine, RiLineChartLine } from "@remixicon/react";
import { format } from "date-fns";
import Link from "next/link";

// This is just a mock type for demonstration
type Dashboard = {
  id: string;
  title: string;
  description: string;
  charts_count: number;
  created_at: string;
};

// Mock data for demonstration
const mockDashboards: Dashboard[] = [
  {
    id: "1",
    title: "Revenue Overview",
    description: "Monthly revenue trends and projections",
    charts_count: 4,
    created_at: new Date().toISOString(),
  },
  {
    id: "2",
    title: "User Analytics",
    description: "User engagement and activity metrics",
    charts_count: 3,
    created_at: new Date().toISOString(),
  },
  {
    id: "3",
    title: "Cost Distribution",
    description: "Breakdown of operational costs",
    charts_count: 5,
    created_at: new Date().toISOString(),
  }
];

export default function ListDashboards() {
  if (mockDashboards.length === 0) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-background">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="rounded-full bg-muted p-4">
            <RiDashboardLine className="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 className="mt-6 text-lg font-medium text-foreground">
            No dashboards yet
          </h3>
          <p className="mt-2 text-sm text-muted-foreground">
            Create your first dashboard to visualize data from your integrations.
          </p>
        </div>
      </Card>
    );
  }

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      {mockDashboards.map((dashboard) => (
        <Link 
          key={dashboard.id} 
          href={`/dashboards/${dashboard.id}`}
          className="block transition-transform hover:scale-[1.02]"
        >
          <Card className="flex flex-col p-6 space-y-4 cursor-pointer hover:shadow-md transition-shadow">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-lg font-medium">{dashboard.title}</h4>
                <p className="text-sm text-muted-foreground">{dashboard.description}</p>
              </div>
            </div>

            <div className="flex-1">
              <div className="h-40 bg-muted rounded-md flex items-center justify-center">
                <div className="grid grid-cols-2 gap-2 text-muted-foreground">
                  <RiLineChartLine className="h-8 w-8" />
                  <RiPieChartLine className="h-8 w-8" />
                  <RiBarChart2Line className="h-8 w-8" />
                  <RiDashboardLine className="h-8 w-8" />
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <RiDashboardLine className="h-4 w-4" />
                <span>{dashboard.charts_count} charts</span>
              </div>
              <span>Created {format(new Date(dashboard.created_at), "MMM d, yyyy")}</span>
            </div>
          </Card>
        </Link>
      ))}
    </div>
  );
}
