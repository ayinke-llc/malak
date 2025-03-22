import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
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
import type { MalakIntegrationChart } from "@/client/Api";

interface AddDataPointDialogProps {
  chart: MalakIntegrationChart;
}

export function AddDataPointDialog({ chart }: AddDataPointDialogProps) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState("");
  const [timestamp, setTimestamp] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Here you would typically make an API call to add the data point
    // For now, we're just closing the dialog as requested
    setOpen(false);
    setValue("");
    setTimestamp("");
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <RiAddLine className="h-4 w-4" />
          Add Data Point
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Data Point</DialogTitle>
          <DialogDescription>
            Add a new data point to {chart.name}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="value" className="text-right">
                Value
              </Label>
              <Input
                id="value"
                type="number"
                step="any"
                value={value}
                onChange={(e) => setValue(e.target.value)}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="timestamp" className="text-right">
                Timestamp
              </Label>
              <Input
                id="timestamp"
                type="datetime-local"
                value={timestamp}
                onChange={(e) => setTimestamp(e.target.value)}
                className="col-span-3"
                required
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit">Add Data Point</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
} 