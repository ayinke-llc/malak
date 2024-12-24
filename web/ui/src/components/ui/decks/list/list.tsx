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
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";

type Deck = {
  id: number;
  name: string;
  uploadedAt: Date;
  size: string;
  url: string;
  reference: string;
};

// This is a placeholder until we implement the actual data fetching
const mockDecks: Deck[] = [
  {
    id: 1,
    name: "Company Overview 2024 - Q1 Financial Results and Future Projections.pdf",
    uploadedAt: new Date(),
    size: "2.4 MB",
    url: "https://example.com/decks/very-long-url-that-needs-to-be-truncated-but-still-copyable-123456789",
    reference: "company-overview-2024",
  },
  {
    id: 2,
    name: "Ctester.pdf",
    uploadedAt: new Date(),
    size: "2.4 MB",
    url: "httpsdfffexample.com/decks/very-long-url-that-needs-to-be-truncated-but-still-copyable-123456789",
    reference: "ctester",
  },
];

export default function ListDecks() {
  const [copiedId, setCopiedId] = useState<number | null>(null);

  const copyToClipboard = async (text: string, id: number) => {
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

  const columnHelper = createColumnHelper<Deck>();

  const columns = useMemo(
    () => [
      columnHelper.accessor("name", {
        header: () => <span className="text-zinc-400">Name</span>,
        cell: (info) => (
          <div className="max-w-[300px]">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Link 
                    href={`/decks/${info.row.original.reference}`}
                    className="block hover:opacity-80"
                  >
                    <span className="font-medium text-zinc-100 truncate block">
                      {truncateText(info.getValue(), 40)}
                    </span>
                  </Link>
                </TooltipTrigger>
                <TooltipContent>
                  <p>{info.getValue()}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>
        ),
      }),
      columnHelper.accessor("size", {
        header: () => <span className="text-zinc-400">Size</span>,
        cell: (info) => <span className="text-zinc-100">{info.getValue()}</span>,
      }),
      columnHelper.accessor("url", {
        header: () => <span className="text-zinc-400">URL</span>,
        cell: (info) => (
          <div className="flex items-center gap-2 max-w-[300px]">
            <span className="text-zinc-100 truncate block">
              {truncateText(info.getValue(), 40)}
            </span>
            <Button
              variant="ghost"
              size="sm"
              className="shrink-0 h-8 w-8 p-0 text-zinc-400 hover:text-zinc-100"
              onClick={() => copyToClipboard(info.getValue(), info.row.original.id)}
            >
              {copiedId === info.row.original.id ? (
                <RiCheckLine className="h-4 w-4" />
              ) : (
                <RiFileCopyLine className="h-4 w-4" />
              )}
            </Button>
          </div>
        ),
      }),
      columnHelper.accessor("uploadedAt", {
        header: () => <span className="text-zinc-400">Uploaded</span>,
        cell: (info) => (
          <span className="text-zinc-100 whitespace-nowrap">
            {format(info.getValue(), "MMM d, yyyy")}
          </span>
        ),
      }),
    ],
    [copiedId]
  );

  const table = useReactTable({
    data: mockDecks,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  if (mockDecks.length === 0) {
    return (
      <Card className="flex flex-col items-center justify-center py-16 px-4 bg-zinc-900/50">
        <div className="flex flex-col items-center justify-center text-center max-w-sm">
          <div className="rounded-full bg-zinc-800 p-4">
            <RiUploadCloud2Line className="h-8 w-8 text-zinc-400" />
          </div>
          <h3 className="mt-6 text-lg font-medium text-zinc-100">
            No decks uploaded yet
          </h3>
          <p className="mt-2 text-sm text-zinc-400/80">
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
    <Card className="bg-zinc-900/50">
      <div className="relative overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
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
                className="border-t border-zinc-800 hover:bg-zinc-800/50 transition-colors"
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