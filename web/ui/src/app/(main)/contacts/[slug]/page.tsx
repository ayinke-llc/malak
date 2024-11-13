"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import CreateContactModal from "@/components/ui/contacts/modal";
import ManageListModal from "@/components/ui/contacts/new-list-modal";
import ContactDetails from "@/components/ui/contacts/single/view";
import { useParams, useRouter } from "next/navigation";

export default function Page() {

  const params = useParams();

  const router = useRouter();

  const reference = params.slug as string;

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
                Viewing details of contact
              </h3>
            </div>

            <div className="flex justify-center gap-2">
              <ManageListModal />
              <CreateContactModal />
            </div>
          </div>
        </section>

        <section>
          <ContactDetails />
        </section>
      </div>
    </>
  );
}
