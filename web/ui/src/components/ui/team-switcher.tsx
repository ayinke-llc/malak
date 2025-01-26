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
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { ServerAPIStatus } from "@/client/Api"
import client from "@/lib/client"
import { SWITCH_WORKSPACE } from "@/lib/query-constants"
import { useMutation } from "@tanstack/react-query"
import { AxiosError, AxiosResponse } from "axios"
import { toast } from "sonner"

export function TeamSwitcher() {
  const { isMobile } = useSidebar();
  const { current, workspaces, setCurrent } = useWorkspacesStore();

  const router = useRouter();

  const mutation = useMutation({
    mutationKey: [SWITCH_WORKSPACE],
    mutationFn: (reference: string) =>
      client.workspaces.switchworkspace(reference),
    onSuccess: ({ data }) => {
      setCurrent(data.workspace);
      toast.success(data.message);
      window.location.reload();
    },
    onError(err: AxiosResponse<ServerAPIStatus>) {
      toast.error(err?.data?.message);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

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
                onClick={() => {
                  mutation.mutate(workspace?.reference as string)
                }}
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
