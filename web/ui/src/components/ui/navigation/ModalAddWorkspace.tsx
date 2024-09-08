import { Badge } from "@/components/Badge"
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
import { RadioCardGroup, RadioCardItem } from "@/components/RadioCard"

export const plans: {
  label: string
  value: string
  description: string
  isRecommended: boolean
}[] = [
    {
      label: "Core plan",
      value: "core",
      description: "Up to 200 investors",
      isRecommended: true,
    },
    {
      label: "Scale",
      value: "scale",
      description: "Up to 1,000 investors",
      isRecommended: false,
    },
  ]

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
                    For best performance, choose a region closest to your
                    application.
                  </p>
                </div>
              </div>
              <div className="mt-4">
                <Label htmlFor="database" className="font-medium">
                  Choose your plan
                </Label>
                <RadioCardGroup
                  defaultValue={plans[0].value}
                  className="mt-2 grid grid-cols-1 gap-4 text-sm md:grid-cols-2"
                >
                  {plans.map((database) => (
                    <RadioCardItem key={database.value} value={database.value}>
                      <div className="flex items-start gap-3">
                        <div>
                          {database.isRecommended ? (
                            <div className="flex items-center gap-2">
                              <span className="leading-5">
                                {database.label}
                              </span>
                              <Badge>Recommended</Badge>
                            </div>
                          ) : (
                            <span>{database.label}</span>
                          )}
                          <p className="mt-1 text-xs text-gray-500">
                            1/8 vCPU, 1 GB RAM
                          </p>
                        </div>
                      </div>
                    </RadioCardItem>
                  ))}
                </RadioCardGroup>
                <p className="mt-2 text-xs text-gray-500">
                  Each plan comes with a 10 day trial
                </p>
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
