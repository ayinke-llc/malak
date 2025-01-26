"use client"

import { ChevronsUpDown, Plus } from "lucide-react"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { Avatar } from "./custom/avatar/avatar"
import useWorkspacesStore from "@/store/workspace"
import { cn } from "@/lib/utils"
import { ModalAddWorkspace } from "./navigation/ModalAddWorkspace"

export function TeamSwitcher() {
  const { isMobile } = useSidebar();
  const { current, workspaces, setCurrent } = useWorkspacesStore();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu modal={false}>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <Avatar className="size-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">
                  {current?.workspace_name}
                </span>
                <span className="truncate text-xs">{"Pro"}</span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            align="start"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-xs text-muted-foreground">
              Teams
            </DropdownMenuLabel>
            {workspaces.map((workspace, index) => (
              <DropdownMenuItem
                key={workspace.reference}
                onClick={() => setCurrent(workspace)}
                className="gap-2 p-2 hover:cursor-pointer"
              >
                <div className="flex size-6 items-center justify-center rounded-sm border">
                  <Avatar className="size-4 shrink-0" />
                </div>
                <span className={cn(
                  workspace?.id == current?.id ? "font-bold" : ""
                )}>
                  {workspace.workspace_name}
                </span>
                <DropdownMenuShortcut>⌘{index + 1}</DropdownMenuShortcut>
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className="gap-2 p-2 hover:cursor-pointer">
              <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                <Plus className="size-4" />
              </div>
              <ModalAddWorkspace
                onSelect={() => { }}
                onOpenChange={() => { }}
                itemName="Create workspace"
              />
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
