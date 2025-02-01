"use client"

import { Label } from "@/components/ui/label"
import { useState } from "react"
import { Controller, useForm } from "react-hook-form"
import * as yup from "yup"
import { yupResolver } from "@hookform/resolvers/yup"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Textarea } from "@/components/ui/textarea"
import { Input } from "@/components/ui/input"
import { RiPencilLine } from "@remixicon/react"
import useWorkspacesStore from "@/store/workspace"
import { MalakContact } from "@/client/Api"

const schema = yup.object().shape({
  company: yup.string().min(5).max(100),
  notes: yup.string().min(5).max(3000),
  phone: yup.string(),
  first_name: yup.string(),
  last_name: yup.string(),
  address: yup.string().min(5).max(300),
})

type FormData = yup.InferType<typeof schema>;

export function EditContactDialog({ contact }: { contact: MalakContact }) {
  const [open, setOpen] = useState(false)

  if (!contact) {
    return null
  }

  const current = useWorkspacesStore(state => state.current)

  const {
    control,
    handleSubmit,
    formState: { errors },
    setValue,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: {
      company: contact?.company || "",
      first_name: contact?.first_name || "",
      last_name: contact?.last_name || "",
      notes: contact?.notes || "",
      phone: contact?.phone || "",
      address: contact?.city || "",
    },
  });

  function onSubmit(values: FormData) {
    // Simulate form submission
    console.log(values)
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          className="h-9 w-9"
        >
          <RiPencilLine className="h-4 w-4" />
        </Button>
      </DialogTrigger>
      <DialogContent >
        <DialogHeader>
          <DialogTitle>Update contact details</DialogTitle>
          <DialogDescription>Update contact's details. You cannot update the contact's email</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>

          <div className="flex space-x-4 mt-4">
            <div className="flex-1 space-y-2">
              <Label htmlFor="first_name">First name:</Label>
              <Controller
                name="first_name"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    id="first_name"
                    className="mt-2"
                  />
                )}
              />
              {errors.first_name && (
                <p className="text-sm text-red-500">{errors.first_name.message}</p>
              )}
            </div>

            <div className="flex-1 space-y-2">
              <Label htmlFor="last_name">Last name:</Label>
              <Controller
                name="last_name"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    id="last_name"
                    className="mt-2"
                  />
                )}
              />
              {errors.first_name && (
                <p className="text-sm text-red-500">{errors.first_name.message}</p>
              )}
            </div>

          </div>
          <div className="space-y-2 mt-4">
            <Label htmlFor="companyName">Company name:</Label>
            <Controller
              name="company"
              control={control}
              render={({ field }) => (
                <Input
                  {...field}
                  id="companyName"
                  placeholder="Enter company name"
                  className="mt-2"
                />
              )}
            />
            {errors.company && (
              <p className="text-sm text-red-500">{errors.company.message}</p>
            )}
          </div>

          <div className="space-y-2 mt-4">
            <Label htmlFor="address">Address:</Label>
            <Controller
              name="address"
              control={control}
              render={({ field }) => (
                <Input
                  {...field}
                  id="address"
                  placeholder="Enter address"
                  className="mt-2"
                />
              )}
            />
            {errors.address && (
              <p className="text-sm text-red-500">{errors.address.message}</p>
            )}
          </div>

          <div className="space-y-2 mt-4">
            <Label htmlFor="notes">Notes:</Label>
            <Controller
              name="notes"
              control={control}
              render={({ field }) => (
                <Textarea
                  {...field}
                  id="notes"
                  placeholder="Enter notes for this contact"
                  className="mt-2"
                />
              )}
            />
            {errors.notes && (
              <p className="text-sm text-red-500">{errors.notes.message}</p>
            )}
          </div>
          <DialogFooter className="mt-5">
            <Button type="submit">Submit</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

