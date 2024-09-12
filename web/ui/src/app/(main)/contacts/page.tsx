"use client"

import CreateContactModal from "@/components/ui/contacts/modal"

export default function Settings() {
  return (
    <>
      <div className="pt-6">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3 id="existing-contacts" className="scroll-mt-10 font-semibold text-gray-900 dark:text-gray-50">
                Stored Contacts
              </h3>
              <p className="text-sm leading-6 text-gray-500">
                View and manage all your existing contacts
              </p>
            </div>
            <CreateContactModal onOpenChange={() => {
              console.log("oops")
            }} />
          </div>
        </section>
      </div>
    </>
  )
}
