import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/providers";
import { siteConfig } from "./siteConfig";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
});

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
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        suppressHydrationWarning
        className={`${inter.className} overflow-y-scroll scroll-auto antialiased 
selection:bg-indigo-100 selection:text-indigo-700 dark:bg-gray-950`}
      >
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
