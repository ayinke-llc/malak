import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Activity } from "../../../types";

interface AddActivityDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (activity: Partial<Activity>) => void;
}

export function AddActivityDialog({
  open,
  onOpenChange,
  onSubmit
}: AddActivityDialogProps) {
  const [activity, setActivity] = useState<Partial<Activity>>({
    type: 'email',
    title: '',
    description: '',
    content: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(activity);
    onOpenChange(false);
    setActivity({
      type: 'email',
      title: '',
      description: '',
      content: ''
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Add New Activity</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Type</label>
            <Select
              value={activity.type}
              onValueChange={(value) => setActivity({ ...activity, type: value as Activity['type'] })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select activity type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="email">Email</SelectItem>
                <SelectItem value="meeting">Meeting</SelectItem>
                <SelectItem value="note">Note</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Title</label>
            <Input
              value={activity.title}
              onChange={(e) => setActivity({ ...activity, title: e.target.value })}
              placeholder={activity.type === 'email' ? 'Note title' : 'Activity title'}
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">{activity.type === 'email' ? 'Note Content' : 'Description'}</label>
            {activity.type === 'email' ? (
              <Textarea
                value={activity.content}
                onChange={(e) => setActivity({ ...activity, content: e.target.value, description: e.target.value.slice(0, 100) + '...' })}
                placeholder="Write your note here"
                className="min-h-[200px]"
                required
              />
            ) : (
              <Input
                value={activity.description}
                onChange={(e) => setActivity({ ...activity, description: e.target.value })}
                placeholder="Brief description"
                required
              />
            )}
          </div>

          {activity.type !== 'email' && (
            <div className="space-y-2">
              <label className="text-sm font-medium">Content (optional)</label>
              <Textarea
                value={activity.content}
                onChange={(e) => setActivity({ ...activity, content: e.target.value })}
                placeholder="Additional details or content"
              />
            </div>
          )}

          <div className="flex justify-end gap-2">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Add {activity.type === 'email' ? 'Note' : 'Activity'}</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 
