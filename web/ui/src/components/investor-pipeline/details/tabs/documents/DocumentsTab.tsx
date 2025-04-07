import { RiFileTextLine } from "@remixicon/react";

interface DocumentsTabProps {
  isReadOnly: boolean;
}

export function DocumentsTab({ isReadOnly }: DocumentsTabProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center space-y-4">
      <RiFileTextLine className="w-12 h-12 text-muted-foreground" />
      <div>
        <h3 className="text-lg font-semibold">Documents Coming Soon</h3>
        <p className="text-sm text-muted-foreground">
          The documents feature is currently under development and will be available soon.
        </p>
      </div>
    </div>
  );
} 
