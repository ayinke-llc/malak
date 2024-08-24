import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/providers";
import Shell from "@/components/shell/dashboard";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Malak",
  description: "OSS Investors' hub for founders",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Providers>
          <Shell>
            {children}
          </Shell>
        </Providers>
      </body>
    </html>
  );
}
