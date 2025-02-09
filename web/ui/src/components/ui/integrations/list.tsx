import { MalakWorkspaceIntegration } from "@/client/Api";
import { IntegrationCard } from "./card";
import Skeleton from "@/components/ui/custom/loader/skeleton";

interface IntegrationsListProps {
  integrations: MalakWorkspaceIntegration[];
  isLoading: boolean;
}

export function IntegrationsList({ integrations, isLoading }: IntegrationsListProps) {
  if (isLoading) {
    return <Skeleton count={10} />
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
