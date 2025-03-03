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
    .split("/")
    .filter(Boolean)
    .map((segment, index, array) => {
      const isLast = index === array.length - 1;
      const href = isLast ? undefined : `/${array.slice(0, index + 1).join("/")}`;
      return {
        label: segment.charAt(0).toUpperCase() + segment.slice(1),
        href,
      };
    });

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                {breadcrumbs.map((item, index) => (
                  <BreadcrumbItem key={index} className="hidden md:block">
                    {item.href ? (
                      <BreadcrumbLink asChild>
                        <Link href={item.href}>
                          {item.label}
                        </Link>
                      </BreadcrumbLink>
                    ) : (
                      <BreadcrumbPage>{item.label}</BreadcrumbPage>
                    )}
                    {index !== breadcrumbs.length - 1 && <BreadcrumbSeparator className="hidden md:block" />}
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
