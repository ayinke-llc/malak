"use client";

import * as React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { type MalakContact } from "@/client/Api";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { LIST_CONTACTS } from "@/lib/query-constants";
import client from "@/lib/client";
import { toast } from "sonner";
import {
  Search,
  Users,
  Mail,
  Building2,
  MapPin,
  MoreHorizontal,
  Plus,
  ChevronLeft,
  ChevronRight,
  ArrowUpDown,
} from "lucide-react";
import Link from "next/link";

export default function ContactsTable() {
  const [page, setPage] = React.useState(1);
  const [perPage, setPerPage] = React.useState(20);
  const [search, setSearch] = React.useState("");
  const [sortBy, setSortBy] = React.useState<"name" | "email" | "company" | "created">("created");
  const [sortOrder, setSortOrder] = React.useState<"asc" | "desc">("desc");
  
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: [LIST_CONTACTS, page, perPage],
    queryFn: async () => {
      const response = await client.contacts.contactsList({
        page,
        per_page: perPage,
      });
      return response.data;
    },
  });

  const deleteContactMutation = useMutation({
    mutationFn: async (reference: string) => {
      await client.contacts.deleteContact(reference);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [LIST_CONTACTS] });
      toast.success("Contact deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete contact");
    },
  });

  const contacts = data?.contacts || [];
  const totalPages = data?.meta.paging.total
    ? Math.ceil(data.meta.paging.total / perPage)
    : 0;

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div className="w-72 h-10 bg-muted animate-pulse rounded-md" />
          <div className="w-40 h-10 bg-muted animate-pulse rounded-md" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i} className="p-6 animate-pulse">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 rounded-full bg-muted" />
                <div className="flex-1 space-y-2">
                  <div className="h-5 bg-muted rounded w-3/4" />
                  <div className="h-4 bg-muted rounded w-1/2" />
                </div>
              </div>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1 flex items-center gap-2">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search contacts..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9"
            />
          </div>
          <Select value={sortBy} onValueChange={(value: any) => setSortBy(value)}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="name">Name</SelectItem>
              <SelectItem value="email">Email</SelectItem>
              <SelectItem value="company">Company</SelectItem>
              <SelectItem value="created">Created</SelectItem>
            </SelectContent>
          </Select>
          <Button
            variant="outline"
            size="icon"
            onClick={() => setSortOrder(prev => prev === "asc" ? "desc" : "asc")}
          >
            <ArrowUpDown className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {contacts.length > 0 ? (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {contacts.map((contact) => {
              const fullName = `${contact.first_name || ''} ${contact.last_name || ''}`.trim();
              const initials = `${(contact.first_name?.[0] || '')}${(contact.last_name?.[0] || '')}`.toUpperCase();

              return (
                <Card key={contact.reference} className="group">
                  <div className="p-6">
                    <div className="flex items-start gap-4">
                      <Avatar className="h-10 w-10 border">
                        <AvatarFallback className={initials ? "bg-primary/10 text-primary" : "bg-muted"}>
                          {initials || "?"}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2">
                          <div>
                            <h3 className="font-medium text-base truncate group-hover:text-primary transition-colors">
                              <Link href={`/contacts/${contact.reference}`}>
                                {fullName || contact.email}
                              </Link>
                            </h3>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground mt-0.5">
                              <Mail className="h-3.5 w-3.5" />
                              <span className="truncate">{contact.email}</span>
                            </div>
                          </div>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
                              >
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem
                                onClick={() => {
                                  navigator.clipboard.writeText(contact.email || "");
                                  toast.success("Email copied to clipboard");
                                }}
                              >
                                Copy email
                              </DropdownMenuItem>
                              <DropdownMenuItem asChild>
                                <Link href={`/contacts/${contact.reference}`}>
                                  View details
                                </Link>
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem
                                className="text-destructive focus:text-destructive"
                                onClick={() => {
                                  if (contact.reference) {
                                    deleteContactMutation.mutate(contact.reference);
                                  }
                                }}
                              >
                                Delete contact
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                        {(contact.company || contact.city) && (
                          <div className="flex flex-wrap gap-3 mt-3">
                            {contact.company && (
                              <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                <Building2 className="h-3.5 w-3.5" />
                                <span>{contact.company}</span>
                              </div>
                            )}
                            {contact.city && (
                              <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                <MapPin className="h-3.5 w-3.5" />
                                <span>{contact.city}</span>
                              </div>
                            )}
                          </div>
                        )}
                        {contact.lists && contact.lists.length > 0 && (
                          <div className="flex flex-wrap gap-1.5 mt-3">
                            {contact.lists.slice(0, 2).map((list) => (
                              <Badge
                                key={list.id}
                                variant="secondary"
                                className="bg-muted/50 text-xs font-normal"
                              >
                                {list.list?.title}
                              </Badge>
                            ))}
                            {contact.lists.length > 2 && (
                              <Badge
                                variant="secondary"
                                className="bg-muted/30 text-xs font-normal"
                              >
                                +{contact.lists.length - 2}
                              </Badge>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </Card>
              );
            })}
          </div>

          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <p className="text-sm text-muted-foreground">
                Showing {(page - 1) * perPage + 1} to{" "}
                {Math.min(page * perPage, data?.meta.paging.total || 0)} of{" "}
                {data?.meta.paging.total || 0} contacts
              </p>
              <Select
                value={String(perPage)}
                onValueChange={(value) => setPerPage(Number(value))}
              >
                <SelectTrigger className="w-[110px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {[10, 20, 50, 100].map((size) => (
                    <SelectItem key={size} value={String(size)}>
                      {size} per page
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page === 1}
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <div className="flex items-center gap-1">
                <Input
                  type="number"
                  min={1}
                  max={totalPages}
                  value={page}
                  onChange={e => {
                    const value = parseInt(e.target.value);
                    if (value >= 1 && value <= totalPages) {
                      setPage(value);
                    }
                  }}
                  className="w-14 h-8 text-center"
                />
                <span className="text-sm text-muted-foreground">
                  of {totalPages}
                </span>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </>
      ) : (
        <Card className="py-16">
          <div className="flex flex-col items-center justify-center text-center">
            <Users className="h-8 w-8 text-muted-foreground/50" />
            <h3 className="mt-4 text-lg font-medium">No contacts found</h3>
            <p className="mt-2 text-sm text-muted-foreground">
              {search
                ? "Try adjusting your search or filters"
                : "Get started by adding your first contact"}
            </p>
            {!search && (
              <Button className="mt-4">
                <Plus className="h-4 w-4 mr-2" />
                Add contact
              </Button>
            )}
          </div>
        </Card>
      )}
    </div>
  );
}
