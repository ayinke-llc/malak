/* eslint-disable @next/next/no-img-element */
"use client"

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { RiSettings4Line } from "@remixicon/react";
import { MalakIntegrationType, MalakWorkspaceIntegration, ServerAPIStatus } from "@/client/Api";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { toast } from "sonner";
import { useState } from "react";
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
import { RiEyeLine, RiEyeOffLine } from "@remixicon/react";
import { useForm, Controller } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { DISABLE_INTEGRATION, ENABLE_INTEGRATION, LIST_INTEGRATIONS, PING_INTEGRATION, UPDATE_INTEGRATION_SETTINGS } from "@/lib/query-constants";
import client from "@/lib/client";
import { AxiosError } from "axios";

interface IntegrationCardProps {
  integration: MalakWorkspaceIntegration;
}

const apiKeySchema = yup.object({
  apiKey: yup.string().required("API key is required"),
});

type ApiKeyFormData = yup.InferType<typeof apiKeySchema>;

const getConnectionTypeBadge = (type: MalakIntegrationType) => {
  switch (type) {
    case MalakIntegrationType.IntegrationTypeOauth2:
      return <Badge variant="default">OAuth2</Badge>
    case MalakIntegrationType.IntegrationTypeApiKey:
      return <Badge variant="default">API Key</Badge>
  }
}

const getStatusBadge = (integration: MalakWorkspaceIntegration) => {
  if (!integration?.integration?.is_enabled) {
    return <Badge variant="secondary">Coming Soon</Badge>
  }
  return null;
}

