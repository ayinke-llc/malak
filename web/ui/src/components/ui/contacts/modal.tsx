import { Button } from "@/components/Button"
import { RiAddLine } from "@remixicon/react"
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
import { Input } from "@/components/Input"
import { Label } from "@/components/Label"
import { useForm, SubmitHandler } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import * as yup from "yup"
import { useState } from "react"
import { useMutation } from "@tanstack/react-query"
import { toast } from "sonner"
import { useRouter } from "next/navigation"
import { ServerAPIStatus } from "@/client/Api"
import client from "@/lib/client"
import { AxiosError } from "axios"
import useAuthStore from "@/store/auth"

export type ModalProps = {
  itemName: string
  onOpenChange: (open: boolean) => void
}

type CreateWorkspaceInput = {
  name: string
}

const schema = yup
  .object({
    name: yup.string().min(5).max(50).required(),
  })
  .required()

export default function CreateContactModal({
  itemName,
  onOpenChange,
}: ModalProps) {

  const [loading, setLoading] = useState<boolean>(false)

  const setWorkspace = useAuthStore.getState().setWorkspace

  const router = useRouter()

  const mutation = useMutation({
    mutationKey: ["create-workspace"],
    mutationFn: (data: CreateWorkspaceInput) => client.workspaces.workspacesCreate(data),
    onSuccess: ({ data }) => {
      setWorkspace(data.workspace)
      toast.success(data.message)
      onOpenChange(false)
      router.push("/")
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response?.data.message
      }
      toast.error(msg)
    },
    retry: false,
    gcTime: Infinity,
    onSettled: () => setLoading(false),
  })

  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm({
    resolver: yupResolver(schema),
  })

  const onSubmit: SubmitHandler<CreateWorkspaceInput> = (data) => {
    setLoading(true)
    mutation.mutate(data)
  }

  return (
    <>
      <Dialog onOpenChange={onOpenChange}>
        <DialogTrigger className="w-full text-right">
          <Button type="button"
            variant="primary"
            className="items-center justify-center whitespace-nowrap">
            <RiAddLine />
            Add User
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-2xl max-w-xs">
          <form onSubmit={handleSubmit(onSubmit)}>
            <DialogHeader>
              <DialogTitle>Add new workspace</DialogTitle>
              <DialogDescription className="mt-1 text-sm leading-6">
                Get started with connecting and building relationships with your investors
              </DialogDescription>
              <div className="mt-4 grid grid-cols-1 gap-4">
                <div>
                  <Label htmlFor="workspace-name" className="font-medium">
                    Contact email
                  </Label>
                  <Input
                    id="workspace-name"
                    name="workspace-name"
                    placeholder="my_workspace"
                    className="mt-2"
                  />
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
              <Button type="submit" className="w-full sm:w-fit">
                Add workspace
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog >
    </>
  )
}
