"use client";

import ListDecks from "@/components/ui/decks/list/list";
import UploadDeckModal from "@/components/ui/decks/modal";

export default function Decks() {
  return (
    <>
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3
                id="company-decks"
                className="text-lg font-medium"
              >
                Company Decks
              </h3>
              <p className="text-sm text-muted-foreground">
                View and manage your investors&apos; contacts
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