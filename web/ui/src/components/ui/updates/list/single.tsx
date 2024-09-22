import { MalakUpdate } from "@/client/Api";
import { Button } from "@/components/Button";
import { Divider } from "@/components/Divider";
import { RiDeleteBin2Line, RiFileCopyLine, RiMoreLine, RiPushpinLine } from "@remixicon/react";
import UpdateBadge from "../../custom/update/badge";
import {
  Dialog, DialogTrigger, DialogContent, DialogTitle,
  DialogHeader, DialogFooter, DialogDescription,
  DialogClose
} from "@/components/Dialog";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/Popover";

const SingleUpdate = (update: MalakUpdate) => {
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
              <Dialog>
                <DialogTrigger asChild>
                  <Button
                    variant="ghost"
                    className="w-full justify-start rounded-none px-2 py-1.5 text-sm"
                  >
                    <RiFileCopyLine className="mr-2 h-4 w-4" />
                    Duplicate
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Duplicate update</DialogTitle>
                    <DialogDescription className="mt-2">
                      Are you sure you want to duplicate this investor update?
                      A new update containing the exact content of this update will created.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter className="mt-4">
                    <DialogClose asChild>
                      <Button variant="secondary">Cancel</Button>
                    </DialogClose>
                    <Button>Confirm</Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
              <Dialog>
                <DialogTrigger asChild>
                  <Button
                    variant="ghost"
                    className="w-full justify-start rounded-none px-2 py-1.5 text-sm text-red-600"
                  >
                    <RiDeleteBin2Line className="mr-2 h-4 w-4" />
                    Delete
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Confirm Deletion</DialogTitle>
                    <DialogDescription className="mt-2">
                      Are you sure you want to delete this investor update? This action cannot be undone.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter className="mt-4">
                    <DialogClose asChild>
                      <Button variant="secondary">Cancel</Button>
                    </DialogClose>
                    <Button variant="destructive">Delete</Button>
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
