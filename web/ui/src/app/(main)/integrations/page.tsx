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
import { RiSettings4Line, RiEyeLine, RiEyeOffLine } from "@remixicon/react";
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
import { useForm, Controller } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

const apiKeySchema = yup.object({
  apiKey: yup.string().required("API key is required"),
});

type ApiKeyFormData = yup.InferType<typeof apiKeySchema>;

export default function Integrations() {
  const [apiKeyDialogOpen, setApiKeyDialogOpen] = useState(false);
  const [oauth2SettingsOpen, setOauth2SettingsOpen] = useState(false);
  const [selectedIntegration, setSelectedIntegration] = useState<any>(null);
  const [isTestingConnection, setIsTestingConnection] = useState(false);
  const [isConnectionTested, setIsConnectionTested] = useState(false);
  const [isConnectionValid, setIsConnectionValid] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  const {
    control,
    handleSubmit,
    formState: { errors },
    watch,
    reset,
    setValue,
  } = useForm<ApiKeyFormData>({
    resolver: yupResolver(apiKeySchema),
    defaultValues: {
      apiKey: "",
    },
  });

  const apiKeyValue = watch("apiKey");

  const { data, isLoading, error } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  useEffect(() => {
    if (error) {
      toast.error("Error occurred while fetching communication preferences");
    }
  }, [error]);

  const handleSwitchToggle = (integration: any, checked: boolean) => {
    if (!checked) {
      // Handle disabling integration here if needed
      return;
    }

    if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeApiKey) {
      setSelectedIntegration(integration);
      setIsEditing(false);
      setApiKeyDialogOpen(true);
    } else if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeOauth2) {
      toast.info(`Redirecting you to authenticate with ${integration?.integration?.integration_name}...`);
      setTimeout(() => {
        window.location.href = "https://google.com"; // Replace with actual OAuth URL
      }, 1500); // Give user time to see the toast
    }
  };

  const handleSettingsClick = (integration: any) => {
    if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeApiKey) {
      setSelectedIntegration(integration);
      setIsEditing(true);
      setApiKeyDialogOpen(true);
    } else if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeOauth2) {
      setSelectedIntegration(integration);
      if (!integration?.is_active) {
        toast.info(`Redirecting you to reconnect with ${integration?.integration?.integration_name}...`);
        setTimeout(() => {
          window.location.href = "https://google.com"; // Replace with actual OAuth URL
        }, 1500);
      } else {
        setOauth2SettingsOpen(true);
      }
    }
  };

  const handleTestConnection = async () => {
    try {
      setIsTestingConnection(true);
      // TODO: Implement actual API test connection
      // const response = await client.workspaces.testIntegrationConnection({ 
      //   integration_id: selectedIntegration.integration.id,
      //   api_key: apiKeyValue
      // });
      
      // Simulated delay for now
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setIsConnectionTested(true);
      setIsConnectionValid(true);
      toast.success("Connection test successful!");
    } catch (error) {
      setIsConnectionTested(true);
      setIsConnectionValid(false);
      toast.error("Connection test failed. Please check your API key and try again.");
    } finally {
      setIsTestingConnection(false);
    }
  };

  const handleApiKeySubmit = async (data: ApiKeyFormData) => {
    if (!isConnectionTested || !isConnectionValid) {
      toast.error("Please test the connection first");
      return;
    }

    try {
      // TODO: Implement API key submission
      toast.success(isEditing ? "API key updated successfully" : "API key saved successfully");
      setApiKeyDialogOpen(false);
      resetDialogState();
    } catch (error) {
      toast.error(isEditing ? "Failed to update API key" : "Failed to save API key");
    }
  };

  const resetDialogState = () => {
    reset();
    setIsEditing(false);
    setIsTestingConnection(false);
    setIsConnectionTested(false);
    setIsConnectionValid(false);
    setShowApiKey(false);
  };

  const getConnectionTypeBadge = (type: MalakIntegrationType) => {
    switch (type) {
      case MalakIntegrationType.IntegrationTypeOauth2:
        return <Badge variant="default">OAuth2</Badge>
      case MalakIntegrationType.IntegrationTypeApiKey:
        return <Badge variant="default">API Key</Badge>
    }
  }

  const handleDisconnectOAuth = async () => {
    try {
      // TODO: Implement actual disconnect/revoke
      // await client.workspaces.revokeIntegrationAccess(selectedIntegration.integration.id);
      toast.success(`Disconnected from ${selectedIntegration?.integration?.integration_name}`);
      setOauth2SettingsOpen(false);
    } catch (error) {
      toast.error("Failed to disconnect integration");
    }
  };

  const handleReconnectOAuth = () => {
    setOauth2SettingsOpen(false);
    toast.info(`Redirecting you to reconnect with ${selectedIntegration?.integration?.integration_name}...`);
    setTimeout(() => {
      window.location.href = "https://google.com"; // Replace with actual OAuth URL
    }, 1500);
  };

  return (
    <>
      <Dialog open={oauth2SettingsOpen} onOpenChange={setOauth2SettingsOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>OAuth2 Connection Settings</DialogTitle>
            <DialogDescription>
              Manage your connection with {selectedIntegration?.integration?.integration_name}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <div className="flex items-center gap-2 text-sm">
                <div className={`h-2 w-2 rounded-full ${selectedIntegration?.is_active ? 'bg-green-500' : 'bg-red-500'}`} />
                <span>{selectedIntegration?.is_active ? 'Connected and Active' : 'Connection Inactive'}</span>
              </div>
              <p className="text-sm text-muted-foreground">
                {selectedIntegration?.is_active 
                  ? "Your integration is currently active and working properly. You can disconnect if needed."
                  : "Your integration is currently inactive. You may need to reconnect to restore functionality."
                }
              </p>
            </div>
          </div>
          <DialogFooter className="flex gap-2">
            {!selectedIntegration?.is_active && (
              <Button 
                variant="outline"
                onClick={handleReconnectOAuth}
              >
                Reconnect
              </Button>
            )}
            <Button 
              variant="destructive"
              onClick={handleDisconnectOAuth}
            >
              Disconnect
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={apiKeyDialogOpen} onOpenChange={(open) => {
        setApiKeyDialogOpen(open);
        if (!open) {
          resetDialogState();
        }
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{isEditing ? "Update API Key" : "Enter API Key"}</DialogTitle>
            <DialogDescription>
              {isEditing 
                ? `Update your API key for ${selectedIntegration?.integration?.integration_name}`
                : `Please provide your API key for ${selectedIntegration?.integration?.integration_name}`
              }
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit(handleApiKeySubmit)}>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="api-key">API Key</Label>
                <div className="relative">
                  <Controller
                    name="apiKey"
                    control={control}
                    render={({ field }) => (
                      <Input
                        {...field}
                        id="api-key"
                        type={showApiKey ? "text" : "password"}
                        placeholder={isEditing ? "Enter new API key" : "Enter your API key"}
                        onChange={(e) => {
                          field.onChange(e);
                          setIsConnectionTested(false);
                          setIsConnectionValid(false);
                        }}
                      />
                    )}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                    onClick={() => setShowApiKey(!showApiKey)}
                  >
                    {showApiKey ? (
                      <RiEyeOffLine className="h-4 w-4 text-muted-foreground" />
                    ) : (
                      <RiEyeLine className="h-4 w-4 text-muted-foreground" />
                    )}
                  </Button>
                </div>
                {errors.apiKey && (
                  <p className="text-sm text-red-500">{errors.apiKey.message}</p>
                )}
                {isConnectionTested && (
                  <div className={`text-sm ${isConnectionValid ? 'text-green-500' : 'text-red-500'}`}>
                    {isConnectionValid 
                      ? "✓ Connection verified successfully" 
                      : "✗ Connection test failed"}
                  </div>
                )}
              </div>
            </div>
            <DialogFooter className="flex gap-2">
              <Button 
                type="button"
                variant="outline"
                onClick={handleTestConnection}
                disabled={!apiKeyValue || isTestingConnection}
              >
                {isTestingConnection ? "Testing..." : "Test Connection"}
              </Button>
              <Button 
                type="submit"
                disabled={!isConnectionTested || !isConnectionValid}
              >
                {isEditing ? "Update" : "Save"}
              </Button>
            </DialogFooter>
          </form>
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
                          onCheckedChange={(checked) => handleSwitchToggle(integration, checked)}
                        />
                        {integration?.is_enabled && (
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button 
                                  variant="ghost" 
                                  size="icon" 
                                  className="h-8 w-8"
                                  onClick={() => handleSettingsClick(integration)}
                                >
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
