"use client"

import BlockNoteJSEditor from "@/components/ui/updates/editor/blocknote";
import SendUpdateButton from "@/components/ui/updates/button/send";
import SendTestButton from "@/components/ui/updates/button/send-test";
import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";
import Skeleton from "@/components/ui/custom/loader/skeleton";

export default function Page() {

  const params = useParams()

  const router = useRouter()

  const reference = params.slug as string

  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_SINGLE_UPDATE],
    queryFn: () => client.workspaces.fetchUpdate(reference),
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  })

  if (error) {
    toast.error("an error occurred while fetching this update")
    router.push("/updates")
    return
  }

  return (
    <div className="pt-6">
      <section>
        <div className="sm:flex sm:items-center sm:justify-between">
          <div>
            <h3 id="existing-contacts" className="scroll-mt-10 font-semibold text-gray-900 dark:text-gray-50">
              Create a new update
            </h3>
            <p className="text-sm leading-6 text-gray-500">
              Sending a new update to your investors
            </p>
          </div>
          <div className="flex flex-wrap justify-center gap-1">
            <SendTestButton />
            <SendUpdateButton />
          </div>
        </div>

        <div className="mt-5">
          {isLoading ? (
            <Skeleton count={10} />
          ) : (
            <BlockNoteJSEditor reference={reference}
              loading={isLoading}
              update={data?.data.update}
            />
          )}
        </div>
      </section>
    </div>
  )
}
