import { RiApps2Line, RiArrowRightSLine, RiCheckLine, RiCloseLine } from "@remixicon/react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { MalakWorkspaceIntegration } from "@/client/Api";

export function IntegrationCard({ integration, isSelected, onClick }: { 
  integration: MalakWorkspaceIntegration; 
  isSelected: boolean;
  onClick: () => void;
}) {
  const isDisabled = !integration.is_enabled || !integration.integration?.is_enabled;

  return (
    <Button
      key={integration.id}
      variant="outline"
      className={cn(
        "w-full justify-between p-4 h-auto hover:bg-muted hover:border-muted-foreground/20",
        isDisabled && "opacity-50 cursor-not-allowed",
        isSelected && "bg-muted border-primary hover:border-primary",
        "text-left"
      )}
      onClick={onClick}
      disabled={isDisabled}
    >
      <div className="flex items-center gap-3">
        {integration.integration?.logo_url ? (
          <img 
            src={integration.integration.logo_url} 
            alt={integration.integration.integration_name || "Integration"} 
            className="w-6 h-6 rounded"
          />
        ) : (
          <RiApps2Line className="w-6 h-6" />
        )}
        <div className="text-left">
          <div className="font-medium text-foreground">{integration.integration?.integration_name}</div>
          <div className="text-xs text-muted-foreground truncate max-w-[180px]">
            {integration.integration?.description || "No description available"}
          </div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Badge variant={integration.is_active ? "default" : "destructive"} className="h-5 shrink-0">
          {integration.is_active ? (
            <RiCheckLine className="h-3 w-3 mr-1" />
          ) : (
            <RiCloseLine className="h-3 w-3 mr-1" />
          )}
          {integration.is_active ? "Active" : "Inactive"}
        </Badge>
        <RiArrowRightSLine className={cn(
          "h-4 w-4 transition-transform shrink-0",
          isSelected && "transform rotate-90"
        )} />
      </div>
    </Button>
  );
} 