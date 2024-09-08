import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
})

import { Sidebar } from "@/components/ui/navigation/sidebar"
import { siteConfig } from "./siteConfig"
import Providers from "@/components/providers/providers"

export const metadata: Metadata = {
  metadataBase: new URL("https://malak.vc"),
  title: siteConfig.name,
  description: siteConfig.description,
  keywords: [],
  authors: [
    {
      name: "Ayinke LLC",
      url: "https://ayinke.ventures",
    },
  ],
  creator: "Ayinke Ventures",
  openGraph: {
    type: "website",
    locale: "en_US",
    url: siteConfig.url,
    title: siteConfig.name,
    description: siteConfig.description,
    siteName: siteConfig.name,
  },
  twitter: {
    card: "summary_large_image",
    title: "Ayinke Ventures",
    creator: "@ayinkeventures",
  },
  icons: {
    icon: "/favicon.ico",
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body
        className={`${inter.className} overflow-y-scroll scroll-auto antialiased selection:bg-indigo-100 selection:text-indigo-700 dark:bg-gray-950`}
        suppressHydrationWarning
      >
        <Providers>
          <div className="mx-auto max-w-screen-2xl">
            <Sidebar />
            <main className="lg:pl-72">{children}</main>
          </div>
        </Providers>
      </body>
    </html>
  )
}
