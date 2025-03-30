"use client";

import ListContacts from "@/components/ui/contacts/list/list";
import CreateContactModal from "@/components/ui/contacts/modal";
import CSVUploadModal from "@/components/ui/contacts/csv-upload-modal";
import { Button } from "@/components/ui/button";
import { RiSettings4Line } from "@remixicon/react";
import Link from "next/link";

export default function Contacts() {
  return (
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
              View and manage the contact details of your investors
            </p>
          </div>

          <div className="flex gap-2">
            <Link href="/contacts/lists">
              <Button variant="default" className="whitespace-nowrap gap-1">
                <RiSettings4Line />
                Manage lists
              </Button>
            </Link>
            <CSVUploadModal />
            <CreateContactModal />
          </div>
        </div>
      </section>

      <section className="mt-10 sm:mt-4">
        <ListContacts />
      </section>
    </div>
  );
} 