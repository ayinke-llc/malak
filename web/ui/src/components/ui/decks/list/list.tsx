"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { format } from "date-fns";
import { RiUploadCloud2Line, RiFileCopyLine, RiCheckLine } from "@remixicon/react";
import UploadDeckModal from "../modal";
import { useState, useMemo } from "react";
import { toast } from "sonner";
import Link from "next/link";
import {
  Tooltip,
  TooltipContent, TooltipTrigger
} from "@/components/ui/tooltip";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { LIST_DECKS } from "@/lib/query-constants";
import client from "@/lib/client";
import { useQuery } from "@tanstack/react-query";
import type { ServerFetchDecksResponse, MalakDeck } from "@/client/Api";

export default function ListDecks() {
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const { data, isLoading } = useQuery<ServerFetchDecksResponse>({
    queryKey: [LIST_DECKS],
    queryFn: () => client.decks.decksList().then(res => res.data),
  });

  const decks = useMemo(() => data?.decks ?? [], [data]);

  const copyToClipboard = async (text: string, id: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopiedId(id);
      toast.success("Link copied to clipboard", {
        description: "The deck URL has been copied to your clipboard.",
      });
      setTimeout(() => setCopiedId(null), 2000);
    } catch (error) {
      toast.error("Failed to copy link", {
        description: "Please try copying the link again.",
      });
    }
  };

  const truncateText = (text: string, maxLength: number) => {
    if (text.length <= maxLength) return text;
    return text.slice(0, maxLength) + "...";
  };

  const columnHelper = createColumnHelper<MalakDeck>();

  const columns = useMemo(
    () => [
      columnHelper.accessor("title", {
        header: () => <span className="text-muted-foreground">Name</span>,
        cell: (info) => {
          const title = info.getValue() ?? "";
          return (
            <div className="max-w-[300px]">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Link
                    href={`/decks/${info.row.original.reference}`}
                    className="block text-foreground hover:opacity-80"
                  >
                    <span className="font-medium truncate block">
                      {truncateText(title, 40)}
                    </span>
                  </Link>
                </TooltipTrigger>
                <TooltipContent>
                  <p>{title}</p>
                </TooltipContent>
              </Tooltip>
            </div>
          );
        },
      }),
      columnHelper.accessor(() => "N/A" as const, {
        id: "size",
        header: () => <span className="text-muted-foreground">Size</span>,
        cell: (info) => <span className="text-foreground">{info.getValue()}</span>,
      }),
      columnHelper.accessor(
        (row) => {
          const shortLink = row.short_link ?? "";
          return `${window.location.origin}/d/${shortLink}`;
        },
        {
          id: "url",
          header: () => <span className="text-muted-foreground">URL</span>,
          cell: (info) => {
            const url = info.getValue();
            return (
              <div className="flex items-center gap-2 max-w-[300px]">
                <span className="truncate block text-foreground">
                  {truncateText(url, 40)}
                </span>
                <Button
                  variant="ghost"
                  size="sm"
                  className="shrink-0 h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
                  onClick={() => copyToClipboard(url, info.row.original.id ?? "")}
                >
                  {copiedId === info.row.original.id ? (
                    <RiCheckLine className="h-4 w-4" />
                  ) : (
                    <RiFileCopyLine className="h-4 w-4" />
                  )}
                </Button>
              </div>
            );
          },
        }
      ),
      columnHelper.accessor("created_at", {
        header: () => <span className="text-muted-foreground">Uploaded</span>,
        cell: (info) => {
          const date = info.getValue();
          return (
            <span className="whitespace-nowrap text-foreground">
              {date ? format(new Date(date), "MMM d, yyyy") : "N/A"}
            </span>
          );
        },
      }),
    ],
    [copiedId]
  );

  const table = useReactTable({
    data: decks,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  if (isLoading) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-zinc-50/50 dark:bg-zinc-900">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-zinc-500 dark:border-zinc-300 border-t-transparent" />
          <p className="mt-4 text-sm text-zinc-500 dark:text-zinc-300">Loading decks...</p>
        </div>
      </Card>
    );
  }

  if (decks.length === 0) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-zinc-50/50 dark:bg-zinc-900">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="rounded-full bg-zinc-100 dark:bg-zinc-800 p-4">
            <RiUploadCloud2Line className="h-8 w-8 text-zinc-500 dark:text-zinc-300" />
          </div>
          <h3 className="mt-6 text-lg font-medium text-zinc-900 dark:text-zinc-50">
            No decks uploaded yet
          </h3>
          <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
            Upload your company decks and PDFs to keep them organized and easily accessible in one place.
          </p>
          <div className="mt-6">
            <UploadDeckModal />
          </div>
        </div>
      </Card>
    );
  }

  return (
    <Card className="bg-background dark:bg-zinc-900 text-foreground">
      <div className="relative overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id} className="border-b border-zinc-200/10 dark:border-zinc-800">
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    className="px-6 py-4 font-medium"
                  >
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row) => (
              <tr
                key={row.id}
                className="border-t border-zinc-200 dark:border-zinc-800 cursor-pointer"
              >
                {row.getVisibleCells().map((cell) => (
                  <td
                    key={cell.id}
                    className="px-6 py-4"
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Card>
  );
} 
