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
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RiAddLine } from "@remixicon/react";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError, AxiosResponse } from "axios";
import client from "@/lib/client";
import { ContentType } from "@/client/Api";

type CreateDashboardInput = {
  name: string;
  description?: string;
};

interface CreateDashboardResponse {
  message: string;
  dashboard: {
    id: string;
    name: string;
    description?: string;
  };
}

const schema = yup.object().shape({
  name: yup.string().required("Dashboard name is required"),
  description: yup.string(),
});

export default function CreateDashboardModal() {
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateDashboardInput>({
    resolver: yupResolver(schema),
  });

  const createMutation = useMutation<
    AxiosResponse<CreateDashboardResponse>,
    AxiosError<ServerAPIStatus>,
    CreateDashboardInput
  >({
    mutationKey: ["CREATE_DASHBOARD"],
    mutationFn: async (data) => {
      return client.request<CreateDashboardResponse, ServerAPIStatus>({
        path: `/dashboards`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
      });
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setOpen(false);
      reset();
    },
    onError(err) {
      let msg = err.message;
      if (err.response?.data) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onSettled: () => setLoading(false),
  });

  const onSubmit: SubmitHandler<CreateDashboardInput> = (data) => {
    setLoading(true);
    createMutation.mutate(data);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          type="button"
          variant="default"
          className="whitespace-nowrap gap-1"
        >
          <RiAddLine />
          New dashboard
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-lg">
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col gap-y-4"
        >
          <DialogHeader>
            <DialogTitle>Create new dashboard</DialogTitle>
            <DialogDescription className="mt-1 text-sm leading-6">
              Create a new dashboard to visualize and track your data
            </DialogDescription>
          </DialogHeader>

          <div>
            <Label htmlFor="name">Dashboard name</Label>
            <Input
              id="name"
              placeholder="Financial Metrics"
              className="mt-2"
              {...register("name")}
            />
            {errors.name && (
              <p className="mt-1 text-xs text-red-600 dark:text-red-500">
                <span className="font-medium">{errors.name.message}</span>
              </p>
            )}
          </div>

          <div>
            <Label htmlFor="description">Description</Label>
            <Input
              id="description"
              placeholder="Track key financial metrics and KPIs"
              className="mt-2"
              {...register("description")}
            />
          </div>

          <DialogFooter className="mt-6">
            <DialogClose asChild>
              <Button
                className="mt-2 w-full sm:mt-0 sm:w-fit"
                variant="secondary"
              >
                Cancel
              </Button>
            </DialogClose>
            <Button
              type="submit"
              className="w-full sm:w-fit"
              disabled={loading}
            >
              Create dashboard
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
} 