"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
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

type ContactListFormData = {
  name: string;
};

const contactListSchema = yup.object().shape({
  name: yup
    .string()
    .required("List name is required")
    .min(3, "List name must be at least 3 characters")
    .max(50, "List name must be at most 50 characters")
    .matches(/^[a-zA-Z0-9\s-_]+$/, "List name can only contain letters, numbers, spaces, hyphens and underscores"),
});

interface ContactListItem {
  list: MalakContactList;
  mappings: MalakContactListMappingWithContact[];
}

interface EditContactListDialogProps {
  list: MalakContactList;
  mappingsCount: number;
  open: boolean;
  onClose: () => void;
}

const EditContactListDialog = ({ list, mappingsCount, open, onClose }: EditContactListDialogProps) => {
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty, isSubmitting },
    reset,
  } = useForm<ContactListFormData>({
    resolver: yupResolver(contactListSchema),
    defaultValues: {
      name: list.title ?? "",
    },
  });

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTACT_LIST],
    mutationFn: async (data: ContactListFormData) => {
      if (!list.reference) throw new Error("List reference is required");
      const response = await client.contacts.editContactList(list.reference, data);
      return response.data;
    },
    onSuccess: () => {
      onClose();
      queryClient.invalidateQueries({ queryKey: [LIST_CONTACT_LISTS] });
      toast.success("List updated successfully");
      reset();
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred");
    },
    onSettled() {
      setLoading(false);
    },
  });

  const onSubmit: SubmitHandler<ContactListFormData> = (data) => {
    setLoading(true);
    mutation.mutate(data);
  };

  return (
    <Dialog open={open} onOpenChange={(open) => {
      if (!open) {
        reset();
        onClose();
      }
    }}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Contact List</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-4">
            <div>
              <Input
                className="w-full"
                {...register("name")}
                placeholder="List name"
                aria-invalid={errors.name ? "true" : "false"}
              />
              {errors.name && (
                <p className="mt-1 text-xs text-destructive">{errors.name.message}</p>
              )}
            </div>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <RiTeamLine className="h-4 w-4" />
              <span>{mappingsCount} contacts</span>
            </div>
            <p className="text-xs text-muted-foreground font-mono">
              {list.reference}
            </p>
          </div>
          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="ghost"
              onClick={() => {
                reset();
                onClose();
              }}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={!isDirty || isSubmitting}
            >
              {isSubmitting ? "Saving..." : "Save changes"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};

interface DeleteContactListDialogProps {
  list: MalakContactList;
  open: boolean;
  onClose: () => void;
  onConfirm: (reference: string) => void;
}

const DeleteContactListDialog = ({ list, open, onClose, onConfirm }: DeleteContactListDialogProps) => {
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = () => {
    if (!list.reference) return;
    setIsDeleting(true);
    onConfirm(list.reference);
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Contact List</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete the list "{list.title}"? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant="ghost"
            onClick={onClose}
            disabled={isDeleting}
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={isDeleting || !list.reference}
          >
            {isDeleting ? "Deleting..." : "Delete List"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

interface ContactListCardProps {
  item: ContactListItem;
  onEdit: (id: string) => void;
  onDelete: (list: MalakContactList) => void;
}

const ContactListCard = ({ item, onEdit, onDelete }: ContactListCardProps) => {
  return (
    <div className="group relative rounded-lg border bg-card p-4 transition-all hover:shadow-sm">
      <div className="flex items-start justify-between">
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <h3 className="font-medium text-sm">{item.list.title}</h3>
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <RiTeamLine className="h-3.5 w-3.5" />
              <span>{item.mappings?.length || 0} contacts</span>
            </div>
          </div>
          <p className="text-[10px] text-muted-foreground font-mono">
            {item.list.reference}
          </p>
        </div>
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="sm"
            className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={() => item.list.id && onEdit(item.list.id)}
          >
            <RiPencilLine className="h-3.5 w-3.5" />
            <span className="sr-only">Edit list</span>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity hover:text-red-500 hover:bg-red-50"
            onClick={() => onDelete(item.list)}
          >
            <RiDeleteBinLine className="h-3.5 w-3.5" />
            <span className="sr-only">Delete list</span>
          </Button>
        </div>
      </div>
    </div>
  );
};

interface CreateNewContactListProps {
  onCreate: () => void;
  onCancel: () => void;
}

const CreateNewContactList = ({ onCreate, onCancel }: CreateNewContactListProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<ContactListFormData>({
    resolver: yupResolver(contactListSchema),
  });

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationKey: [CREATE_CONTACT_LIST],
    mutationFn: async (data: ContactListFormData) => {
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
  });

  const onSubmit: SubmitHandler<ContactListFormData> = (data) => {
    mutation.mutate(data);
  };

  return (
    <Dialog open onOpenChange={onCancel}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New List</DialogTitle>
          <DialogDescription>
            Create a new contact list to organize your contacts.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <Input
              placeholder="List name"
              className="w-full"
              {...register("name")}
              aria-invalid={errors.name ? "true" : "false"}
            />
            {errors.name && (
              <p className="mt-1 text-xs text-destructive">{errors.name.message}</p>
            )}
          </div>
          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="ghost"
              onClick={onCancel}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Creating..." : "Create List"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default function ContactListsPage() {
  const [editingList, setEditingList] = useState<ContactListItem | null>(null);
  const [deletingList, setDeletingList] = useState<MalakContactList | null>(null);
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
      setDeletingList(null);
      queryClient.invalidateQueries({ 
        queryKey: [LIST_CONTACT_LISTS],
        exact: true,
      });
      toast.success("List deleted successfully");
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "Failed to delete list");
    },
  });

  const handleEdit = (id: string) => {
    const list = data?.lists.find(item => item.list.id === id);
    if (list) {
      setEditingList(list);
    }
  };

  const handleDelete = (list: MalakContactList) => {
    if (!list.reference) {
      toast.error("Cannot delete list: missing reference");
      return;
    }
    setDeletingList(list);
  };

  const handleConfirmDelete = (reference: string) => {
    deleteMutation.mutate(reference);
  };

  return (
    <div className="space-y-8 p-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/contacts">
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <RiArrowLeftLine className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-xl font-semibold">Contact Lists</h1>
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

      {error ? (
        <div className="rounded-lg border border-destructive/20 bg-destructive/10 p-4">
          <div className="flex items-center gap-2 text-destructive">
            <RiErrorWarningLine className="h-5 w-5" />
            <p className="text-sm font-medium">Failed to load contact lists</p>
          </div>
        </div>
      ) : !data ? (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="rounded-lg border p-4 animate-pulse">
              <div className="flex items-start justify-between">
                <div className="space-y-2 flex-1">
                  <div className="h-4 w-1/2 bg-muted rounded" />
                  <div className="h-3 w-1/4 bg-muted rounded" />
                </div>
                <div className="flex gap-1">
                  <div className="h-7 w-7 bg-muted rounded" />
                  <div className="h-7 w-7 bg-muted rounded" />
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {data.lists && data.lists.length > 0 ? (
            data.lists.map((item: ContactListItem) => (
              <ContactListCard
                key={item.list.id}
                item={item}
                onEdit={handleEdit}
                onDelete={(list) => handleDelete(list)}
              />
            ))
          ) : (
            <div className="col-span-full">
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
            </div>
          )}
        </div>
      )}

      {showNewListForm && (
        <CreateNewContactList
          onCreate={() => setShowNewListForm(false)}
          onCancel={() => setShowNewListForm(false)}
        />
      )}

      {editingList && (
        <EditContactListDialog
          list={editingList.list}
          mappingsCount={editingList.mappings?.length || 0}
          open={!!editingList}
          onClose={() => setEditingList(null)}
        />
      )}

      {deletingList && (
        <DeleteContactListDialog
          list={deletingList}
          open={!!deletingList}
          onClose={() => setDeletingList(null)}
          onConfirm={handleConfirmDelete}
        />
      )}
    </div>
  );
} 