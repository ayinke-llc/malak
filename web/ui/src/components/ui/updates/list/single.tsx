"use client";

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
import {
  EVENT_TOGGLE_PINNED_STATE,
  EVENT_UPDATE_DELETE,
  EVENT_UPDATE_DUPLICATE,
} from "@/lib/analytics-constansts";
import client from "@/lib/client";
import {
  DELETE_UPDATE,
  DUPLICATE_UPDATE,
  LIST_UPDATES,
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
      posthog?.capture(EVENT_UPDATE_DUPLICATE);
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
      posthog?.capture(EVENT_UPDATE_DELETE);
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
      posthog?.capture(EVENT_TOGGLE_PINNED_STATE);
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
      toast.success(`${resp.data.message}. refresh page to view pinned items`);
    },
  });


  return (
    <>
      <div
        key={update.id}
        className="flex items-center justify-between p-2"
      >
        <div className="flex flex-col space-y-1">
          <div className="flex items-center space-x-2">
            <Link href={`/updates/${update.reference}`}>
              <h3 className="font-semibold text-foreground">
                {update.title ? decodeHtmlEntities(update.title as string) : ''}
              </h3>
            </Link>
            <UpdateBadge status={update.status as string} />
          </div>
          <p className="text-sm text-muted-foreground">
            {formatInTimeZone(
              new Date(update?.created_at as string),
              current?.timezone || 'UTC',
              "EEEE, MMMM do, yyyy"
            )}
          </p>
        </div>
        <div className="flex space-x-2">
          <Button
            variant="ghost"
            size="icon"
            aria-label="Pin update"
            loading={loading}
            onClick={() => {
              togglePinnedStatus.mutate(update.reference as string);
            }}
          >
            <RiPushpinLine
              className="h-4 w-4"
              color={update?.is_pinned as boolean ? "red" : "currentColor"} />
          </Button>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="icon" aria-label="More options">
                <RiMoreLine className="h-4 w-4 text-muted-foreground" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-40 p-0">
              <Dialog open={duplicateDialogOpen}>
                <Button
                  variant="ghost"
                  className="w-full justify-start rounded-none px-2 py-1.5 text-sm text-foreground"
                  onClick={() => {
                    setDuplicateDialogOpen(true);
                  }}
                >
                  <RiFileCopyLine className="mr-2 h-4 w-4 text-muted-foreground" />
                  Duplicate
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
                  className="w-full justify-start rounded-none px-2 py-1.5 text-sm text-red-600"
                  onClick={() => {
                    setDeleteDialogOpen(true);
                  }}
                >
                  <RiDeleteBin2Line className="mr-2 h-4 w-4" />
                  Delete
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
            </PopoverContent>
          </Popover>
        </div>
      </div>
      <Separator />
    </>
  );
};

export default SingleUpdate;
