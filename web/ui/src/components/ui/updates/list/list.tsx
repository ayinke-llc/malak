import { MalakUpdate } from "@/client/Api"
import client from "@/lib/client"
import { LIST_UPDATES } from "@/lib/query-constants"
import { useInfiniteQuery } from "@tanstack/react-query"
import SingleUpdate from "./single"
import { Button } from "@/components/Button"

export type UpdateListTableProps = {
  data: MalakUpdate[]
}

const UpdatesListTable = () => {

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

  return (
    <div>
      {data?.pages?.map((value, idx) => {
        return value?.data?.updates?.map((update, idx) => {
          return <SingleUpdate {...update} key={idx} />
        })
      })}

      {hasNextPage && <Button
        variant="primary"
        onClick={() => fetchNextPage()}>
        Load more
      </Button>}
    </div>
  )
}

export default UpdatesListTable;
