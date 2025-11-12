"use client";

import { Sidebar } from "@/components/ui/navigation/sidebar";
import EmailVerificationBadge from "@/components/EmailVerificationNotice";

export default function Layout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <Sidebar>
      <main className="min-h-screen w-full lg:pl-20">
        <div className="px-2 py-4 sm:px-6 sm:pb-10 sm:pt-10 lg:px-10 lg:pt-7">
          <EmailVerificationBadge />
          {children}
        </div>
      </main>
    </Sidebar>
  );
}
