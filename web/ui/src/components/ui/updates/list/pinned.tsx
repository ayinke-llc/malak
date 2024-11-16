import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RiCalendarLine, RiFlipHorizontal2Line, RiFlipHorizontalLine, RiMoreLine, RiPushpinLine, RiUnpinLine } from "@remixicon/react"
import { MalakUpdate } from "@/client/Api"
import Link from "next/link"
import UpdateBadge from "../../custom/update/badge"
import Skeleton from "../../custom/loader/skeleton"
import { useMutation, useQuery } from "@tanstack/react-query"
import { LIST_PINNED_UPDATES } from "@/lib/query-constants"
import client from "@/lib/client"
import { toast } from "sonner";
import { format } from "date-fns"
import { Badge } from "@/components/ui/badge"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

const PinnedList = () => {

  const {
    data,
    error,
    isFetching,
  } = useQuery({
    queryKey: [LIST_PINNED_UPDATES],
    queryFn: () => {
      return client.workspaces.updatesPinsList();
    },
    retry: false,
  });

  if (error) {
    toast.error(error.message);
  }
  return (
    <>
      {
        isFetching ? (
          <div className="pb-10">
            <Skeleton count={10} />
          </div>
        )
          : (
            <div className="grid grid-cols-4 md:grid-cols-4 lg:grid-cols-4 gap-4 mb-6">
              {data?.data?.updates?.map((update) => {
                return <Item {...update} />
              })}
            </div>
          )
      }

    </>
  )
}

const Item = (update: MalakUpdate) => {

  const togglePinnedStatus = (reference: string) => {
  }

  return (

    <Card key={update?.id as string} className="bg-white border-none shadow-md hover:shadow-lg transition-shadow duration-300">
      <CardContent className="p-4">
        <div className="flex items-start justify-between mb-3">
          <Badge variant={update?.status as string === 'published' ? 'default' : 'secondary'} className="text-xs">
            {update?.status}
          </Badge>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <RiMoreLine className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => { }}>
                Unpin
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Link href={`/updates/${update.reference}`}>View Details</Link>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <Link href={`/updates/${update?.reference as string}`} className="block group">
          <h3 className="text-lg font-medium text-gray-800 group-hover:text-primary transition-colors duration-200 mb-2">
            {update?.title as string}
          </h3>
        </Link>
        <div className="flex items-center text-sm text-gray-500">
          <RiCalendarLine className="h-4 w-4 mr-2" />
          {format(update?.created_at as string, "EEEE, MMMM do, yyyy")}
        </div>
      </CardContent>
    </Card>
  )
}

export default PinnedList;
