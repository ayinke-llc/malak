import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  RiMailLine,
  RiPhoneLine,
  RiMapPinLine,
  RiBuilding2Line,
  RiCalendarLine,
  RiDeleteBinLine,
  RiTeamLine,
  RiBarChartLine,
  RiTimeLine,
  RiFileTextLine,
  RiDashboardLine,
  RiFolderOpenLine,
  RiAddLine,
  RiSearchLine,
  RiArrowUpSLine,
  RiArrowDownSLine
} from '@remixicon/react';
import { Badge } from "@/components/ui/badge";
import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis } from "recharts";
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
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { MalakContact, MalakContactShareItem, ServerAPIStatus } from "@/client/Api";
import { fullName } from "@/lib/custom";
import { format, formatDistanceToNow } from "date-fns";
import Skeleton from "../../custom/loader/skeleton";
import { useMutation } from "@tanstack/react-query";
import { DELETE_CONTACT } from "@/lib/query-constants";
import client from "@/lib/client";
import { toast } from "sonner";
import { AxiosError } from "axios";
import Link from "next/link";
import { EditContactDialog } from "./form";
import { ContactListsView } from "../lists/contact-lists-view";

type TimePeriod = 'days' | 'weeks' | 'months';

interface ContactDetailsProps {
  reference: string;
  contact: MalakContact
  shared_items: MalakContactShareItem[]
  isLoading: boolean
}

interface AddToListDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  availableLists: Array<{
    id: number;
    title: string;
    description: string;
    count: number;
    color: string;
  }>;
  onAddToList: (listId: number) => void;
}

