import { useState } from "react";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { RiAddLine } from "@remixicon/react";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { MalakIntegrationChartType } from "@/client/Api";
import type { MalakWorkspaceIntegration, ServerAPIStatus } from "@/client/Api";
import { CREATE_CHART, LIST_CHARTS, FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import client from "@/lib/client";
import { AxiosError } from "axios";

interface CreateChartFormData {
  title: string;
  type: MalakIntegrationChartType;
}

const createChartSchema = yup.object({
  title: yup.string().required("Chart title is required"),
  type: yup.string().oneOf([MalakIntegrationChartType.IntegrationChartTypeBar, MalakIntegrationChartType.IntegrationChartTypePie], "Invalid chart type").required("Chart type is required"),
});

export function CreateChartDialog({ integration }: { integration: MalakWorkspaceIntegration }) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();
  
  const form = useForm<CreateChartFormData>({
    resolver: yupResolver(createChartSchema),
    defaultValues: {
      title: "",
      type: MalakIntegrationChartType.IntegrationChartTypeBar
    }
  });

  const { mutate: createChart, isPending } = useMutation({
    mutationKey: [CREATE_CHART],
    mutationFn: async (data: CreateChartFormData) => {
      return client.workspaces.integrationsChartsCreate(integration.reference!, {
        title: data.title,
        chart_type: data.type
      });
    },
    onSuccess: ({ data }) => {
      toast.success("Chart created successfully");
      queryClient.invalidateQueries({ queryKey: [LIST_CHARTS] });
      queryClient.invalidateQueries({ queryKey: [FETCH_CHART_DATA_POINTS] });
      setOpen(false);
      form.reset();
    },
    onError: (error: AxiosError<ServerAPIStatus>) => {
      toast.error(error.response?.data?.message || "Failed to create chart. Please try again.");
    }
  });

  const handleCreateChart = (data: CreateChartFormData) => {
    createChart(data);
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => {
      setOpen(isOpen);
      if (!isOpen) {
        form.reset();
      }
    }}>
      <DialogTrigger asChild>
        <Button className="gap-2">
          <RiAddLine className="h-4 w-4" />
          Create Chart
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Chart</DialogTitle>
        </DialogHeader>
        <form onSubmit={form.handleSubmit(handleCreateChart)} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="title">Chart Title</Label>
            <Input
              id="title"
              placeholder="Enter chart title"
              {...form.register("title")}
            />
            {form.formState.errors.title && (
              <p className="text-sm text-destructive mt-1">{form.formState.errors.title.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="type">Chart Type</Label>
            <Select 
              value={form.watch("type")} 
              onValueChange={(value: MalakIntegrationChartType) => form.setValue("type", value)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={MalakIntegrationChartType.IntegrationChartTypeBar}>Bar Chart</SelectItem>
                <SelectItem value={MalakIntegrationChartType.IntegrationChartTypePie}>Pie Chart</SelectItem>
              </SelectContent>
            </Select>
            {form.formState.errors.type && (
              <p className="text-sm text-destructive mt-1">{form.formState.errors.type.message}</p>
            )}
          </div>
          <div className="flex gap-2 justify-end">
            <Button 
              type="button" 
              variant="destructive" 
              onClick={() => {
                setOpen(false);
                form.reset();
              }}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? "Creating..." : "Create Chart"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 