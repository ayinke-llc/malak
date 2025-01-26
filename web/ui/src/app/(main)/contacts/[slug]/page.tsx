"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import ContactDetails from "@/components/ui/contacts/single/view";
import client from "@/lib/client";
import { FETCH_CONTACT } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";

export default function Page() {
  const params = useParams();
  const router = useRouter();
  const reference = params.slug as string;

  const { data, error, isLoading } = useQuery({
    queryKey: [FETCH_CONTACT],
    queryFn: () => client.contacts.contactsDetail(reference),
  });

  if (error) {
    toast.error("an error occurred while fetching this contact");
    router.push("/contacts");
  }

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
