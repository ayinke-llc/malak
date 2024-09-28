"use client";

import { Button } from "@/components/Button";
import ListUpdatesTable from "@/components/ui/updates/list/list";
import { RiAddLine } from "@remixicon/react";
import Link from "next/link";

export default function Page() {
  return (
    <>
      <div className="pt-6">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3 className="scroll-mt-10 font-semibold text-gray-900 dark:text-gray-50">
                Investor updates
              </h3>
              <p className="text-sm leading-6 text-gray-500">
                View and manage your previously sent updates. Or send a new one
              </p>
            </div>

            <div className="w-full text-right">
              <Button
                type="button"
                asChild
                variant="primary"
                className="whitespace-nowrap"
              >
                <Link href={"/updates/new"}>
                  <RiAddLine />
                  New update
                </Link>
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
