import { MalakWorkspaceIntegration } from "@/client/Api";
import { IntegrationCard } from "./card";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import { Card } from "@/components/ui/card";
import { RiApps2Line } from "@remixicon/react";

interface IntegrationsListProps {
  integrations: MalakWorkspaceIntegration[];
  isLoading: boolean;
}

export function IntegrationsList({ integrations, isLoading }: IntegrationsListProps) {
  if (isLoading) {
    return <Skeleton count={10} />
  }

  if (!integrations.length) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-background">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="rounded-full bg-muted p-4">
            <RiApps2Line className="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 className="mt-6 text-lg font-medium text-foreground">
            No integrations available
          </h3>
          <p className="mt-2 text-sm text-muted-foreground">
            There are currently no integrations available for your workspace. Check back later for updates.
          </p>
        </div>
      </Card>
    );
  }

  // Sort integrations by integration.is_enabled
  const sortedIntegrations = [...integrations].sort((a, b) => {
    // Convert boolean to number for sorting (true = 1, false = 0)
    const aEnabled = a.integration?.is_enabled ? 1 : 0;
    const bEnabled = b.integration?.is_enabled ? 1 : 0;
    // Sort in descending order (enabled first)
    return bEnabled - aEnabled;
  });

  return (
    <div className="grid gap-6 md:grid-cols-3 lg:grid-cols-6">
      {sortedIntegrations.map((integration, index) => (
        <IntegrationCard
          key={index}
          integration={integration}
        />
      ))}
    </div>
  );
} 
