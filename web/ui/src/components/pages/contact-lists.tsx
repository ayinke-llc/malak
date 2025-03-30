"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import client from "@/lib/client";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RiAddLine,
  RiCheckLine,
  RiPencilLine,
  RiDeleteBinLine,
  RiArrowLeftLine,
  RiErrorWarningLine,
  RiCloseLine,
  RiTeamLine,
} from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";
import type {
  MalakContactList,
  MalakContactListMappingWithContact,
  ServerAPIStatus,
  ServerFetchContactListsResponse,
} from "@/client/Api";
import Link from "next/link";

const LIST_CONTACT_LISTS = "list-contact-lists";
const CREATE_CONTACT_LIST = "create-contact-list";
const UPDATE_CONTACT_LIST = "update-contact-list";

type CreateContactListInput = {
  name: string;
};

const schema = yup
  .object({
    name: yup.string().required().min(3).max(50),
  })
  .required();

interface EditListProps {
  list: MalakContactList;
  onEdited: () => void;
}

interface ContactListItem {
  list: MalakContactList;
  mappings: MalakContactListMappingWithContact[];
}

const EditList = ({ list, onEdited }: EditListProps) => {
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
  } = useForm<CreateContactListInput>({
    resolver: yupResolver(schema),
    defaultValues: {
      name: list.title ?? "",
    },
  });

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTACT_LIST],
    mutationFn: async (data: CreateContactListInput) => {
      if (!list.reference) throw new Error("List reference is required");
      const response = await client.contacts.editContactList(list.reference, data);
      return response.data;
    },
    onSuccess: () => {
      onEdited();
      queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
      toast.success("List updated successfully");
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred");
    },
    onSettled() {
      setLoading(false);
    },
  });

  const onSubmit: SubmitHandler<CreateContactListInput> = (data) => {
    setLoading(true);
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="w-full">
      <div className="flex items-center gap-2">
        <Input
          className="flex-grow bg-transparent border-gray-200 text-gray-900 focus-visible:ring-0 focus-visible:border-gray-300"
          {...register("name")}
        />
        <Button
          type="submit"
          size="sm"
          variant="ghost"
          disabled={!isDirty || loading}
          className="h-8 w-8 p-0 text-gray-500 hover:text-emerald-400 hover:bg-emerald-500/10 disabled:opacity-50"
        >
          <RiCheckLine className="h-4 w-4" />
          <span className="sr-only">Save changes</span>
        </Button>
      </div>
      {errors.name && (
        <p className="mt-1 text-xs text-red-500">{errors.name.message}</p>
      )}
    </form>
  );
};

interface CreateNewContactListProps {
  onCreate: () => void;
}

const CreateNewContactList = ({ onCreate }: CreateNewContactListProps) => {
  const [loading, setLoading] = useState(false);
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateContactListInput>({
    resolver: yupResolver(schema),
  });

  const mutation = useMutation({
    mutationKey: [CREATE_CONTACT_LIST],
    mutationFn: async (data: CreateContactListInput) => {
      const response = await client.contacts.createContactList(data);
      return response.data;
    },
    onSuccess: () => {
      onCreate();
      queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
      reset();
      toast.success("List created successfully");
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred");
    },
    onSettled() {
      setLoading(false);
    },
  });

  const onSubmit: SubmitHandler<CreateContactListInput> = (data) => {
    setLoading(true);
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="List name"
          className="flex-grow bg-transparent border-gray-200 text-gray-900 focus-visible:ring-0 focus-visible:border-gray-300"
          {...register("name")}
        />
        <Button
          type="submit"
          size="sm"
          variant="ghost"
          disabled={loading}
          className="h-8 w-8 p-0 text-gray-500 hover:text-emerald-400 hover:bg-emerald-500/10"
        >
          <RiCheckLine className="h-4 w-4" />
          <span className="sr-only">Create list</span>
        </Button>
      </div>
      {errors.name && (
        <p className="text-xs text-red-500">{errors.name.message}</p>
      )}
    </form>
  );
};

