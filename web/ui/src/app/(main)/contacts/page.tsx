"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import CreateContactModal from "@/components/ui/contacts/modal";
import ManageListModal from "@/components/ui/contacts/new-list-modal";

export default function Page() {
  return (
    <>
      <div className="pt-6">
        <section>
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-base font-medium text-gray-200">
                Saved contacts
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                View and manage your investors' contacts
              </p>
            </div>

            <div className="flex items-center gap-2">
              <ManageListModal />
              <CreateContactModal />
            </div>
          </div>
        </section>

        <section className="mt-6">
          <ListContacts />
        </section>
      </div>
    </>
  );
}
