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
import { MalakIntegrationChartType, MalakIntegrationDataPointType } from "@/client/Api";
import type { MalakWorkspaceIntegration, ServerAPIStatus } from "@/client/Api";
import { CREATE_CHART, LIST_CHARTS, FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import client from "@/lib/client";
import { AxiosError } from "axios";

interface CreateChartFormData {
  title: string;
  type: MalakIntegrationChartType.IntegrationChartTypeBar;
  datapoint: MalakIntegrationDataPointType;
}

const createChartSchema = yup.object({
  title: yup.string().required("Chart title is required"),
  type: yup.string().oneOf([MalakIntegrationChartType.IntegrationChartTypeBar], "Invalid chart type").required("Chart type is required"),
  datapoint: yup.string().oneOf([MalakIntegrationDataPointType.IntegrationDataPointTypeCurrency, MalakIntegrationDataPointType.IntegrationDataPointTypeOthers], "Invalid datapoint type").required("Datapoint type is required"),
});

export function CreateChartDialog({ integration }: { integration: MalakWorkspaceIntegration }) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const form = useForm<CreateChartFormData>({
    resolver: yupResolver(createChartSchema),
    defaultValues: {
      title: "",
      type: MalakIntegrationChartType.IntegrationChartTypeBar,
      datapoint: "" as MalakIntegrationDataPointType
    }
  });

  const { mutate: createChart, isPending } = useMutation({
    mutationKey: [CREATE_CHART],
    mutationFn: async (data: CreateChartFormData) => {
      return client.workspaces.integrationsChartsCreate(integration.reference!, {
        title: data.title,
        chart_type: data.type,
        datapoint: data.datapoint
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

    if (!Object.values(MalakIntegrationChartType).includes(data.type)) {
      toast.error("Invalid chart type selected");
      return;
    }


    if (!Object.values(MalakIntegrationDataPointType).includes(data.datapoint)) {
      toast.error("Invalid datapoint type selected");
      return;
    }

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
              onValueChange={(value) => form.setValue("type", MalakIntegrationChartType.IntegrationChartTypeBar)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={MalakIntegrationChartType.IntegrationChartTypeBar}>Bar Chart</SelectItem>
              </SelectContent>
            </Select>
            {form.formState.errors.type && (
              <p className="text-sm text-destructive mt-1">{form.formState.errors.type.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="datapoint">Datapoint Type</Label>
            <Select
              value={form.watch("datapoint")}
              onValueChange={(value: MalakIntegrationDataPointType) => form.setValue("datapoint", value)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={MalakIntegrationDataPointType.IntegrationDataPointTypeCurrency}>Currency</SelectItem>
                <SelectItem value={MalakIntegrationDataPointType.IntegrationDataPointTypeOthers}>Others</SelectItem>
              </SelectContent>
            </Select>
            {form.formState.errors.datapoint && (
              <p className="text-sm text-destructive mt-1">{form.formState.errors.datapoint.message}</p>
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
