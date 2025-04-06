import { Button } from "@/components/ui/button";
import { RiCalendarLine, RiFileTextLine } from "@remixicon/react";

interface SuggestedActionsProps {
  onAddActivity: () => void;
  onUploadDocument: () => void;
}

export function SuggestedActions({ onAddActivity, onUploadDocument }: SuggestedActionsProps) {
  return (
    <div className="bg-card rounded-lg p-4 border">
      <h3 className="font-medium mb-3">Suggested Actions</h3>
      <div className="space-y-2">
        <Button
          variant="outline"
          className="w-full justify-start"
          onClick={onAddActivity}
        >
          <RiCalendarLine className="w-4 h-4 mr-2" />
          Add Activity or Note
        </Button>
        <Button
          variant="outline"
          className="w-full justify-start"
          onClick={onUploadDocument}
        >
          <RiFileTextLine className="w-4 h-4 mr-2" />
          Upload Documents
        </Button>
      </div>
    </div>
  );
} 
