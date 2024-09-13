"use client"

import Editor from "@/components/ui/updates/editor/editor"
import { useState } from "react";



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
            <span>oops</span>
          </div>

          <div className="mt-5">
            <Editor />
            <h1>FUCK</h1>
          </div>
        </section>
      </div>
    </>
  )
}
