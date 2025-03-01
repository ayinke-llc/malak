"use client";

import { ServerAPIStatus, ServerFetchUpdateReponse, MalakBlock } from "@/client/Api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import client from "@/lib/client";
import { RiArrowLeftLine } from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import { AxiosError, AxiosResponse } from "axios";
import { format } from "date-fns";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

// This would come from your API
const MALAK_TEMPLATES = [
  {
    id: "weekly-update",
    title: "Weekly Update",
    description: "A template for sharing weekly progress, challenges, and next steps.",
    usageCount: 124, // This would come from the API
    content: `# Weekly Update

## Accomplishments
- 
- 
- 

## Challenges
- 
- 

## Next Steps
- 
- 
`
  },
  {
    id: "project-milestone",
    title: "Project Milestone",
    description: "A template for announcing project milestones and achievements.",
    usageCount: 89, // This would come from the API
    content: `# Project Milestone Update

## Milestone Achievement
- 

## Impact
- 

## Next Phase
- 
`
  },
  {
    id: "status-report",
    title: "Status Report",
    description: "A template for detailed status reports with metrics and timelines.",
    usageCount: 56, // This would come from the API
    content: `# Status Report

## Project Status
- Overall Status: [On Track/At Risk/Blocked]
- Timeline: [On Schedule/Delayed]
- Budget: [Within Budget/Over Budget]

## Key Metrics
- 
- 

## Action Items
- 
- 
`
  }
];

export default function TemplatesPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<"malak" | "custom">("malak");

  const templates = MALAK_TEMPLATES;

  const mutation = useMutation({
    mutationFn: async (templateId: string) => {
      const template = MALAK_TEMPLATES.find(t => t.id === templateId);
      if (!template) throw new Error("Template not found");

      // First create the update
      const createResp = await client.workspaces.updatesCreate({
        title: `${format(new Date(), "EEEE, MMMM do, yyyy")} Update`,
      });

      // Then update its content
      await client.workspaces.updateContent(createResp.data.update.reference!, {
        title: createResp.data.update.title!,
        update: [
          {
            type: "markdown",
            content: template.content
          }
        ] as MalakBlock[]
      });

      return createResp;
    },
    gcTime: 0,
    onError: (err: AxiosError<ServerAPIStatus>): void => {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    onSuccess: (resp: AxiosResponse<ServerFetchUpdateReponse>) => {
      router.push(`/updates/${resp.data.update.reference}`);
    },
    onMutate: () => setIsLoading(true),
    onSettled: () => setIsLoading(false)
  });

  const handleTemplateSelect = (templateId: string) => {
    mutation.mutate(templateId);
  };

  return (
    <div className="container mx-auto py-6">
      <div className="mb-6">
        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="gap-2"
          disabled={isLoading}
        >
          <RiArrowLeftLine />
          Back to Updates
        </Button>
      </div>

      <div className="mb-8">
        <h1 className="text-2xl font-bold">Choose a Template</h1>
        <p className="text-muted-foreground">Select a template to start your update</p>
      </div>

      <Tabs defaultValue="malak" className="mb-8" onValueChange={(value) => setActiveTab(value as "malak" | "custom")}>
        <TabsList>
          <TabsTrigger value="malak">Malak Templates</TabsTrigger>
          <TabsTrigger value="custom">My Templates</TabsTrigger>
        </TabsList>
        <TabsContent value="malak">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {templates.map((template) => (
              <Card 
                key={template.id} 
                className={`cursor-pointer transition-colors ${isLoading ? 'opacity-50 pointer-events-none' : 'hover:bg-accent/5'}`}
                onClick={() => handleTemplateSelect(template.id)}
              >
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle>{template.title}</CardTitle>
                    <Badge variant="secondary" className="ml-2">
                      Used {template.usageCount} times
                    </Badge>
                  </div>
                  <CardDescription>{template.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <pre className="max-h-40 overflow-hidden text-sm text-muted-foreground">
                    {template.content.slice(0, 150)}...
                  </pre>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
        <TabsContent value="custom">
          <div className="text-center py-12">
            <p className="text-muted-foreground mb-2">Custom templates coming soon</p>
            <p className="text-sm text-muted-foreground">You'll be able to create and save your own templates here</p>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
} 