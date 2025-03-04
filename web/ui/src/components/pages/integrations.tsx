"use client"

import { useEffect } from "react";
import client from "@/lib/client";
import { LIST_INTEGRATIONS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { IntegrationsList } from "@/components/ui/integrations/list";
import { Card } from "@/components/ui/card";
import { RiErrorWarningLine } from "@remixicon/react";
import { Button } from "@/components/ui/button";

export default function Integrations() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  useEffect(() => {
    if (error) {
      toast.error("Error occurred while fetching integrations");
    }
  }, [error]);

  if (error) {
    return (
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3 id="integrations" className="text-lg font-medium">
                Available Integrations
              </h3>
              <p className="text-sm text-muted-foreground">View and manage these integrations on your workspace</p>
            </div>
          </div>
        </section>

        <section className="mt-10">
          <Card className="flex flex-col items-center justify-center py-16 px-4 bg-background">
            <div className="flex flex-col items-center justify-center text-center max-w-sm">
              <div className="rounded-full bg-destructive/10 p-4">
                <RiErrorWarningLine className="h-8 w-8 text-destructive" />
              </div>
              <h3 className="mt-6 text-lg font-medium text-foreground">
                Error loading integrations
              </h3>
              <p className="mt-2 text-sm text-muted-foreground">
                We could not load your integrations. Please try again.
              </p>
              <Button
                onClick={() => refetch()}
                className="mt-6"
                variant="outline"
              >
                Try Again
              </Button>
            </div>
          </Card>
        </section>
      </div>
    );
  }

  return (
    <div className="pt-6 bg-background">
      <section>
        <div className="sm:flex sm:items-center sm:justify-between">
          <div>
            <h3 id="integrations" className="text-lg font-medium">
              Available Integrations
            </h3>
            <p className="text-sm text-muted-foreground">View and manage these integrations on your workspace</p>
          </div>
        </div>
      </section>

      <section className="mt-10">
        <IntegrationsList
          integrations={data?.data?.integrations || []}
          isLoading={isLoading}
        />
      </section>
    </div>
  );
} 