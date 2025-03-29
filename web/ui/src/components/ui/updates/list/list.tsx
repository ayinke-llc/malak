import type { MalakUpdate } from "@/client/Api";
import { Button } from "@/components/ui/button";
import client from "@/lib/client";
import { LIST_UPDATES } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import Skeleton from "../../custom/loader/skeleton";
import SingleUpdate from "./single";
import { useEffect, useState } from "react";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { RiArrowDownSLine } from "@remixicon/react";
import { AlertCircle } from "lucide-react";

export type UpdateListTableProps = {
  data: MalakUpdate[];
};

type FilterStatus = "all" | "draft" | "sent";

const ListUpdatesTable = () => {
  const [pageIndex, setPageIndex] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [filterStatus, setFilterStatus] = useState<FilterStatus>("all");

  const {
    data,
    error,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: [LIST_UPDATES, pageIndex, pageSize, filterStatus],
    queryFn: () => {
      return client.workspaces.updatesList({
        page: pageIndex,
        per_page: pageSize,
        status: filterStatus === "all" ? undefined : filterStatus,
      });
    },
    retry: false,
  });

  const updates = data?.data?.updates || [];
  const totalPages = Math.ceil((data?.data?.meta?.paging?.total || 0) / pageSize) || 1;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 bg-muted p-1 rounded-lg">
          <Button
            variant={filterStatus === "all" ? "default" : "ghost"}
            size="sm"
            className="h-8"
            onClick={() => setFilterStatus("all")}
          >
            All Updates
          </Button>
          <Button
            variant={filterStatus === "draft" ? "default" : "ghost"}
            size="sm"
            className="h-8"
            onClick={() => setFilterStatus("draft")}
          >
            Drafts
          </Button>
          <Button
            variant={filterStatus === "sent" ? "default" : "ghost"}
            size="sm"
            className="h-8"
            onClick={() => setFilterStatus("sent")}
          >
            Sent
          </Button>
        </div>
      </div>

      <div className="rounded-lg border bg-background">
        {isLoading ? (
          <div className="p-6 space-y-6">
            <Skeleton count={5} />
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center p-6 text-center">
            <AlertCircle className="h-12 w-12 text-destructive mb-4" />
            <h3 className="text-lg font-semibold mb-2">Failed to load updates</h3>
            <p className="text-muted-foreground mb-4">{error.message}</p>
            <Button
              variant="outline"
              onClick={() => refetch()}
            >
              Try again
            </Button>
          </div>
        ) : updates.length > 0 ? (
          <div className="divide-y divide-border/50">
            {updates.map((update) => (
              <div
                key={update.reference}
                className="p-4 hover:bg-muted transition-colors"
              >
                <SingleUpdate {...update} />
              </div>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center p-6 text-center">
            <p className="text-muted-foreground">No updates found.</p>
          </div>
        )}
      </div>

      {!isLoading && (
        <div className="flex items-center justify-between py-2">
          <div className="flex items-center gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="outline"
                  size="sm"
                  className="h-8"
                >
                  {pageSize} per page <RiArrowDownSLine className="ml-2 h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {[10, 20, 30, 50].map((size) => (
                  <DropdownMenuItem
                    key={size}
                    onClick={() => setPageSize(size)}
                  >
                    Show {size} rows
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div className="flex items-center gap-1.5">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPageIndex(1)}
              disabled={pageIndex === 1}
              className="h-8 w-8 p-0"
            >
              <span className="sr-only">Go to first page</span>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
                <path fillRule="evenodd" d="M15.79 14.77a.75.75 0 01-1.06.02L10.25 10.5l4.48-4.29a.75.75 0 111.04 1.08L11.78 10l4.47 4.23a.75.75 0 01-.02 1.06zm-6 0a.75.75 0 01-1.06.02L4.25 10.5l4.48-4.29a.75.75 0 111.04 1.08L5.78 10l4.47 4.23a.75.75 0 01-.02 1.06z" clipRule="evenodd" />
              </svg>
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPageIndex(pageIndex - 1)}
              disabled={pageIndex === 1}
              className="h-8 w-8 p-0"
            >
              <span className="sr-only">Go to previous page</span>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
                <path fillRule="evenodd" d="M12.79 14.77a.75.75 0 01-1.06.02L7.25 10.5l4.48-4.29a.75.75 0 111.04 1.08L8.78 10l4.47 4.23a.75.75 0 01-.02 1.06z" clipRule="evenodd" />
              </svg>
            </Button>

            <div className="flex items-center gap-1">
              <span className="text-sm text-muted-foreground">Page</span>
              <Input
                type="number"
                min={1}
                max={totalPages}
                value={pageIndex}
                onChange={e => {
                  const page = e.target.value ? Number(e.target.value) : 1;
                  if (page >= 1 && page <= totalPages) {
                    setPageIndex(page);
                  }
                }}
                className="w-12 h-8 text-center [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
              />
              <span className="text-sm text-muted-foreground">of {totalPages}</span>
            </div>

            <Button
              variant="outline"
              size="sm"
              onClick={() => setPageIndex(pageIndex + 1)}
              disabled={pageIndex === totalPages}
              className="h-8 w-8 p-0"
            >
              <span className="sr-only">Go to next page</span>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
                <path fillRule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clipRule="evenodd" />
              </svg>
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPageIndex(totalPages)}
              disabled={pageIndex === totalPages}
              className="h-8 w-8 p-0"
            >
              <span className="sr-only">Go to last page</span>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
                <path fillRule="evenodd" d="M4.21 14.77a.75.75 0 01.02-1.06L8.168 10 4.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l4.5-4.25a.75.75 0 01-1.06-.02z" clipRule="evenodd" />
              </svg>
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ListUpdatesTable;
