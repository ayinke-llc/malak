"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { format } from "date-fns";
import { RiFileList3Line, RiUploadCloud2Line } from "@remixicon/react";

// This is a placeholder until we implement the actual data fetching
const mockDecks: any[] = [];

export default function ListDecks() {
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
          <Button className="mt-6 bg-zinc-100 text-zinc-900 hover:bg-zinc-200" size="lg">
            <RiUploadCloud2Line className="mr-2 h-4 w-4" />
            Upload your first deck
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <Card className="bg-zinc-900/50">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="text-zinc-400">Name</TableHead>
            <TableHead className="text-zinc-400">Company</TableHead>
            <TableHead className="text-zinc-400">Size</TableHead>
            <TableHead className="text-zinc-400">Uploaded</TableHead>
            <TableHead className="w-[100px] text-zinc-400">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {mockDecks.map((deck) => (
            <TableRow key={deck.id}>
              <TableCell className="font-medium text-zinc-100">{deck.name}</TableCell>
              <TableCell className="text-zinc-100">{deck.company}</TableCell>
              <TableCell className="text-zinc-100">{deck.size}</TableCell>
              <TableCell className="text-zinc-100">{format(deck.uploadedAt, "MMM d, yyyy")}</TableCell>
              <TableCell>
                <div className="flex gap-2">
                  <Button variant="ghost" size="sm" className="text-zinc-100 hover:text-zinc-200 hover:bg-zinc-800">
                    Download
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Card>
  );
} 