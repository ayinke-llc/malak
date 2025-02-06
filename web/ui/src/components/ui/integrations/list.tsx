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

  return (
    <div className="grid gap-6 md:grid-cols-3 lg:grid-cols-6">
      {integrations?.map((integration, index) => (
        <IntegrationCard
          key={index}
          integration={integration}
        />
      ))}
    </div>
  );
} 
