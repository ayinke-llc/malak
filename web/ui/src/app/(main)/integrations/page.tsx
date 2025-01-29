"use client"

import { Button } from "@/components/ui/button";
import { useState } from "react";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Github,
  Slack,
  Twitter,
  Figma,
  NotebookIcon as Notion,
  DropletIcon as Dropbox,
  Trello,
  Linkedin,
  DollarSign,
  ChartBar,
} from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

export default function Integrations() {
  const [integrations, setIntegrations] = useState(initialIntegrations)
  const [apiKey, setApiKey] = useState("")

  const handleConnect = (index: number) => {
    const updatedIntegrations = [...integrations]
    updatedIntegrations[index].status = "Connected"
    setIntegrations(updatedIntegrations)
    setApiKey("") // Clear API key after connecting
  }

  const handleDisconnect = (index: number) => {
    const updatedIntegrations = [...integrations]
    updatedIntegrations[index].status = "Available"
    setIntegrations(updatedIntegrations)
  }

  const getStatusBadge = (status: IntegrationStatus) => {
    switch (status) {
      case "Connected":
        return <Badge>Connected</Badge>
      case "Available":
        return <Badge variant="secondary">Available</Badge>
      case "Coming Soon":
        return <Badge variant="outline">Coming Soon</Badge>
    }
  }

  const getConnectionTypeBadge = (type: ConnectionType) => {
    switch (type) {
      case "OAuth2":
        return <Badge variant="default">OAuth2</Badge>
      case "API Key":
        return <Badge variant="secondary">API Key</Badge>
    }
  }

  return (
    <>
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3
                id="company-decks"
                className="text-lg font-medium"
              >
                Available Integrations
              </h3>
              <p className="text-sm text-muted-foreground">
                View and manage integrations on your workspace
              </p>
            </div>
          </div>
        </section>

        <section className="mt-10">
          <div className="grid gap-6 md:grid-cols-3 lg:grid-cols-6">
            {integrations.map((integration, index) => (
              <Card key={index} className={integration.status === "Coming Soon" ? "opacity-70" : ""}>
                <CardHeader className="flex flex-row items-center gap-4">
                  <integration.icon className="w-8 h-8" />
                  <div>
                    <CardTitle>{integration.name}</CardTitle>
                    <CardDescription>{integration.description}</CardDescription>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-2">
                    {getStatusBadge(integration.status)}
                    {getConnectionTypeBadge(integration.connectionType)}
                  </div>
                </CardContent>
                <CardFooter>
                  {integration.status === "Connected" && (
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button variant="outline" className="w-full">
                          Disconnect
                        </Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Deactivate {integration.name} Integration</DialogTitle>
                          <DialogDescription>
                            Are you sure you want to deactivate the {integration.name} integration? This action will remove
                            all associated data and settings.
                          </DialogDescription>
                        </DialogHeader>
                        <DialogFooter>
                          <Button variant="outline" onClick={() => { }}>
                            Cancel
                          </Button>
                          <Button variant="destructive" onClick={() => handleDisconnect(index)}>
                            Deactivate
                          </Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>
                  )}
                  {integration.status === "Available" && (
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button className="w-full">Connect</Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Connect to {integration.name}</DialogTitle>
                          <DialogDescription>
                            {integration.connectionType === "OAuth2"
                              ? "Click the button below to authorize the connection."
                              : "Enter your API key to connect."}
                          </DialogDescription>
                        </DialogHeader>
                        {integration.connectionType === "OAuth2" ? (
                          <Button onClick={() => handleConnect(index)}>Authorize with {integration.name}</Button>
                        ) : (
                          <div className="grid gap-4 py-4">
                            <div className="grid grid-cols-4 items-center gap-4">
                              <Label htmlFor="api-key" className="text-right">
                                API Key
                              </Label>
                              <div className="col-span-3 space-y-2">
                                <Input
                                  id="api-key"
                                  value={apiKey}
                                  onChange={(e) => setApiKey(e.target.value)}
                                  className="col-span-3"
                                />
                                <p className="text-sm text-muted-foreground">
                                  Your API key is securely encrypted and stored.
                                </p>
                              </div>
                            </div>
                          </div>
                        )}
                        <DialogFooter>
                          {integration.connectionType === "API Key" && (
                            <Button onClick={() => handleConnect(index)}>Connect</Button>
                          )}
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>
                  )}
                  {integration.status === "Coming Soon" && (
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button disabled className="w-full">
                            Coming Soon
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>This integration is not available yet. Stay tuned!</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  )}
                </CardFooter>
              </Card>
            ))}
          </div>
        </section>
      </div>
    </>
  );
}

type IntegrationStatus = "Connected" | "Available" | "Coming Soon"
type ConnectionType = "OAuth2" | "API Key"

interface Integration {
  name: string
  description: string
  icon: React.ElementType
  status: IntegrationStatus
  connectionType: ConnectionType
}

// Mock data for integrations
const initialIntegrations: Integration[] = [
  {
    name: "Brex",
    description: "Connect your Brex accounts and show your investors your up to date balances in a single click",
    icon: DollarSign,
    status: "Connected",
    connectionType: "OAuth2",
  },
  {
    name: "Stripe",
    description: "Connect your financial account and display your balances and payouts",
    icon: DollarSign,
    status: "Available",
    connectionType: "OAuth2",
  },
  {
    name: "Mercury",
    description: "Connect your Mercury accounts and show your investors your up to date balances in a single click",
    icon: DollarSign,
    status: "Available",
    connectionType: "OAuth2",
  },
  {
    name: "Mono",
    description: "Connect your Nigerian bank accounts and show your investors your up to date balances in a single click",
    icon: DollarSign,
    status: "Available",
    connectionType: "OAuth2",
  },
  {
    name: "Google analytics",
    description: "Show your website metrics in one click in your investors' update",
    icon: ChartBar,
    status: "Connected",
    connectionType: "OAuth2",
  },
  {
    name: "Quickbooks",
    description: "Show your receivables in a single click",
    icon: DollarSign,
    status: "Available",
    connectionType: "OAuth2",
  }
]
