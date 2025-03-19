"use client"

import { useState } from "react";
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Badge } from "@/components/ui/badge";
import { X } from "lucide-react";
import {
  RiAddLine,
  RiEditLine,
  RiDeleteBinLine,
  RiFileCopyLine,
  RiLoader4Line,
  RiEyeLine,
  RiEyeOffLine,
} from "@remixicon/react";
import { toast } from "sonner";
import { format, addHours, addDays } from "date-fns";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import CopyToClipboard from "react-copy-to-clipboard";
import { cn } from "@/lib/utils";
import { MultiSelect } from "@/components/ui/multi-select";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  ColumnDef,
} from "@tanstack/react-table";

// Mock data for initial list
const mockApiKeys = [
  {
    id: "1",
    name: "Production API Key",
    key: "pk_live_123456789",
    created_at: new Date().toISOString(),
    expires_at: null,
  },
  {
    id: "2",
    name: "Development API Key",
    key: "pk_test_987654321",
    created_at: new Date().toISOString(),
    expires_at: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
  },
];

type ApiKey = {
  id: string;
  name: string;
  key: string;
  created_at: string;
  expires_at: string | null;
};

const apiKeySchema = yup.object().shape({
  name: yup.string().required("Name is required"),
});

type ApiKeyFormData = yup.InferType<typeof apiKeySchema>;

export default function ApiKeys() {
  const [isCreating, setIsCreating] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedKey, setSelectedKey] = useState<ApiKey | null>(null);
  const [apiKeys, setApiKeys] = useState<ApiKey[]>(mockApiKeys);
  const [revealedKeys, setRevealedKeys] = useState<Set<string>>(new Set());
  const [revokeExpiration, setRevokeExpiration] = useState<string>("immediately");
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<string | null>(null);

  const form = useForm<ApiKeyFormData>({
    resolver: yupResolver(apiKeySchema) as any,
    defaultValues: {
      name: "",
    },
  });

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

  const handleCreate = async (data: ApiKeyFormData) => {
    setIsLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      const newKey = `pk_${Math.random().toString(36).substr(2, 9)}`;
      const newApiKey: ApiKey = {
        id: Math.random().toString(36).substr(2, 9),
        name: data.name,
        key: newKey,
        created_at: new Date().toISOString(),
        expires_at: null,
      };
      setApiKeys(prev => [...prev, newApiKey]);
      setNewlyCreatedKey(newKey);
      toast.success("API key created successfully");
      form.reset();
      setIsCreating(false);
    } catch (error) {
      toast.error("Failed to create API key");
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdate = async (data: ApiKeyFormData) => {
    if (!selectedKey) return;
    setIsLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      setApiKeys(prev => prev.map(key => 
        key.id === selectedKey.id 
          ? {
              ...key,
              name: data.name,
            }
          : key
      ));
      toast.success("API key updated successfully");
      setIsEditing(false);
      setSelectedKey(null);
      form.reset();
    } catch (error) {
      toast.error("Failed to update API key");
    } finally {
      setIsLoading(false);
    }
  };

  const handleEdit = (key: ApiKey) => {
    setSelectedKey(key);
    form.reset({
      name: key.name,
    });
    setIsEditing(true);
  };

  const handleRevoke = async (id: string) => {
    setIsLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      if (revokeExpiration === "immediately") {
        setApiKeys(prev => prev.filter(key => key.id !== id));
        toast.success("API key revoked immediately");
      } else {
        const hours = revokeExpiration === "24h" ? 24 : 168; // 7 days = 168 hours
        const expiresAt = new Date(Date.now() + hours * 60 * 60 * 1000).toISOString();
        setApiKeys(prev => prev.map(key => 
          key.id === id 
            ? { ...key, expires_at: expiresAt }
            : key
        ));
        toast.success(`API key will expire in ${hours === 24 ? "24 hours" : "7 days"}`);
      }
    } catch (error) {
      toast.error("Failed to revoke API key");
    } finally {
      setIsLoading(false);
    }
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

  const columns: ColumnDef<ApiKey>[] = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "created_at",
      header: "Created At",
      cell: ({ row }) => format(new Date(row.original.created_at), "PPp"),
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
            <Dialog open={isEditing && selectedKey?.id === key.id} onOpenChange={setIsEditing}>
              <DialogTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleEdit(key)}
                >
                  <RiEditLine className="w-4 h-4" />
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Edit API Key</DialogTitle>
                  <DialogDescription>
                    Update your API key settings
                  </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                  <form onSubmit={form.handleSubmit(handleUpdate)} className="space-y-4">
                    <FormField
                      control={form.control}
                      name="name"
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
                        onClick={() => setIsEditing(false)}
                      >
                        Cancel
                      </Button>
                      <Button type="submit">
                        Update Key
                      </Button>
                    </div>
                  </form>
                </Form>
              </DialogContent>
            </Dialog>
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
                      onClick={() => handleRevoke(key.id)}
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

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RiLoader4Line className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-lg font-medium">API Keys</h2>
        <Dialog open={isCreating} onOpenChange={setIsCreating}>
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
                  name="name"
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
                    onClick={() => setIsCreating(false)}
                  >
                    Cancel
                  </Button>
                  <Button type="submit">
                    Create Key
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
            <CopyToClipboard text={newlyCreatedKey}>
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