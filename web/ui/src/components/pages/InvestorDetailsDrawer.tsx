import { useState, useEffect } from "react";
import type { MalakContact, MalakFundraiseContactDealDetails } from "@/client/Api";
import { fullName } from "@/lib/custom";
import {
  Sheet,
  SheetContent,
  SheetHeader
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle, DialogFooter
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  RiMailLine, RiCalendarLine,
  RiFileTextLine, RiStarFill,
  RiStarLine,
  RiArchiveFill, RiEditLine
} from "@remixicon/react";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { NumericFormat } from "react-number-format";
import { format, fromUnixTime, isValid, parseISO } from "date-fns";
import type { Card, Activity } from "@/components/investor-pipeline/types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { FETCH_FUNDRAISING_PIPELINE, UPDATE_INVESTOR_IN_PIPELINE } from "@/lib/query-constants";
import client from "@/lib/client";
import { toast } from "sonner";
import type { AxiosError } from "axios";
import type { ServerAPIStatus } from "@/client/Api";
import { ActivityList } from "../investor-pipeline/details/tabs/activity/ActivityList";
import Copy from "../ui/custom/copy";
import { DocumentsTab } from "../investor-pipeline/details/tabs/documents/DocumentsTab";

interface InvestorDetailsDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  isArchived?: boolean;
  contact?: MalakContact;
  deal?: MalakFundraiseContactDealDetails;
  slug: string;
}

