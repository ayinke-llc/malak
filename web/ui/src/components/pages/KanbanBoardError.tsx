import { RiErrorWarningLine } from "@remixicon/react";
import { Button } from "@/components/ui/button";

interface KanbanBoardErrorProps {
  error: Error;
  resetErrorBoundary: () => void;
}

export function KanbanBoardError({ error, resetErrorBoundary }: KanbanBoardErrorProps) {
  return (
    <div className="flex flex-col items-center justify-center h-screen p-4 text-center">
      <RiErrorWarningLine className="w-12 h-12 text-destructive mb-4" />
      <h2 className="text-2xl font-semibold mb-2">Something went wrong</h2>
      <p className="text-muted-foreground mb-4 max-w-md">
        {error.message || "An error occurred while loading the fundraising pipeline."}
      </p>
      <Button onClick={resetErrorBoundary}>Try again</Button>
    </div>
  );
} 