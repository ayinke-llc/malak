"use client"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/Dropdown"
import { cx, focusInput } from "@/lib/utils"
import { RiArrowRightSLine, RiExpandUpDownLine } from "@remixicon/react"
import React, { useState } from "react"
import { ModalAddWorkspace } from "./ModalAddWorkspace"
import useWorkspacesStore from "@/store/workspace"
import { useMutation } from "@tanstack/react-query"
import { SWITCH_WORKSPACE } from "@/lib/query-constants"
import client from "@/lib/client"
import { ServerAPIStatus } from "@/client/Api"
import { AxiosError } from "axios"
import { toast } from "sonner"
import Skeleton from "../custom/loader/skeleton"


export const WorkspacesDropdownDesktop = () => {
  const [dropdownOpen, setDropdownOpen] = React.useState(false)
  const [hasOpenDialog, setHasOpenDialog] = React.useState(false)
  const dropdownTriggerRef = React.useRef<null | HTMLButtonElement>(null)
  const focusRef = React.useRef<null | HTMLButtonElement>(null)

  const [loading, setLoading] = useState<boolean>(false)

  const workspaces = useWorkspacesStore.getState().workspaces
  const current = useWorkspacesStore.getState().current
  const setCurrent = useWorkspacesStore.getState().setCurrent

  const handleDialogItemSelect = () => {
    focusRef.current = dropdownTriggerRef.current
  }

  const handleDialogItemOpenChange = (open: boolean) => {
    setHasOpenDialog(open)
    if (!open) {
      setDropdownOpen(false)
    }
  }

  const mutation = useMutation({
    mutationKey: [SWITCH_WORKSPACE],
    mutationFn: (reference: string) => client.workspaces.switchworkspace(reference),
    onSuccess: ({ data }) => {
      setCurrent(data.workspace)
      toast.success(data.message)
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response?.data.message
      }
      toast.error(msg)
    },
    retry: false,
    gcTime: Infinity,
    onSettled: () => setLoading(false),
  })

  return (
    <div suppressHydrationWarning={true}>
      <DropdownMenu
        open={dropdownOpen}
        onOpenChange={setDropdownOpen}
        modal={false}
      >
        <DropdownMenuTrigger asChild>
          <button
            className={cx(
              "flex w-full items-center gap-x-2.5 rounded-md border border-gray-300 bg-white p-2 text-sm shadow-sm transition-all hover:bg-gray-50 dark:border-gray-800 dark:bg-gray-950 hover:dark:bg-gray-900",
              focusInput,
            )}
          >
            <span
              className="uppercase flex aspect-square size-8 items-center justify-center rounded bg-indigo-600 p-2 text-xs font-medium text-white dark:bg-indigo-500"
              aria-hidden="true"
            >
              {current?.workspace_name?.split(' ')
                .slice(0, 2)
                .map((name) => name[0])
                .join('')}
            </span>
            <div className="flex w-full items-center justify-between gap-x-4 truncate">
              <div className="truncate">
                <p className="truncate whitespace-nowrap text-sm font-medium text-indigo-600 capitalize">
                  {current?.workspace_name}
                </p>
                <p className="whitespace-nowrap text-left text-xs text-indigo-600 capitalize">
                  Admin
                </p>
              </div>
              <RiExpandUpDownLine
                className="size-5 shrink-0 text-gray-500"
                aria-hidden="true"
              />
            </div>
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          hidden={hasOpenDialog}
          onCloseAutoFocus={(event) => {
            if (focusRef.current) {
              focusRef.current.focus()
              focusRef.current = null
              event.preventDefault()
            }
          }}
        >
          <DropdownMenuGroup>
            <DropdownMenuLabel>
              Workspaces ({workspaces.length})
            </DropdownMenuLabel>
            {loading ? <Skeleton count={2} /> : workspaces.map((workspace) => (
              <DropdownMenuItem key={workspace.reference} onClick={() => {
                mutation.mutate(workspace.reference as string)
              }}>
                <div className="flex w-full items-center gap-x-2.5">
                  <span
                    className={cx(
                      "bg-indigo-600 dark:bg-indigo-500",
                      "uppercase flex aspect-square size-8 items-center justify-center rounded p-2 text-xs font-medium text-white",
                    )}
                    aria-hidden="true"
                  >
                    {workspace.workspace_name?.split(' ')
                      .slice(0, 2)
                      .map((name) => name[0])
                      .join('')}
                  </span>
                  <div className={cx(
                    workspace.reference === current?.reference && "text-indigo-600 dark:text-indigo-400",
                  )}>
                    <p className="text-sm font-medium">
                      {workspace.workspace_name}
                    </p>
                    <p className="text-xs">
                      {workspace.reference}
                    </p>
                  </div>
                </div>
              </DropdownMenuItem>
            ))}
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <ModalAddWorkspace
            onSelect={handleDialogItemSelect}
            onOpenChange={handleDialogItemOpenChange}
            itemName="Add workspace"
          />
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}

