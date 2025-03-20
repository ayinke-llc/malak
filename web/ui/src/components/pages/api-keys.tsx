"use client"

import type { MalakAPIKey } from "@/client/Api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import client from "@/lib/client";
import { AnalyticsEvent } from "@/lib/events";
import { CREATE_API_KEY, LIST_API_KEYS } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RiAddLine,
  RiDeleteBinLine,
  RiFileCopyLine,
  RiLoader4Line
} from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { addDays, addHours, format } from "date-fns";
import { X } from "lucide-react";
import { usePostHog } from "posthog-js/react";
import { useState } from "react";
import CopyToClipboard from "react-copy-to-clipboard";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";


const apiKeySchema = yup.object().shape({
  key_name: yup.string()
    .required("Name is required")
    .min(3, "Name must be at least 3 characters")
    .max(20, "Name must not exceed 20 characters"),
});

type ApiKeyFormData = yup.InferType<typeof apiKeySchema>;

export default function ApiKeys() {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [selectedKey, setSelectedKey] = useState<MalakAPIKey | null>(null);
  const [revealedKeys, setRevealedKeys] = useState<Set<string>>(new Set());
  const [revokeExpiration, setRevokeExpiration] = useState<string>("immediately");
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const posthog = usePostHog();

  const form = useForm<ApiKeyFormData>({
    resolver: yupResolver(apiKeySchema) as any,
    defaultValues: {
      key_name: "",
    },
  });

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [LIST_API_KEYS],
    queryFn: async () => {
      const response = await client.developers.keysList();
      return response.data.keys;
    },
  });

  const apiKeys = data || [];

  const createMutation = useMutation({
    mutationKey: [CREATE_API_KEY],
    mutationFn: (data: ApiKeyFormData) => {
      return client.developers.keysCreate({
        title: data.key_name,
      });
    },
    onSuccess: ({ data }) => {
      setNewlyCreatedKey(data.value);
      toast.success(data.message);
      form.reset();
      setIsCreateDialogOpen(false);
      posthog?.capture(AnalyticsEvent.CreateApiKey);
      queryClient.invalidateQueries({ queryKey: [LIST_API_KEYS] });
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message || "Failed to create API key");
    },
  });

  const updateMutation = useMutation({
    mutationKey: ["UPDATE_API_KEY"],
    mutationFn: async (data: ApiKeyFormData) => {
      if (!selectedKey) throw new Error("No key selected");
      // TODO: Replace with actual API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { data: { message: "API key updated successfully" } };
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setSelectedKey(null);
      form.reset();
      queryClient.invalidateQueries({ queryKey: [LIST_API_KEYS] });
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message || "Failed to update API key");
    },
  });

  const revokeMutation = useMutation({
    mutationKey: ["REVOKE_API_KEY"],
    mutationFn: async ({ id, expiration }: { id: string, expiration: string }) => {
      // TODO: Replace with actual API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { data: { message: "API key revoked successfully" } };
    },
    onSuccess: ({ data }, variables) => {
      if (variables.expiration === "immediately") {
        toast.success("API key revoked immediately");
      } else {
        const hours = variables.expiration === "24h" ? 24 : 168;
        toast.success(`API key will expire in ${hours === 24 ? "24 hours" : "7 days"}`);
      }
      queryClient.invalidateQueries({ queryKey: [LIST_API_KEYS] });
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message || "Failed to revoke API key");
    },
  });

  const handleCreate = async (data: ApiKeyFormData) => {
    createMutation.mutate(data);
  };

  const toggleKeyVisibility = (keyId: string) => {
    setRevealedKeys(prev => {
      const newSet = new Set(prev);
      if (newSet.has(keyId)) {
        newSet.delete(keyId);
      } else {
        newSet.add(keyId);
      }
      return newSet;
    });
  };

  const handleRevoke = async (id: string) => {
    revokeMutation.mutate({ id, expiration: revokeExpiration });
  };

  const getRevokeTimeLabel = (value: string) => {
    const now = new Date();
    switch (value) {
      case "immediately":
        return "Immediately. This action cannot be undone.";
      case "24h":
        return `At ${format(addHours(now, 24), "PPp")}`;
      case "7d":
        return `At ${format(addDays(now, 7), "PPp")}`;
      default:
        return "Immediately. This action cannot be undone.";
    }
  };

  const columns: ColumnDef<MalakAPIKey>[] = [
    {
      accessorKey: "key_name",
      header: "Name",
    },
    {
      accessorKey: "created_at",
      header: "Created At",
      cell: ({ row }) => row.original.created_at ? format(new Date(row.original.created_at), "PPp") : "-",
    },
    {
      accessorKey: "expires_at",
      header: "Expires At",
      cell: ({ row }) => row.original.expires_at ? format(new Date(row.original.expires_at), "PPp") : "Never",
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => {
        const key = row.original;
        return (
          <div className="flex items-center space-x-2">
            <Dialog>
              <DialogTrigger asChild>
                <Button variant="ghost" size="icon">
                  <RiDeleteBinLine className="w-4 h-4" />
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Revoke API Key</DialogTitle>
                  <DialogDescription>
                    Choose when to revoke this API key
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Revoke</label>
                    <Select
                      value={revokeExpiration}
                      onValueChange={setRevokeExpiration}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select when to revoke" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="immediately">Immediately</SelectItem>
                        <SelectItem value="24h">In 24 hours</SelectItem>
                        <SelectItem value="7d">In 7 days</SelectItem>
                      </SelectContent>
                    </Select>
                    <p className="text-sm text-muted-foreground">
                      {getRevokeTimeLabel(revokeExpiration)}
                    </p>
                  </div>
                  <div className="flex justify-end space-x-2">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => setRevokeExpiration("immediately")}
                    >
                      Cancel
                    </Button>
                    <Button
                      variant="destructive"
                      onClick={() => handleRevoke(key.id || "")}
                    >
                      Revoke
                    </Button>
                  </div>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        );
      },
    },
  ];

  const table = useReactTable({
    columns,
    data: apiKeys,
    getCoreRowModel: getCoreRowModel(),
  });

  const isMutating = createMutation.isPending || updateMutation.isPending || revokeMutation.isPending;

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <Skeleton className="h-7 w-32" />
          <Skeleton className="h-10 w-[140px]" />
        </div>

        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <Skeleton className="h-4 w-24" />
                </TableHead>
                <TableHead>
                  <Skeleton className="h-4 w-32" />
                </TableHead>
                <TableHead>
                  <Skeleton className="h-4 w-32" />
                </TableHead>
                <TableHead>
                  <Skeleton className="h-4 w-20" />
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {Array.from({ length: 3 }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell>
                    <Skeleton className="h-5 w-32" />
                  </TableCell>
                  <TableCell>
                    <Skeleton className="h-5 w-40" />
                  </TableCell>
                  <TableCell>
                    <Skeleton className="h-5 w-40" />
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      <Skeleton className="h-8 w-8 rounded" />
                      <Skeleton className="h-8 w-8 rounded" />
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    );
  }

  if (isMutating) {
    return (
      <div className="flex items-center justify-center h-64">
        <RiLoader4Line className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-[400px] p-8">
        <div className="max-w-2xl mx-auto bg-background border rounded-lg shadow-sm">
          <div className="p-6">
            <div className="flex items-center space-x-3">
              <div className="h-10 w-10 rounded-full bg-red-50 flex items-center justify-center">
                <RiDeleteBinLine className="h-5 w-5 text-red-500" />
              </div>
              <div>
                <h3 className="text-base font-semibold">Error Loading API Keys</h3>
                <p className="text-sm text-muted-foreground">
                  We couldn't retrieve your API keys at this time
                </p>
              </div>
            </div>

            <div className="mt-4 bg-muted/50 rounded-md p-4">
              <p className="text-sm text-muted-foreground font-mono">
                {error instanceof Error ? (
                  <span className="text-destructive">{error.message}</span>
                ) : (
                  "An unexpected error occurred while fetching your API keys"
                )}
              </p>
            </div>

            <div className="mt-6">
              <div className="flex items-center justify-between p-4 bg-muted/30 rounded-md">
                <div className="space-y-1">
                  <p className="text-sm font-medium">Retry the request</p>
                  <p className="text-sm text-muted-foreground">
                    Attempt to fetch your API keys again
                  </p>
                </div>
                <Button
                  variant="outline"
                  onClick={() => refetch()}
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <>
                      <RiLoader4Line className="w-4 h-4 mr-2 animate-spin" />
                      Retrying...
                    </>
                  ) : (
                    <>
                      <RiLoader4Line className="w-4 h-4 mr-2" />
                      Try Again
                    </>
                  )}
                </Button>
              </div>
            </div>

            <div className="mt-6 flex items-center gap-2 text-sm text-muted-foreground">
              <p>
                Error Code: <code className="text-xs bg-muted px-1 py-0.5 rounded">ERR_API_KEYS_FETCH</code>
              </p>
              <span>•</span>
              <p>
                Timestamp: {new Date().toLocaleTimeString()}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-lg font-medium">API Keys</h2>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <RiAddLine className="w-4 h-4 mr-2" />
              Create API Key
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create API Key</DialogTitle>
              <DialogDescription>
                Create a new API key for your application
              </DialogDescription>
            </DialogHeader>
            <Form {...form}>
              <form onSubmit={form.handleSubmit(handleCreate)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="key_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Name</FormLabel>
                      <FormControl>
                        <Input placeholder="Enter key name" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className="flex justify-end space-x-2">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setIsCreateDialogOpen(false)}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    disabled={createMutation.isPending || !form.formState.isValid || !form.formState.isDirty}
                  >
                    {createMutation.isPending ? (
                      <>
                        <RiLoader4Line className="mr-2 h-4 w-4 animate-spin" />
                        Creating...
                      </>
                    ) : (
                      "Create Key"
                    )}
                  </Button>
                </div>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
      </div>

      {newlyCreatedKey && (
        <div className="max-w-2xl mx-auto p-4 bg-muted rounded-lg border">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium">New API Key Created</h3>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setNewlyCreatedKey(null)}
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
          <p className="text-sm text-muted-foreground mb-3">
            Make sure to copy your API key now. You won't be able to see it again!
          </p>
          <div className="flex items-center space-x-2 p-3 bg-background rounded-md">
            <code className="font-mono text-sm flex-1">{newlyCreatedKey}</code>
            <CopyToClipboard text={newlyCreatedKey} onCopy={() => toast.success("API key copied to clipboard")}>
              <Button variant="ghost" size="icon">
                <RiFileCopyLine className="w-4 h-4" />
              </Button>
            </CopyToClipboard>
          </div>
        </div>
      )}

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No API keys found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
} 
