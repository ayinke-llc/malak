import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  RiRefreshLine,
  RiLoader4Line, RiLinkM,
  RiMailLine,
  RiArrowLeftLine,
  RiArrowRightLine,
  RiSearchLine,
  RiFileCopyLine,
  RiArrowUpLine,
  RiArrowDownLine,
} from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { format } from "date-fns";
import { Input } from "@/components/ui/input";
import CopyToClipboard from "react-copy-to-clipboard";
import { GENERATE_ACCESS_LINK, FETCH_ACCESS_LINKS } from "@/lib/query-constants";
import client from "@/lib/client";
import { ServerAPIStatus, MalakDashboardLink } from "@/client/Api";
import { AxiosError } from "axios";
import { MALAK_APP_URL } from "@/lib/config";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
  getSortedRowModel,
  SortingState,
  getFilteredRowModel,
  ColumnFiltersState,
  getPaginationRowModel,
} from "@tanstack/react-table";

interface AccessManagementProps {
  reference: string;
  shareLink: string;
  onLinkChange: (s: string) => void
}

const columns: ColumnDef<ReturnType<typeof transformLink>>[] = [
  {
    accessorKey: "email",
    header: "Email",
    cell: ({ row }) => (
      <div className="font-medium">
        {row.original.email}
      </div>
    ),
  },
  {
    accessorKey: "createdAt",
    header: ({ column }) => (
      <Button
        variant="ghost"
        className="p-0 hover:bg-transparent"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Added
        {column.getIsSorted() === "asc" ? (
          <RiArrowUpLine className="ml-2 h-4 w-4" />
        ) : column.getIsSorted() === "desc" ? (
          <RiArrowDownLine className="ml-2 h-4 w-4" />
        ) : null}
      </Button>
    ),
    cell: ({ row }) => format(row.original.createdAt, "MMM d, yyyy"),
  },
  {
    id: "shareLink",
    header: "Share Link",
    cell: ({ row }) => {
      const shareUrl = `${MALAK_APP_URL}/shared/dashboards/${row.original.token}`;
      
      return (
        <div className="flex items-center gap-4">
          <span className="text-sm text-muted-foreground w-[300px] truncate">
            {truncateUrl(shareUrl)}
          </span>
          <CopyToClipboard
            text={shareUrl}
            onCopy={() => toast.success("Share link copied to clipboard")}
          >
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
            >
              <RiFileCopyLine className="h-4 w-4" />
            </Button>
          </CopyToClipboard>
        </div>
      );
    },
  },
  {
    id: "actions",
    header: "",
    cell: ({ row }) => {
      const email = row.original.email;
      return (
        email && row.original.status === "active" && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => row.original.onRevoke?.(email)}
            disabled={row.original.isRevoking}
            className="text-red-600 hover:text-red-700 hover:bg-red-50"
          >
            {row.original.isRevoking ? (
              <RiLoader4Line className="h-4 w-4 animate-spin" />
            ) : (
              "Revoke"
            )}
          </Button>
        )
      );
    },
  },
];

const transformLink = (link: MalakDashboardLink) => ({
  type: "email" as const,
  email: link.contact?.email ?? "",
  token: link.token ?? "",
  createdAt: new Date(link.created_at ?? ""),
  lastAccess: link.updated_at ? new Date(link.updated_at) : undefined,
  accessCount: 0,
  status: "active" as const,
  onRevoke: undefined as ((email: string) => void) | undefined,
  isRevoking: false,
});

const truncateUrl = (url: string) => {
  if (url.length <= 45) return url;
  return `${url.slice(0, 30)}...${url.slice(-12)}`;
};

