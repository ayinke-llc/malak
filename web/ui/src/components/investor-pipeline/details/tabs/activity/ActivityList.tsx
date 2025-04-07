import { useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { RiAddLine } from "@remixicon/react";
import { Activity } from "../../../types";
import { ActivityItem } from "./ActivityItem";
import { AddActivityDialog } from "./AddActivityDialog";

interface ActivityListProps {
  activities: Activity[];
  isLoading: boolean;
  hasMore: boolean;
  isArchived?: boolean;
  onAddActivity: (activity: Partial<Activity>) => void;
}

function ActivitySkeleton() {
  return (
    <div className="relative">
      <div className="absolute -left-[27px] bg-background p-1 border rounded-full">
        <div className="h-4 w-4 rounded-full bg-muted animate-pulse" />
      </div>
      <div className="bg-card rounded-lg p-4 border">
        <div className="flex items-center justify-between mb-2">
          <div className="h-4 w-32 bg-muted animate-pulse rounded" />
          <div className="h-3 w-24 bg-muted animate-pulse rounded" />
        </div>
        <div className="h-4 w-full bg-muted animate-pulse rounded mb-3" />
        <div className="h-16 w-full bg-muted animate-pulse rounded" />
      </div>
    </div>
  );
}

export function ActivityList({
  activities,
  isLoading,
  hasMore,
  isArchived = false,
  onAddActivity
}: ActivityListProps) {
  const [isAddingActivity, setIsAddingActivity] = useState(false);
  const observerTarget = useRef<HTMLDivElement>(null);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="text-sm text-muted-foreground">
          Showing {activities.length} of {activities.length >= 250 ? 'maximum ' : ''}{250} activities
        </div>
        {!isArchived && (
          <Button
            onClick={() => setIsAddingActivity(true)}
            size="sm"
            className="gap-2"
          >
            <RiAddLine className="w-4 h-4" />
            Add Activity or Note
          </Button>
        )}
      </div>

      <div className="relative space-y-6 pl-8 before:absolute before:left-3 before:top-2 before:bottom-0 before:w-[2px] before:bg-border">
        {activities.map((activity) => (
          <ActivityItem key={activity.id} activity={activity} />
        ))}

        {/* Loading state */}
        {isLoading && (
          <>
            <ActivitySkeleton />
            <ActivitySkeleton />
            <ActivitySkeleton />
          </>
        )}

        {/* Intersection observer target */}
        <div ref={observerTarget} className="h-4" />

        {/* End of list message */}
        {!hasMore && activities.length > 0 && (
          <div className="text-center text-sm text-muted-foreground py-4">
            {activities.length >= 250
              ? `Maximum limit of ${250} activities reached`
              : "No more activities to load"}
          </div>
        )}
      </div>

      {!isArchived && (
        <AddActivityDialog
          open={isAddingActivity}
          onOpenChange={setIsAddingActivity}
          onSubmit={onAddActivity}
        />
      )}
    </div>
  );
} 
