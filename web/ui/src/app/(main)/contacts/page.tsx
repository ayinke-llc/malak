"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import CreateContactModal from "@/components/ui/contacts/modal";
import ManageListModal from "@/components/ui/contacts/new-list-modal";

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
            View and manage your investors' contacts
            </p>
          </div>

          <div>
            
            <ManageListModal />
            <CreateContactModal />
          </div>
        </div>
      </section>

      <section className="mt-10">
       
        <ListContacts />
      </section>
    </div>
  </>

  );
}
