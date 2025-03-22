import { useState } from "react";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import client from "@/lib/client";
import { MalakIntegrationChart, MalakWorkspaceIntegration } from "@/client/Api";
import { FETCH_CHART_DATA_POINTS } from "@/lib/query-constants";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

interface AddDataPointDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  chart: MalakIntegrationChart;
  workspaceIntegration: MalakWorkspaceIntegration;
}

const schema = yup.object({
  value: yup
    .number()
    .typeError("Please enter a valid number")
    .min(0, "Value must be a non-negative number")
    .required("Value is required"),
});

type FormData = yup.InferType<typeof schema>;

export function AddDataPointDialog({
  open,
  onOpenChange,
  chart,
  workspaceIntegration,
}: AddDataPointDialogProps) {
  const queryClient = useQueryClient();

  const form = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: {
      value: undefined,
    },
  });

  const { mutate: addDataPoint, isPending } = useMutation({
    mutationFn: async (data: FormData) => {
      return client.workspaces.integrationsChartsPointsCreate(
        workspaceIntegration.reference as string,
        chart.reference as string,
        { value: data.value }
      );
    },
    onSuccess: () => {
      toast.success("Data point added successfully");
      onOpenChange(false);
      form.reset();
      // Invalidate queries to refresh the data
      queryClient.invalidateQueries({ queryKey: [FETCH_CHART_DATA_POINTS, chart.reference] });
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to add data point");
    },
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add Data Point</DialogTitle>
          <DialogDescription>
            Add a new data point to {chart.user_facing_name}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit((data) => addDataPoint(data))} className="space-y-4">
            <FormField
              control={form.control}
              name="value"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Value</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      min="0"
                      {...field}
                      onChange={(e) => field.onChange(e.target.valueAsNumber)}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit" disabled={isPending}>
                {isPending ? "Adding..." : "Add Data Point"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
} 