"use client"

import { useEffect } from "react";
import client from "@/lib/client";
import { LIST_INTEGRATIONS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { IntegrationsList } from "@/components/ui/integrations/list";

export default function Integrations() {
  const { data, isLoading, error } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  useEffect(() => {
    if (error) {
      toast.error("Error occurred while fetching communication preferences");
    }
  }, [error]);

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
