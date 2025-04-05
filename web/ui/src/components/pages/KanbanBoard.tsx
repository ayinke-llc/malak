"use client";

import { useState } from "react";
import { DragDropContext, Droppable, Draggable, DropResult } from "@hello-pangea/dnd";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  RiTimeLine, RiAddLine, RiSettings4Line, RiArchiveLine, RiArchiveFill,
  RiInboxUnarchiveLine, RiInformationLine
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
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";

import { ShareSettingsDialog } from "@/components/investor-pipeline/ShareSettingsDialog";
import { AddInvestorDialog as AddInvestorDialogComponent } from "@/components/investor-pipeline/AddInvestorDialog";
import type { SearchResult } from "@/components/investor-pipeline/AddInvestorDialog";
import type { Card as InvestorCard, Columns, Board, ShareSettings } from "@/components/investor-pipeline/types";
import { initialBoard } from "@/components/investor-pipeline/mock-data";
import { Input } from "@/components/ui/input";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export default function KanbanBoard() {
  const [board, setBoard] = useState<Board>(initialBoard);
  const [selectedInvestor, setSelectedInvestor] = useState<InvestorCard | null>(null);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [isAddInvestorOpen, setIsAddInvestorOpen] = useState(false);
  const [isShareSettingsOpen, setIsShareSettingsOpen] = useState(false);
  const [isArchiveConfirmOpen, setIsArchiveConfirmOpen] = useState(false);
  const [isUnarchiveConfirmOpen, setIsUnarchiveConfirmOpen] = useState(false);
  const [shareSettings, setShareSettings] = useState<ShareSettings>({
    isEnabled: false,
    shareLink: "",
    requireEmail: false,
    requirePassword: false,
  });

  const handleDragEnd = (result: DropResult) => {
    if (!result.destination || board.isArchived) return;

    const { source, destination } = result;
    
    // Don't do anything if dropped in the same position
    if (
      source.droppableId === destination.droppableId &&
      source.index === destination.index
    ) {
      return;
    }

    const sourceColumn = board.columns[source.droppableId];
    const destColumn = board.columns[destination.droppableId];
    
    if (sourceColumn && destColumn) {
      const sourceCards = Array.from(sourceColumn.cards);
      const destCards = source.droppableId === destination.droppableId ? sourceCards : Array.from(destColumn.cards);
      const [removed] = sourceCards.splice(source.index, 1);
      destCards.splice(destination.index, 0, removed);

      setBoard({
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
      });
    }
  };

  const handleArchive = () => {
    setBoard({
      ...board,
      isArchived: true
    });
    setIsArchiveConfirmOpen(false);
    toast.success("Pipeline archived successfully");
  };

  const handleUnarchive = () => {
    setBoard({
      ...board,
      isArchived: false
    });
    setIsUnarchiveConfirmOpen(false);
    toast.success("Pipeline unarchived successfully");
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

    const newCard: InvestorCard = {
      id: investor.id,
      title: investor.company,
      amount: "TBD",
      stage: "Initial Contact",
      dueDate: new Date().toISOString().split('T')[0],
      contact: {
        name: investor.name,
        image: investor.image,
      },
      roundDetails: {
        raising: "TBD",
        type: "TBD",
        ownership: "TBD"
      },
      checkSize: investor.checkSize,
      initialContactDate: investor.initialContactDate,
      isLeadInvestor: investor.isLeadInvestor,
      rating: investor.rating
    };

    setBoard({
      ...board,
      columns: {
        ...board.columns,
        backlog: {
          ...board.columns.backlog,
          cards: [...board.columns.backlog.cards, newCard],
        },
      }
    });
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
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsShareSettingsOpen(true)}
          >
            <RiSettings4Line className="mr-1 h-4 w-4" />
            Share Settings
          </Button>
          {board.isArchived ? (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsUnarchiveConfirmOpen(true)}
            >
              <RiInboxUnarchiveLine className="mr-1 h-4 w-4" />
              Unarchive
            </Button>
          ) : (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsArchiveConfirmOpen(true)}
            >
              <RiArchiveLine className="mr-1 h-4 w-4" />
              Archive
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
              {Object.entries(board.columns).map(([columnId, column]) => (
                <div key={columnId} className="flex flex-col w-[280px] shrink-0 rounded-lg bg-muted/20">
                  <div className="p-2 mb-1">
                    <div className="flex items-center justify-between px-2 py-1">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <div className="flex items-center gap-1 cursor-help">
                              <h2 className="font-medium text-sm">{column.title}</h2>
                              <RiInformationLine className="h-4 w-4 text-muted-foreground" />
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p className="max-w-xs">{column.description}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <Badge variant="secondary" className="text-xs">
                        {column.cards.length}
                      </Badge>
                    </div>
                  </div>

                  <Droppable droppableId={columnId} key={columnId}>
                    {(provided, snapshot) => (
                      <div
                        ref={provided.innerRef}
                        {...provided.droppableProps}
                        className={`flex-1 p-2 space-y-2 overflow-y-auto
                          ${snapshot.isDraggingOver ? 'bg-muted/30' : ''}
                          scrollbar-thin scrollbar-thumb-muted-foreground/20 scrollbar-track-transparent
                          hover:scrollbar-thumb-muted-foreground/30`}
                        style={{
                          height: 'calc(100vh - 8rem)',
                          maxHeight: 'calc(100vh - 8rem)',
                          overflowY: 'auto',
                          overflowX: 'hidden'
                        }}
                      >
                        {column.cards.map((card, index) => (
                          <Draggable
                            key={card.id}
                            draggableId={card.id}
                            index={index}
                            isDragDisabled={board.isArchived}
                          >
                            {(provided, snapshot) => (
                              <Card
                                ref={provided.innerRef}
                                {...provided.draggableProps}
                                {...provided.dragHandleProps}
                                className={`cursor-pointer border-none shadow-sm hover:shadow-md transition-shadow
                                  ${snapshot.isDragging ? 'opacity-50 shadow-lg ring-2 ring-primary' : ''}`}
                                onClick={() => {
                                  setSelectedInvestor(card);
                                  setIsDetailsOpen(true);
                                }}
                              >
                                <CardContent className="p-3">
                                  <div className="flex items-center gap-3">
                                    <Avatar className="h-8 w-8">
                                      <AvatarImage
                                        src={card.contact.image}
                                        alt={card.contact.name}
                                      />
                                      <AvatarFallback>
                                        {card.contact.name
                                          .split(" ")
                                          .map((n) => n[0])
                                          .join("")}
                                      </AvatarFallback>
                                    </Avatar>
                                    <div className="min-w-0 flex-1">
                                      <h4 className="truncate font-medium text-sm">
                                        {card.title}
                                      </h4>
                                      <p className="truncate text-xs text-muted-foreground">
                                        {card.contact.name}
                                      </p>
                                      <div className="mt-1 flex items-center gap-2">
                                        <Badge variant="secondary" className="text-xs">
                                          {card.amount}
                                        </Badge>
                                        <div className="flex items-center text-xs text-muted-foreground">
                                          <RiTimeLine className="mr-1 h-3 w-3" />
                                          {card.dueDate}
                                        </div>
                                      </div>
                                    </div>
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
      />

      <AddInvestorDialogComponent
        open={isAddInvestorOpen}
        onOpenChange={setIsAddInvestorOpen}
        onAddInvestor={handleAddInvestor}
      />

      <ShareSettingsDialog
        open={isShareSettingsOpen}
        onOpenChange={setIsShareSettingsOpen}
        settings={shareSettings}
        onSettingsChange={setShareSettings}
      />

      <AlertDialog open={isArchiveConfirmOpen} onOpenChange={setIsArchiveConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Archive Pipeline</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to archive this pipeline? This will make it read-only
              and prevent any further changes.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleArchive}>
              Archive Pipeline
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={isUnarchiveConfirmOpen} onOpenChange={setIsUnarchiveConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Unarchive Pipeline</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to unarchive this pipeline? This will make it editable
              and allow changes to be made.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleUnarchive}>
              Unarchive Pipeline
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
} 
