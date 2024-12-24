"use client";

import { Button } from "@/components/ui/button";
import { RiAddLine } from "@remixicon/react";
import ListDecks from "@/components/ui/decks/list/list";

export default function Page() {
  return (
    <>
      <div className="pt-6">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3
                id="company-decks"
                className="text-lg font-medium text-zinc-100"
              >
                Company Decks
              </h3>
              <p className="text-sm text-zinc-400/80">
                View and manage your uploaded company decks and PDFs
              </p>
            </div>

            <div>
              <Button
                type="button"
                variant="default"
                className="whitespace-nowrap bg-zinc-100 text-zinc-900 hover:bg-zinc-200"
              >
                <RiAddLine className="mr-2 h-4 w-4" />
                Upload Deck
              </Button>
            </div>
          </div>
        </section>

        <section className="mt-10">
          <ListDecks />
        </section>
      </div>
    </>
  );
}
