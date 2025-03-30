import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { RiSearchLine, RiTeamLine } from '@remixicon/react';
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ADD_CONTACT_TO_LIST, FETCH_CONTACT_LISTS, FETCH_CONTACT } from "@/lib/query-constants";
import client from "@/lib/client";
import { toast } from "sonner";
import { AxiosError, AxiosResponse } from "axios";
import { MalakContactList, ServerFetchContactListsResponse } from "@/client/Api";

interface AddToListDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  contactReference: string;
  currentListIds: string[];
}

export function AddToListDialog({ open, onOpenChange, contactReference, currentListIds }: AddToListDialogProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedListReference, setSelectedListReference] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const { data: response, error } = useQuery<AxiosResponse<ServerFetchContactListsResponse>>({
    queryKey: [FETCH_CONTACT_LISTS],
    queryFn: () => client.contacts.fetchContactLists(),
  });

  const allLists = response?.data.lists.map((item: { list: MalakContactList }) => item.list) ?? [];

  const availableLists = allLists.filter((list: MalakContactList): list is MalakContactList & { reference: string } => {
    return list.reference !== undefined && !currentListIds.includes(list.id || "");
  });

  const filteredLists = availableLists.filter((list: MalakContactList) =>
    list.title?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const addToListMutation = useMutation({
    mutationKey: [ADD_CONTACT_TO_LIST],
    mutationFn: async (listReference: string) => {
      return client.contacts.addEmailToContactList(listReference, { reference: contactReference });
    },
    onSuccess: () => {
      toast.success("Contact added to list successfully");
      onOpenChange(false);
      setSelectedListReference(null);
      setSearchQuery("");
      queryClient.invalidateQueries({ queryKey: [FETCH_CONTACT, contactReference] });
    },
    onError: (err: AxiosError<{ message?: string }>) => {
      toast.error(err.response?.data?.message || "Failed to add contact to list");
    },
  });

  const handleAddToList = () => {
    if (selectedListReference !== null) {
      addToListMutation.mutate(selectedListReference);
    }
  };

  if (error) {
    toast.error("Failed to fetch contact lists");
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Add to List</DialogTitle>
          </DialogHeader>
          <div className="py-6 text-center text-muted-foreground">
            Failed to load contact lists. Please try again later.
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

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
                  key={list.reference}
                  className={`flex items-center gap-4 p-3 rounded-lg border cursor-pointer transition-colors ${
                    selectedListReference === list.reference ? 'bg-primary/5 border-primary' : 'hover:bg-accent/5'
                  }`}
                  onClick={() => setSelectedListReference(list.reference)}
                >
                  <div className={`h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center text-white`}>
                    <RiTeamLine className="h-4 w-4" />
                  </div>
                  <div className="flex-1">
                    <h4 className="text-sm font-medium">{list.title}</h4>
                  </div>
                </div>
              ))}

              {filteredLists.length === 0 && (
                <div className="text-center py-8">
                  <p className="text-sm text-muted-foreground">
                    {searchQuery ? "No matching lists found" : "No available lists"}
                  </p>
                </div>
              )}
            </div>
          </ScrollArea>
        </div>
        <DialogFooter>
          <Button 
            variant="outline" 
            onClick={() => onOpenChange(false)}
            disabled={addToListMutation.isPending}
          >
            Cancel
          </Button>
          <Button 
            onClick={handleAddToList} 
            disabled={selectedListReference === null || addToListMutation.isPending}
          >
            {addToListMutation.isPending ? "Adding..." : "Add to List"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
} 