"use client";

import BlockNoteJSEditor from "@/components/ui/updates/editor/blocknote";
import SendUpdateButton from "@/components/ui/updates/button/send";
import SendTestButton from "@/components/ui/updates/button/send-test";
import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import Analytics from "@/components/ui/updates/analytics/analytics";

export default function Page() {
  const params = useParams();

  const router = useRouter();

  const reference = params.slug as string;

  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_SINGLE_UPDATE],
    queryFn: () => client.workspaces.fetchUpdate(reference),
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  if (error) {
    toast.error("an error occurred while fetching this update");
    router.push("/updates");
    return;
  }

  return (
    <div>
      <section>

        <div className="mt-2">
          {isLoading ? (
            <Skeleton count={20} />
          ) : (
            <>
              <Analytics reference={reference} />
              <div className="flex flex-col sm:flex-row justify-end gap-2">
                <SendTestButton reference={reference} />
                <SendUpdateButton reference={reference} />
              </div>
              <BlockNoteJSEditor
                reference={reference}
                loading={isLoading}
                update={data?.data.update}
              />
            </>
          )}
        </div>
      </section>
    </div>
  );
}
