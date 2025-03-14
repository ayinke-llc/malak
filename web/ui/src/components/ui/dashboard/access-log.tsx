import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { format } from "date-fns";
import { 
  RiDownloadLine, 
  RiLoader4Line, 
  RiMailLine, 
  RiGlobalLine,
  RiArrowLeftLine,
  RiArrowRightLine,
} from "@remixicon/react";

interface AccessLogEntry {
  id: string;
  timestamp: Date;
  userEmail?: string;
  accessMethod: 'link' | 'email';
  ipAddress: string;
  userAgent: string;
}

interface AccessLogProps {
  reference: string;
}

const ITEMS_PER_PAGE = 10;

export function AccessLog({ reference }: AccessLogProps) {
  const [filter, setFilter] = useState("");
  const [accessType, setAccessType] = useState<"all" | "link" | "email">("all");
  const [currentPage, setCurrentPage] = useState(1);

  // Simulate fetching access logs - in real app, fetch based on page
  const logs: AccessLogEntry[] = Array.from({ length: 50 }, (_, i) => ({
    id: `${i + 1}`,
    timestamp: new Date(Date.now() - i * 3600000),
    userEmail: i % 2 === 0 ? `user${i}@example.com` : undefined,
    accessMethod: i % 2 === 0 ? "email" : "link",
    ipAddress: `192.168.1.${i + 1}`,
    userAgent: i % 2 === 0 ? "Chrome/Windows" : "Safari/Mac"
  }));

  const exportMutation = useMutation({
    mutationFn: async () => {
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { success: true };
    }
  });

  const filteredLogs = logs.filter(log => {
    if (accessType !== "all" && log.accessMethod !== accessType) return false;
    if (!filter) return true;
    
    return (
      log.userEmail?.toLowerCase().includes(filter.toLowerCase()) ||
      log.ipAddress.includes(filter)
    );
  });

  const totalPages = Math.ceil(filteredLogs.length / ITEMS_PER_PAGE);
  const paginatedLogs = filteredLogs.slice(
    (currentPage - 1) * ITEMS_PER_PAGE,
    currentPage * ITEMS_PER_PAGE
  );

  const handleExport = () => {
    exportMutation.mutate();
  };

  return (
    <div className="space-y-6">
      <div className="sticky top-0 z-10 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex items-center gap-4 bg-muted/50 p-4 rounded-lg border shadow-sm">
          <div className="flex-1">
            <Input
              placeholder="Filter by email or IP..."
              value={filter}
              onChange={(e) => {
                setFilter(e.target.value);
                setCurrentPage(1);
              }}
              className="w-full bg-background"
            />
          </div>
          <Select
            value={accessType}
            onValueChange={(value: "all" | "link" | "email") => {
              setAccessType(value);
              setCurrentPage(1);
            }}
          >
            <SelectTrigger className="w-[180px] bg-background">
              <SelectValue placeholder="Access Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Access</SelectItem>
              <SelectItem value="link">Link Access</SelectItem>
              <SelectItem value="email">Email Access</SelectItem>
            </SelectContent>
          </Select>
          <Button
            variant="outline"
            onClick={handleExport}
            disabled={exportMutation.isPending}
            className="bg-background"
          >
            {exportMutation.isPending ? (
              <RiLoader4Line className="h-4 w-4 animate-spin" />
            ) : (
              <RiDownloadLine className="h-4 w-4" />
            )}
            <span className="ml-2">Export</span>
          </Button>
        </div>

        <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground px-1.5 pb-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div>
            Showing {paginatedLogs.length} of {filteredLogs.length} entries
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="h-8 w-8 p-0"
              >
                <RiArrowLeftLine className="h-4 w-4" />
              </Button>
              <span>
                Page {currentPage} of {totalPages}
              </span>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                disabled={currentPage === totalPages}
                className="h-8 w-8 p-0"
              >
                <RiArrowRightLine className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="rounded-lg border overflow-hidden mt-0">
        <div className="overflow-auto">
          <Table>
            <TableHeader>
              <TableRow className="bg-muted/50 hover:bg-muted/50">
                <TableHead className="w-[180px] font-semibold">Time</TableHead>
                <TableHead className="font-semibold min-w-[200px]">User</TableHead>
                <TableHead className="w-[120px] font-semibold">Access Type</TableHead>
                <TableHead className="font-semibold min-w-[120px]">IP Address</TableHead>
                <TableHead className="font-semibold min-w-[200px]">Device</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {paginatedLogs.map((log) => (
                <TableRow key={log.id} className="hover:bg-muted/50">
                  <TableCell className="align-middle py-4 whitespace-nowrap">
                    {format(log.timestamp, "MMM d, yyyy HH:mm")}
                  </TableCell>
                  <TableCell className="align-middle">
                    <div className="flex items-center gap-2">
                      {log.userEmail || "Anonymous"}
                    </div>
                  </TableCell>
                  <TableCell className="align-middle whitespace-nowrap">
                    <div className="flex items-center gap-2">
                      {log.accessMethod === "email" ? (
                        <RiMailLine className="h-4 w-4 text-primary" />
                      ) : (
                        <RiGlobalLine className="h-4 w-4 text-muted-foreground" />
                      )}
                      <span className={log.accessMethod === "email" ? "text-primary" : ""}>
                        {log.accessMethod}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="align-middle font-mono text-sm">
                    {log.ipAddress}
                  </TableCell>
                  <TableCell className="align-middle">
                    {log.userAgent}
                  </TableCell>
                </TableRow>
              ))}
              {paginatedLogs.length === 0 && (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center">
                    <div className="flex flex-col items-center justify-center text-muted-foreground">
                      <p className="text-sm">No access logs found</p>
                      <p className="text-xs mt-1">Try adjusting your filters</p>
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
} 