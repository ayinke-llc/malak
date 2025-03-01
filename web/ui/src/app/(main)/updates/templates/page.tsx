"use client";

import { ServerAPIStatus, ServerFetchUpdateReponse, MalakBlock, MalakSystemTemplate } from "@/client/Api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Alert, AlertDescription } from "@/components/ui/alert";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import client from "@/lib/client";
import { RiArrowLeftLine, RiErrorWarningLine } from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { AxiosError, AxiosResponse } from "axios";
import { format } from "date-fns";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";
import { LIST_TEMPLATES } from "@/lib/query-constants";

function TemplateCard({ template, isLoading, onClick }: {
  template: MalakSystemTemplate;
  isLoading: boolean;
  onClick: () => void;
}) {
  return (
    <Card
      className={`cursor-pointer transition-colors ${isLoading ? 'opacity-50 pointer-events-none' : 'hover:bg-accent/5'}`}
      onClick={onClick}
    >
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>{template.title}</CardTitle>
          <Badge variant="secondary" className="ml-2">
            Used {template.number_of_uses || 0} times
          </Badge>
        </div>
        <CardDescription>{template.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <pre className="max-h-40 overflow-hidden text-sm text-muted-foreground">
          {template.content?.[0]?.content?.slice(0, 150)}...
        </pre>
      </CardContent>
    </Card>
  );
}

function TemplateCardSkeleton() {
  return (
    <Card className="cursor-pointer">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="w-32">
            <Skeleton count={1} />
          </div>
          <div className="w-24 ml-2">
            <Skeleton count={1} />
          </div>
        </div>
        <div className="mt-2 w-full">
          <Skeleton count={1} />
        </div>
      </CardHeader>
      <CardContent>
        <div className="w-full">
          <Skeleton count={4} />
        </div>
      </CardContent>
    </Card>
  );
}

export default function TemplatesPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<"malak" | "custom">("malak");

  const { data: templatesData, isLoading: isLoadingTemplates, error } = useQuery({
    queryKey: [LIST_TEMPLATES],
    queryFn: () => client.workspaces.updatesTemplatesList(),
  });

  const templates = activeTab === "malak"
    ? templatesData?.data.templates.system || []
    : templatesData?.data.templates.workspace || [];

  const mutation = useMutation({
    mutationFn: async (templateId: string) => {
      const template = [...(templatesData?.data.templates.system || []), ...(templatesData?.data.templates.workspace || [])]
        .find(t => t.id === templateId);
      if (!template) throw new Error("Template not found");

      const createResp = await client.workspaces.updatesCreate({
        title: `${format(new Date(), "EEEE, MMMM do, yyyy")} Update`,
      });

      await client.workspaces.updateContent(createResp.data.update.reference!, {
        title: createResp.data.update.title!,
        update: template.content || []
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

      {error ? (
        <Alert variant="destructive">
          <RiErrorWarningLine className="h-4 w-4" />
          <AlertDescription>
            Failed to load templates. Please try again later.
          </AlertDescription>
        </Alert>
      ) : (
        <Tabs defaultValue="malak" className="mb-8" onValueChange={(value) => setActiveTab(value as "malak" | "custom")}>
          <TabsList>
            <TabsTrigger value="malak">Malak Templates</TabsTrigger>
            <TabsTrigger value="custom">My Templates</TabsTrigger>
          </TabsList>
          <TabsContent value="malak">
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {isLoadingTemplates ? (
                <>
                  <TemplateCardSkeleton />
                  <TemplateCardSkeleton />
                  <TemplateCardSkeleton />
                </>
              ) : templates.length > 0 ? (
                templates.map((template: MalakSystemTemplate) => (
                  <TemplateCard
                    key={template.id}
                    template={template}
                    isLoading={isLoading}
                    onClick={() => handleTemplateSelect(template.id!)}
                  />
                ))
              ) : (
                <div className="col-span-full text-center py-12">
                  <p className="text-muted-foreground">No templates available</p>
                </div>
              )}
            </div>
          </TabsContent>
          <TabsContent value="custom">
            <div className="text-center py-12">
              <p className="text-muted-foreground mb-2">Custom templates coming soon</p>
              <p className="text-sm text-muted-foreground">You'll be able to create and save your own templates here</p>
            </div>
          </TabsContent>
        </Tabs>
      )}
    </div>
  );
} 
