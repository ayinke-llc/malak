"use client";

import BlockNoteJSEditor from "@/components/ui/updates/editor/blocknote";
import SendUpdateButton from "@/components/ui/updates/button/send";
import SendTestButton from "@/components/ui/updates/button/send-test";
import client from "@/lib/client";
import { FETCH_SINGLE_UPDATE } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import Skeleton from "@/components/ui/custom/loader/skeleton";
import Analytics from "@/components/ui/updates/analytics/analytics";
import { RiErrorWarningLine } from "@remixicon/react";

export default function UpdateDetailsPage({ reference }: { reference: string }) {
  const router = useRouter();

  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_SINGLE_UPDATE, reference],
    queryFn: () => client.workspaces.fetchUpdate(reference),
    retry: false,
  });

  if (isLoading) {
    return (
      <div className="flex flex-col space-y-4 p-4">
        <Skeleton count={1} />
        <div className="space-y-2">
          <Skeleton count={15} />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] p-4">
        <div className="text-center space-y-4">
          <RiErrorWarningLine className="h-12 w-12 text-red-500 mx-auto" />
          <h3 className="text-lg font-semibold text-gray-900">Failed to load update</h3>
          <p className="text-gray-500 max-w-md">
            We couldn't load this update. This might be because it was deleted or you don't have permission to view it.
          </p>
          <button
            onClick={() => router.push("/updates")}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Return to Updates
          </button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <section>
        <div className="mb-6">
          <button
            onClick={() => router.push("/updates")}
            className="inline-flex items-center text-sm text-gray-500 hover:text-gray-700 gap-1"
          >
            <svg className="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z" clipRule="evenodd" />
            </svg>
            Back to Updates
          </button>
        </div>
        <div className="mt-2">
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
        </div>
      </section>
    </div>
  );
} 
