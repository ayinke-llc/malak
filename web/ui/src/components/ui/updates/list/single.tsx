"use client";

import { AnalyticsEvent } from "@/lib/events";
import type {
  MalakUpdate,
  ServerAPIStatus,
  ServerCreatedUpdateResponse,
} from "@/client/Api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import client from "@/lib/client";
import {
  DELETE_UPDATE,
  DUPLICATE_UPDATE,
  LIST_UPDATES,
  LIST_PINNED_UPDATES,
  TOGGLE_PINNED_STATE,
} from "@/lib/query-constants";
import {
  RiDeleteBin2Line,
  RiFileCopyLine,
  RiMoreLine,
  RiPushpinLine,
} from "@remixicon/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AxiosError, AxiosResponse } from "axios";
import { useRouter } from "next/navigation";
import { usePostHog } from "posthog-js/react";
import { useState } from "react";
import { toast } from "sonner";
import UpdateBadge from "../../custom/update/badge";
import Link from "next/link";
import { Separator } from "../../separator";
import { formatInTimeZone } from 'date-fns-tz';
import useWorkspacesStore from "@/store/workspace";

// Function to decode HTML entities
const decodeHtmlEntities = (text: string) => {
  const textarea = document.createElement('textarea');
  textarea.innerHTML = text;
  return textarea.value;
};

