import { useState } from "react";
import { Button } from "@/components/ui/button";
import { RiAddLine, RiTeamLine, RiArrowUpSLine, RiArrowDownSLine } from '@remixicon/react';
import { toast } from "sonner";
import { AddToListDialog } from "./add-to-list-dialog";
import { MalakContact, MalakContactListMapping } from "@/client/Api";

interface ContactListsViewProps {
  contact: MalakContact;
}

export function ContactListsView({ contact }: ContactListsViewProps) {
  const [showAddToListDialog, setShowAddToListDialog] = useState(false);
  const [showAllLists, setShowAllLists] = useState(false);

  const currentLists = contact.lists || [];

  return (
    <div className="mt-6">
      <div className="grid gap-6">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold">Contact Lists</h3>
            <p className="text-sm text-muted-foreground mt-1">
              Member of {currentLists.length} {currentLists.length === 1 ? 'list' : 'lists'}
            </p>
          </div>
          <Button 
            variant="outline" 
            size="sm" 
            className="flex items-center gap-2"
            onClick={() => setShowAddToListDialog(true)}
          >
            <RiAddLine className="h-4 w-4" />
            Add to List
          </Button>
        </div>

        {currentLists.length === 0 ? (
          <div className="text-center py-8 border rounded-lg bg-muted/10">
            <RiTeamLine className="h-8 w-8 text-muted-foreground mx-auto mb-3" />
            <p className="text-sm text-muted-foreground">
              This contact is not a member of any lists yet
            </p>
            <Button
              variant="link"
              className="mt-2"
              onClick={() => setShowAddToListDialog(true)}
            >
              Add to a list
            </Button>
          </div>
        ) : (
          <>
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Current Lists
                </span>
              </div>
            </div>

            <div className="grid gap-4">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {currentLists
                  .slice(0, showAllLists ? undefined : 6)
                  .map((list: MalakContactListMapping) => (
                    <div
                      key={list.list_id}
                      className="flex flex-col p-4 rounded-lg border bg-card hover:bg-accent/5 transition-colors"
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex items-center gap-3">
                          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center text-white shrink-0">
                            <RiTeamLine className="h-4 w-4" />
                          </div>
                          <div className="space-y-1">
                            <h4 className="font-medium text-foreground line-clamp-1">{list.list?.title}</h4>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
              </div>

              {currentLists.length > 6 && (
                <Button
                  variant="outline"
                  className="w-full mt-2"
                  onClick={() => setShowAllLists(!showAllLists)}
                >
                  {showAllLists ? (
                    <span className="flex items-center gap-2">
                      Show Less Lists
                      <RiArrowUpSLine className="h-4 w-4" />
                    </span>
                  ) : (
                    <span className="flex items-center gap-2">
                      Show {currentLists.length - 6} More Lists
                      <RiArrowDownSLine className="h-4 w-4" />
                    </span>
                  )}
                </Button>
              )}
            </div>
          </>
        )}
      </div>

      <AddToListDialog
        open={showAddToListDialog}
        onOpenChange={setShowAddToListDialog}
        contactReference={contact.reference || ""}
        currentListIds={currentLists
          .map(list => list.list_id ? Number(list.list_id) : undefined)
          .filter((id): id is number => id !== undefined)}
      />
    </div>
  );
} 