function AddActivityDialog({
  open,
  onOpenChange,
  onSubmit
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (activity: Partial<Activity>) => void;
}) {
  const [activity, setActivity] = useState<Partial<Activity>>({
    type: 'email',
    title: '',
    description: '',
    content: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(activity);
    onOpenChange(false);
    setActivity({
      type: 'email',
      title: '',
      description: '',
      content: ''
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Add New Activity</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Type</label>
            <Select
              value={activity.type}
              onValueChange={(value) => setActivity({ ...activity, type: value as Activity['type'] })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select activity type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="email">Email</SelectItem>
                <SelectItem value="meeting">Meeting</SelectItem>
                <SelectItem value="document">Document</SelectItem>
                <SelectItem value="team">Team</SelectItem>
                <SelectItem value="stage_change">Stage Change</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Title</label>
            <Input
              value={activity.title}
              onChange={(e) => setActivity({ ...activity, title: e.target.value })}
              placeholder="Activity title"
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <Input
              value={activity.description}
              onChange={(e) => setActivity({ ...activity, description: e.target.value })}
              placeholder="Brief description"
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Content (optional)</label>
            <Textarea
              value={activity.content}
              onChange={(e) => setActivity({ ...activity, content: e.target.value })}
              placeholder="Additional details or content"
            />
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Add Activity</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}



function EditInvestorDialog({
  open,
  onOpenChange,
  investor,
  onSave,
  deal,
  isLoading
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  onSave: (updatedDeal: Partial<MalakFundraiseContactDealDetails>) => void;
  deal?: MalakFundraiseContactDealDetails;
  isLoading: boolean;
}) {
  const [editedDeal, setEditedDeal] = useState<Partial<MalakFundraiseContactDealDetails>>({
    check_size: deal?.check_size ?? 0,
    can_lead_round: deal?.can_lead_round ?? false,
    rating: deal?.rating ?? 0
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(editedDeal);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Investor Details</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Check Size</Label>
            <NumericFormat
              value={editedDeal.check_size ? (editedDeal.check_size / 100).toLocaleString() : ""}
              onValueChange={(values) => {
                const { value } = values;
                setEditedDeal({ ...editedDeal, check_size: Number(value) * 100 });
              }}
              thousandSeparator
              prefix="$"
              customInput={Input}
              placeholder="Enter check size"
              className="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
          </div>

          <div className="space-y-2">
            <Label>Can Lead Round</Label>
            <div className="flex items-center space-x-2">
              <Switch
                checked={editedDeal.can_lead_round}
                onCheckedChange={(checked) => setEditedDeal({ ...editedDeal, can_lead_round: checked })}
              />
              <span className="text-sm text-muted-foreground">
                {editedDeal.can_lead_round ? 'Yes' : 'No'}
              </span>
            </div>
          </div>

          <div className="space-y-2">
            <Label>Rating</Label>
            <div className="flex items-center gap-1">
              {[1, 2, 3, 4, 5].map((star) => (
                <Button
                  key={star}
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="hover:bg-transparent"
                  onClick={() => setEditedDeal({ ...editedDeal, rating: star })}
                >
                  {star <= (editedDeal.rating || 0) ? (
                    <RiStarFill className="w-6 h-6 text-yellow-400" />
                  ) : (
                    <RiStarLine className="w-6 h-6 text-muted-foreground" />
                  )}
                </Button>
              ))}
              <span className="ml-2 text-sm text-muted-foreground">
                {editedDeal.rating} of 5
              </span>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" loading={isLoading}>
              Save Changes
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

const formatSafeDate = (dateValue: string | number | undefined | null, formatStr: string = 'MMM d, yyyy') => {
  if (!dateValue) return "Not set";
  try {
    // Handle different date formats
    const date = typeof dateValue === 'number'
      ? fromUnixTime(dateValue) // Unix timestamp in seconds
      : parseISO(dateValue); // ISO string

    return isValid(date) ? format(date, formatStr) : "Invalid date";
  } catch (e) {
    return "Invalid date";
  }
};

export function InvestorDetailsDrawer({
  open,
  onOpenChange,
  investor,
  isArchived = false,
  contact,
  deal,
  slug
}: InvestorDetailsDrawerProps) {
  const [activeTab, setActiveTab] = useState("overview");
  const [activities, setActivities] = useState<Activity[]>([]);
  const [isAddingActivity, setIsAddingActivity] = useState(false);
  const [isEditingInvestor, setIsEditingInvestor] = useState(false);
  const queryClient = useQueryClient();

  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const updateInvestorMutation = useMutation({
    mutationKey: [UPDATE_INVESTOR_IN_PIPELINE, slug],
    mutationFn: async (updatedDeal: Partial<MalakFundraiseContactDealDetails>) => {
      if (!contact?.reference) {
        throw new Error("No contact reference provided")
      }

      const response = await client.pipelines.contactsPartialUpdate(slug, investor?.dataID as string, {
        check_size: updatedDeal.check_size ?? 0,
        can_lead_round: updatedDeal.can_lead_round ?? false,
        rating: updatedDeal.rating ?? 0
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FETCH_FUNDRAISING_PIPELINE, slug] });
      toast.success("Investor details updated successfully");
      setIsEditingInvestor(false);
      onOpenChange(false);
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err.response?.data.message ?? "Failed to update investor details");
    },
  });

  // reset states when investor changes
  useEffect(() => {
    if (investor) {
      setActivities([]);
      setPage(0);
      setHasMore(true);
    }
  }, [investor?.id]);

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

  const handleSaveInvestor = (updatedDeal: Partial<MalakFundraiseContactDealDetails>) => {
    updateInvestorMutation.mutate(updatedDeal);
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
              <div className="space-x-2">
                <Badge variant="outline">{investor.stage}</Badge>
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
            <Tabs value={activeTab} onValueChange={setActiveTab} >
              <TabsList className="w-full">
                <TabsTrigger value="overview" className="flex-1">Overview</TabsTrigger>
                <TabsTrigger value="activity" className="flex-1">Activities</TabsTrigger>
                <TabsTrigger value="documents" className="flex-1">Documents</TabsTrigger>
              </TabsList>

              <TabsContent value="overview" className="mt-6">
                <div className="space-y-6">
                  <div className="bg-card rounded-lg p-4 border">
                    <h3 className="font-medium mb-3">Contact Information</h3>
                    <div className="space-y-3">
                      <div className="flex items-center gap-2">
                        <div className="flex-1 min-w-0">
                          <h3 className="text-lg font-semibold truncate">
                            {contact ? fullName(contact) : investor?.contact?.name}
                          </h3>
                          {contact?.company ? (
                            <p className="text-sm text-muted-foreground">
                              {contact.company}
                            </p>
                          ) : investor?.contact?.company ? (
                            <p className="text-sm text-muted-foreground">
                              {investor.contact.company}
                            </p>
                          ) : (
                            <p className="text-sm text-muted-foreground italic">
                              No company
                            </p>
                          )}
                          {contact?.email && (
                            <div className="flex items-center gap-2 mt-2">
                              <RiMailLine className="w-4 h-4 text-muted-foreground" />
                              <span className="text-sm">{contact.email}</span>
                              <Copy
                                text={contact.email}
                                onCopyText={"Email copied to clipboard"}
                                tooltipText="Copy email"
                              />
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>

                  <div className="bg-card rounded-lg p-4 border">
                    <div className="flex items-center justify-between mb-3">
                      <h3 className="font-medium">Deal Information</h3>
                      {!isArchived && (
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => setIsEditingInvestor(true)}
                        >
                          <RiEditLine className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                    <div className="space-y-3 text-sm">
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Expected Check Size</span>
                        <span className="font-medium">
                          {deal?.check_size !== undefined
                            ? `$${(Number(deal.check_size) / 100).toLocaleString()}`
                            : investor?.checkSize || "TBD"}
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Initial Contact</span>
                        <span className="font-medium">
                          {formatSafeDate(deal?.initial_contact) || formatSafeDate(investor?.initialContactDate)}
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Lead Investor</span>
                        <Badge variant={deal?.can_lead_round || investor?.isLeadInvestor ? "default" : "outline"}>
                          {deal?.can_lead_round || investor?.isLeadInvestor ? "Yes" : "No"}
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Rating</span>
                        <div className="flex items-center">
                          {[1, 2, 3, 4, 5].map((star) => (
                            star <= (deal?.rating || investor?.rating || 0)
                              ? <RiStarFill key={star} className="w-4 h-4 text-yellow-400" />
                              : <RiStarLine key={star} className="w-4 h-4 text-muted-foreground" />
                          ))}
                        </div>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Stage</span>
                        <Badge variant="outline">{investor?.stage || "Not set"}</Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Due Date</span>
                        <span className="font-medium">
                          {formatSafeDate(investor?.dueDate)}
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="bg-card rounded-lg p-4 border">
                    <h3 className="font-medium mb-3">Suggested Actions</h3>
                    <div className="space-y-2">
                      <Button
                        variant="outline"
                        className="w-full justify-start"
                        disabled={true}
                      >
                        <RiCalendarLine className="w-4 h-4 mr-2" />
                        Add Activity or Note
                      </Button>
                      <Button
                        variant="outline"
                        className="w-full justify-start"
                        disabled={true}
                      >
                        <RiFileTextLine className="w-4 h-4 mr-2" />
                        Upload Documents
                      </Button>
                    </div>
                  </div>
                </div>
              </TabsContent>

              <TabsContent value="activity">
                <ActivityList />
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
                deal={deal}
                isLoading={updateInvestorMutation.isPending}
              />

              <AddActivityDialog
                open={isAddingActivity}
                onOpenChange={setIsAddingActivity}
                onSubmit={handleAddActivity}
              />
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
} 
