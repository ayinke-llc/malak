"use client"

import type { ServerAPIStatus } from "@/client/Api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import client from "@/lib/client";
import { CREATE_WORKSPACE } from "@/lib/query-constants";
import useWorkspacesStore from "@/store/workspace";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiAddLargeLine } from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useState, useEffect } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";

export type ModalProps = {
  itemName: string;
  onSelect: () => void;
  onOpenChange: (open: boolean) => void;
  forceOpen?: boolean;
};

type CreateWorkspaceInput = {
  name: string;
};

const schema = yup
  .object({
    name: yup.string().min(5).max(50).required(),
  })
  .required();

export function ModalAddWorkspace({
  itemName,
  onSelect,
  onOpenChange,
  forceOpen = false,
}: ModalProps) {

  const [loading, setLoading] = useState<boolean>(false);
  const appendWorkspace = useWorkspacesStore(state => state.appendWorkspaceAfterCreation)
  const [open, setOpen] = useState(forceOpen);
  const setCurrent = useWorkspacesStore(state => state.setCurrent)

  useEffect(() => {
    setOpen(forceOpen);
  }, [forceOpen]);

  useEffect(() => {
    if (open || forceOpen) {
      // Small delay to ensure the modal is rendered
      const timeoutId = setTimeout(() => {
        const input = document.getElementById('workspace-name');
        input?.focus();
      }, 100);

      return () => clearTimeout(timeoutId);
    }
  }, [open, forceOpen]);

  const handleOpenChange = (isOpen: boolean) => {
    if (!forceOpen) {
      setOpen(isOpen);
      onOpenChange(isOpen);
    }
  };


  const mutation = useMutation({
    mutationKey: [CREATE_WORKSPACE],
    mutationFn: (data: CreateWorkspaceInput) =>
      client.workspaces.workspacesCreate(data),
    onSuccess: ({ data }) => {
      setCurrent(data.workspace)
      appendWorkspace(data.workspace);
      toast.success(data.message);
      onOpenChange(false);
      setOpen(false);
      window.location.reload();
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response?.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onSettled: () => setLoading(false),
  });

  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm({
    resolver: yupResolver(schema),
  });

  const onSubmit: SubmitHandler<CreateWorkspaceInput> = (data) => {
    setLoading(true);
    mutation.mutate(data);
  };

  return (
    <Dialog onOpenChange={handleOpenChange} open={forceOpen || open}>
      {!forceOpen && (
        <DialogTrigger className="w-full text-left" asChild>
          <DropdownMenuItem
            className="hover:cursor-pointer"
            onSelect={(event) => {
              event.preventDefault();
              onSelect?.();
            }}
          >
            <div className="flex size-6 items-center justify-center rounded-md border bg-background">
              <RiAddLargeLine className="size-4" />
            </div>
            {itemName}
          </DropdownMenuItem>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-lg">
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col gap-y-1"
        >
          <DialogHeader>
            <DialogTitle>Add new workspace</DialogTitle>
            <DialogDescription className="mt-1 text-sm leading-6">
              Get started with connecting and building relationships with your
              investors
            </DialogDescription>
            <div className="mt-4 grid grid-cols-2 gap-4">
              <div className="col-span-full">
                <Label htmlFor="workspace-name" className="font-medium">
                  Workspace name
                </Label>
                <Input
                  id="workspace-name"
                  placeholder="Ayinke Ventures"
                  className="mt-2"
                  {...register("name")}
                />
                <p className="mt-2 text-xs text-gray-500">
                  Please provide the name of your product, startup or company
                </p>
                {errors.name && (
                  <p className="mt-4 text-xs text-red-600 dark:text-red-500">
                    <span className="font-medium">{errors.name.message}</span>
                  </p>
                )}
              </div>
            </div>
          </DialogHeader>
          <DialogFooter className="mt-6">
            {!forceOpen && (
              <DialogClose asChild>
                <Button
                  className="mt-2 w-full sm:mt-0 sm:w-fit"
                  variant="secondary"
                >
                  Go back
                </Button>
              </DialogClose>
            )}
            <Button type="submit" className="w-full sm:w-fit" disabled={loading}>
              Add workspace
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
