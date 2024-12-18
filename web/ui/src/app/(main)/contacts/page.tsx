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
                id="existing-contacts"
                className="scroll-mt-10 font-semibold text-gray-900 dark:text-gray-50"
              >
                Stored Contacts
              </h3>
              <p className="text-sm leading-6 text-gray-500">
                View and manage all your existing contacts
              </p>
            </div>

            <div className="flex justify-center gap-2">
              <ManageListModal />
              <CreateContactModal />
            </div>
          </div>
        </section>

        <section>
          <ListContacts />
        </section>
      </div>
    </>
  );
}