const AddToListDialog = ({ open, onOpenChange, availableLists, onAddToList }: AddToListDialogProps) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedListId, setSelectedListId] = useState<number | null>(null);

  const filteredLists = availableLists.filter(list =>
    list.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    list.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleAddToList = () => {
    if (selectedListId !== null) {
      onAddToList(selectedListId);
      onOpenChange(false);
      setSelectedListId(null);
      setSearchQuery("");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add to List</DialogTitle>
          <DialogDescription>
            Select a list to add this contact to
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <div className="relative">
            <RiSearchLine className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search lists..."
              className="pl-8"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <ScrollArea className="h-[300px] pr-4">
            <div className="space-y-2">
              {filteredLists.map((list) => (
                <div
                  key={list.id}
                  className={`flex items-center gap-4 p-3 rounded-lg border cursor-pointer transition-colors ${
                    selectedListId === list.id ? 'bg-primary/5 border-primary' : 'hover:bg-accent/5'
                  }`}
                  onClick={() => setSelectedListId(list.id)}
                >
                  <div className={`h-8 w-8 rounded-lg ${list.color} flex items-center justify-center text-white`}>
                    <RiTeamLine className="h-4 w-4" />
                  </div>
                  <div className="flex-1">
                    <h4 className="text-sm font-medium">{list.title}</h4>
                    <p className="text-xs text-muted-foreground">{list.description}</p>
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {list.count} contacts
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleAddToList} disabled={selectedListId === null}>
            Add to List
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

const ContactDetails = ({ isLoading, reference, contact, shared_items }: ContactDetailsProps) => {
  const router = useRouter();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const deleteMutation = useMutation({
    mutationKey: [DELETE_CONTACT],
    mutationFn: async (reference: string) => client.contacts.deleteContact(reference),
    onSuccess: () => {
      toast.success("Contact deleted successfully");
      router.back();
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred while deleting contact");
    },
  });

  const handleDelete = async () => {
    deleteMutation.mutate(reference)
  };

  return (
    <div className="mt-6 space-y-6">
      {isLoading ? (
        <Skeleton count={20} />
      ) : (
        <>
          <Card className="shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-7">
              <div className="space-y-1">
                <div className="flex items-center gap-3">
                  <div className="h-16 w-16 rounded-full bg-primary/10 flex items-center justify-center">
                    <span className="text-xl font-semibold text-primary">
                      {contact.email ? (contact.email[0].toUpperCase()) : '?'}
                    </span>
                  </div>
                  <div>
                    <CardTitle className="text-2xl font-bold tracking-tight text-foreground">
                      {contact.first_name || contact.last_name ? fullName(contact) : contact.email}
                    </CardTitle>
                    <CardDescription className="text-base text-muted-foreground">
                      {contact.company && (
                        <span className="flex items-center gap-1">
                          <RiBuilding2Line className="h-4 w-4" />
                          {contact.company}
                        </span>
                      )}
                    </CardDescription>
                  </div>
                </div>
              </div>
              <div className="flex gap-3">
                <EditContactDialog contact={contact} />
                <Button
                  variant="destructive"
                  size="icon"
                  onClick={() => setShowDeleteDialog(true)}
                  disabled={isLoading}
                  className="h-9 w-9"
                >
                  <RiDeleteBinLine className="h-4 w-4" />
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <Tabs defaultValue="details" className="w-full">
                <TabsList className="mb-4">
                  <TabsTrigger value="details" className="text-sm">Details</TabsTrigger>
                  <TabsTrigger value="lists" className="text-sm">Lists</TabsTrigger>
                  <TabsTrigger value="notes" className="text-sm">Notes</TabsTrigger>
                </TabsList>

                <TabsContent value="details">
                  <div className="grid gap-8 mt-6">
                    {/* Quick Actions */}
                    <div className="grid grid-cols-2 gap-4 max-w-md">
                      <Button 
                        variant="outline" 
                        className="h-auto py-4 flex flex-col gap-1 w-full"
                        onClick={() => window.location.href = `mailto:${contact?.email}`}
                        disabled={!contact?.email}
                      >
                        <RiMailLine className="h-5 w-5" />
                        <span className="text-xs">Send Email</span>
                      </Button>
                      <Button 
                        variant="outline" 
                        className="h-auto py-4 flex flex-col gap-1 w-full"
                        onClick={() => window.location.href = `tel:${contact?.phone}`}
                        disabled={!contact?.phone}
                      >
                        <RiPhoneLine className="h-5 w-5" />
                        <span className="text-xs">Call</span>
                      </Button>
                    </div>

                    {/* Contact Information */}
                    <div className="grid gap-6">
                      <h3 className="text-lg font-semibold">Contact Information</h3>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-4">
                          <div className="flex items-center gap-3">
                            <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center">
                              <RiMailLine className="h-4 w-4 text-primary" />
                            </div>
                            <div>
                              <p className="text-sm text-muted-foreground">Email</p>
                              <p className="font-medium">{contact?.email || "N/A"}</p>
                            </div>
                          </div>
                          <div className="flex items-center gap-3">
                            <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center">
                              <RiPhoneLine className="h-4 w-4 text-primary" />
                            </div>
                            <div>
                              <p className="text-sm text-muted-foreground">Phone</p>
                              <p className="font-medium">{contact?.phone || "N/A"}</p>
                            </div>
                          </div>
                        </div>
                        <div className="space-y-4">
                          <div className="flex items-center gap-3">
                            <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center">
                              <RiBuilding2Line className="h-4 w-4 text-primary" />
                            </div>
                            <div>
                              <p className="text-sm text-muted-foreground">Company</p>
                              <p className="font-medium">{contact?.company || "N/A"}</p>
                            </div>
                          </div>
                          <div className="flex items-center gap-3">
                            <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center">
                              <RiMapPinLine className="h-4 w-4 text-primary" />
                            </div>
                            <div>
                              <p className="text-sm text-muted-foreground">Location</p>
                              <p className="font-medium">{contact?.city || "N/A"}</p>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Timeline Section */}
                    <div className="grid gap-4">
                      <h3 className="text-lg font-semibold">Timeline</h3>
                      <div className="space-y-4">
                        <div className="flex items-center gap-3">
                          <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center">
                            <RiCalendarLine className="h-4 w-4 text-primary" />
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Created</p>
                            <p className="font-medium">
                              {format(contact?.created_at as string || new Date(), "EEEE, MMMM do, yyyy")}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-3">
                          <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center">
                            <RiTimeLine className="h-4 w-4 text-primary" />
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Last Updated</p>
                            <p className="font-medium">
                              {format(contact?.updated_at as string || new Date(), "EEEE, MMMM do, yyyy")}
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </TabsContent>

                <TabsContent value="lists">
                  <ContactListsView contact={contact} />
                </TabsContent>

                <TabsContent value="notes">
                  <div className="mt-6">
                    {contact?.notes ? (
                      <div className="prose prose-sm max-w-none">
                        <p className="text-muted-foreground">{contact.notes}</p>
                      </div>
                    ) : (
                      <div className="text-center py-8 border rounded-lg bg-muted/10">
                        <p className="text-sm text-muted-foreground">No notes available for this contact</p>
                      </div>
                    )}
                  </div>
                </TabsContent>
              </Tabs>
            </CardContent>
          </Card>

          {/* Latest Shared Items Section */}
          <div className="space-y-8 mt-10">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground">Latest shared items</h2>

            {/* Updates Section */}
            <div>
              <div className="flex items-center mb-4">
                <div className="flex items-center gap-2">
                  <RiFileTextLine className="h-5 w-5 text-muted-foreground" />
                  <h3 className="text-lg font-semibold text-foreground">Updates</h3>
                </div>
              </div>
              {(!shared_items || shared_items.filter((value) => value.item_type === "update").length === 0) ? (
                <div className="p-8 text-center border rounded-lg bg-muted/10">
                  <p className="text-sm text-muted-foreground">No updates have been shared with this contact yet</p>
                </div>
              ) : (
                <div className="space-y-2">
                  {shared_items.
                    filter((value) => value.item_type === "update").
                    map((item) => (
                      <div key={item.item_reference} className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent/5 hover:text-accent-foreground cursor-pointer">
                        <Link href={`/updates/${item?.item_reference}`} className="text-sm font-medium text-foreground">
                          {item?.title}
                        </Link>
                        <span className="text-sm text-muted-foreground">
                          Sent {formatDistanceToNow(item?.shared_at as string, { addSuffix: true })}
                        </span>
                      </div>
                    ))}
                </div>
              )}
            </div>

            {/* Dashboards Section */}
            <div>
              <div className="flex items-center mb-4">
                <div className="flex items-center gap-2">
                  <RiDashboardLine className="h-5 w-5 text-muted-foreground" />
                  <h3 className="text-lg font-semibold text-foreground">Dashboards</h3>
                </div>
              </div>
              {(!shared_items || shared_items.filter((value) => value.item_type === "dashboard").length === 0) ? (
                <div className="p-8 text-center border rounded-lg bg-muted/10">
                  <p className="text-sm text-muted-foreground">No dashboards have been shared with this contact yet</p>
                </div>
              ) : (
                <div className="space-y-2">
                  {shared_items.
                    filter((value) => value.item_type === "dashboard").
                    map((item) => (
                      <div key={item.item_reference} className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent/5 hover:text-accent-foreground cursor-pointer">
                        <Link href={`/dashboards/${item?.item_reference}`} className="text-sm font-medium text-foreground">
                          {item?.title}
                        </Link>
                        <span className="text-sm text-muted-foreground">
                          Sent {formatDistanceToNow(item?.shared_at as string, { addSuffix: true })}
                        </span>
                      </div>
                    ))}
                </div>
              )}
            </div>

            {/* Data Rooms Section */}
            <div>
              <div className="flex items-center mb-4">
                <div className="flex items-center gap-2">
                  <RiFolderOpenLine className="h-5 w-5 text-muted-foreground" />
                  <h3 className="text-lg font-semibold text-foreground">Data rooms</h3>
                </div>
              </div>
              {(!shared_items || shared_items.filter((value) => value.item_type === "deck").length === 0) ? (
                <div className="p-8 text-center border rounded-lg bg-muted/10">
                  <p className="text-sm text-muted-foreground">No data rooms have been shared with this contact yet</p>
                </div>
              ) : (
                <div className="space-y-2">
                  {shared_items.
                    filter((value) => value.item_type === "deck").
                    map((item) => (
                      <div key={item.item_reference} className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent/5 hover:text-accent-foreground cursor-pointer">
                        <Link href={`/data-rooms/${item?.item_reference}`} className="text-sm font-medium text-foreground">
                          {item?.title}
                        </Link>
                        <span className="text-sm text-muted-foreground">
                          Sent {formatDistanceToNow(item?.shared_at as string, { addSuffix: true })}
                        </span>
                      </div>
                    ))}
                </div>
              )}
            </div>
          </div>
        </>
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the contact
              and remove their data from our servers.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteMutation.isPending}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteMutation.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};

export default ContactDetails;
