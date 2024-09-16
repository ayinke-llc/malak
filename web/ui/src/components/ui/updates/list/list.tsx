import { MalakUpdate } from "@/client/Api"
import client from "@/lib/client"
import { LIST_UPDATES } from "@/lib/query-constants"
import { useInfiniteQuery } from "@tanstack/react-query"
import SingleUpdate from "./single"
import { Button } from "@/components/Button"
import { toast } from "sonner"
import Skeleton from "../../custom/loader/skeleton"

export type UpdateListTableProps = {
  data: MalakUpdate[]
}

const ListUpdatesTable = () => {

  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage
  } = useInfiniteQuery({
    queryKey: [LIST_UPDATES],
    queryFn: ({ pageParam }) => {
      console.log(pageParam)
      return client.workspaces.updatesList({
        page: pageParam,
      })
    },
    getNextPageParam: (lastPage) => {
      if (lastPage.data?.updates == undefined) {
        return undefined
      }

      return (lastPage.data.meta.paging?.page as number) + 1
    },
    retry: false,
    initialPageParam: 1,
  })

  if (error) {
    toast.error(error.message)
  }

  return (
    <div>
      {isFetching ?
        <Skeleton count={30} /> : (
          data?.pages?.map((value) => {
            return value?.data?.updates?.map((update, idx) => {
              return <SingleUpdate {...update} key={idx} />
            })
          }))}

      {isFetchingNextPage && <Skeleton count={30} />}

      {hasNextPage && (
        <div className="mt-5">
          <Button
            variant="secondary"
            onClick={() => fetchNextPage()}>
            Load more
          </Button>
        </div>
      )}
    </div>
  )
}

export default ListUpdatesTable;