export function AccessManagement({ onLinkChange, reference, shareLink }: AccessManagementProps) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [currentLink, setCurrentLink] = useState(shareLink);
  const [page, setPage] = useState(1);
  const [perPage] = useState(10);

  const { data: accessLinksData, isLoading } = useQuery({
    queryKey: [FETCH_ACCESS_LINKS, reference, page, perPage],
    queryFn: () => client.dashboards.accessControlDetail(reference, { page, per_page: perPage }),
  });

  const regenerateLinkMutation = useMutation({
    mutationKey: [GENERATE_ACCESS_LINK],
    mutationFn: () => client.dashboards.accessControlLinkCreate(reference, {}),
    onSuccess: ({ data }) => {
      toast.success("link generated")
      const fullShareLink = MALAK_APP_URL + "/shared/dashboards/" + data?.link?.token as string;
      setCurrentLink(fullShareLink)
      onLinkChange(fullShareLink)
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "Could not generate link");
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  const revokeAccessMutation = useMutation({
    mutationFn: async (email: string) => {
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { success: true };
    },
    onSuccess: () => {
      toast.success("Access revoked successfully");
    },
    onError: () => {
      toast.error("Failed to revoke access");
    }
  });

  const handleRegenerateLink = () => {
    regenerateLinkMutation.mutate();
  };

  const handleRevokeAccess = (email: string) => {
    revokeAccessMutation.mutate(email);
  };

  const data = (accessLinksData?.data?.links ?? []).map((link: MalakDashboardLink) => {
    const transformed = transformLink(link);
    transformed.onRevoke = handleRevokeAccess;
    transformed.isRevoking = revokeAccessMutation.isPending;
    return transformed;
  });

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    state: {
      sorting,
      columnFilters,
    },
    initialState: {
      pagination: {
        pageSize: perPage,
      },
    },
    pageCount: Math.ceil((accessLinksData?.data?.meta?.paging?.total ?? 0) / perPage),
    manualPagination: true,
  });

  return (
    <div className="space-y-8">
      <div className="sticky top-0 z-10 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 space-y-4 pb-4">
        <div className="rounded-lg border bg-muted/50 p-6 space-y-6 shadow-sm">
          <div className="flex items-center justify-between">
            <div className="space-y-1.5">
              <div className="flex items-center gap-2">
                <RiLinkM className="h-5 w-5 text-primary" />
                <h4 className="text-lg font-semibold">Link Sharing</h4>
              </div>
              <p className="text-sm text-muted-foreground">
                Anyone with the link can view this dashboard
              </p>
            </div>
          </div>

          <div className="flex items-center gap-4 pt-2">
            <div className="flex-1 flex items-center gap-2">
              <Input
                value={currentLink}
                readOnly
                className="bg-background flex-1"
              />
              <CopyToClipboard
                text={currentLink}
                onCopy={() => toast.success("Link copied to clipboard")}
              >
                <Button
                  variant="outline"
                  size="icon"
                  className="bg-background"
                >
                  <RiFileCopyLine className="h-4 w-4" />
                </Button>
              </CopyToClipboard>
            </div>
            <Button
              variant="outline"
              onClick={handleRegenerateLink}
              disabled={regenerateLinkMutation.isPending}
              className="bg-background shrink-0"
            >
              {regenerateLinkMutation.isPending ? (
                <RiLoader4Line className="h-4 w-4 animate-spin" />
              ) : (
                <RiRefreshLine className="h-4 w-4" />
              )}
              <span className="ml-2">Regenerate Link</span>
            </Button>
          </div>
        </div>

        {/* Access List Header */}
        <div className="flex items-center justify-between px-1.5 pt-2">
          <h4 className="text-lg font-semibold">Access List</h4>
          <div className="relative">
            <RiSearchLine className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Filter by email..."
              value={(table.getColumn("email")?.getFilterValue() as string) ?? ""}
              onChange={(e) => table.getColumn("email")?.setFilterValue(e.target.value)}
              className="pl-9 w-[200px]"
            />
          </div>
        </div>

        <div className="flex items-center justify-between text-sm text-muted-foreground px-1.5 border-b">
          <div>
            Showing {data.length} of {accessLinksData?.data?.meta?.paging?.total ?? 0} entries
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setPage(page => Math.max(1, page - 1))}
                disabled={page === 1}
                className="h-8 w-8 p-0"
              >
                <RiArrowLeftLine className="h-4 w-4" />
              </Button>
              <span>
                Page {page} of{" "}
                {Math.ceil((accessLinksData?.data?.meta?.paging?.total ?? 0) / perPage)}
              </span>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setPage(page => page + 1)}
                disabled={page >= Math.ceil((accessLinksData?.data?.meta?.paging?.total ?? 0) / perPage)}
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
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id} className="bg-muted/50 hover:bg-muted/50">
                  {headerGroup.headers.map((header) => (
                    <TableHead key={header.id} className="font-semibold">
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                    </TableHead>
                  ))}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={columns.length} className="h-24">
                    <div className="flex items-center justify-center">
                      <RiLoader4Line className="h-6 w-6 animate-spin" />
                    </div>
                  </TableCell>
                </TableRow>
              ) : data.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow key={row.id}>
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    No results.
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
