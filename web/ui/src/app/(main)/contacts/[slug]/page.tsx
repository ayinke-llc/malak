"use client";

import { MalakContact, MalakContactShareItem } from "@/client/Api";
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
        <ContactDetails
          reference={reference}
          shared_items={data?.data?.shared_items as MalakContactShareItem[]}
          contact={data?.data?.contact as MalakContact}
        />
      </section>
    </div>
  );
}
