'use client'

import {
  MalakContact,
  MalakUpdateRecipient,
  MalakUpdateStat
} from "@/client/Api"
import {
  Card,
  CardContent,
  CardHeader
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
  RiThumbUpLine,
  RiArrowDownSLine,
  RiArrowUpSLine
} from "@remixicon/react"
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"
import { useVirtualizer } from "@tanstack/react-virtual"
import { format } from "date-fns"
import { useRef, useState } from 'react'
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"

const columns: ColumnDef<MalakUpdateRecipient>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => (
      <div className='pl-3 mt-3 mb-3'>
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
  const [isExpanded, setIsExpanded] = useState(recipientStats.length <= 6)
  const tableContainerRef = useRef<HTMLDivElement>(null)

  const progressPercentage = !update?.total_opens ? 0 : (update?.unique_opens as number / update?.total_opens as number) * 100

  const table = useReactTable({
    columns,
    data: recipientStats,
    getCoreRowModel: getCoreRowModel(),
  })

  const { rows } = table.getRowModel()

  const rowVirtualizer = useVirtualizer({
    count: rows?.length,
    getScrollElement: () => tableContainerRef.current,
    estimateSize: () => 50,
    overscan: 5,
  })

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="flex items-center gap-4 w-full">
          <Button
            variant="ghost"
            size="sm"
            className="flex items-center gap-2"
            onClick={() => setIsExpanded(!isExpanded)}
          >
            {isExpanded ? (
              <RiArrowUpSLine className="h-4 w-4" />
            ) : (
              <RiArrowDownSLine className="h-4 w-4" />
            )}
            <span>Analytics {isExpanded ? 'Collapse' : 'Expand'}</span>
          </Button>

          {!isExpanded && (
            <div className="flex items-center gap-6 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <RiEye2Line className="h-4 w-4 text-blue-600" />
                <span>{(update?.unique_opens ?? 0).toLocaleString()} opens</span>
              </div>
              <Separator orientation="vertical" className="h-4" />
              <div className="flex items-center gap-2">
                <RiThumbUpLine className="h-4 w-4 text-yellow-600" />
                <span>{update?.total_reactions?.toLocaleString() ?? 0} reactions</span>
              </div>
              <Separator orientation="vertical" className="h-4" />
              <div className="flex items-center gap-2">
                <RiMouseLine className="h-4 w-4 text-purple-600" />
                <span>{update?.total_clicks?.toLocaleString() ?? 0} clicks</span>
              </div>
            </div>
          )}
        </div>
      </CardHeader>
      {isExpanded && (
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
            <div className="flex items-center space-x-4">
              <div className="p-2 bg-blue-100 rounded-full">
                <RiEye2Line className="h-6 w-6 text-blue-600" />
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Unique Opens</p>
                <h3 className="text-2xl font-bold">{(update?.unique_opens ?? 0).toLocaleString()}</h3>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="p-2 bg-green-100 rounded-full">
                <RiBarChart2Line className="h-6 w-6 text-green-600" />
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Opens</p>
                <h3 className="text-2xl font-bold">{(update?.total_opens ?? 0).toLocaleString()}</h3>
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
            style={{ height: 'auto', overflowY: 'auto' }}
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
      )}
    </Card>
  )
}
