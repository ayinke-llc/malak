'use client'

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"
import { RiBarChart2Line, RiEye2Line, RiMouseLine, RiThumbUpLine } from "@remixicon/react"
import { MalakUpdateStat } from "@/client/Api"

type Recipient = {
  name: string
  email: string
  delivered: boolean
  lastOpened: Date | null
  reacted: string | null
  clicks: number
}

const columns: ColumnDef<Recipient>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => (
      <div>
        <div className="font-medium">{row.original.name}</div>
        <div className="text-sm text-muted-foreground">{row.original.email}</div>
      </div>
    ),
  },
  {
    accessorKey: "delivered",
    header: "Delivered",
    cell: ({ row }) => (
      <div className="flex justify-center">
        {row.original.delivered ? (
          <span className="text-green-500 text-xl">✓</span>
        ) : (
          <span className="text-red-500 text-xl">✗</span>
        )}
      </div>
    ),
  },
  {
    accessorKey: "lastOpened",
    header: "Last Opened",
    cell: ({ row }) => (
      <div>
        {row.original.lastOpened
          ? row.original.lastOpened.toLocaleString()
          : "Not opened"}
      </div>
    ),
  },
  {
    accessorKey: "reacted",
    header: "Reacted",
    cell: ({ row }) => (
      <div className="text-center text-xl">{row.original.reacted || "-"}</div>
    ),
  },
  {
    accessorKey: "clicks",
    header: "Clicks",
    cell: ({ row }) => <div className="text-center">{row.original.clicks}</div>,
  },
]

const data: Recipient[] = [
  {
    name: "Lanre Adelowo",
    email: "yo@lanre.wtf",
    delivered: true,
    lastOpened: new Date("2023-06-15T10:30:00"),
    reacted: "👍",
    clicks: 3,
  },
  {
    name: "Bobby bro",
    email: "bob@example.com",
    delivered: true,
    lastOpened: null,
    reacted: null,
    clicks: 0,
  },
  {
    name: "Charlie XCX",
    email: "charlie@example.com",
    delivered: false,
    lastOpened: null,
    reacted: null,
    clicks: 0,
  },
  {
    name: "Elon Musk",
    email: "elon@example.com",
    delivered: true,
    lastOpened: new Date("2023-06-14T15:45:00"),
    reacted: "❤️",
    clicks: 5,
  },
  {
    name: "PG - Paul george or Graham",
    email: "pg@example.com",
    delivered: true,
    lastOpened: new Date("2023-06-15T09:00:00"),
    reacted: "😊",
    clicks: 2,
  },
]

export interface Props {
  update: MalakUpdateStat
}

export default function View(
  { update }: Props
) {
  const progressPercentage = (update?.unique_opens as number / update?.total_opens as number) * 100

  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="text-2xl font-bold">Update Analytics</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
          <div className="flex items-center space-x-4">
            <div className="p-2 bg-blue-100 rounded-full">
              <RiEye2Line className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm font-medium text-muted-foreground">Unique Opens</p>
              <h3 className="text-2xl font-bold">{update?.unique_opens?.toLocaleString()}</h3>
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <div className="p-2 bg-green-100 rounded-full">
              <RiBarChart2Line className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm font-medium text-muted-foreground">Total Opens</p>
              <h3 className="text-2xl font-bold">{update?.total_opens?.toLocaleString()}</h3>
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <div className="p-2 bg-yellow-100 rounded-full">
              <RiThumbUpLine className="h-6 w-6 text-yellow-600" />
            </div>
            <div>
              <p className="text-sm font-medium text-muted-foreground">Total Reactions</p>
              <h3 className="text-2xl font-bold">{update?.total_reactions?.toLocaleString() ?? 0}</h3>
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <div className="p-2 bg-purple-100 rounded-full">
              <RiMouseLine className="h-6 w-6 text-purple-600" />
            </div>
            <div>
              <p className="text-sm font-medium text-muted-foreground">Total Clicks</p>
              <h3 className="text-2xl font-bold">{update?.total_clicks?.toLocaleString() ?? 0}</h3>
            </div>
          </div>
        </div>
        <div className="space-y-2 mb-6">
          <div className="flex justify-between text-sm text-muted-foreground">
            <span>Total Opens Progress</span>
            <span>{progressPercentage.toFixed(1)}%</span>
          </div>
          <Progress value={progressPercentage} className="h-2" />
        </div>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <TableHead key={header.id}>
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
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={columns.length} className="h-24 text-center">
                    No results.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}
