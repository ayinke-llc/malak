import { Button } from "@/components/Button";
import { MalakUpdate, ServerAPIStatus, ServerCreatedUpdateResponse } from "@/client/Api";
import { Divider } from "@/components/Divider";
import { RiDeleteBin2Line, RiFileCopyLine, RiMoreLine, RiPushpinLine } from "@remixicon/react";
import UpdateBadge from "../../custom/update/badge";
import {
  Dialog, DialogContent, DialogTitle,
  DialogHeader, DialogFooter, DialogDescription
} from "@/components/Dialog";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/Popover";
import { useMutation } from "@tanstack/react-query";
import { DELETE_UPDATE, DUPLICATE_UPDATE } from "@/lib/query-constants";
import { useState } from "react";
import client from "@/lib/client";
import { AxiosError, AxiosResponse } from "axios";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { usePostHog } from "posthog-js/react";
import { EVENT_UPDATE_DELETE, EVENT_UPDATE_DUPLICATE } from "@/lib/analytics-constansts";

const SingleUpdate = (update: MalakUpdate) => {

  const [loading, setLoading] = useState<boolean>(false)
  const [deleted, setDeleted] = useState<boolean>(false)
  const [duplicateDialogOpen, setDuplicateDialogOpen] = useState<boolean>(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false)

  const router = useRouter()

  const posthog = usePostHog()

  const duplicateMutation = useMutation({
    mutationKey: [DUPLICATE_UPDATE],
    retry: false,
    gcTime: Infinity,
    onMutate: () => {
      posthog?.capture(EVENT_UPDATE_DUPLICATE)
    },
    onSettled: () => setLoading(false),
    mutationFn: (reference: string) => client.workspaces.duplicateUpdate(reference),
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response.data.message
      }
      toast.error(msg)
    },
    onSuccess: (resp: AxiosResponse<ServerCreatedUpdateResponse>) => {
      toast.success(resp.data.message)
      setDuplicateDialogOpen(false)
      router.push(`/updates/${resp.data.update.reference}`)
    }
  })

  const deletionMutation = useMutation({
    mutationKey: [DELETE_UPDATE],
    retry: false,
    gcTime: Infinity,
    onMutate: () => {
      posthog?.capture(EVENT_UPDATE_DELETE)
    },
    onSettled: () => setLoading(false),
    mutationFn: (reference: string) => client.workspaces.deleteUpdate(reference),
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response.data.message
      }
      toast.error(msg)
    },
    onSuccess: (resp: AxiosResponse<ServerAPIStatus>) => {
      setDeleteDialogOpen(false)
      setDeleted(true)
      toast.success(resp.data.message)
    }
  })

  if (deleted) {
    // essentially remove it from the list
    return null
  }

  return (
    <>
      <div key={update.id}
        className="flex items-center justify-between p-2 hover:bg-accent rounded-lg transition-colors">
        <div className="flex flex-col space-y-1">
          <div className="flex items-center space-x-2">
            <h3 className="font-semibold">{update.reference}</h3>
            <UpdateBadge status={update.status as string} />
          </div>
          <p className="text-sm text-muted-foreground">{update.created_at}</p>
        </div>
        <div className="flex space-x-2">
          <Button variant="ghost" size="icon" aria-label="Pin update">
            <RiPushpinLine className="h-4 w-4" />
          </Button>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="icon" aria-label="More options">
                <RiMoreLine className="h-4 w-4" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-40 p-0">
              <Dialog open={duplicateDialogOpen}>
                <Button
                  variant="ghost"
                  className="w-full justify-start rounded-none px-2 py-1.5 text-sm"
                  onClick={() => {
                    setDuplicateDialogOpen(true)
                  }}
                >
                  <RiFileCopyLine className="mr-2 h-4 w-4" />
                  Duplicate
                </Button>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Duplicate update</DialogTitle>
                    <DialogDescription className="mt-2">
                      Are you sure you want to duplicate this investor update?
                      A new update containing the exact content of this update will created.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter className="mt-4">
                    <Button
                      variant="secondary"
                      isLoading={loading}
                      onClick={() => {
                        setDuplicateDialogOpen(false)
                      }}
                    >
                      Cancel
                    </Button>
                    <Button
                      loadingText="Duplicating"
                      isLoading={loading}
                      onClick={() => {
                        setLoading(true)
                        duplicateMutation.mutate(update.reference as string)
                      }}>
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
                    setDeleteDialogOpen(true)
                  }}
                >
                  <RiDeleteBin2Line className="mr-2 h-4 w-4" />
                  Delete
                </Button>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Confirm Deletion</DialogTitle>
                    <DialogDescription className="mt-2">
                      Are you sure you want to delete this investor update? This action cannot be undone.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter className="mt-4">
                    <Button
                      variant="secondary"
                      isLoading={loading}
                      onClick={() => {
                        setDeleteDialogOpen(false)
                      }}
                    >
                      Cancel
                    </Button>
                    <Button variant="destructive"
                      loadingText="Deleting"
                      isLoading={loading}
                      onClick={() => {
                        setLoading(true)
                        deletionMutation.mutate(update.reference as string)
                      }}>Delete</Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </PopoverContent>
          </Popover>
        </div>
      </div>
      <Divider />
    </>
  )
}

export default SingleUpdate;
