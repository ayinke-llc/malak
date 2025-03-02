"use client";

import { MalakSystemTemplate, ServerAPIStatus, ServerFetchUpdateReponse } from "@/client/Api";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import client from "@/lib/client";
import { LIST_TEMPLATES } from "@/lib/query-constants";
import { RiArrowLeftLine, RiErrorWarningLine } from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { AxiosError, AxiosResponse } from "axios";
import { format } from "date-fns";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useCreateBlockNote } from "@blocknote/react";
import Markdown from "react-markdown";
import { defaultEditorContent } from "@/components/ui/updates/editor/default-value";
import { Loader2 } from "lucide-react";

function TemplateCard({ template, isLoading, onClick }: {
  template: MalakSystemTemplate;
  isLoading: boolean;
  onClick: () => void;
}) {
  const editor = useCreateBlockNote();
  const [previewText, setPreviewText] = useState<string>('');
  
  useEffect(() => {
    async function parseContent() {
      if (template.content) {
        const markdown = await editor.blocksToMarkdownLossy(template.content as any);
        setPreviewText(markdown.slice(0, 150));
      }
    }
    parseContent();
  }, [editor, template.content]);

  return (
    <Card
      className={`relative cursor-pointer transition-colors ${isLoading ? 'opacity-70 pointer-events-none' : 'hover:bg-accent/5'}`}
      onClick={onClick}
    >
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-background/50 z-50 rounded-lg">
          <Loader2 className="h-6 w-6 animate-spin" />
        </div>
      )}
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
        <div className="max-h-40 overflow-hidden text-sm text-muted-foreground rounded-md bg-muted/50 p-3">
          <div className="prose prose-sm prose-neutral dark:prose-invert max-w-none">
            <Markdown>{previewText + (previewText ? '...' : '')}</Markdown>
          </div>
        </div>
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
  const [activeTemplateId, setActiveTemplateId] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<"malak" | "custom">("malak");

  const { data: templatesData, isLoading: isLoadingTemplates, error } = useQuery({
    queryKey: [LIST_TEMPLATES],
    queryFn: () => client.workspaces.updatesTemplatesList(),
  });

  const templates = activeTab === "malak"
    ? templatesData?.data.templates.system || []
    : templatesData?.data.templates.workspace || [];

  const mutation = useMutation({
    mutationFn: async (templateReference: string) => {
      const template = [...(templatesData?.data.templates.system || []), ...(templatesData?.data.templates.workspace || [])]
        .find(t => t.reference === templateReference);
      if (!template) throw new Error("Template not found");

      const createResp = await client.workspaces.updatesCreate({
        title: `${format(new Date(), "EEEE, MMMM do, yyyy")} Update`,
        template: {
          is_system_template: true,
          reference: template.reference
        }
      });

      const initialContent = defaultEditorContent(createResp.data.update.reference!);
      const templateContent = template.content || [];
      const combinedContent = [...initialContent, ...templateContent];

      await client.workspaces.updateContent(createResp.data.update.reference!, {
        title: createResp.data.update.title!,
        update: combinedContent
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
    onMutate: (templateReference: string) => setActiveTemplateId(templateReference),
    onSettled: () => setActiveTemplateId(null)
  });

  const handleTemplateSelect = (templateReference: string) => {
    toast.promise(mutation.mutateAsync(templateReference), {
      loading: 'Creating update from template...',
      success: 'Update created successfully!',
      error: 'Failed to create update'
    });
  };

  return (
    <div className="container mx-auto py-6">
      <div className="mb-6">
        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="gap-2"
          disabled={!!activeTemplateId}
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
            <TabsTrigger value="malak" disabled={!!activeTemplateId}>Malak Templates</TabsTrigger>
            <TabsTrigger value="custom" disabled={!!activeTemplateId}>My Templates</TabsTrigger>
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
                    isLoading={activeTemplateId === template.reference}
                    onClick={() => handleTemplateSelect(template.reference!)}
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
              <p className="text-sm text-muted-foreground">You&apos;ll be able to create and save your own templates here</p>
            </div>
          </TabsContent>
        </Tabs>
      )}
    </div>
  );
} 
