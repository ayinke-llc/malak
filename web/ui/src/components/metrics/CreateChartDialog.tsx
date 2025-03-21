import { useState } from "react";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { RiAddLine } from "@remixicon/react";
import { toast } from "sonner";
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
import type { MalakWorkspaceIntegration } from "@/client/Api";

interface CreateChartFormData {
  title: string;
  type: "bar" | "pie";
}

const createChartSchema = yup.object({
  title: yup.string().required("Chart title is required"),
  type: yup.string().oneOf(["bar", "pie"], "Invalid chart type").required("Chart type is required"),
});

export function CreateChartDialog({ integration }: { integration: MalakWorkspaceIntegration }) {
  const [open, setOpen] = useState(false);
  
  const form = useForm<CreateChartFormData>({
    resolver: yupResolver(createChartSchema),
    defaultValues: {
      title: "",
      type: "bar"
    }
  });

  const handleCreateChart = (data: CreateChartFormData) => {
    // This will be implemented later
    toast.success("Chart creation will be implemented soon");
    setOpen(false);
    form.reset();
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
              onValueChange={(value: "bar" | "pie") => form.setValue("type", value)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="bar">Bar Chart</SelectItem>
                <SelectItem value="pie">Pie Chart</SelectItem>
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
            <Button type="submit">
              Create Chart
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 