import { Button } from "@/components/Button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/Dialog"
import { DropdownMenuItem } from "@/components/Dropdown"
import { Input } from "@/components/Input"
import { Label } from "@/components/Label"

export type ModalProps = {
  itemName: string
  onSelect: () => void
  onOpenChange: (open: boolean) => void
}

export function ModalAddWorkspace({
  itemName,
  onSelect,
  onOpenChange,
}: ModalProps) {
  return (
    <>
      <Dialog onOpenChange={onOpenChange}>
        <DialogTrigger className="w-full text-left">
          <DropdownMenuItem
            onSelect={(event) => {
              event.preventDefault()
              onSelect && onSelect()
            }}
          >
            {itemName}
          </DropdownMenuItem>
        </DialogTrigger>
        <DialogContent className="sm:max-w-2xl">
          <form>
            <DialogHeader>
              <DialogTitle>Add new workspace</DialogTitle>
              <DialogDescription className="mt-1 text-sm leading-6">
                Get started with connecting and building relationships with your investors
              </DialogDescription>
              <div className="mt-4 grid grid-cols-2 gap-4">
                <div className="col-span-full">
                  <Label htmlFor="workspace-name" className="font-medium">
                    Workspace name
                  </Label>
                  <Input
                    id="workspace-name"
                    name="workspace-name"
                    placeholder="my_workspace"
                    className="mt-2"
                  />
                  <p className="mt-2 text-xs text-gray-500">
                    Please provide the name of your product, startup or company
                  </p>
                </div>
              </div>
            </DialogHeader>
            <DialogFooter className="mt-6">
              <DialogClose asChild>
                <Button
                  className="mt-2 w-full sm:mt-0 sm:w-fit"
                  variant="secondary"
                >
                  Go back
                </Button>
              </DialogClose>
              <DialogClose asChild>
                <Button type="submit" className="w-full sm:w-fit">
                  Add workspace
                </Button>
              </DialogClose>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  )
}
