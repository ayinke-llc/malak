"use client"

import { usePathname } from "next/navigation"

export default function Layout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {

  const pathName = usePathname()

  return (
    <div className="p-4 sm:px-6 sm:pb-10 sm:pt-10 lg:px-10 lg:pt-7">
      <h1 className="text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50 capitalize">
        {pathName.replace("/", "")}
      </h1>
      <div>{children}</div>
    </div>
  )
}