export function IntegrationCard({ integration }: IntegrationCardProps) {

  const [apiKeyDialogOpen, setApiKeyDialogOpen] = useState(false);
  const [oauth2SettingsOpen, setOauth2SettingsOpen] = useState(false);
  const [isTestingConnection, setIsTestingConnection] = useState(false);
  const [isConnectionTested, setIsConnectionTested] = useState(false);
  const [isConnectionValid, setIsConnectionValid] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [localEnabled, setLocalEnabled] = useState(integration?.is_enabled || false);

  const queryClient = useQueryClient();

  const {
    control,
    handleSubmit,
    formState: { errors },
    watch,
    reset,
  } = useForm<ApiKeyFormData>({
    resolver: yupResolver(apiKeySchema),
    defaultValues: {
      apiKey: "",
    },
  });

  const mutation = useMutation({
    mutationKey: [PING_INTEGRATION],
    mutationFn: (data: ApiKeyFormData) => {
      setIsTestingConnection(true);
      return client.workspaces.integrationsPingCreate(integration.reference as string, {
        api_key: data.apiKey,
      })
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setIsConnectionValid(true);
      setIsConnectionTested(true);
      setIsTestingConnection(false);
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message || "An error occurred while pinging this integration");
      setIsConnectionValid(false);
      setIsConnectionTested(true);
      setIsTestingConnection(false);
    },
  });

  const enableMutation = useMutation({
    mutationKey: [ENABLE_INTEGRATION],
    mutationFn: (data: ApiKeyFormData) => {
      return client.workspaces.integrationsCreate(integration?.reference as string, {
        api_key: data.apiKey,
      })
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setApiKeyDialogOpen(false);
      resetDialogState();
      setLocalEnabled(true);
      queryClient.invalidateQueries({ queryKey: [LIST_INTEGRATIONS] });
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message || "An error occurred while enabling this integration");
    },
  });

  const disableMutation = useMutation({
    mutationKey: [DISABLE_INTEGRATION],
    mutationFn: () => {
      return client.workspaces.integrationsDelete(integration?.reference as string)
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setApiKeyDialogOpen(false);
      resetDialogState();
      setLocalEnabled(false);
      queryClient.invalidateQueries({ queryKey: [LIST_INTEGRATIONS] });
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message || "An error occurred while disabling this integration");
      setLocalEnabled(true);
    },
  });

  const updateMutationSettings = useMutation({
    mutationKey: [UPDATE_INTEGRATION_SETTINGS],
    mutationFn: (data: ApiKeyFormData) => {
      return client.workspaces.integrationsUpdate(integration?.reference as string, {
        api_key: data.apiKey,
      })
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setApiKeyDialogOpen(false);
      resetDialogState();
      setLocalEnabled(true);
      queryClient.invalidateQueries({ queryKey: [LIST_INTEGRATIONS] });
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message || "An error occurred while updating this integration");
    },
  });

  const apiKeyValue = watch("apiKey");

  const handleToggleIntegration = (checked: boolean) => {
    if (!checked) {
      disableMutation.mutate();
      return;
    }

    setLocalEnabled(checked);
    if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeApiKey) {
      setIsEditing(false);
      setApiKeyDialogOpen(true);
    } else if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeOauth2) {
      handleOAuth2Redirect();
    }
  };

  const handleSettingsClick = () => {
    if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeApiKey) {
      setIsEditing(true);
      setApiKeyDialogOpen(true);
    } else if (integration?.integration?.integration_type === MalakIntegrationType.IntegrationTypeOauth2) {
      if (!integration?.is_active) {
        handleOAuth2Redirect();
      } else {
        setOauth2SettingsOpen(true);
      }
    }
  };

  const handleOAuth2Redirect = () => {
    toast.info(`Redirecting you to authenticate with ${integration?.integration?.integration_name}...`);
    setTimeout(() => {
      window.location.href = "https://google.com"; // Replace with actual OAuth URL
    }, 1500);
  };

  const handleTestConnection = async () => {
    mutation.mutate({ apiKey: apiKeyValue });
  };

  const handleApiKeySubmit = async (data: ApiKeyFormData) => {
    if (!isConnectionTested || !isConnectionValid) {
      toast.error("Please test the connection first");
      return;
    }

    if (isEditing) {
      updateMutationSettings.mutate(data);
    } else {
      enableMutation.mutate(data);
    }
  };

  const resetDialogState = () => {
    reset();
    setIsEditing(false);
    setIsTestingConnection(false);
    setIsConnectionTested(false);
    setIsConnectionValid(false);
    setShowApiKey(false);
    setLocalEnabled(integration?.is_enabled || false);
  };

  const handleDisconnectOAuth = async () => {
    try {
      toast.success(`Disconnected from ${integration?.integration?.integration_name}`);
      setOauth2SettingsOpen(false);
    } catch (error) {
      toast.error("Failed to disconnect integration");
    }
  };

  return (
    <>
      <Dialog open={oauth2SettingsOpen} onOpenChange={setOauth2SettingsOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>OAuth2 Connection Settings</DialogTitle>
            <DialogDescription>
              Manage your connection with {integration?.integration?.integration_name}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <div className="flex items-center gap-2 text-sm">
                <div className={`h-2 w-2 rounded-full ${integration?.is_active ? 'bg-green-500' : 'bg-red-500'}`} />
                <span>{integration?.is_active ? 'Connected and Active' : 'Connection Inactive'}</span>
              </div>
              <p className="text-sm text-muted-foreground">
                {integration?.is_active
                  ? "Your integration is currently active and working properly. You can disconnect if needed."
                  : "Your integration is currently inactive. You may need to reconnect to restore functionality."
                }
              </p>
            </div>
          </div>
          <DialogFooter className="flex gap-2">
            {!integration?.is_active && (
              <Button
                variant="outline"
                onClick={handleOAuth2Redirect}
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
          queryClient.invalidateQueries({ queryKey: [LIST_INTEGRATIONS] });
        }
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{isEditing ? "Update API Key" : "Enter API Key"}</DialogTitle>
            <DialogDescription>
              {isEditing
                ? `Update your API key for ${integration?.integration?.integration_name}`
                : `Please provide your API key for ${integration?.integration?.integration_name}`
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

      <Card className="drop-shadow-lg shadow-lg flex flex-col h-[280px]">
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
              {getStatusBadge(integration)}
            </div>
            <div className="flex items-center gap-2">
              {integration?.is_enabled && integration?.integration?.is_enabled && (
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
                checked={localEnabled}
                disabled={!integration?.integration?.is_enabled}
                onCheckedChange={handleToggleIntegration}
              />
              {integration?.is_enabled && integration?.integration?.is_enabled && (
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={handleSettingsClick}
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
    </>
  );
}
