"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import ContactDetails from "@/components/ui/contacts/single/view";
import { useParams, useRouter } from "next/navigation";

export default function Page() {
  const params = useParams();
  const router = useRouter();
  const reference = params.slug as string;

  return (
    <div className="pt-6">
      <section>
        <div className="sm:flex sm:items-center sm:justify-between">
          <div>
            <h3
              id="company-decks"
              className="text-lg font-medium text-zinc-100"
            >
              Contact details
            </h3>
            <p className="text-sm text-zinc-400/80">
              Viewing contact's details
            </p>
          </div>
        </div>
      </section>

      <section>
        <ContactDetails reference={reference} />
      </section>
    </div>
  );
}
