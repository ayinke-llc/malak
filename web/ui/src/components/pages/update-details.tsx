"use client";

import BlockNoteJSEditor from "@/components/ui/updates/editor/blocknote";
import SendUpdateButton from "@/components/ui/updates/button/send";
import SendTestButton from "@/components/ui/updates/button/send-test";
import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import Analytics from "@/components/ui/updates/analytics/analytics";

export default function UpdateDetailsPage({ reference }: { reference: string }) {
  const router = useRouter();

  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_SINGLE_UPDATE, reference],
    queryFn: () => client.workspaces.fetchUpdate(reference),
    retry: false,
  });

  if (error) {
    toast.error("an error occurred while fetching this update");
    router.push("/updates");
    return null;
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
              <div className="flex flex-col sm:flex-row justify-end gap-2 mt-8">
                {data?.data?.update?.status === "draft" && <SendTestButton reference={reference} />}
                <SendUpdateButton reference={reference} isSent={data?.data?.update?.status === "sent"} />
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
