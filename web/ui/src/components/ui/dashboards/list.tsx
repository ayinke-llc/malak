"use client"

import { Card } from "@/components/ui/card";
import { RiDashboardLine } from "@remixicon/react";
import { format } from "date-fns";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import client from "@/lib/client";
import { LIST_DASHBOARDS } from "@/lib/query-constants";
import { Button } from "@/components/ui/button";
import { RiArrowLeftLine, RiArrowRightLine } from "@remixicon/react";
import { useState } from "react";
import type { MalakDashboard, ServerListDashboardResponse } from "@/client/Api";

export default function ListDashboards() {
  const [page, setPage] = useState(1);
  const perPage = 12;

  const { data, isLoading, isError } = useQuery<ServerListDashboardResponse>({
    queryKey: [LIST_DASHBOARDS, page],
    queryFn: async () => {
      const response = await client.dashboards.dashboardsList({
        page,
        per_page: perPage,
      });
      return response.data;
    },
  });

  if (isLoading) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-background">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="rounded-full bg-muted p-4">
            <RiDashboardLine className="h-8 w-8 text-muted-foreground animate-spin" />
          </div>
          <h3 className="mt-6 text-lg font-medium text-foreground">
            Loading dashboards...
          </h3>
        </div>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-background">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="rounded-full bg-destructive/10 p-4">
            <RiDashboardLine className="h-8 w-8 text-destructive" />
          </div>
          <h3 className="mt-6 text-lg font-medium text-foreground">
            Error loading dashboards
          </h3>
          <p className="mt-2 text-sm text-muted-foreground">
            Please try again later.
          </p>
        </div>
      </Card>
    );
  }

  if (!data?.dashboards?.length) {
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

  const totalPages = Math.ceil(data.meta.paging.total / perPage);

  return (
    <div className="space-y-6">
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {data.dashboards.map((dashboard: MalakDashboard) => (
          <Link
            key={dashboard.id}
            href={`/dashboards/${dashboard.reference}`}
            className="block transition-transform hover:scale-[1.02]"
          >
            <Card className="flex flex-col p-6 space-y-4 cursor-pointer hover:shadow-md transition-shadow">
              <div>
                <h4 className="text-lg font-medium">{dashboard.title}</h4>
                <p className="text-sm text-muted-foreground">{dashboard.description}</p>
              </div>

              <div className="flex items-center justify-between text-sm text-muted-foreground">
                <div className="flex items-center gap-2">
                  <RiDashboardLine className="h-4 w-4" />
                  <span>{dashboard.chart_count} charts</span>
                </div>
                <span>Created {format(new Date(dashboard.created_at!), "MMM d, yyyy")}</span>
              </div>
            </Card>
          </Link>
        ))}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
          >
            <RiArrowLeftLine className="h-4 w-4" />
          </Button>
          <span className="text-sm text-muted-foreground">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
          >
            <RiArrowRightLine className="h-4 w-4" />
          </Button>
        </div>
      )}
    </div>
  );
}
