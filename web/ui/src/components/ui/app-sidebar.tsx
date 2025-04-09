"use client"

import * as React from "react"
import Link from "next/link"
import { NavMain } from "@/components/ui/nav-main"
import { NavUser } from "@/components/ui/nav-user"
import { TeamSwitcher } from "@/components/ui/team-switcher"
import { Sparkles } from "lucide-react"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  useSidebar,
} from "@/components/ui/sidebar"
import { links } from "./navigation/navlist"
import useAuthStore from "@/store/auth"
import useWorkspacesStore from "@/store/workspace"
import { RiLinksFill } from "@remixicon/react"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const user = useAuthStore(state => state.user)
  const workspace = useWorkspacesStore(state => state.current)
  const { state } = useSidebar()

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        {state === "expanded" && (
          <Link
            href="https://www.youtube.com/watch?v=BdFBOPWKRO4"
            target="_blank"
            rel="noopener noreferrer"
            className="group block mx-2 mb-4"
          >
            <div className="relative overflow-hidden rounded-lg border border-primary/20 bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5 p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/20 text-primary group-hover:bg-primary/30 transition-colors">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    className="w-4 h-4"
                  >
                    <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" />
                  </svg>
                </div>
                <div className="flex-1">
                  <h4 className="font-medium text-primary mb-0.5">Dashboard walkthrough</h4>
                  <p className="text-xs text-muted-foreground">Learn how to use Malak effectively</p>
                </div>
              </div>
              <div className="absolute -right-2 -top-2 h-16 w-16 rotate-12 bg-gradient-to-br from-primary/30 to-transparent opacity-20 blur-2xl group-hover:opacity-30 transition-opacity" />
            </div>
          </Link>
        )}
        <NavMain items={links} />
        {state === "expanded" && (
          <div className="mt-4 px-2">
            <div className="text-xs font-medium text-muted-foreground mb-2 px-2">Resources</div>
            <Link
              href="https://changelog.malak.vc"
              target="_blank"
              rel="noopener noreferrer"
              className="group flex items-center gap-2 rounded-md px-2 py-1.5 text-sm text-muted-foreground hover:text-foreground hover:bg-accent/50 transition-colors"
            >
              <RiLinksFill className="w-4 h-4" />
              Changelog
            </Link>
          </div>
        )}
      </SidebarContent>
      <SidebarFooter>
        {!workspace?.is_subscription_active && (
          <Link
            href="/settings?tab=billing"
            className="group block mx-2 mb-4"
          >
            <div className="relative overflow-hidden rounded-lg border border-primary/20 bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5 p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/20 text-primary group-hover:bg-primary/30 transition-colors">
                  <Sparkles size={16} />
                </div>
                <div className="flex-1">
                  <h4 className="font-medium text-primary mb-0.5">Upgrade Now</h4>
                  <p className="text-xs text-muted-foreground">Get access to premium features and more</p>
                </div>
              </div>
              <div className="absolute -right-2 -top-2 h-16 w-16 rotate-12 bg-gradient-to-br from-primary/30 to-transparent opacity-20 blur-2xl group-hover:opacity-30 transition-opacity" />
            </div>
          </Link>
        )}
        <NavUser user={{
          email: user?.email as string,
          full_name: user?.full_name as string
        }} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
