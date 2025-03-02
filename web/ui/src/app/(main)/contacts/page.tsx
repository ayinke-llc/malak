"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import CreateContactModal from "@/components/ui/contacts/modal";
import ManageListModal from "@/components/ui/contacts/new-list-modal";

export default function Page() {
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
                Company contacts
              </h3>
              <p className="text-sm text-muted-foreground">
                View and manage your investors&apos; contacts
              </p>
            </div>

            <div>

              <ManageListModal />
              <CreateContactModal />
            </div>
          </div>
        </section>

        <section className="mt-10 sm:mt-4">
          <ListContacts />
        </section>
      </div>
    </>

  );
}
