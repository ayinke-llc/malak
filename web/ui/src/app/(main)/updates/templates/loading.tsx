import { Card, CardContent, CardHeader } from "@/components/ui/card";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import { Button } from "@/components/ui/button";
import { RiArrowLeftLine } from "@remixicon/react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

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

export default function Loading() {
  return (
    <div className="container mx-auto py-6">
      <div className="mb-6">
        <Button
          variant="ghost"
          className="gap-2"
          disabled
        >
          <RiArrowLeftLine />
          Back to Updates
        </Button>
      </div>

      <div className="mb-8">
        <h1 className="text-2xl font-bold">Choose a Template</h1>
        <p className="text-muted-foreground">Select a template to start your update</p>
      </div>

      <Tabs defaultValue="malak" className="mb-8">
        <TabsList>
          <TabsTrigger value="malak" disabled>Malak Templates</TabsTrigger>
          <TabsTrigger value="custom" disabled>My Templates</TabsTrigger>
        </TabsList>
        <TabsContent value="malak">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <TemplateCardSkeleton />
            <TemplateCardSkeleton />
            <TemplateCardSkeleton />
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
} 