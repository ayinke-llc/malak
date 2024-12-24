"use client";

import ListDecks from "@/components/ui/decks/list/list";
import UploadDeckModal from "@/components/ui/decks/modal";

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
              <UploadDeckModal />
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
