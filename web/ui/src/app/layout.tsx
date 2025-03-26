export const dynamic = "force-dynamic";
export const fetchCache = "force-no-store"

import type { Metadata } from "next";
import { GeistSans } from "geist/font/sans";
import "./globals.css";
import Providers from "@/components/providers/providers";
import { siteConfig } from "./siteConfig";
import { ThemeProvider } from "@/components/providers/theme";
import { TooltipProvider } from "@radix-ui/react-tooltip";
import { Suspense } from "react";
import { ErrorBoundary } from "@/components/error-boundary";
import { BetaBanner } from "@/components/ui/beta-banner";

export const metadata: Metadata = {
  metadataBase: new URL("https://app.malak.vc"),
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
};


export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        suppressHydrationWarning={true}
        className={`${GeistSans.className} overflow-y-scroll scroll-auto antialiased dark:bg-gray-950 theme-custom`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          disableTransitionOnChange
        >
          <Suspense>
            <ErrorBoundary>
              <TooltipProvider>
                <BetaBanner />
                <Providers>{children}</Providers>
              </TooltipProvider>
            </ErrorBoundary>
          </Suspense>
        </ThemeProvider>
      </body>
    </html>
  );
}
