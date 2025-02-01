"use client"

import { ChevronsUpDown, Plus } from "lucide-react"

import { ServerAPIStatus } from "@/client/Api"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import client from "@/lib/client"
import { SWITCH_WORKSPACE } from "@/lib/query-constants"
import { cn } from "@/lib/utils"
import useWorkspacesStore from "@/store/workspace"
import { useMutation } from "@tanstack/react-query"
import { AxiosResponse } from "axios"
import { useRouter } from "next/navigation"
import { toast } from "sonner"
import { Avatar } from "./custom/avatar/avatar"
import { ModalAddWorkspace } from "./navigation/ModalAddWorkspace"

export function TeamSwitcher() {
  const { isMobile } = useSidebar();
  const { current, workspaces, setCurrent } = useWorkspacesStore();

  const router = useRouter();

  const mutation = useMutation({
    mutationKey: [SWITCH_WORKSPACE],
    mutationFn: (reference: string) => client.workspaces.switchworkspace(reference),
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
              <div
                className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                {current?.logo_url ? (
                  <img
                    className="size-4 shrink-0"
                    src={current?.logo_url as string}
                    alt={`${current?.workspace_name}'s logo`} />
                ) : <Avatar className="size-4 shrink-0" />}
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
                disabled={mutation.isPending}
              >
                <div className="flex size-6 items-center justify-center rounded-sm border">
                  {workspace?.logo_url ? (
                    <img
                      className="size-4 shrink-0"
                      src={workspace?.logo_url as string}
                      alt={`${workspace?.workspace_name}'s logo`} />
                  ) : <Avatar className="size-4 shrink-0" />}
                </div>
                <span className={cn(
                  workspace?.id == current?.id ? "font-bold" : ""
                )}>
                  {workspace.workspace_name}
                </span>
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className="gap-2 p-2 hover:cursor-pointer" asChild>
              <>
                <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                  <Plus className="size-4" />
                </div>
                <ModalAddWorkspace
                  onSelect={() => { }}
                  onOpenChange={() => { }}
                  itemName="Create workspace"
                />
              </>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
