import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { RiAddLine, RiCloseLine, RiArchiveFill } from "@remixicon/react";
import { Card } from "../types";
import { ActivityList } from "./tabs/activity/ActivityList";
import { DocumentsTab } from "./tabs/documents/DocumentsTab";
import { EditInvestorDialog } from "./dialogs/EditInvestorDialog";
import { Activity } from "../types";
import { OverviewTab } from "./tabs/overview/OverviewTab";

interface InvestorDetailsDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  isArchived?: boolean;
}

export function InvestorDetailsDrawer({
  open,
  onOpenChange,
  investor,
  isArchived = false,
}: InvestorDetailsDrawerProps) {
  const [activeTab, setActiveTab] = useState("overview");
  const [activities, setActivities] = useState<Activity[]>([]);
  const [isAddingActivity, setIsAddingActivity] = useState(false);
  const [isEditingInvestor, setIsEditingInvestor] = useState(false);

  // Infinite scroll states
  const [page, setPage] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const handleAddActivity = (newActivity: Partial<Activity>) => {
    const activity: Activity = {
      ...newActivity,
      id: Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: newActivity.type as Activity['type'],
      title: newActivity.title || '',
      description: newActivity.description || ''
    };

    setActivities(prev => [activity, ...prev]);
  };

  const handleSaveInvestor = (updatedInvestor: Card) => {
    // In a real implementation, you would update the investor in your backend here
    onOpenChange(false);
  };

  if (!investor) return null;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        className="w-full max-w-[90%] lg:max-w-[75%] 2xl:max-w-[1400px] overflow-y-auto"
      >
        <div className="h-full flex flex-col">
          <SheetHeader className="flex-none">
            <div className="flex items-center justify-between">
              <SheetTitle className="text-xl">{investor.title}</SheetTitle>
              <div className="flex items-center gap-2">
                {!isArchived && (
                  <Button variant="ghost" size="icon" onClick={() => setIsAddingActivity(true)}>
                    <RiAddLine className="w-4 h-4" />
                  </Button>
                )}
                <Button variant="ghost" size="icon" onClick={() => onOpenChange(false)}>
                  <RiCloseLine className="w-4 h-4" />
                </Button>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div className="space-x-2">
                <Badge variant="outline">{investor.stage}</Badge>
                <Badge variant="secondary" className="bg-primary/10 text-primary">
                  ${investor.amount}
                </Badge>
                {isArchived && (
                  <Badge variant="outline" className="text-muted-foreground">
                    <RiArchiveFill className="w-3 h-3 mr-1 inline" />
                    Archived
                  </Badge>
                )}
              </div>
            </div>
          </SheetHeader>

          <div className="mt-6">
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList className="w-full">
                <TabsTrigger value="overview" className="flex-1">Overview</TabsTrigger>
                <TabsTrigger value="activity" className="flex-1">Activity</TabsTrigger>
                <TabsTrigger value="documents" className="flex-1">Documents</TabsTrigger>
              </TabsList>

              <TabsContent value="overview" className="mt-6">
                <OverviewTab
                  investor={investor}
                  onAddActivity={() => setActiveTab("activity")}
                  onUploadDocument={() => setActiveTab("documents")}
                />
              </TabsContent>

              <TabsContent value="activity">
                <ActivityList
                  activities={activities}
                  isLoading={isLoading}
                  hasMore={hasMore}
                  isArchived={isArchived}
                  onAddActivity={handleAddActivity}
                />
              </TabsContent>

              <TabsContent value="documents">
                <DocumentsTab isReadOnly={isArchived} />
              </TabsContent>
            </Tabs>
          </div>

          {!isArchived && (
            <>
              <EditInvestorDialog
                open={isEditingInvestor}
                onOpenChange={setIsEditingInvestor}
                investor={investor}
                onSave={handleSaveInvestor}
              />
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
} 