const SingleUpdate = (update: MalakUpdate) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [duplicateDialogOpen, setDuplicateDialogOpen] =
    useState<boolean>(false);

  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);

  const queryClient = useQueryClient();

  const router = useRouter();

  const posthog = usePostHog();

  const current = useWorkspacesStore((state) => state.current);

  const duplicateMutation = useMutation({
    mutationKey: [DUPLICATE_UPDATE],
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onMutate: () => {
      posthog?.capture(AnalyticsEvent.DuplicateUpdate);
    },
    onSettled: () => setLoading(false),
    mutationFn: (reference: string) =>
      client.workspaces.duplicateUpdate(reference),
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    onSuccess: (resp: AxiosResponse<ServerCreatedUpdateResponse>) => {
      toast.success(resp.data.message);
      setDuplicateDialogOpen(false);
      router.push(`/updates/${resp.data.update.reference}`);
    },
  });

  const deletionMutation = useMutation({
    mutationKey: [DELETE_UPDATE],
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onMutate: () => {
      posthog?.capture(AnalyticsEvent.DeleteUpdate);
    },
    onSettled: () => setLoading(false),
    mutationFn: (reference: string) =>
      client.workspaces.deleteUpdate(reference),
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    onSuccess: (resp: AxiosResponse<ServerAPIStatus>) => {
      queryClient.invalidateQueries({ queryKey: [LIST_UPDATES] });
      setDeleteDialogOpen(false);
      toast.success(resp.data.message);
    },
  });

  const togglePinnedStatus = useMutation({
    mutationKey: [TOGGLE_PINNED_STATE],
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onMutate: () => {
      posthog?.capture(AnalyticsEvent.TogglePinnedState);
    },
    onSettled: () => setLoading(false),
    mutationFn: (reference: string) =>
      client.workspaces.toggleUpdatePin(reference),
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    onSuccess: (resp: AxiosResponse<ServerCreatedUpdateResponse>) => {
      queryClient.invalidateQueries({ queryKey: [LIST_UPDATES] });
      queryClient.invalidateQueries({ queryKey: [LIST_PINNED_UPDATES] });
      toast.success(resp.data.message);
    },
  });

  return (
    <>
      <div
        key={update.id}
        className="group relative flex items-start justify-between p-4 hover:bg-muted/50 transition-all duration-200 rounded-lg"
      >
        <div className="flex flex-col space-y-2 flex-1">
          <div className="flex items-center gap-3">
            <Link 
              href={`/updates/${update.reference}`}
              className="flex-1 min-w-0"
            >
              <h3 className="font-medium text-foreground truncate hover:text-primary transition-colors">
                {update.title ? decodeHtmlEntities(update.title as string) : ''}
              </h3>
            </Link>
            <UpdateBadge status={update.status as string} />
          </div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
              <path fillRule="evenodd" d="M5 1.5a1.5 1.5 0 00-1.5 1.5v11.5A1.5 1.5 0 005 16h10a1.5 1.5 0 001.5-1.5V3A1.5 1.5 0 0015 1.5H5zM2.5 3a3 3 0 013-3h10a3 3 0 013 3v11.5a3 3 0 01-3 3H5a3 3 0 01-3-3V3z" clipRule="evenodd" />
              <path fillRule="evenodd" d="M12 6.5a.5.5 0 01.5-.5h2a.5.5 0 010 1h-2a.5.5 0 01-.5-.5zm-6 0a.5.5 0 01.5-.5h2a.5.5 0 010 1h-2a.5.5 0 01-.5-.5zm6 3a.5.5 0 01.5-.5h2a.5.5 0 010 1h-2a.5.5 0 01-.5-.5zm-6 0a.5.5 0 01.5-.5h2a.5.5 0 010 1h-2a.5.5 0 01-.5-.5zm6 3a.5.5 0 01.5-.5h2a.5.5 0 010 1h-2a.5.5 0 01-.5-.5zm-6 0a.5.5 0 01.5-.5h2a.5.5 0 010 1h-2a.5.5 0 01-.5-.5z" clipRule="evenodd" />
            </svg>
            <span>
              {formatInTimeZone(
                new Date(update?.created_at as string),
                current?.timezone || 'UTC',
                "EEEE, MMMM do, yyyy"
              )}
            </span>
          </div>
        </div>
        <div className="flex items-center gap-1 ml-4 opacity-0 group-hover:opacity-100 transition-opacity">
          <Button
            variant="ghost"
            size="icon"
            aria-label="Pin update"
            loading={loading}
            onClick={() => {
              togglePinnedStatus.mutate(update.reference as string);
            }}
            className="h-8 w-8 hover:bg-muted"
          >
            <RiPushpinLine
              className="h-4 w-4"
              color={update?.is_pinned as boolean ? "red" : "currentColor"} />
          </Button>
          <Popover>
            <PopoverTrigger asChild>
              <Button 
                variant="ghost" 
                size="icon" 
                aria-label="More options"
                className="h-8 w-8 hover:bg-muted"
              >
                <RiMoreLine className="h-4 w-4 text-muted-foreground" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-48 p-0" align="end">
              <div className="py-1">
                <Dialog open={duplicateDialogOpen}>
                  <Button
                    variant="ghost"
                    className="w-full justify-start px-3 py-2 text-sm text-foreground hover:bg-muted"
                    onClick={() => {
                      setDuplicateDialogOpen(true);
                    }}
                  >
                    <RiFileCopyLine className="mr-2 h-4 w-4 text-muted-foreground" />
                    Duplicate update
                  </Button>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Duplicate update</DialogTitle>
                      <DialogDescription className="mt-2 text-muted-foreground">
                        Are you sure you want to duplicate this investor update? A
                        new update containing the exact content of this update
                        will created.
                      </DialogDescription>
                    </DialogHeader>
                    <DialogFooter className="mt-4">
                      <Button
                        variant="secondary"
                        loading={loading}
                        onClick={() => {
                          setDuplicateDialogOpen(false);
                        }}
                      >
                        Cancel
                      </Button>
                      <Button
                        loading={loading}
                        onClick={() => {
                          setLoading(true);
                          duplicateMutation.mutate(update.reference as string);
                        }}
                      >
                        Duplicate
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
                <Dialog open={deleteDialogOpen}>
                  <Button
                    variant="ghost"
                    className="w-full justify-start px-3 py-2 text-sm text-red-600 hover:bg-muted"
                    onClick={() => {
                      setDeleteDialogOpen(true);
                    }}
                    disabled={update.status === "sent"}
                  >
                    <RiDeleteBin2Line className="mr-2 h-4 w-4" />
                    Delete update
                  </Button>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Confirm Deletion</DialogTitle>
                      <DialogDescription className="mt-2 text-muted-foreground">
                        Are you sure you want to delete this investor update? This
                        action cannot be undone.
                      </DialogDescription>
                    </DialogHeader>
                    <DialogFooter className="mt-4">
                      <Button
                        variant="secondary"
                        loading={loading}
                        onClick={() => {
                          setDeleteDialogOpen(false);
                        }}
                      >
                        Cancel
                      </Button>
                      <Button
                        variant="destructive"
                        loading={loading}
                        onClick={() => {
                          setLoading(true);
                          deletionMutation.mutate(update.reference as string);
                        }}
                      >
                        Delete
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </PopoverContent>
          </Popover>
        </div>
      </div>
      <Separator className="opacity-50" />
    </>
  );
};

export default SingleUpdate;
