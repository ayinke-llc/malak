"use client"

import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import client from "@/lib/client";
import { LIST_INTEGRATIONS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { RiSettings4Line } from "@remixicon/react";
import { MalakIntegrationType } from "@/client/Api";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export default function Integrations() {
  const [apiKeyDialogOpen, setApiKeyDialogOpen] = useState(false);
  const [selectedIntegration, setSelectedIntegration] = useState<any>(null);
  const [apiKey, setApiKey] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  useEffect(() => {
    if (error) {
      toast.error("Error occurred while fetching communication preferences");
    }
  }, [error]);

  const handleSwitchToggle = (integration: any) => {
    if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeApiKey) {
      setSelectedIntegration(integration);
      setApiKeyDialogOpen(true);
    } else if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeOauth2) {
      toast.info(`Redirecting you to authenticate with ${integration?.integration?.integration_name}...`);
      setTimeout(() => {
        window.location.href = "https://google.com"; // Replace with actual OAuth URL
      }, 1500); // Give user time to see the toast
    }
  };

  const handleApiKeySubmit = async () => {
    try {
      // TODO: Implement API key submission
      toast.success("API key saved successfully");
      setApiKeyDialogOpen(false);
      setApiKey("");
    } catch (error) {
      toast.error("Failed to save API key");
    }
  };

  const getConnectionTypeBadge = (type: MalakIntegrationType) => {
    switch (type) {
      case MalakIntegrationType.IntegrationTypeOauth2:
        return <Badge variant="default">OAuth2</Badge>
      case MalakIntegrationType.IntegrationTypeApiKey:
        return <Badge variant="default">API Key</Badge>
    }
  }

  return (
    <>
      <Dialog open={apiKeyDialogOpen} onOpenChange={setApiKeyDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Enter API Key</DialogTitle>
            <DialogDescription>
              Please provide your API key for {selectedIntegration?.integration?.integration_name}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="api-key">API Key</Label>
              <Input
                id="api-key"
                type="password"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                placeholder="Enter your API key"
              />
            </div>
          </div>
          <DialogFooter>
            <Button onClick={handleApiKeySubmit}>Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

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
          {isLoading ? (
            <Skeleton count={10} />
          ) : (
            <div className="grid gap-6 md:grid-cols-3 lg:grid-cols-6">
              {data?.data?.integrations.map((integration, index) => (
                <Card key={index} className="drop-shadow-lg shadow-lg flex flex-col h-[280px]">
                  <CardHeader className="flex-shrink-0">
                    <div className="flex items-center gap-4">
                      <img
                        className="w-8 h-8 rounded-md"
                        src={integration?.integration?.logo_url}
                        alt={`${integration?.integration?.integration_name} logo`}
                      />
                      <CardTitle>{integration?.integration?.integration_name}</CardTitle>
                    </div>
                  </CardHeader>
                  <CardContent className="flex-grow flex flex-col justify-between">
                    <CardDescription className="mt-2 h-12 overflow-hidden text-ellipsis">
                      {integration?.integration?.description}
                    </CardDescription>
                    <div className="flex justify-between items-center mt-4">
                      <div className="flex gap-2">
                        {getConnectionTypeBadge(integration?.integration?.integration_type as MalakIntegrationType)}
                      </div>
                      <div className="flex items-center gap-2">
                        {integration?.is_enabled && (
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger>
                                <div className={`h-2 w-2 rounded-full ${integration?.is_active ? 'bg-green-500' : 'bg-red-500'}`} />
                              </TooltipTrigger>
                              <TooltipContent>
                                <p>{integration?.is_active ? 'Connected' : 'Connection Failed'}</p>
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        )}
                        <Switch 
                          checked={integration?.is_enabled} 
                          disabled={!integration?.integration?.is_enabled}
                          onCheckedChange={() => handleSwitchToggle(integration)}
                        />
                        {integration?.is_enabled && (
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button variant="ghost" size="icon" className="h-8 w-8">
                                  <RiSettings4Line className="h-4 w-4" />
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p>Integration Settings</p>
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        )}
                      </div>
                    </div>
                  </CardContent>
                  <CardFooter className="flex-shrink-0 flex justify-end">
                  </CardFooter>
                </Card>
              ))}
            </div>
          )}
        </section>
      </div>
    </>
  );
}
