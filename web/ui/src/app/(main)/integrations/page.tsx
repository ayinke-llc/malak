"use client"

import { useEffect } from "react";
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


  const getStatusBadge = (status: string) => {
    switch (status) {
      case "Connected":
        return <Badge>Connected</Badge>
      case "Available":
        return <Badge variant="secondary">Available</Badge>
      case "Coming Soon":
        return <Badge variant="outline">Coming Soon</Badge>
    }
  }

  const getConnectionTypeBadge = (type: string) => {
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
                id="integrations"
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
          {isLoading ? <Skeleton count={10} /> : (
            <div className="grid gap-6 md:grid-cols-3 lg:grid-cols-6">
              {data?.data?.integrations.map((integration, index) => (
                <Card key={index} className="drop-shadow-lg shadow-lg">
                  <CardHeader className="flex flex-row items-center gap-4">
                    <img className="w-8 h-8 rounded-md"
                      src={integration?.integration?.logo_url} />
                    <div>
                      <CardTitle>{integration?.integration?.integration_name}</CardTitle>
                      <CardDescription className="mt-2">{integration?.integration?.description}</CardDescription>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex gap-2">
                      {getStatusBadge("")}
                      {getConnectionTypeBadge("")}
                    </div>
                  </CardContent>
                  <CardFooter className="flex justify-end items-end">
                    <Switch
                      checked={integration?.is_enabled}
                      disabled={!integration?.integration?.is_enabled} />

                    {integration?.is_enabled && <Button variant="ghost" size="icon">
                      <RiSettings4Line className="h-4 w-4" />
                    </Button>}
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
