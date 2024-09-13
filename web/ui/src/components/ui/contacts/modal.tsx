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
import { ServerAPIStatus } from "@/client/Api"
import client from "@/lib/client"
import { AxiosError } from "axios"
import { CREATE_CONTACT_MUTATION } from "@/lib/query-constants"

export type ModalProps = {
}

type CreateContactInput = {
  email: string
}

const schema = yup
  .object({
    email: yup.string().required().email(),
  })
  .required()

export default function CreateContactModal({
}: ModalProps) {

  const [loading, setLoading] = useState<boolean>(false)

  const [hasOpenDialog, setHasOpenDialog] = useState(false)

  const handleDialogItemOpenChange = (open: boolean) => {
    setHasOpenDialog(open)
  }

  const mutation = useMutation({
    mutationKey: [CREATE_CONTACT_MUTATION],
    mutationFn: (data: CreateContactInput) => client.contacts.contactsCreate(data),
    onSuccess: ({ data }) => {
      toast.success(data.message)
      handleDialogItemOpenChange(false)
      reset()
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response.data.message
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
    reset,
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: {
      email: "",
    }
  })

  const onSubmit: SubmitHandler<CreateContactInput> = (data) => {
    setLoading(true)
    mutation.mutate(data)
  }

  return (
    <div className="flex justify-center">
      <Dialog onOpenChange={handleDialogItemOpenChange} open={hasOpenDialog}>
        <DialogTrigger asChild>
          <div className="w-full text-right" >
            <Button type="button"
              variant="primary"
              className="whitespace-nowrap">
              <RiAddLine />
              Add User
            </Button>
          </div>
        </DialogTrigger>
        <DialogContent className="sm:max-w-lg">
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-y-1">
            <DialogHeader>
              <DialogTitle>Add a new contact</DialogTitle>
              <DialogDescription className="mt-1 text-sm leading-6">
                Get started with connecting and building relationships with a
                specific investor
              </DialogDescription>
              <div className="mt-4">
                <Label htmlFor="workspace-name" className="font-medium">
                  Email address
                </Label>
                <Input
                  id="email"
                  placeholder="yo@lanre.wtf"
                  className="mt-2"
                  type="email"
                  {...register("email")}
                />
                {errors.email && (
                  < p className="mt-4 text-xs text-red-600 dark:text-red-500">
                    <span className="font-medium">
                      {errors.email.message}
                    </span>
                  </p>
                )}
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
              <Button type="submit" className="w-full sm:w-fit"
                isLoading={loading}>
                Add Contact
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog >
    </div>
  )
}
