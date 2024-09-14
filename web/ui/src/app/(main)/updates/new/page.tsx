"use client"

import { Button } from "@/components/Button";
import NovelEditor from "@/components/ui/updates/editor/editor"

export default function Page() {

  return (
    <>
      <div className="pt-6">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3 id="existing-contacts" className="scroll-mt-10 font-semibold text-gray-900 dark:text-gray-50">
                Create a new update
              </h3>
              <p className="text-sm leading-6 text-gray-500">
                Sending a new update to your investors
              </p>
            </div>
            <Button type="submit">
              Send
            </Button>
          </div>

          <div className="mt-5">
            <NovelEditor />
          </div>
        </section>
      </div>
    </>
  )
}