export default function ContactListsPage() {
  const [editingId, setEditingId] = useState<string>();
  const [showNewListForm, setShowNewListForm] = useState(false);
  const queryClient = useQueryClient();

  const { data, error } = useQuery<ServerFetchContactListsResponse>({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: async () => {
      const response = await client.contacts.fetchContactLists();
      return response.data;
    },
  });

  const deleteMutation = useMutation({
    mutationKey: [UPDATE_CONTACT_LIST],
    mutationFn: async (reference: string) => {
      const response = await client.contacts.deleteContactList(reference);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
      toast.success("List deleted successfully");
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred");
    },
  });

  const handleEdit = (id: string) => {
    if (id) {
      setEditingId(id);
    }
  };

  const handleDelete = (reference: string) => {
    if (reference) {
      deleteMutation.mutate(reference);
    }
  };

  return (
    <div className="pt-6 bg-background">
      <section>
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-4">
            <Link href="/contacts">
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <RiArrowLeftLine className="h-4 w-4" />
              </Button>
            </Link>
            <div>
              <h1 className="text-2xl font-semibold">Contact Lists</h1>
              <p className="text-sm text-muted-foreground">
                Create and manage your contact lists
              </p>
            </div>
          </div>
          <Button
            variant="default"
            size="sm"
            className="gap-2"
            onClick={() => setShowNewListForm(true)}
          >
            <RiAddLine className="h-4 w-4" />
            New List
          </Button>
        </div>
      </section>

      <section className="mt-6">
        {error ? (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4">
            <div className="flex items-center gap-2 text-red-600">
              <RiErrorWarningLine className="h-5 w-5" />
              <p className="text-sm font-medium">Failed to load contact lists</p>
            </div>
          </div>
        ) : !data ? (
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="rounded-lg border p-6 animate-pulse">
                <div className="flex items-start justify-between">
                  <div className="space-y-3 flex-1">
                    <div className="h-5 w-2/3 bg-muted rounded" />
                    <div className="h-4 w-1/3 bg-muted rounded" />
                  </div>
                  <div className="h-8 w-8 bg-muted rounded" />
                </div>
              </div>
            ))}
          </div>
        ) : (
          <>
            {showNewListForm && (
              <div className="mb-6 rounded-lg border bg-card p-4">
                <div className="mb-3 flex items-center justify-between">
                  <h3 className="text-lg font-medium">Create New List</h3>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0"
                    onClick={() => setShowNewListForm(false)}
                  >
                    <RiCloseLine className="h-4 w-4" />
                  </Button>
                </div>
                <CreateNewContactList
                  onCreate={() => {
                    setShowNewListForm(false);
                    queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
                  }}
                />
              </div>
            )}

            {data.lists && data.lists.length > 0 ? (
              <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                {data.lists.map((item: ContactListItem) => (
                  <div
                    key={item.list.id}
                    className="group relative rounded-lg border bg-card p-6 transition-all hover:shadow-md"
                  >
                    {item.list.id === editingId ? (
                      <EditList
                        list={item.list}
                        onEdited={() => {
                          setEditingId(undefined);
                          queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
                        }}
                      />
                    ) : (
                      <>
                        <div className="flex items-start justify-between">
                          <div className="space-y-1">
                            <h3 className="font-medium line-clamp-1">{item.list.title}</h3>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                              <div className="flex items-center gap-1">
                                <RiTeamLine className="h-4 w-4" />
                                <span>{item.mappings.length} contacts</span>
                              </div>
                            </div>
                          </div>
                          <div className="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100">
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-8 w-8 p-0"
                              onClick={() => item.list.id && handleEdit(item.list.id)}
                            >
                              <RiPencilLine className="h-4 w-4" />
                              <span className="sr-only">Edit list</span>
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-8 w-8 p-0 hover:text-red-500 hover:bg-red-50"
                              onClick={() => item.list.reference && handleDelete(item.list.reference)}
                            >
                              <RiDeleteBinLine className="h-4 w-4" />
                              <span className="sr-only">Delete list</span>
                            </Button>
                          </div>
                        </div>
                        <div className="mt-4">
                          <p className="text-xs text-muted-foreground font-mono">
                            {item.list.reference}
                          </p>
                        </div>
                      </>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <div className="rounded-lg border border-dashed p-8 text-center">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  <RiTeamLine className="h-6 w-6 text-primary" />
                </div>
                <h3 className="mt-4 text-lg font-medium">No lists yet</h3>
                <p className="mt-1 text-sm text-muted-foreground">
                  Create your first contact list to start organizing your contacts
                </p>
                <Button
                  variant="default"
                  size="sm"
                  className="mt-4 gap-2"
                  onClick={() => setShowNewListForm(true)}
                >
                  <RiAddLine className="h-4 w-4" />
                  Create your first list
                </Button>
              </div>
            )}
          </>
        )}
      </section>
    </div>
  );
} 