export const WorkspacesDropdownMobile = () => {
  const [dropdownOpen, setDropdownOpen] = React.useState(false)
  const [hasOpenDialog, setHasOpenDialog] = React.useState(false)
  const dropdownTriggerRef = React.useRef<null | HTMLButtonElement>(null)
  const focusRef = React.useRef<null | HTMLButtonElement>(null)

  const workspaces = useWorkspacesStore.getState().workspaces
  const current = useWorkspacesStore.getState().current
  const setCurrent = useWorkspacesStore.getState().setCurrent

  const [loading, setLoading] = useState<boolean>(false)

  const handleDialogItemSelect = () => {
    focusRef.current = dropdownTriggerRef.current
  }

  const handleDialogItemOpenChange = (open: boolean) => {
    setHasOpenDialog(open)
    if (open === false) {
      setDropdownOpen(false)
    }
  }
  return (
    <div suppressHydrationWarning={true}>
      <DropdownMenu
        open={dropdownOpen}
        onOpenChange={setDropdownOpen}
        modal={false}
      >
        <DropdownMenuTrigger asChild>
          <button className="flex items-center gap-x-1.5 rounded-md p-2 hover:bg-gray-100 focus:outline-none hover:dark:bg-gray-900">
            <span
              className={cx(
                "uppercase",
                "flex aspect-square size-7 items-center justify-center rounded bg-indigo-600 p-2 text-xs font-medium text-white dark:bg-indigo-500",
              )}
              aria-hidden="true"
            >
              {current?.workspace_name?.split(' ')
                .slice(0, 2)
                .map((name) => name[0])
                .join('')}
            </span>
            <RiArrowRightSLine
              className="size-4 shrink-0 text-gray-500"
              aria-hidden="true"
            />
            <div className="flex w-full items-center justify-between gap-x-3 truncate">
              <p className="truncate whitespace-nowrap text-sm font-medium text-indigo-600 capitalize">
                {current?.workspace_name}
              </p>
              <RiExpandUpDownLine
                className="size-4 shrink-0 text-gray-500"
                aria-hidden="true"
              />
            </div>
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          className="!min-w-72"
          hidden={hasOpenDialog}
          onCloseAutoFocus={(event) => {
            if (focusRef.current) {
              focusRef.current.focus()
              focusRef.current = null
              event.preventDefault()
            }
          }}
        >
          <DropdownMenuGroup>
            <DropdownMenuLabel>
              Workspaces ({workspaces.length})
            </DropdownMenuLabel>
            {loading ? <Skeleton count={3} /> : workspaces.map((workspace) => (
              <DropdownMenuItem key={workspace.reference}>
                <div className="flex w-full items-center gap-x-2.5">
                  <span
                    className={cx(
                      "bg-indigo-600 dark:bg-indigo-500",
                      "uppercase flex size-8 items-center justify-center rounded p-2 text-xs font-medium text-white",
                    )}
                    aria-hidden="true"
                  >
                    {workspace.workspace_name?.split(' ')
                      .slice(0, 2)
                      .map((name) => name[0])
                      .join('')}
                  </span>
                  <div className={cx(
                    workspace.reference === current?.reference && "text-indigo-600 dark:text-indigo-400",
                  )}>
                    <p className="text-sm font-medium text-gray-900 dark:text-gray-50">
                      {workspace.workspace_name}
                    </p>
                    <p className="text-xs text-gray-700 dark:text-gray-300">
                      {workspace.reference}
                    </p>
                  </div>
                </div>
              </DropdownMenuItem>
            ))}
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <ModalAddWorkspace
            onSelect={handleDialogItemSelect}
            onOpenChange={handleDialogItemOpenChange}
            itemName="Add workspace"
          />
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
