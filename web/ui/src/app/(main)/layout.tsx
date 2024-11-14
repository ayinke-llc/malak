"use client";

import { Sidebar } from "@/components/ui/navigation/sidebar";

export default function Layout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {

  return (
    <>
      <Sidebar>
        <main className="lg:pl-20">
          <div className="p-4 sm:px-6 sm:pb-10 sm:pt-10 lg:px-10 lg:pt-7">
            {children}
          </div>
        </main>
      </Sidebar>
    </>
  );
}
