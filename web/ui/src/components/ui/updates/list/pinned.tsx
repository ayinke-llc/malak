import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RiCalendarLine, RiMoreLine } from "@remixicon/react"
import { MalakUpdate, ServerAPIStatus, ServerCreatedUpdateResponse, } from "@/client/Api"
import Link from "next/link"
import type { AxiosError, AxiosResponse } from "axios";
import Skeleton from "../../custom/loader/skeleton"
import { useMutation, useQuery } from "@tanstack/react-query"
import { LIST_PINNED_UPDATES, TOGGLE_PINNED_STATE } from "@/lib/query-constants"
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
import { EVENT_TOGGLE_PINNED_STATE } from "@/lib/analytics-constansts"
import { usePostHog } from "posthog-js/react"
import { useState } from "react"
import UpdateBadge from "../../custom/update/badge"

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

  const posthog = usePostHog()

  const [loading, setLoading] = useState<boolean>(false);
  const [deleted, setDeleted] = useState<boolean>(false);

  const togglePinnedStatusMutation = useMutation({
    mutationKey: [TOGGLE_PINNED_STATE],
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onMutate: () => {
      posthog?.capture(EVENT_TOGGLE_PINNED_STATE);
      setLoading(true)
    },
    onSettled: () => setLoading(false),
    mutationFn: (reference: string) =>
      client.workspaces.toggleUpdatePin(reference),
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    onSuccess: (resp: AxiosResponse<ServerCreatedUpdateResponse>) => {
      toast.success(`Update has been Unpinned`);
      setDeleted(resp.data?.update?.is_pinned as boolean)
    },
  });

  if (deleted) {
    return null
  }

  if (loading) {
    return <Skeleton count={5} />
  }

  return (
    <Card key={update?.id as string} className="shadow-md hover:shadow-lg transition-shadow duration-300">
      <CardContent className="p-4">
        <div className="flex items-start justify-between mb-3">
          <UpdateBadge status={update?.status as string} />
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <RiMoreLine className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                onClick={() => togglePinnedStatusMutation.mutate(update?.reference as string)}>
                Unpin
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Link href={`/updates/${update.reference}`}>View Details</Link>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <Link href={`/updates/${update?.reference as string}`} className="block group">
          <h3 className="text-lg font-medium mb-2">
            {update?.title as string}
          </h3>
        </Link>
        <div className="flex items-center text-sm ">
          <RiCalendarLine className="h-4 w-4 mr-2" />
          {format(update?.created_at as string, "EEEE, MMMM do, yyyy")}
        </div>
      </CardContent>
    </Card>
  )
}

export default PinnedList;
