import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE_ANALYTICS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import Skeleton from "../../custom/loader/skeleton";
import View from "./view";
import { MalakUpdateRecipientStat, MalakUpdateStat } from "@/client/Api";
import { useState } from "react";

type Props = {
  reference: string
}

const Analytics = (props: Props) => {

  const [showAll, setShowAll] = useState(false)

  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_SINGLE_UPDATE_ANALYTICS],
    queryFn: () => client.workspaces.fetchUpdateAnalytics(props.reference),
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  if (error) {
    toast.error("an error occurred while fetching analytics for this update");
  }

  return (
    <div>
      <section>

        <div className="mt-2">
          {isLoading ? (
            <Skeleton count={20} />
          ) : (
            <View
              showAll={showAll}
              toggleShowAll={() => {
                setShowAll(!showAll)
              }}
              update={data?.data?.update as MalakUpdateStat}
              recipientStats={data?.data?.recipients as MalakUpdateRecipientStat[]} />
          )}
        </div>
      </section>
    </div>
  );
}

export default Analytics;
