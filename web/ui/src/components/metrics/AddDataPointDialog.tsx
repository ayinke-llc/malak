import { useState } from "react";
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

interface AddDataPointDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  chart: MalakIntegrationChart;
  workspaceIntegration: MalakWorkspaceIntegration;
}

export function AddDataPointDialog({
  open,
  onOpenChange,
  chart,
  workspaceIntegration,
}: AddDataPointDialogProps) {
  const [value, setValue] = useState("");
  const queryClient = useQueryClient();

  const { mutate: addDataPoint, isPending } = useMutation({
    mutationFn: async () => {
      const numericValue = parseInt(value, 10);
      if (isNaN(numericValue) || numericValue < 0) {
        throw new Error("Please enter a valid non-negative number");
      }
      return client.workspaces.integrationsChartsPointsCreate(
        workspaceIntegration.reference as string,
        chart.reference as string,
        { value: numericValue }
      );
    },
    onSuccess: () => {
      toast.success("Data point added successfully");
      onOpenChange(false);
      setValue("");
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
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="value" className="text-right">
              Value
            </Label>
            <Input
              id="value"
              type="number"
              min="0"
              value={value}
              onChange={(e) => setValue(e.target.value)}
              className="col-span-3"
            />
          </div>
        </div>
        <DialogFooter>
          <Button
            type="submit"
            onClick={() => addDataPoint()}
            disabled={isPending || !value}
          >
            {isPending ? "Adding..." : "Add Data Point"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
} 