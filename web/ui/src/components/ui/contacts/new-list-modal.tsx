"use client"

import type {
  MalakContactList,
  MalakContactListMappingWithContact,
  ServerAPIStatus, ServerFetchContactListsResponse
} from "@/client/Api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import client from "@/lib/client";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RiAddLine,
  RiCheckLine, RiPencilLine,
  RiDeleteBinLine, RiSettings4Line
} from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";

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

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTACT_LIST],
    mutationFn: async (data: CreateContactListInput) => {
      if (!list.reference) throw new Error("List reference is required");
      const response = await client.contacts.editContactList(list.reference, data);
      return response.data;
    },
    onSuccess: () => {
      onEdited();
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
          className="flex-grow bg-transparent border-gray-800/50 text-gray-300 focus-visible:ring-0 focus-visible:border-gray-700"
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
          className="flex-grow bg-transparent focus-visible:ring-0 focus-visible:border-gray-700"
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

export default function ManageListModal() {
  const [hasOpenDialog, setHasOpenDialog] = useState(false);
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

  const handleDialogItemOpenChange = (open: boolean) => {
    setHasOpenDialog(open);
    if (!open) {
      setEditingId(undefined);
      setShowNewListForm(false);
    }
  };

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
    <Dialog open={hasOpenDialog} onOpenChange={handleDialogItemOpenChange}>
      <DialogTrigger asChild>
        <div className="w-full text-right mb-4">
          <Button
            type="button"
            variant="default"
            className="whitespace-nowrap gap-1"
          >
            <RiSettings4Line />
            Manage lists
          </Button>
        </div>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Contact lists</DialogTitle>
          <DialogDescription className="text-gray-500">
            Create and manage your contact lists.
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          {error ? (
            <div className="text-sm text-red-500">
              Failed to load contact lists
            </div>
          ) : !data ? (
            <div className="text-sm text-gray-500">Loading lists...</div>
          ) : (
            <div className="space-y-4">
              {data.lists && data.lists.length > 0 ? (
                <div className="space-y-4">
                  {data.lists.map((item: ContactListItem) => (
                    <div
                      key={item.list.id}
                      className="flex items-center justify-between gap-2 rounded-md border border-gray-800/50 p-3"
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
                          <div>
                            <h4 className="font-medium">
                              {item.list.title}
                            </h4>
                            <p className="text-sm">
                              {item.list.reference}
                            </p>
                          </div>
                          <div className="flex items-center gap-2">
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
                              className="h-8 w-8 p-0 hover:text-red-400 hover:bg-red-500/10"
                              onClick={() => item.list.reference && handleDelete(item.list.reference)}
                            >
                              <RiDeleteBinLine className="h-4 w-4" />
                              <span className="sr-only">Delete list</span>
                            </Button>
                          </div>
                        </>
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-sm">No lists created yet</div>
              )}

              {showNewListForm ? (
                <CreateNewContactList
                  onCreate={() => {
                    setShowNewListForm(false);
                    queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
                  }}
                />
              ) : (
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full border-gray-800/50 bg-transparent hover:bg-800/30 hover:text-gray-500"
                  onClick={() => setShowNewListForm(true)}
                >
                  <RiAddLine className="h-4 w-4 mr-2" />
                  Create new list
                </Button>
              )}
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
