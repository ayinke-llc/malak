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
  RiFileCopyLine
} from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { format } from "date-fns";
import { Input } from "@/components/ui/input";
import CopyToClipboard from "react-copy-to-clipboard";
import { GENERATE_ACCESS_LINK } from "@/lib/query-constants";
import client from "@/lib/client";
import { ServerAPIStatus } from "@/client/Api";
import { AxiosError } from "axios";
import { MALAK_APP_URL } from "@/lib/config";

interface ShareAccess {
  type: 'link' | 'email';
  email?: string;
  createdAt: Date;
  lastAccess?: Date;
  accessCount: number;
  status: 'active' | 'revoked';
}

interface AccessManagementProps {
  reference: string;
  shareLink: string;
  onLinkChange: (s: string) => void
}

const ITEMS_PER_PAGE = 10;

export function AccessManagement({ onLinkChange, reference, shareLink }: AccessManagementProps) {

  const [currentPage, setCurrentPage] = useState(1);
  const [filter, setFilter] = useState("");
  const [typeFilter, setTypeFilter] = useState<"all" | "link" | "email">("all");
  const [currentLink, setCurrentLink] = useState(shareLink);

  const accessList: ShareAccess[] = Array.from({ length: 50 }, (_, i) => ({
    type: i % 3 === 0 ? "link" : "email",
    email: i % 3 === 0 ? undefined : `user${i}@example.com`,
    createdAt: new Date(Date.now() - i * 86400000),
    lastAccess: i % 4 === 0 ? undefined : new Date(Date.now() - i * 3600000),
    accessCount: Math.floor(Math.random() * 20),
    status: i % 5 === 0 ? "revoked" : "active"
  }));

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

  const filteredList = accessList.filter(access => {
    if (typeFilter !== "all" && access.type !== typeFilter) return false;
    if (!filter) return true;

    return access.email?.toLowerCase().includes(filter.toLowerCase());
  });

  const totalPages = Math.ceil(filteredList.length / ITEMS_PER_PAGE);
  const paginatedList = filteredList.slice(
    (currentPage - 1) * ITEMS_PER_PAGE,
    currentPage * ITEMS_PER_PAGE
  );

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
          <div className="flex items-center gap-4">
            <div className="relative">
              <RiSearchLine className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Filter by email..."
                value={filter}
                onChange={(e) => {
                  setFilter(e.target.value);
                  setCurrentPage(1);
                }}
                className="pl-9 w-[200px]"
              />
            </div>
            <Select
              value={typeFilter}
              onValueChange={(value: "all" | "link" | "email") => {
                setTypeFilter(value);
                setCurrentPage(1);
              }}
            >
              <SelectTrigger className="w-[140px]">
                <SelectValue placeholder="Type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Types</SelectItem>
                <SelectItem value="link">Link</SelectItem>
                <SelectItem value="email">Email</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <div className="flex items-center justify-between text-sm text-muted-foreground px-1.5 border-b">
          <div>
            Showing {paginatedList.length} of {filteredList.length} entries
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
                <TableHead className="w-[120px] font-semibold">Type</TableHead>
                <TableHead className="font-semibold min-w-[200px]">User</TableHead>
                <TableHead className="w-[120px] font-semibold">Added</TableHead>
                <TableHead className="w-[140px] font-semibold">Last Access</TableHead>
                <TableHead className="w-[100px] font-semibold text-center">Access Count</TableHead>
                <TableHead className="w-[100px] font-semibold">Status</TableHead>
                <TableHead className="w-[100px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {paginatedList.map((access, index) => (
                <TableRow key={index} className="hover:bg-muted/50">
                  <TableCell className="align-middle py-4 whitespace-nowrap">
                    <div className="flex items-center gap-2">
                      {access.type === "email" ? (
                        <RiMailLine className="h-4 w-4 text-primary" />
                      ) : (
                        <RiLinkM className="h-4 w-4 text-muted-foreground" />
                      )}
                      <span className={access.type === "email" ? "text-primary" : ""}>
                        {access.type}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="align-middle font-medium">
                    {access.email || "Anyone with link"}
                  </TableCell>
                  <TableCell className="align-middle whitespace-nowrap">
                    {format(access.createdAt, "MMM d, yyyy")}
                  </TableCell>
                  <TableCell className="align-middle whitespace-nowrap">
                    {access.lastAccess
                      ? format(access.lastAccess, "MMM d, yyyy HH:mm")
                      : "Never"}
                  </TableCell>
                  <TableCell className="align-middle text-center">
                    {access.accessCount}
                  </TableCell>
                  <TableCell className="align-middle">
                    <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${access.status === "active"
                      ? "bg-green-50 text-green-700"
                      : "bg-red-50 text-red-700"
                      }`}>
                      {access.status}
                    </span>
                  </TableCell>
                  <TableCell className="align-middle text-right">
                    {access.type === "email" && access.status === "active" && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRevokeAccess(access.email!)}
                        disabled={revokeAccessMutation.isPending}
                        className="text-red-600 hover:text-red-700 hover:bg-red-50"
                      >
                        {revokeAccessMutation.isPending ? (
                          <RiLoader4Line className="h-4 w-4 animate-spin" />
                        ) : (
                          "Revoke"
                        )}
                      </Button>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
} 
