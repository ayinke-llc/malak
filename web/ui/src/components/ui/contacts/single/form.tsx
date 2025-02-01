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
import { Input } from "@/components/ui/input"
import { RiPencilLine } from "@remixicon/react"

const schema = yup.object().shape({
  company: yup.string().min(5).max(100),
  notes: yup.string().min(5).max(3000),
  phone: yup.string(),
  first_name: yup.string(),
  last_name: yup.string(),
})

type FormData = yup.InferType<typeof schema>;

export function EditContactDialog() {
  const [open, setOpen] = useState(false)

  const {
    control,
    handleSubmit,
    formState: { errors },
    setValue,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: {
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
          <DialogFooter className="mt-5">
            <Button type="submit">Submit</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

