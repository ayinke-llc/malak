import { Button } from "@/components/Button"

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
            <Button type="button" variant="primary"
              className="relative inline-flex items-center justify-center whitespace-nowrap rounded-md border px-3 py-2 text-center text-sm font-medium shadow-sm transition-all duration-100 ease-in-out disabled:pointer-events-none disabled:shadow-none outline outline-offset-2 outline-0 focus-visible:outline-2 outline-indigo-500 dark:outline-indigo-500 border-transparent text-white dark:text-gray-900 bg-indigo-600 dark:bg-indigo-500 hover:bg-indigo-500 dark:hover:bg-indigo-600 disabled:bg-indigo-100 disabled:text-gray-400 disabled:dark:bg-indigo-800 disabled:dark:text-indigo-400 mt-4 w-full gap-2 sm:mt-0 sm:w-fit">
              Add User
            </Button>
          </div>
        </section>
      </div>
    </>
  )
}
