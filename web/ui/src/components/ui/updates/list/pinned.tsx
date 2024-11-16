import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RiPushpinLine, RiUnpinLine } from "@remixicon/react"
import { MalakUpdate } from "@/client/Api"
import Link from "next/link"
import UpdateBadge from "../../custom/update/badge"
import Skeleton from "../../custom/loader/skeleton"
import { useMutation, useQuery } from "@tanstack/react-query"
import { LIST_PINNED_UPDATES } from "@/lib/query-constants"
import client from "@/lib/client"
import { toast } from "sonner";
import { format } from "date-fns"

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

    <Card key={update.id} className="bg-primary/5">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">
          <Link href={`/updates/${update.reference}`}>
            {update?.title as string}
          </Link>
        </CardTitle>
        <UpdateBadge status={update?.status as string} />
      </CardHeader>
      <CardContent>
        <p className="text-xs text-muted-foreground">
          {format(update?.created_at as string, "EEEE, MMMM do, yyyy")}
        </p>
        <div className="flex justify-end mt-2">
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={() => togglePinnedStatus(update?.reference as string)}
          >
            <RiUnpinLine className="h-4 w-4" color="red" />
            <span className="sr-only">Unpin update</span>
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

export default PinnedList;
