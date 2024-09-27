"use client"

import { ServerAPIStatus } from "@/client/Api";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import SendUpdateButton from "@/components/ui/updates/button/send";
import SendTestButton from "@/components/ui/updates/button/send-test";
import BlockNoteJSEditor from "@/components/ui/updates/editor/blocknote";
import NovelEditor from "@/components/ui/updates/editor/editor";
import client from "@/lib/client";
import { CREATE_UPDATE } from "@/lib/query-constants";
import { useMutation } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function Page() {

  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [reference, setReference] = useState<string | undefined>(undefined)

  const router = useRouter()

  const mutation = useMutation({
    mutationKey: [CREATE_UPDATE],
    mutationFn: () => client.workspaces.updatesCreate(),
    onSuccess: ({ data }) => {
      setReference(data.update.reference)
      toast.success("Your update have been created now. As you type, we will sync and save your changes")
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response.data.message
      }
      toast.error(msg)
      router.push("/updates")
    },
    retry: false,
    gcTime: Infinity,
    onSettled: () => setIsLoading(false),
  })

  useEffect(() => {
    mutation.mutate()
  }, [])

  return (
    <div className="pt-6">
      {isLoading ? (<div className="mt-10">
        <Skeleton count={40} />
      </div>) : (
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
            <BlockNoteJSEditor reference={reference} />
          </div>
        </section>
      )}
    </div>
  )
}
