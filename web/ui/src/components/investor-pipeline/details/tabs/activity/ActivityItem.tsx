import { RiMailLine, RiCalendarLine, RiFileTextLine } from "@remixicon/react";
import { Activity } from "../../../types";

interface ActivityItemProps {
  activity: Activity;
}

export function getActivityIcon(type: Activity['type']) {
  switch (type) {
    case 'email':
      return <RiMailLine className="w-4 h-4 text-primary" />;
    case 'meeting':
      return <RiCalendarLine className="w-4 h-4 text-primary" />;
    case 'stage_change':
      return <RiFileTextLine className="w-4 h-4 text-primary" />;
    default:
      return <RiMailLine className="w-4 h-4 text-primary" />;
  }
}

export function ActivityItem({ activity }: ActivityItemProps) {
  return (
    <div className="relative">
      <div className="absolute -left-[27px] bg-background p-1 border rounded-full">
        {getActivityIcon(activity.type)}
      </div>
      <div className="bg-card rounded-lg p-4 border">
        <div className="flex items-center justify-between mb-2">
          <h4 className="font-medium">{activity.title}</h4>
          <span className="text-xs text-muted-foreground">
            {new Date(activity.timestamp).toLocaleString()}
          </span>
        </div>
        <p className="text-sm text-muted-foreground mb-3">
          {activity.description}
        </p>
        {activity.content && (
          <div className="bg-muted/50 rounded-md p-3 text-sm">
            <p className="whitespace-pre-wrap">{activity.content}</p>
          </div>
        )}
      </div>
    </div>
  );
} 
