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
} from "@/components/ui/sidebar"
import { links } from "./navigation/navlist"
import useAuthStore from "@/store/auth"
import useWorkspacesStore from "@/store/workspace"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  const user = useAuthStore(state => state.user)
  const workspace = useWorkspacesStore(state => state.current)

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={links} />
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
