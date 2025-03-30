"use client";

import { MalakContact, MalakContactShareItem } from "@/client/Api";
import ContactDetails from "@/components/ui/contacts/single/view";
import client from "@/lib/client";
import { FETCH_CONTACT } from "@/lib/query-constants";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { RiErrorWarningLine } from "@remixicon/react";

interface ContactDetailsPageProps {
  reference: string;
}

export default function ContactDetailsPage({ reference }: ContactDetailsPageProps) {
  const router = useRouter();

  const { data, error, isLoading, refetch } = useQuery({
    queryKey: [FETCH_CONTACT, reference],
    queryFn: () => client.contacts.contactsDetail(reference),
  });

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] p-4">
        <div className="text-center space-y-4">
          <RiErrorWarningLine className="h-12 w-12 text-destructive mx-auto" />
          <h3 className="text-lg font-semibold text-foreground">Failed to load contact</h3>
          <p className="text-muted-foreground max-w-md">
            We couldn't load this contact. This might be because it was deleted or you don't have permission to view it.
          </p>
          <div className="space-x-4">
            <button
              onClick={() => refetch()}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
            >
              Try Again
            </button>
            <button
              onClick={() => router.push("/contacts")}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-foreground bg-muted hover:bg-muted/80 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-ring"
            >
              Return to Contacts
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="pt-6">
      <section>
        <ContactDetails
          reference={reference}
          shared_items={data?.data?.shared_items as MalakContactShareItem[]}
          contact={data?.data?.contact as MalakContact}
          isLoading={isLoading}
        />
      </section>
    </div>
  );
} 
