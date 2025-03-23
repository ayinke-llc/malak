"use client"

import { useEffect, useState } from "react";
import client from "@/lib/client";
import { LIST_INTEGRATIONS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { IntegrationsList } from "@/components/ui/integrations/list";
import { Card } from "@/components/ui/card";
import { RiErrorWarningLine, RiInformationLine } from "@remixicon/react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { RiCloseLine } from "@remixicon/react";

export default function Integrations() {
  const [showBanner, setShowBanner] = useState(true);
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [LIST_INTEGRATIONS],
    queryFn: () => client.workspaces.integrationsList(),
  });

  useEffect(() => {
    if (error) {
      toast.error("Error occurred while fetching integrations");
    }
  }, [error]);

  const InfoBanner = () => (
    <Alert variant="default" className="mb-6 border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-900/20">
      <div className="flex items-center justify-between">
        <div className="flex gap-3">
          <RiInformationLine className="h-5 w-5 text-blue-600 dark:text-blue-400" />
          <div>
            <AlertTitle className="text-blue-800 dark:text-blue-300">Security Information</AlertTitle>
            <AlertDescription className="text-sm text-blue-700 dark:text-blue-400">
              Learn about our secure secrets storage in our{" "}
              <a
                href="https://docs.malak.vc/self-hosting/secrets"
                target="_blank"
                rel="noopener noreferrer"
                className="font-medium underline underline-offset-4 hover:text-blue-800 dark:hover:text-blue-300"
              >
                documentation
              </a>
              .
            </AlertDescription>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-6 w-6 text-blue-700 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
          onClick={() => setShowBanner(false)}
        >
          <RiCloseLine className="h-4 w-4" />
          <span className="sr-only">Dismiss</span>
        </Button>
      </div>
    </Alert>
  );

  if (error) {
    return (
      <div className="pt-6 bg-background">
        {showBanner && <InfoBanner />}
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
                onClick={() => refetch()}
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
    <div className="pt-6 bg-background">
      {showBanner && <InfoBanner />}
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
        <IntegrationsList
          integrations={data?.data?.integrations || []}
          isLoading={isLoading}
        />
      </section>
    </div>
  );
} 