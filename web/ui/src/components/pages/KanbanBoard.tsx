"use client";

import { useState } from "react";
import { DragDropContext, Droppable, Draggable, DropResult } from "@hello-pangea/dnd";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  RiTimeLine, RiAddLine, RiSettings4Line, RiArchiveLine,
  RiInboxUnarchiveLine, RiInformationLine, RiCalendarLine, RiErrorWarningLine,
  RiCloseLine, RiMailLine, RiPhoneLine, RiMoneyDollarCircleLine
} from "@remixicon/react";
import { InvestorDetailsDrawer } from "./InvestorDetailsDrawer";
import { toast } from "sonner";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { ShareSettingsDialog } from "@/components/investor-pipeline/ShareSettingsDialog";
import { AddInvestorDialog as AddInvestorDialogComponent } from "@/components/investor-pipeline/AddInvestorDialog";
import type { SearchResult } from "@/components/investor-pipeline/AddInvestorDialog";
import type { Card as InvestorCard, Board, ShareSettings } from "@/components/investor-pipeline/types";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  Tabs,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { FETCH_FUNDRAISING_PIPELINE, CLOSE_FUNDRAISING_PIPELINE, ADD_INVESTOR_TO_PIPELINE, SEARCH_CONTACTS } from "@/lib/query-constants";
import client from "@/lib/client";
import type { ServerFetchBoardResponse } from "@/client/Api";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import type { AxiosError } from "axios";
import type { ServerAPIStatus } from "@/client/Api";
import { fullName } from "@/lib/custom";
import type { MalakContact } from "@/client/Api";
import type { MalakFundraiseContactDealDetails } from "@/client/Api";

interface KanbanBoardProps {
  slug: string;
}

