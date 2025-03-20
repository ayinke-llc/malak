"use client"

import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger
} from "@/components/ui/tabs";
import { useSearchParams } from "next/navigation";
import ApiKeys from "./api-keys";
import Billing from "./billing";
import GeneralSettings from "./general-settings";
import Soon from "./soon";

export default function Settings() {

  const searchParams = useSearchParams();

  const getDefaultValue = (): string => {
    const value = searchParams.get("tab")

    if (!value) {
      return "general"
    }

    if (["general", "billing", "team", "webhook", "api"].includes(value.toLowerCase())) {
      return value.toLowerCase()
    }

    return "general"
  }

  return (
    <>
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3
                id="company-settings"
                className="text-lg font-medium"
              >
                Company Preferences
              </h3>
              <p className="text-sm text-muted-foreground">
                View and manage your company&apos;s preferences
              </p>
            </div>

            <div>
            </div>
          </div>
        </section>

        <section className="mt-10">
          <Tabs defaultValue={getDefaultValue()} className="space-y-6">
            <TabsList className="w-full justify-start border-b pb-px mb-4" variant="vercel">
              <TabsTrigger value="general">General</TabsTrigger>
              <TabsTrigger value="billing">Billing</TabsTrigger>
              <TabsTrigger value="team">Team</TabsTrigger>
              <TabsTrigger value="api">API Key</TabsTrigger>
              <TabsTrigger value="webhooks">Webhooks</TabsTrigger>
            </TabsList>
            <TabsContent value="general">
              <GeneralSettings />
            </TabsContent>
            <TabsContent value="billing">
              <Billing />
            </TabsContent>
            <TabsContent value="team">
              <Soon feature="your team" />
            </TabsContent>
            <TabsContent value="api">
              <ApiKeys />
            </TabsContent>
            <TabsContent value="webhooks">
              <Soon feature="your webhooks" />
            </TabsContent>
          </Tabs>
        </section>
      </div>
    </>
  );
}
