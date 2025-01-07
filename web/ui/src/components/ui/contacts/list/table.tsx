"use client";

import * as React from "react";
import {
  CaretSortIcon,
  ChevronDownIcon,
} from "@radix-ui/react-icons";
import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TableHeader,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { RiMoreLine } from "@remixicon/react";
import Link from "next/link";

export type Contact = {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  company?: string;
  city?: string;
  phone?: string;
  notes?: string;
  reference: string;
  lists: ContactList[];
  created_at: string;
  updated_at: string;
};

interface ContactList {
  id: string;
  title: string;
  reference: string;
}

const data: Contact[] = [
  {
    id: "1",
    email: "john@example.com",
    first_name: "John",
    last_name: "Doe",
    company: "Example Corp",
    city: "San Francisco",
    phone: "+1234567890",
    reference: "contact_123",
    lists: [
      { id: "1", title: "Investors", reference: "list_1" },
      { id: "2", title: "VCs", reference: "list_2" },
      { id: "3", title: "Angels", reference: "list_3" },
    ],
    created_at: "2024-01-15T10:00:00Z",
    updated_at: "2024-01-15T10:00:00Z",
  },
  {
    id: "2",
    email: "jane@example.com",
    first_name: "Jane",
    last_name: "Smith",
    company: "Tech Corp",
    city: "New York",
    phone: "+0987654321",
    reference: "contact_456",
    lists: [
      { id: "1", title: "Investors", reference: "list_1" },
    ],
    created_at: "2024-01-14T10:00:00Z",
    updated_at: "2024-01-14T10:00:00Z",
  },
];

export const columns: ColumnDef<Contact>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
        className="border-gray-700/50"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
        className="border-gray-700/50"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "email",
    header: ({ column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            className="p-0 hover:bg-transparent text-gray-400 hover:text-gray-300 font-normal -ml-4"
          >
            Contact
            <CaretSortIcon className="ml-2 h-4 w-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => {
      const contact = row.original;
      const fullName = `${contact.first_name} ${contact.last_name}`;

      return (
        <div className="min-w-[320px]">
          <Link
            href={`contacts/${contact.reference}`}
            className="group block"
          >
            <div className="flex items-start gap-3">
              <div className="mt-0.5 flex-shrink-0 w-8 h-8 rounded-full bg-gray-800/50 flex items-center justify-center text-gray-400 text-sm font-medium uppercase">
                {contact.first_name[0]}{contact.last_name[0]}
              </div>
              <div>
                <h3 className="font-medium text-gray-200 group-hover:text-white transition-colors">
                  {fullName}
                </h3>
                <div className="flex flex-col gap-0.5">
                  <p className="text-sm text-gray-400">
                    {contact.email}
                  </p>
                  {contact.company && (
                    <p className="text-sm text-gray-500">
                      {contact.company}
                    </p>
                  )}
                </div>
              </div>
            </div>
          </Link>
        </div>
      );
    },
  },
  {
    accessorKey: "lists",
    header: () => (
      <div className="text-left text-gray-400 font-normal">Lists</div>
    ),
    cell: ({ row }) => {
      const lists = row.getValue("lists") as ContactList[];

      if (lists && lists.length > 0) {
        return (
          <div className="flex gap-1.5 flex-wrap min-w-[200px]">
            {lists.slice(0, 3).map((list) => (
              <Badge
                key={list.id}
                variant="secondary"
                className="bg-gray-800/30 text-gray-300 hover:bg-gray-800/50 transition-colors border border-gray-800/50 px-2 py-0.5 text-xs font-normal"
              >
                {list.title}
              </Badge>
            ))}
            {lists.length > 3 && (
              <Badge
                variant="secondary"
                className="bg-transparent border border-gray-800/50 text-gray-500 px-2 py-0.5 text-xs font-normal"
              >
                +{lists.length - 3} more
              </Badge>
            )}
          </div>
        );
      }
      return null;
    },
  },
  {
    accessorKey: "created_at",
    header: () => (
      <div className="text-left text-gray-400 font-normal">Created</div>
    ),
    cell: ({ row }) => {
      const date = new Date(row.getValue("created_at"));
      const formatted = new Intl.DateTimeFormat("en-US", {
        dateStyle: "medium",
        timeStyle: "short",
      }).format(date);

      return (
        <div className="text-sm text-gray-500 whitespace-nowrap">
          {formatted}
        </div>
      );
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const contact = row.original;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              className="h-8 w-8 p-0 hover:bg-gray-800/50 text-gray-500 hover:text-gray-400"
            >
              <span className="sr-only">Open menu</span>
              <RiMoreLine className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="w-48 bg-gray-900 border-gray-800/50"
          >
            <DropdownMenuLabel className="text-gray-400 text-xs">Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(contact.email)}
              className="text-gray-300 focus:bg-gray-800 focus:text-gray-200 text-sm"
            >
              Copy email
            </DropdownMenuItem>
            <DropdownMenuSeparator className="bg-gray-800/50" />
            <DropdownMenuItem className="focus:bg-gray-800 focus:text-gray-200 text-sm">
              <Link
                href={`/contacts/${contact.reference}`}
                className="text-gray-300 w-full"
              >
                View details
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem className="text-gray-300 focus:bg-gray-800 focus:text-gray-200 text-sm">
              Add to list
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];

export default function ContactsTable() {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState({});

  const table = useReactTable({
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
  });

  return (
    <div className="w-full space-y-4">
      <div className="flex items-center justify-between">
        <Input
          placeholder="Filter contacts..."
          value={(table.getColumn("email")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("email")?.setFilterValue(event.target.value)
          }
          className="max-w-xs bg-transparent border-gray-800/50 text-sm placeholder:text-gray-600 focus-visible:ring-0 focus-visible:border-gray-700"
        />
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="outline"
                size="sm"
                className="h-8 border-gray-800/50 bg-transparent text-gray-400 hover:bg-gray-800/30 hover:text-gray-300"
              >
                View <ChevronDownIcon className="ml-2 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-40 bg-gray-900 border-gray-800/50">
              {table
                .getAllColumns()
                .filter((column) => column.getCanHide())
                .map((column) => {
                  return (
                    <DropdownMenuCheckboxItem
                      key={column.id}
                      className="capitalize text-gray-300 text-sm"
                      checked={column.getIsVisible()}
                      onCheckedChange={(value) =>
                        column.toggleVisibility(!!value)
                      }
                    >
                      {column.id}
                    </DropdownMenuCheckboxItem>
                  );
                })}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <div className="rounded-lg border border-gray-800/50 bg-gray-900/30">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow
                key={headerGroup.id}
                className="border-b border-gray-800/50 hover:bg-transparent"
              >
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead
                      key={header.id}
                      className="text-gray-400 h-10"
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  className="border-b border-gray-800/50 hover:bg-gray-800/20"
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className="py-3"
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center text-gray-500"
                >
                  No contacts found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between py-2">
        <div className="text-sm text-gray-500">
          {table.getFilteredSelectedRowModel().rows.length} of{" "}
          {table.getFilteredRowModel().rows.length} selected
        </div>
        <div className="flex items-center space-x-2 text-sm">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
            className="h-8 border-gray-800/50 bg-transparent text-gray-400 hover:bg-gray-800/30 hover:text-gray-300 disabled:opacity-50"
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
            className="h-8 border-gray-800/50 bg-transparent text-gray-400 hover:bg-gray-800/30 hover:text-gray-300 disabled:opacity-50"
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
