import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE_ANALYTICS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import Skeleton from "../../custom/loader/skeleton";
import View from "./view";
import { MalakUpdateRecipientStat, MalakUpdateStat } from "@/client/Api";

type Props = {
  reference: string
}

const Analytics = (props: Props) => {
  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_SINGLE_UPDATE_ANALYTICS],
    queryFn: () => client.workspaces.fetchUpdateAnalytics(props.reference),
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  if (error) {
    toast.error("an error occurred while fetching analytics for this update");
    return null;
  }

  if (!data?.data?.update) {
    return null;
  }

  // Don't render if all stats are 0 or undefined
  const hasAnalytics = data.data.update.unique_opens || 
                      data.data.update.total_opens || 
                      data.data.update.total_reactions || 
                      data.data.update.total_clicks || 
                      (data.data.recipients && data.data.recipients.length > 0);

  if (!hasAnalytics) {
    return null;
  }

  return (
    <div>
      <section>
        <div className="mt-2">
          {isLoading ? (
            <Skeleton count={20} />
          ) : (
            <View
              update={data.data.update as MalakUpdateStat}
              recipientStats={data.data.recipients ?? []} />
          )}
        </div>
      </section>
    </div>
  );
}

export default Analytics;
