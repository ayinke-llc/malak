"use client";

import { usePathname } from "next/navigation";
import { AppSidebar } from "@/components/ui/app-sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import Link from "next/link";

export function Sidebar({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  const breadcrumbs = pathname
    ? pathname
        .split("/")
        .filter(Boolean)
        .map((segment, index, array) => {
          const isLast = index === array.length - 1;
          const href = isLast ? undefined : `/${array.slice(0, index + 1).join("/")}`;
          return {
            label: isLast ? segment : segment.charAt(0).toUpperCase() + segment.slice(1),
            href,
          };
        })
    : [];

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="sticky top-0 z-10 flex h-14 sm:h-16 shrink-0 items-center gap-2 bg-background">
          <div className="flex items-center gap-2 px-2 sm:px-4">
            <SidebarTrigger className="sm:-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4 hidden sm:block" />
            <Breadcrumb>
              <BreadcrumbList className="overflow-hidden">
                <BreadcrumbItem>
                  <BreadcrumbLink asChild>
                    <Link href="/">Home</Link>
                  </BreadcrumbLink>
                  {breadcrumbs.length > 0 && <BreadcrumbSeparator />}
                </BreadcrumbItem>
                {breadcrumbs.map((item, index) => (
                  <BreadcrumbItem key={index}>
                    {item.href ? (
                      <BreadcrumbLink asChild>
                        <Link href={item.href} className="truncate">{item.label}</Link>
                      </BreadcrumbLink>
                    ) : (
                      <BreadcrumbPage className="truncate">{item.label}</BreadcrumbPage>
                    )}
                    {index !== breadcrumbs.length - 1 && <BreadcrumbSeparator />}
                  </BreadcrumbItem>
                ))}
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
