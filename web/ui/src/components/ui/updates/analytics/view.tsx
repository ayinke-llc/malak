'use client'

import { useRef, useEffect, useState } from 'react'
import { format } from "date-fns"
import {
  MalakContact,
  MalakUpdateRecipient,
  MalakUpdateStat
} from "@/client/Api"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import { fullName } from "@/lib/custom"
import {
  RiBarChart2Line,
  RiEye2Line,
  RiMouseLine,
  RiThumbUpLine
} from "@remixicon/react"
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"
import { useVirtualizer } from "@tanstack/react-virtual"

const columns: ColumnDef<MalakUpdateRecipient>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => (
      <div className='mt-3 mb-3'>
        <div className="font-medium">
          {fullName(row.original.contact as MalakContact)}
        </div>
        <div className="text-sm text-muted-foreground">{row.original.contact?.email}</div>
      </div>
    ),
  },
  {
    accessorKey: "delivered",
    header: "Delivered",
    cell: ({ row }) => (
      <div className="flex justify-right">
        {row.original.update_recipient_stat?.is_delivered ? (
          "✅"
        ) : (
          "❌"
        )}
      </div>
    ),
  },
  {
    accessorKey: "lastOpened",
    header: "Last Opened",
    cell: ({ row }) => (
      <div>
        {row.original.update_recipient_stat?.last_opened_at
          ? format(row.original.update_recipient_stat?.last_opened_at, "EEEE, MMMM do, yyyy")
          : "Not opened"}
      </div>
    ),
  },
  {
    accessorKey: "reacted",
    header: "Reacted",
    cell: ({ row }) => (
      <div className="text-left text-xl">{row.original.update_recipient_stat?.has_reaction ? "👍" : "❌"}</div>
    ),
  },
]

export interface Props {
  update: MalakUpdateStat
  recipientStats: MalakUpdateRecipient[]
}

export default function View(
  { update, recipientStats }: Props
) {

  const [showAll, setShowAll] = useState(false)

  const tableContainerRef = useRef<HTMLDivElement>(null)
  const [tableHeight, setTableHeight] = useState(400) // Default height

  const progressPercentage = (update?.unique_opens as number / (update?.total_opens as number)) * 100

  const table = useReactTable({
    columns,
    data: recipientStats,
    getCoreRowModel: getCoreRowModel(),
  })

  const { rows } = table.getRowModel()

  const rowVirtualizer = useVirtualizer({
    count: rows?.length,
    getScrollElement: () => tableContainerRef.current,
    estimateSize: () => 50, // Adjust this value based on your row height
    overscan: 5,
  })

  useEffect(() => {
    if (tableContainerRef.current) {
      setTableHeight(tableContainerRef.current.offsetHeight)
    }
  }, [showAll])

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
          <Progress value={progressPercentage} className="h-2 [&>div]:bg-green-500" />
        </div>
        <div
          className="rounded-md border"
          ref={tableContainerRef}
          style={{ height: showAll ? 'auto' : '700px', overflowY: 'auto' }}
        >
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
              {recipientStats.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={columns.length} className="h-24 text-center">
                    No results.
                  </TableCell>
                </TableRow>
              ) : (
                rowVirtualizer.getVirtualItems().map((virtualRow) => {
                  const row = rows?.[virtualRow.index]
                  if (!row) return null // Skip rendering if row is undefined
                  return (
                    <TableRow
                      key={row.id}
                      data-index={virtualRow.index}
                      ref={rowVirtualizer.measureElement}
                    >
                      {row.getVisibleCells().map((cell) => (
                        <TableCell key={cell.id}>
                          {flexRender(
                            cell.column.columnDef.cell,
                            cell.getContext()
                          )}
                        </TableCell>
                      ))}
                    </TableRow>
                  )
                })
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}