export default function KanbanBoard({ slug }: KanbanBoardProps) {
  const [selectedInvestor, setSelectedInvestor] = useState<InvestorCard | null>(null);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [isAddInvestorOpen, setIsAddInvestorOpen] = useState(false);
  const [isShareSettingsOpen, setIsShareSettingsOpen] = useState(false);
  const [isCloseConfirmOpen, setIsCloseConfirmOpen] = useState(false);
  const [shareSettings, setShareSettings] = useState<ShareSettings>({
    isEnabled: false,
    shareLink: "",
    requireEmail: false,
    requirePassword: false,
  });
  const [activeTab, setActiveTab] = useState("overview");

  const queryClient = useQueryClient();

  const { data: boardData, isLoading, error } = useQuery<ServerFetchBoardResponse>({
    queryKey: [FETCH_FUNDRAISING_PIPELINE, slug],
    queryFn: async () => {
      const response = await client.pipelines.boardDetail(slug);
      return response.data;
    },
  });

  const updateBoardMutation = useMutation({
    mutationFn: async (updatedBoard: Board) => {
      await new Promise(resolve => setTimeout(resolve, 500));
      return updatedBoard;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FETCH_FUNDRAISING_PIPELINE, slug] });
      toast.success("Pipeline updated successfully");
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err.response?.data.message ?? "Failed to update pipeline");
    },
  });

  const closeBoardMutation = useMutation({
    mutationKey: [CLOSE_FUNDRAISING_PIPELINE, slug],
    mutationFn: async () => {
      const response = await client.pipelines.pipelinesDelete(slug);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FETCH_FUNDRAISING_PIPELINE, slug] });
      toast.success("Pipeline closed successfully");
      setIsCloseConfirmOpen(false);
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err.response?.data.message ?? "Failed to close pipeline");
    },
  });

  const addInvestorMutation = useMutation({
    mutationKey: [ADD_INVESTOR_TO_PIPELINE, slug],
    mutationFn: async (investor: SearchResult & {
      checkSize: string;
      initialContactDate: string;
      isLeadInvestor: boolean;
      rating: number;
    }) => {
      const response = await client.pipelines.contactsCreate(slug, {
        contact_reference: investor.reference,
        check_size: Number(investor.checkSize) * 100, // Convert to cents
        initial_contact: Math.floor(new Date(investor.initialContactDate).getTime() / 1000),
        can_lead_round: investor.isLeadInvestor,
        rating: investor.rating
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [FETCH_FUNDRAISING_PIPELINE, slug] });
      toast.success("Investor added successfully");
      setIsAddInvestorOpen(false);
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err.response?.data.message ?? "Failed to add investor");
    },
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-screen p-4 text-center">
        <RiErrorWarningLine className="w-12 h-12 text-destructive mb-4" />
        <h2 className="text-2xl font-semibold mb-2">Unable to Load Board</h2>
        <p className="text-muted-foreground mb-4 max-w-md">
          An error occurred while loading the board. Please try again.
        </p>
        <Button
          variant="default"
          onClick={() => queryClient.invalidateQueries({ queryKey: [FETCH_FUNDRAISING_PIPELINE, slug] })}
        >
          Retry
        </Button>
      </div>
    );
  }

  if (!boardData) {
    return null;
  }

  const { pipeline = {}, columns = [], contacts = [], positions = [], deals = [] } = boardData;

  // Get all contact IDs from all contacts in the board
  const existingContactIds = contacts.map(contact => contact.contact_id || "").filter(Boolean);
  console.log('All board contacts:', contacts);
  console.log('Extracted contact IDs:', existingContactIds);

  // Transform the data into the format expected by the board
  const board: Board = {
    isArchived: pipeline.is_closed || false,
    columns: columns
      .sort((a, b) => {
        // Hardcoded column order based on title
        const columnOrder = [
          "Backlog",
          "Contacted",
          "Partner Meeting",
          "Passed",
          "Termsheet/SAFE",
          "Closed"
        ];
        
        const aIndex = columnOrder.indexOf(a.title || "");
        const bIndex = columnOrder.indexOf(b.title || "");
        
        // If both columns are in our order list, sort by their position
        if (aIndex !== -1 && bIndex !== -1) {
          return aIndex - bIndex;
        }
        
        // If only one column is in our order list, prioritize it
        if (aIndex !== -1) return -1;
        if (bIndex !== -1) return 1;
        
        // For any columns not in our list, maintain their original order
        return 0;
      })
      .reduce<Board["columns"]>((acc, column) => {
        if (column?.reference) {
          acc[column.reference] = {
            id: column.id || column.reference,
            title: column.title || "",
            description: column.description || "",
            cards: (contacts || [])
              .filter(contact => contact && contact.fundraising_pipeline_column_id === column.id)
              .map(contact => {
                // Find the position for this contact
                const position = positions.find(p => p.fundraising_pipeline_column_contact_id === contact.id);
                const deal = contact.deal_details;
                const contactDetails = contact.contact;
                
                return {
                  id: contact.reference || "",
                  title: contactDetails ? fullName(contactDetails) : "",
                  amount: deal?.check_size ? `$${(deal.check_size / 100).toLocaleString()}` : "",
                  stage: column?.title || "",
                  dueDate: pipeline.expected_close_date || "",
                  contact: {
                    name: contactDetails ? fullName(contactDetails) : "",
                    company: contactDetails?.company || "",
                    email: contactDetails?.email || "",
                    phone: contactDetails?.phone || "",
                    title: contactDetails?.title || "",
                  },
                  roundDetails: {
                    raising: pipeline.target_amount ? `$${(pipeline.target_amount / 100).toLocaleString()}` : "",
                    type: pipeline.stage || "",
                    ownership: "" // This data isn't available in the current API
                  },
                  checkSize: deal?.check_size ? `$${(deal.check_size / 100).toLocaleString()}` : "",
                  initialContactDate: deal?.initial_contact || contact.created_at || "",
                  isLeadInvestor: deal?.can_lead_round || false,
                  rating: deal?.rating || 0,
                  originalContact: contactDetails,
                  originalDeal: deal
                };
              })
          };
        }
        return acc;
      }, {})
  };

  const handleDragEnd = (result: DropResult) => {
    if (!result?.destination || board.isArchived) return;

    const { source, destination } = result;

    if (!source || !destination) return;

    if (
      source.droppableId === destination.droppableId &&
      source.index === destination.index
    ) {
      return;
    }

    const sourceColumn = board.columns[source.droppableId];
    const destColumn = board.columns[destination.droppableId];

    if (!sourceColumn || !destColumn) {
      toast.error("Invalid drag operation: column not found");
      return;
    }

    const sourceCards = Array.from(sourceColumn.cards || []);
    const destCards = source.droppableId === destination.droppableId
      ? sourceCards
      : Array.from(destColumn.cards || []);

    if (sourceCards.length === 0) {
      toast.error("Invalid drag operation: no cards to move");
      return;
    }

    const [removed] = sourceCards.splice(source.index, 1);
    if (!removed) {
      toast.error("Invalid drag operation: card not found");
      return;
    }

    destCards.splice(destination.index, 0, removed);

    const updatedBoard = {
      ...board,
      columns: {
        ...board.columns,
        [source.droppableId]: {
          ...sourceColumn,
          cards: sourceCards,
        },
        [destination.droppableId]: {
          ...destColumn,
          cards: destCards,
        },
      }
    };

    updateBoardMutation.mutate(updatedBoard);
  };

  const handleClose = () => {
    closeBoardMutation.mutate();
  };

  const handleAddInvestor = (investor: SearchResult & {
    checkSize: string;
    initialContactDate: string;
    isLeadInvestor: boolean;
    rating: number;
  }) => {
    if (board.isArchived) {
      toast.error("Cannot add investors to an archived pipeline");
      return;
    }

    addInvestorMutation.mutate(investor);
  };

  return (
    <div className="flex flex-col h-screen max-h-screen overflow-hidden">
      <div className="flex justify-between items-center p-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <h1 className="text-2xl font-semibold">Fundraising Pipeline</h1>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsAddInvestorOpen(true)}
            disabled={board.isArchived}
          >
            <RiAddLine className="mr-1 h-4 w-4" />
            Add Investor
          </Button>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="outline"
                  size="sm"
                  disabled
                >
                  <RiSettings4Line className="mr-1 h-4 w-4" />
                  Share Settings
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Coming soon!</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          {board.isArchived ? (
            <Button
              variant="outline"
              size="sm"
              disabled
            >
              <RiArchiveLine className="mr-1 h-4 w-4" />
              Closed
            </Button>
          ) : (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsCloseConfirmOpen(true)}
            >
              <RiCloseLine className="mr-1 h-4 w-4" />
              Close Board
            </Button>
          )}
        </div>
      </div>

      <div className="flex-1 min-h-0">
        <div
          className="h-full overflow-x-auto"
          style={{
            msOverflowStyle: 'none',
            scrollbarWidth: 'none',
            WebkitOverflowScrolling: 'touch',
          }}
        >
          <style jsx global>{`
            /* Hide scrollbar for Chrome, Safari and Opera */
            .overflow-x-auto::-webkit-scrollbar {
              display: none;
            }
          `}</style>

          <DragDropContext onDragEnd={handleDragEnd}>
            <div className="flex gap-4 p-4 h-full">
              {Object.entries(board.columns || {}).map(([columnId, column]) => (
                <div key={columnId} className="flex flex-col w-[280px] shrink-0 rounded-lg bg-muted/20">
                  <div className="p-2 mb-1">
                    <div className="flex items-center justify-between px-2 py-1">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <div className="flex items-center gap-1 cursor-help">
                              <h2 className="font-medium text-sm">{column?.title || "Unnamed Column"}</h2>
                              <RiInformationLine className="h-4 w-4 text-muted-foreground" />
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p className="max-w-xs">{column?.description || "No description available"}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <Badge variant="secondary" className="text-xs">
                        {(column?.cards || []).length}
                      </Badge>
                    </div>
                  </div>

                  <Droppable droppableId={columnId} key={columnId}>
                    {(provided, snapshot) => (
                      <div
                        ref={provided.innerRef}
                        {...provided.droppableProps}
                        className={`flex-1 p-2 space-y-2 overflow-y-auto
                          ${snapshot?.isDraggingOver ? 'bg-muted/30' : ''}
                          scrollbar-thin scrollbar-thumb-muted-foreground/20 scrollbar-track-transparent
                          hover:scrollbar-thumb-muted-foreground/30`}
                        style={{
                          height: 'calc(100vh - 8rem)',
                          maxHeight: 'calc(100vh - 8rem)',
                          overflowY: 'auto',
                          overflowX: 'hidden'
                        }}
                      >
                        {(column?.cards || []).map((card, index) => (
                          <Draggable
                            key={card?.id || index}
                            draggableId={card?.id || `card-${index}`}
                            index={index}
                            isDragDisabled={board.isArchived}
                          >
                            {(provided, snapshot) => (
                              <Card
                                ref={provided.innerRef}
                                {...provided.draggableProps}
                                {...provided.dragHandleProps}
                                className={`cursor-pointer border-none shadow-sm hover:shadow-md transition-shadow
                                  ${snapshot?.isDragging ? 'opacity-50 shadow-lg ring-2 ring-primary' : ''}`}
                                onClick={() => {
                                  if (card) {
                                    setSelectedInvestor(card);
                                    setIsDetailsOpen(true);
                                  }
                                }}
                              >
                                <CardContent className="p-3 space-y-3">
                                  <div className="flex items-start justify-between gap-2">
                                    <div className="flex-1 min-w-0">
                                      <div className="flex items-center gap-2">
                                        <Avatar className="h-8 w-8">
                                          <AvatarImage src="" />
                                          <AvatarFallback className="text-xs">
                                            {card?.title?.split(' ').map(n => n[0]).join('').toUpperCase()}
                                          </AvatarFallback>
                                        </Avatar>
                                        <div className="space-y-1 min-w-0">
                                          <h4 className="font-semibold text-sm truncate">
                                            {card?.title}
                                          </h4>
                                          {card?.contact?.company && (
                                            <div className="flex items-center gap-1.5">
                                              <p className="text-xs font-medium text-muted-foreground truncate">
                                                {card.contact.company}
                                              </p>
                                              {card?.contact?.title && (
                                                <span className="text-xs text-muted-foreground/60">
                                                  • {card.contact.title}
                                                </span>
                                              )}
                                            </div>
                                          )}
                                        </div>
                                      </div>
                                    </div>
                                    {card?.rating > 0 && (
                                      <div className="flex items-center">
                                        <Badge variant="secondary" className="text-[10px] h-5 px-2">
                                          {Array(card.rating).fill('★').join('')}
                                        </Badge>
                                      </div>
                                    )}
                                  </div>
                                  
                                  <div className="space-y-1.5">
                                    {card?.contact?.email && (
                                      <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                        <RiMailLine className="h-3.5 w-3.5 shrink-0" />
                                        <span className="truncate">{card.contact.email}</span>
                                      </div>
                                    )}
                                    {card?.contact?.phone && (
                                      <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                        <RiPhoneLine className="h-3.5 w-3.5 shrink-0" />
                                        <span className="truncate">{card.contact.phone}</span>
                                      </div>
                                    )}
                                    {card?.checkSize && (
                                      <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                        <RiMoneyDollarCircleLine className="h-3.5 w-3.5 shrink-0" />
                                        <span className="truncate">Check size: {card.checkSize}</span>
                                      </div>
                                    )}
                                  </div>

                                  <div className="flex items-center justify-between text-xs text-muted-foreground pt-1">
                                    {card?.initialContactDate && (
                                      <div className="flex items-center gap-1">
                                        <RiCalendarLine className="h-3.5 w-3.5" />
                                        {new Date(card.initialContactDate).toLocaleDateString()}
                                      </div>
                                    )}
                                    {card?.isLeadInvestor && (
                                      <Badge variant="outline" className="text-[10px] h-5 px-2 border-primary/20">
                                        Lead
                                      </Badge>
                                    )}
                                  </div>
                                </CardContent>
                              </Card>
                            )}
                          </Draggable>
                        ))}
                        {provided.placeholder}
                      </div>
                    )}
                  </Droppable>
                </div>
              ))}
            </div>
          </DragDropContext>
        </div>
      </div>

      <InvestorDetailsDrawer
        investor={selectedInvestor}
        open={isDetailsOpen}
        onOpenChange={setIsDetailsOpen}
        isArchived={board.isArchived}
        contact={selectedInvestor?.originalContact}
        deal={selectedInvestor?.originalDeal}
      />

      <AddInvestorDialogComponent
        open={isAddInvestorOpen}
        onOpenChange={setIsAddInvestorOpen}
        onAddInvestor={handleAddInvestor}
        isLoading={addInvestorMutation.isPending}
        existingContacts={existingContactIds}
      />

      <ShareSettingsDialog
        open={isShareSettingsOpen}
        onOpenChange={setIsShareSettingsOpen}
        settings={shareSettings}
        onSettingsChange={setShareSettings}
      />

      <AlertDialog open={isCloseConfirmOpen} onOpenChange={setIsCloseConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Close Pipeline Board</AlertDialogTitle>
            <AlertDialogDescription className="space-y-4">
              <p>
                Are you sure you want to close this pipeline board? This action is permanent and has the following effects:
              </p>
              <ul className="list-disc pl-4 space-y-2">
                <li>The board will become read-only</li>
                <li>You cannot add new investors or move existing ones</li>
                <li>This action cannot be undone - you cannot reopen the board once closed</li>
              </ul>
              <p className="font-medium text-destructive">
                Please confirm that you understand this is a permanent action.
              </p>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleClose}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={closeBoardMutation.isPending}
            >
              {closeBoardMutation.isPending ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-background mr-2" />
                  Closing...
                </>
              ) : (
                "Close Board Permanently"
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Tabs className="w-full">
        <TabsList className="w-full">
          <TabsTrigger value="overview" className="flex-1">Overview</TabsTrigger>
          <TabsTrigger value="activity" className="flex-1">Activity</TabsTrigger>
          <TabsTrigger value="documents" className="flex-1">Documents</TabsTrigger>
        </TabsList>
      </Tabs>

      <Dialog>
        <DialogTrigger asChild>
          <Button
            variant="outline"
            className="w-full justify-start"
            onClick={() => setActiveTab("activity")}
          >
            <RiCalendarLine className="w-4 h-4 mr-2" />
            Add Activity or Note
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogTitle>Add Activity or Note</DialogTitle>
          <DialogDescription>
            Record an activity or add a note to the pipeline.
          </DialogDescription>
          {/* Add your form content here */}
        </DialogContent>
      </Dialog>
    </div>
  );
} 
