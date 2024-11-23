"use client";

import { ServerAPIStatus, ServerFetchUpdateReponse } from "@/client/Api";
import { Button } from "@/components/ui/button";
import ListUpdatesTable from "@/components/ui/updates/list/list";
import PinnedList from "@/components/ui/updates/list/pinned";
import client from "@/lib/client";
import { RiAddLine } from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import { AxiosError, AxiosResponse } from "axios";
import { format } from "date-fns";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

export default function Page() {
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const router = useRouter();

  const mutation = useMutation({
    mutationFn: () => {
      return client.workspaces.updatesCreate({
        title: `${format(new Date(), "EEEE, MMMM do, yyyy")} Update`,
      });
    },
    gcTime: 0,
    onError: (err: AxiosError<ServerAPIStatus>): void => {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }

      toast.error(msg);
    },
    retry: false,
    onSuccess: (resp: AxiosResponse<ServerFetchUpdateReponse>) => {
      router.push(`/updates/${resp.data.update.reference}`);
    },
    onMutate: () => setIsLoading(true),
    onSettled: () => setIsLoading(false)
  });

  return (
    <>
      <div className="pt-6">

        <section>
          <PinnedList />
        </section>

        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div className="w-full text-right">
              <Button
                type="button"
                variant="default"
                className="whitespace-nowrap"
                loading={isLoading}
                onClick={() => mutation.mutate()}
              >
                <RiAddLine />
                New update
              </Button>
            </div>
          </div>
        </section>
        <div className="mt-10 sm:mt-4">
          <ListUpdatesTable />
        </div>
      </div>
    </>
  );
}
