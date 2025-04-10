import { useState, useRef, useEffect } from "react";
import type { MalakContact, MalakFundraiseContactDealDetails } from "@/client/Api";
import { fullName } from "@/lib/custom";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
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
  RiFileTextLine, RiCloseLine, RiAddLine, RiStarFill,
  RiStarLine,
  RiArchiveFill,
  RiTeamLine,
  RiArrowRightSLine,
  RiEditLine,
  RiClipboardLine
} from "@remixicon/react";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { NumericFormat } from "react-number-format";
import { format, fromUnixTime, isValid, parseISO } from "date-fns";
import type { Card, Activity, Note } from "@/components/investor-pipeline/types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { FETCH_FUNDRAISING_PIPELINE, UPDATE_INVESTOR_IN_PIPELINE } from "@/lib/query-constants";
import client from "@/lib/client";
import { toast } from "sonner";
import type { AxiosError } from "axios";
import type { ServerAPIStatus } from "@/client/Api";
import CopyToClipboard from "react-copy-to-clipboard";
import { ActivityList } from "../investor-pipeline/details/tabs/activity/ActivityList";
import Copy from "../ui/custom/copy";

interface InvestorDetailsDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  isArchived?: boolean;
  contact?: MalakContact;
  deal?: MalakFundraiseContactDealDetails;
  slug: string;
}

const TOTAL_ACTIVITIES_LIMIT = 250;
const ACTIVITIES_PER_PAGE = 25;

// Mock function to fetch activities
const fetchActivities = async (page: number, investorId: string): Promise<Activity[]> => {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1000));

  const startIndex = page * ACTIVITIES_PER_PAGE;
  // If we've reached the limit, return empty array
  if (startIndex >= TOTAL_ACTIVITIES_LIMIT) {
    return [];
  }

  // Calculate how many items to generate (handle last page)
  const itemsToGenerate = Math.min(
    ACTIVITIES_PER_PAGE,
    TOTAL_ACTIVITIES_LIMIT - startIndex
  );

  // Generate activities
  return Array.from({ length: itemsToGenerate }, (_, i) => ({
    id: `${page}-${i}-${Math.random()}`,
    type: ['email', 'meeting', 'document', 'team', 'stage_change'][Math.floor(Math.random() * 5)] as Activity['type'],
    title: `Activity ${startIndex + i + 1}`,
    description: `Description for activity ${startIndex + i + 1}`,
    timestamp: new Date(Date.now() - (startIndex + i) * 24 * 60 * 60 * 1000).toISOString(),
    content: Math.random() > 0.5 ? `Content for activity ${startIndex + i + 1}` : undefined,
  }));
};

function ActivitySkeleton() {
  return (
    <div className="relative">
      <div className="absolute -left-[27px] bg-background p-1 border rounded-full">
        <Skeleton className="h-4 w-4 rounded-full" />
      </div>
      <div className="bg-card rounded-lg p-4 border">
        <div className="flex items-center justify-between mb-2">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-24" />
        </div>
        <Skeleton className="h-4 w-full mb-3" />
        <Skeleton className="h-16 w-full" />
      </div>
    </div>
  );
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

function NoteDialog({
  open,
  onOpenChange,
  onSubmit,
  initialNote
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (note: Partial<Note>) => void;
  initialNote?: Note;
}) {
  const [note, setNote] = useState<Partial<Note>>(() => initialNote || {
    title: '',
    content: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(note);
    onOpenChange(false);
    if (!initialNote) {
      setNote({
        title: '',
        content: ''
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{initialNote ? 'Edit Note' : 'Add New Note'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Title</label>
            <Input
              value={note.title}
              onChange={(e) => setNote({ ...note, title: e.target.value })}
              placeholder="Note title"
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Content</label>
            <Textarea
              value={note.content}
              onChange={(e) => setNote({ ...note, content: e.target.value })}
              placeholder="Note content"
              required
              className="min-h-[200px]"
            />
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">{initialNote ? 'Save Changes' : 'Add Note'}</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function formatFileSize(bytes: number): string {
  const units = ['B', 'KB', 'MB', 'GB'];
  let size = bytes;
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return `${size.toFixed(1)} ${units[unitIndex]}`;
}

function truncateText(text: string, maxLength: number) {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength - 3) + '...';
}

function DocumentsTab({ isReadOnly }: { isReadOnly: boolean }) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center space-y-4">
      <RiFileTextLine className="w-12 h-12 text-muted-foreground" />
      <div>
        <h3 className="text-lg font-semibold">Documents Coming Soon</h3>
        <p className="text-sm text-muted-foreground">
          The documents feature is currently under development and will be available soon.
        </p>
      </div>
    </div>
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

// Add a safe date formatting helper
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

  // Infinite scroll states
  const [page, setPage] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const observerTarget = useRef<HTMLDivElement>(null);

  const updateInvestorMutation = useMutation({
    mutationKey: [UPDATE_INVESTOR_IN_PIPELINE, slug],
    mutationFn: async (updatedDeal: Partial<MalakFundraiseContactDealDetails>) => {
      if (!contact?.reference) {
        throw new Error("No contact reference provided")
      }

      const response = await client.pipelines.contactsPartialUpdate(slug, contact.id as string, {
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

  const loadMoreActivities = async () => {
    if (!investor || isLoading || !hasMore) return;

    setIsLoading(true);
    try {
      const newActivities = await fetchActivities(page, investor.id);
      if (newActivities.length === 0 || activities.length + newActivities.length >= TOTAL_ACTIVITIES_LIMIT) {
        setHasMore(false);
      }
      setActivities(prev => [...prev, ...newActivities]);
      setPage(prev => prev + 1);
    } catch (error) {
      console.error('Error loading activities:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting && hasMore && !isLoading) {
          loadMoreActivities();
        }
      },
      { threshold: 0.1 }
    );

    if (observerTarget.current) {
      observer.observe(observerTarget.current);
    }

    return () => observer.disconnect();
  }, [hasMore, isLoading, page]);

  // Reset states when investor changes
  useEffect(() => {
    if (investor) {
      setActivities([]);
      setPage(0);
      setHasMore(true);
      loadMoreActivities();
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

  const getActivityIcon = (type: Activity['type']) => {
    switch (type) {
      case 'email':
        return <RiMailLine className="w-4 h-4 text-primary" />;
      case 'meeting':
        return <RiCalendarLine className="w-4 h-4 text-primary" />;
      case 'document':
        return <RiFileTextLine className="w-4 h-4 text-primary" />;
      case 'team':
        return <RiTeamLine className="w-4 h-4 text-primary" />;
      case 'stage_change':
        return <RiArrowRightSLine className="w-4 h-4 text-primary" />;
      default:
        return <RiMailLine className="w-4 h-4 text-primary" />;
    }
  };

  const handleSaveInvestor = (updatedDeal: Partial<MalakFundraiseContactDealDetails>) => {
    updateInvestorMutation.mutate(updatedDeal);
  };

  const existingContactIds = [];

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
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger>
                                    <Copy
                                      text={contact.email}
                                      onCopyText={"Email copied to clipboard"}
                                    />
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    Copy email
                                  </TooltipContent>
                                </Tooltip>
                              </TooltipProvider>
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
                        onClick={() => setActiveTab("activity")}
                      >
                        <RiCalendarLine className="w-4 h-4 mr-2" />
                        Add Activity or Note
                      </Button>
                      <Button
                        variant="outline"
                        className="w-full justify-start"
                        onClick={() => setActiveTab("documents")}
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
