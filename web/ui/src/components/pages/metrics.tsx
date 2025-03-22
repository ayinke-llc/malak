"use client"

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { RiErrorWarningLine, RiApps2Line, RiArrowLeftLine } from "@remixicon/react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import client from "@/lib/client";
import { LIST_INTEGRATIONS, LIST_CHARTS } from "@/lib/query-constants";
import type { MalakWorkspaceIntegration, MalakIntegrationChart } from "@/client/Api";
import { MalakIntegrationType } from "@/client/Api";
import { IntegrationCard } from "../metrics/IntegrationCard";
import { ChartCard } from "../metrics/ChartCard";
import { ChartDataView } from "../metrics/ChartDataView";
import { CreateChartDialog } from "../metrics/CreateChartDialog";

export default function Metrics() {
  const [selectedIntegration, setSelectedIntegration] = useState<MalakWorkspaceIntegration | null>(null);
  const [selectedChart, setSelectedChart] = useState<MalakIntegrationChart | null>(null);

  const { data: integrationsData, isLoading: isLoadingIntegrations, error: integrationsError, refetch: refetchIntegrations } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  const { data: chartsData, isLoading: isLoadingCharts } = useQuery({
    queryKey: [LIST_CHARTS],
    queryFn: () => client.dashboards.chartsList(),
    enabled: !!selectedIntegration,
  });

  // Filter charts for the selected integration
  const integrationCharts = chartsData?.data?.charts?.filter(
    (chart: MalakIntegrationChart) => chart.workspace_integration_id === selectedIntegration?.id
  ) || [];

  // Sort integrations by active status and creation date
  const sortedIntegrations = [...(integrationsData?.data?.integrations || [])].sort((a, b) => {
    // First sort by active status (active ones first)
    if (a.is_active !== b.is_active) {
      return b.is_active ? 1 : -1;
    }
    // Then sort by creation date (newest first)
    const dateA = a.created_at ? new Date(a.created_at).getTime() : 0;
    const dateB = b.created_at ? new Date(b.created_at).getTime() : 0;
    return dateB - dateA;
  });

  if (integrationsError) {
    return (
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3 className="text-lg font-medium">Integration Metrics</h3>
              <p className="text-sm text-muted-foreground">View metrics from your connected integrations</p>
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
                onClick={() => refetchIntegrations()}
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
    <div className="pt-6 bg-background min-h-screen">
      <section>
        <div className="sm:flex sm:items-center sm:justify-between mb-8">
          <div>
            <h3 className="text-2xl font-semibold">Integration Metrics</h3>
            <p className="text-sm text-muted-foreground mt-1">View and analyze metrics from your connected integrations</p>
          </div>
        </div>
      </section>

      <section>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-4">
          {/* Integrations List */}
          <div className="lg:col-span-1">
            <Card className="overflow-hidden">
              <div className="p-4 border-b">
                <h4 className="font-medium">Available Integrations</h4>
                <p className="text-sm text-muted-foreground mt-1">Select an integration to view its metrics</p>
              </div>
              <ScrollArea className="h-[calc(100vh-300px)]">
                <div className="p-4 space-y-2">
                  {isLoadingIntegrations ? (
                    <>
                      <Skeleton className="h-20 w-full" />
                      <Skeleton className="h-20 w-full" />
                      <Skeleton className="h-20 w-full" />
                    </>
                  ) : (
                    sortedIntegrations.map((integration) => (
                      <IntegrationCard
                        key={integration.id}
                        integration={integration}
                        isSelected={selectedIntegration?.id === integration.id}
                        onClick={() => {
                          setSelectedIntegration(integration);
                          setSelectedChart(null);
                        }}
                      />
                    ))
                  )}
                </div>
              </ScrollArea>
            </Card>
          </div>

          {/* Charts and Data View */}
          <div className="lg:col-span-3">
            {selectedChart ? (
              <div>
                <div className="mb-4">
                  <Button
                    variant="ghost"
                    className="gap-2"
                    onClick={() => setSelectedChart(null)}
                  >
                    <RiArrowLeftLine className="h-4 w-4" />
                    Back to Charts
                  </Button>
                </div>
                <ChartDataView 
                  chart={selectedChart} 
                  isSystemIntegration={selectedIntegration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeSystem} 
                  workspaceIntegration={selectedIntegration!}
                />
              </div>
            ) : (
              <Card className="h-[calc(100vh-200px)] overflow-hidden">
                <div className="p-4 border-b flex justify-between items-center">
                  <div>
                    <h4 className="font-medium">
                      {selectedIntegration ? (
                        <>Charts for {selectedIntegration.integration?.integration_name}</>
                      ) : (
                        <>Available Charts</>
                      )}
                    </h4>
                    <p className="text-sm text-muted-foreground mt-1">
                      {selectedIntegration ? (
                        <>Select a chart to view its data points</>
                      ) : (
                        <>Select an integration to view its charts</>
                      )}
                    </p>
                  </div>
                  {selectedIntegration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeSystem && (
                    <CreateChartDialog integration={selectedIntegration} />
                  )}
                </div>
                <ScrollArea className="h-[calc(100vh-280px)]">
                  <div className="p-4">
                    {selectedIntegration ? (
                      isLoadingCharts ? (
                        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                          <Skeleton className="h-32 w-full" />
                          <Skeleton className="h-32 w-full" />
                          <Skeleton className="h-32 w-full" />
                        </div>
                      ) : integrationCharts.length > 0 ? (
                        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                          {integrationCharts.map((chart: MalakIntegrationChart) => (
                            <ChartCard 
                              key={chart.id} 
                              chart={chart}
                              isSelected={false}
                              onClick={() => setSelectedChart(chart)}
                            />
                          ))}
                        </div>
                      ) : (
                        <div className="flex flex-col items-center justify-center py-16 text-center">
                          <div className="rounded-full bg-muted p-4 mb-4">
                            <RiApps2Line className="h-6 w-6 text-muted-foreground" />
                          </div>
                          <h4 className="text-lg font-medium">No charts available</h4>
                          <p className="text-sm text-muted-foreground mt-1">
                            This integration doesn't have any charts configured yet
                          </p>
                        </div>
                      )
                    ) : (
                      <div className="flex flex-col items-center justify-center py-16 text-center">
                        <div className="rounded-full bg-muted p-4 mb-4">
                          <RiApps2Line className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <h4 className="text-lg font-medium">Select an integration</h4>
                        <p className="text-sm text-muted-foreground mt-1">
                          Choose an integration from the sidebar to view its charts
                        </p>
                      </div>
                    )}
                  </div>
                </ScrollArea>
              </Card>
            )}
          </div>
        </div>
      </section>
    </div>
  );
} 