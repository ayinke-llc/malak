"use client"

import { ChevronsUpDown } from "lucide-react"
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
import { toast } from "sonner"
import { ModalAddWorkspace } from "./navigation/ModalAddWorkspace"
import { useState, useEffect } from "react"
import { Dialog, DialogContent } from "@/components/ui/dialog"
import useAuthStore from "@/store/auth"

export function TeamSwitcher() {
  const { isMobile } = useSidebar();
  const isAuthenticated = useAuthStore(state => state.isAuthenticated);

  const current = useWorkspacesStore(state => state.current)
  const workspaces = useWorkspacesStore(state => state.workspaces)
  const setCurrent = useWorkspacesStore(state => state.setCurrent);

  const [isDropdownOpen, setDropdownOpen] = useState(false);
  const [showCreateWorkspace, setShowCreateWorkspace] = useState(false);

  useEffect(() => {
    if (!isAuthenticated()) {
      setShowCreateWorkspace(false);
      return;
    }

    if (!workspaces || workspaces.length === 0) {
      setShowCreateWorkspace(true);
    } else {
      setShowCreateWorkspace(false);
    }
  }, [workspaces, isAuthenticated]);

  // Add cleanup effect
  useEffect(() => {
    return () => {
      setShowCreateWorkspace(false);
    };
  }, []);

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

  if (!isAuthenticated()) {
    return null;
  }

  if (showCreateWorkspace) {
    return (
      <Dialog open={showCreateWorkspace} onOpenChange={(open) => !open && workspaces?.length > 0 && setShowCreateWorkspace(false)}>
        <DialogContent className="sm:max-w-[425px]" onInteractOutside={(e) => {
          // Prevent closing when clicking outside if no workspaces exist
          if (!workspaces?.length) {
            e.preventDefault();
          }
        }}>
          <ModalAddWorkspace
            onSelect={() => setShowCreateWorkspace(false)}
            onOpenChange={(open) => !open && workspaces?.length > 0 && setShowCreateWorkspace(false)}
            itemName="Create workspace"
            forceOpen={true}
          />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu modal={false}
          open={isDropdownOpen}
          onOpenChange={setDropdownOpen}
        >
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
              <div
                className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <img
                  className="size-4 shrink-0"
                  src={current?.logo_url || "https://res.cloudinary.com/dwkjke5ea/image/upload/v1742121952/malak/logos/mtnjuwfl0gb9r11pz5qg.svg"}
                  alt={`${current?.workspace_name}'s logo`} />
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
            {workspaces?.map((workspace, index) => (
              <DropdownMenuItem
                key={workspace.reference}
                onClick={() => {
                  mutation.mutate(workspace?.reference as string)
                }}
                className="gap-2 p-2 hover:cursor-pointer"
                disabled={mutation.isPending}
              >
                <div className="flex size-6 items-center justify-center rounded-sm border">
                  <img
                    className="size-4 shrink-0"
                    src={workspace?.logo_url || "https://res.cloudinary.com/dwkjke5ea/image/upload/v1742121952/malak/logos/mtnjuwfl0gb9r11pz5qg.svg"}
                    alt={`${workspace?.workspace_name}'s logo`} />
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
