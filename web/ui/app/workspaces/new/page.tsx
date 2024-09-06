"use client"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useForm, SubmitHandler } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import * as yup from "yup"
import { Skeleton } from "@/components/ui/skeleton"
import { useState } from "react"
import client from "@/lib/client"
import { useMutation } from "@tanstack/react-query"
import { toast } from "sonner"
import { AxiosError } from "axios"
import { ServerAPIStatus } from "@/client/Api"

type CreateWorkspaceInput = {
  name: string
}

const schema = yup
  .object({
    name: yup.string().min(2).max(50).required(),
  })
  .required()

export default function Page() {

  const [loading, setLoading] = useState<boolean>(false)
  const [dialogOpened, setDialogOpened] = useState<boolean>(false)

  const mutation = useMutation({
    mutationKey: ["create-workspace"],
    mutationFn: (data: CreateWorkspaceInput) => client.workspaces.workspacesCreate(data),
    onSuccess: ({ data }) => {
      console.log(data.workspace)
      onDialogTrigger()
      toast.success(data.message)
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

  const onDialogTrigger = () => {
    setDialogOpened(!dialogOpened)
  }

  return (
    <div className="flex flex-col items-center gap-1 text-center">
      {loading ?
        <div className="columns-2">
          <SkeletonCard />
          <SkeletonCard />
        </div> : (
          <>
            <h3 className="text-2xl font-bold tracking-tight">
              Create a new workspace
            </h3>
            <Dialog open={dialogOpened}>
              <Button variant="outline" className="mt-4" onClick={onDialogTrigger}>Create workspace</Button>
              <DialogContent className="sm:max-w-[425px]" onClose={onDialogTrigger}>
                <form onSubmit={handleSubmit(onSubmit)}>
                  <DialogHeader>
                    <DialogTitle>Create workspace</DialogTitle>
                    <DialogDescription>
                      Create a workspace to start building your relationship with investors.
                      Each workspace should ideally represent a startup/company
                    </DialogDescription>
                  </DialogHeader>
                  <div className="grid gap-4 py-4">
                    <div className="grid grid-cols-4 items-center gap-4">
                      <Label htmlFor="name" className="text-right">
                        Name
                      </Label>

                      <Input id="name"
                        placeholder="Your workspace name"
                        className="col-span-3"
                        {...register("name", { required: true, maxLength: 50, minLength: 5 })} />

                      {errors.name && (
                        <p className="col-span-3 mt-2 text-sm text-red-600 dark:text-red-500">
                          <span className="font-medium">
                            {errors.name.message}
                          </span>
                        </p>
                      )}
                    </div>
                  </div>
                  <DialogFooter>
                    <Button type="submit">Create</Button>
                  </DialogFooter>
                </form>
              </DialogContent>
            </Dialog>
          </>
        )}
    </div>
  )
}

const SkeletonCard = () => {
  return (
    <div className="flex flex-col space-y-3">
      <Skeleton className="h-[125px] w-[250px] rounded-xl" />
      <div className="space-y-2">
        <Skeleton className="h-4 w-[250px]" />
        <Skeleton className="h-4 w-[200px]" />
        <Skeleton className="h-4 w-[200px]" />
        <Skeleton className="h-4 w-[200px]" />
        <Skeleton className="h-4 w-[200px]" />
        <Skeleton className="h-4 w-[200px]" />
        <Skeleton className="h-4 w-[200px]" />
      </div>
    </div>
  )
}
