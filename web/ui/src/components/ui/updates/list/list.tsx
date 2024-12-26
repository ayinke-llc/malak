import type { MalakUpdate } from "@/client/Api";
import { Button } from "@/components/ui/button";
import client from "@/lib/client";
import { LIST_UPDATES } from "@/lib/query-constants";
import { useInfiniteQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import Skeleton from "../../custom/loader/skeleton";
import SingleUpdate from "./single";

export type UpdateListTableProps = {
  data: MalakUpdate[];
};

const ListUpdatesTable = () => {
  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
  } = useInfiniteQuery({
    queryKey: [LIST_UPDATES],
    queryFn: ({ pageParam }) => {
      return client.workspaces.updatesList({
        page: pageParam,
      });
    },
    getNextPageParam: (lastPage) => {
      if (lastPage.data?.updates === undefined) {
        return undefined;
      }

      return (lastPage.data.meta.paging?.page as number) + 1;
    },
    retry: false,
    initialPageParam: 1,
  });

  if (error) {
    toast.error(error.message);
  }

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        {isFetching ? (
          <div className="p-6 space-y-6">
            <Skeleton count={5} />
          </div>
        ) : (
          <div className="divide-y divide-border/50">
            {data?.pages?.map((value) => {
              return value?.data?.updates?.map((update) => (
                <div 
                  key={update.reference} 
                  className="p-4 hover:bg-muted/50 transition-colors"
                >
                  <SingleUpdate {...update} />
                </div>
              ));
            })}
          </div>
        )}
      </div>

      {isFetchingNextPage && (
        <div className="rounded-lg border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 p-6">
          <div className="space-y-6">
            <Skeleton count={3} />
          </div>
        </div>
      )}

      {hasNextPage && (
        <div className="flex justify-center py-6">
          <Button
            variant="outline"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
            className="min-w-[200px] bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
          >
            {isFetchingNextPage ? "Loading more updates..." : "Show more updates"}
          </Button>
        </div>
      )}
    </div>
  );
};

export default ListUpdatesTable;
