import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE_ANALYTICS } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import Skeleton from "../../custom/loader/skeleton";
import View from "./view";
import { MalakUpdateStat } from "@/client/Api";

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
    console.log(error)
    toast.error("an error occurred while fetching analytics for this update");
  }


  return (
    <div>
      <section>

        <div className="mt-2">
          {isLoading ? (
            <Skeleton count={20} />
          ) : (
            <>
              <View update={data?.data.update as MalakUpdateStat} />
            </>
          )}
        </div>
      </section>
    </div>
  );
}

export default Analytics;
