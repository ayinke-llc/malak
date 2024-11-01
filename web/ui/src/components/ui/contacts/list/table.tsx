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
import { Button } from "@/components/Button";
import { Checkbox } from "@/components/Checkbox";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/Dropdown";
import { Input } from "@/components/Input";
import {
  Table,
  TableRoot,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TableHeaderCell,
} from "@/components/Table";
import { Badge } from "@/components/Badge";
import { RiMoreLine, RiSettingsLine } from "@remixicon/react";
import Link from "next/link";

const data: Payment[] = [
  {
    reference: "m5gr84i9",
    amount: 316,
    status: "success",
    email: "oops@yahoo.com",
    lists: ["preseed investors", "series a", "war time investors", "good guys", "war room"],
  },
  {
    reference: "3u1reuv4",
    amount: 242,
    status: "success",
    email: "oops@gmail.com",
  },
  {
    reference: "derv1ws0",
    amount: 837,
    status: "processing",
    email: "Modfuhf@gmail.com",
  },
  {
    reference: "5kma53ae",
    amount: 874,
    status: "success",
    email: "fjkfgjf@gmail.com",
  },
  {
    reference: "bhqecj4p",
    amount: 721,
    status: "failed",
    email: "jfhkff@hotmail.com",
  },
  {
    reference: "m5gr84i9",
    amount: 316,
    status: "success",
    email: "oops@yahoo.com",
    lists: ["preseed investors", "series a", "war time investors", "good guys", "war room"],
  },
  {
    reference: "3u1reuv4",
    amount: 242,
    status: "success",
    email: "oops@gmail.com",
  },
  {
    reference: "derv1ws0",
    amount: 837,
    status: "processing",
    email: "Modfuhf@gmail.com",
  },
  {
    reference: "5kma53ae",
    amount: 874,
    status: "success",
    email: "fjkfgjf@gmail.com",
  },
  {
    reference: "bhqecj4p",
    amount: 721,
    status: "failed",
    email: "jfhkff@hotmail.com",
  },
  {
    reference: "m5gr84i9",
    amount: 316,
    status: "success",
    email: "oops@yahoo.com",
    lists: ["preseed investors", "series a", "war time investors", "good guys", "war room"],
  },
  {
    reference: "3u1reuv4",
    amount: 242,
    status: "success",
    email: "oops@gmail.com",
  },
  {
    reference: "derv1ws0",
    amount: 837,
    status: "processing",
    email: "Modfuhf@gmail.com",
  },
  {
    reference: "5kma53ae",
    amount: 874,
    status: "success",
    email: "fjkfgjf@gmail.com",
  },
  {
    reference: "bhqecj4p",
    amount: 721,
    status: "failed",
    email: "jfhkff@hotmail.com",
  },
  {
    reference: "m5gr84i9",
    amount: 316,
    status: "success",
    email: "oops@yahoo.com",
    lists: ["preseed investors", "series a", "war time investors", "good guys", "war room"],
  },
  {
    reference: "3u1reuv4",
    amount: 242,
    status: "success",
    email: "oops@gmail.com",
  },
  {
    reference: "derv1ws0",
    amount: 837,
    status: "processing",
    email: "Modfuhf@gmail.com",
  },
  {
    reference: "5kma53ae",
    amount: 874,
    status: "success",
    email: "fjkfgjf@gmail.com",
  },
  {
    reference: "bhqecj4p",
    amount: 1000,
    status: "failed",
    email: "omomijoor@hotmail.com",
  },
];

export type Payment = {
  reference: string;
  amount: number;
  status: "pending" | "processing" | "success" | "failed";
  email: string;
  lists?: string[];
};

export const columns: ColumnDef<Payment>[] = [
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
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "email",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Email
          <CaretSortIcon className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const id = row.getValue("reference")

      return (
        <div>
          <Link href={`contacts/${id}`}>
            <h3 className="font-semibold capitalize">Lanre Adelowo</h3>
            <p className="text-sm text-muted-foreground lowercase">
              {row.getValue("email")}
            </p>
          </Link>
        </div>
      );
    },
  },
  {
    accessorKey: "lists",
    header: () => <div className="text-left">Lists</div>,
    cell: ({ row }) => {
      const values = row.getValue("lists") as string[];

      if (values !== undefined) {
        return (
          <div className="text-left font-medium gap-2 flex">
            {values.slice(0, 3).map((value) => {
              return (
                <Badge variant={"default"}>{value}</Badge>
              )
            })}
            {values.length > 3 && `+ ${values.length - 3}`}
          </div>
        );
      }
    },
  },
  {
    accessorKey: "amount",
    header: () => <div className="text-left">Created at</div>,
    cell: ({ row }) => {
      return <div className="text-left font-medium">{"2024/01/12 13:09"}</div>;
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const payment = row.original;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <RiMoreLine className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem>View</DropdownMenuItem>
            <DropdownMenuItem>Edit</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(payment.reference)}
            >
              Copy reference
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];

export default function ContactsTable() {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    [],
  );
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
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
      rowSelection,
    },
  });

  return (
    <div className="w-full">
      <div className="flex items-center py-4">
        <Input
          placeholder="Filter emails..."
          value={(table.getColumn("email")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("email")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="ml-auto">
              Columns <ChevronDownIcon className="ml-2 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {table
              .getAllColumns()
              .filter((column) => column.getCanHide())
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className="capitalize"
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
      <div className="rounded-md border">
        <TableRoot>
          <Table>
            <TableHead>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <TableHeaderCell key={header.id}>
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                      </TableHeaderCell>
                    );
                  })}
                </TableRow>
              ))}
            </TableHead>
            <TableBody>
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                    className="cursor-pointer hover:bg-gray-100"
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
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
        </TableRoot>
      </div>
      <div className="flex items-center justify-end space-x-2 py-4">
        <div className="flex-1 text-sm text-muted-foreground">
          {table.getFilteredSelectedRowModel().rows.length} of{" "}
          {table.getFilteredRowModel().rows.length} row(s) selected.
        </div>
        <div className="space-x-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Previous
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
