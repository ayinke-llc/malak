import type { ServerAPIStatus } from "@/client/Api";
import { Button } from "@/components/Button";
import {
  Dialog, DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/Dialog";
import { Input } from "@/components/Input";
import client from "@/lib/client";
import { CREATE_CONTACT_LIST, CREATE_CONTACT_MUTATION, LIST_CONTACT_LISTS } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiAddLine, RiCheckboxLine, RiCheckLine, RiCloseLargeLine, RiDeleteBinLine, RiEyeLine, RiPencilLine, RiPencilRuler2Line, RiPencilRulerLine, RiTwitterXLine } from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";

type CreateContactListInput = {
  title: string;
};

const schema = yup
  .object({
    title: yup.string().required().min(3).max(50),
  })
  .required();

export default function ManageListModal() {

  const [loading, setLoading] = useState<boolean>(false);
  const [hasOpenDialog, setHasOpenDialog] = useState(false);

  const [editingId, setEditingId] = useState<string | null>(null)
  const [editText, setEditText] = useState('')
  const [deleteId, setDeleteId] = useState<number | null>(null)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const [newItemText, setNewItemText] = useState('')

  const { data, error, isLoading } = useQuery({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: () => client.contacts.fetchContactLists(),
  });

  if (error) {
    toast.error("an error occurred while fetching this update");
  }

  const handleEdit = (id: string, text: string) => {
    setEditingId(id)
    setEditText(text)
  }

  const handleSave = () => {
    setEditingId(null)
  }

  const handleDelete = (id: string) => {
    setIsDeleteModalOpen(true)
  }

  const confirmDelete = () => {
    setIsDeleteModalOpen(false)
  }

  const handleDialogItemOpenChange = (open: boolean) => {
    setHasOpenDialog(open);
  };

  const mutation = useMutation({
    mutationKey: [CREATE_CONTACT_MUTATION],
    mutationFn: (data: CreateContactListInput) =>
      client.contacts.createContactList({ name: data.title }),
    onSuccess: ({ data }) => {
      toast.success(data.message);
      handleDialogItemOpenChange(false);
      reset();
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onSettled: () => setLoading(false),
  });

  const {
    register,
    formState: { errors },
    handleSubmit,
    reset,
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: {
      title: "",
    },
  });

  const onSubmit: SubmitHandler<CreateContactListInput> = (data) => {
    setLoading(true);
    mutation.mutate(data);
  };


  return (
    <>
      <Dialog onOpenChange={handleDialogItemOpenChange} open={hasOpenDialog}>
        <DialogTrigger asChild>
          <div className="w-full text-right">
            <Button
              type="button"
              variant="secondary"
              className="whitespace-nowrap gap-1"
            >
              <RiEyeLine />
              Manage lists
            </Button>
          </div>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Contacts List</DialogTitle>
            <DialogDescription>
              View, edit, delete or manage your lists
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <CreateNewContactList />
            {data?.data?.lists?.map((item) => (
              <div
                key={item.id}
                className="flex items-center space-x-2 p-3 rounded-lg border-2 border-gray-200"
              >
                {editingId === item.id ? (
                  <Input
                    value={editText}
                    onChange={(e) => setEditText(e.target.value)}
                    className="flex-grow"
                  />
                ) : (
                  <span className="flex-grow">{item.title}</span>
                )}
                <div className="flex space-x-2">
                  {editingId === item.id ? (
                    <Button onClick={handleSave} size="sm">Save</Button>
                  ) : (
                    <Button
                      onClick={() => handleEdit(item?.id as string, item?.title as string)}
                      size="icon"
                      variant="ghost"
                    >
                      <RiPencilLine className="h-4 w-4" />
                    </Button>
                  )}
                  <Button
                    onClick={() => handleDelete(item?.id as string)}
                    size="icon"
                    variant="ghost"
                  >
                    <RiDeleteBinLine className="h-4 w-4" color="red" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this item?
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="secondary" onClick={() => setIsDeleteModalOpen(false)}
              isLoading={loading}
            >
              Cancel
            </Button>
            <Button variant="destructive" onClick={confirmDelete}
              isLoading={loading}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

const CreateNewContactList = () => {

  const [isAddingItem, setIsAddingItem] = useState(false)

  const [loading, setLoading] = useState<boolean>(false);

  const mutation = useMutation({
    mutationKey: [CREATE_CONTACT_LIST],
    mutationFn: (data: CreateContactListInput) =>
      client.contacts.createContactList({ name: data.title }),
    onSuccess: ({ data }) => {
      toast.success(data.message);
      reset();
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onSettled: () => setLoading(false),
  });

  const {
    register,
    formState: { errors },
    handleSubmit,
    reset,
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: {
      title: "",
    },
  });

  const onSubmit: SubmitHandler<CreateContactListInput> = (data) => {
    setLoading(true);
    mutation.mutate(data);
  };

  return (
    <div className="w-full">
      {isAddingItem ? (
        <div>
          <div className="flex items-center space-x-2 w-full">
            <form onSubmit={handleSubmit(onSubmit)} className="flex items-center space-x-2 w-full">
              <Input
                placeholder="Name of your list"
                className="flex-grow"
                {...register("title")}
              />
              <Button type="submit" size="icon" variant="ghost" isLoading={loading}>
                <RiCheckLine className="h-4 w-4" color="green" />
                <span className="sr-only">Add item</span>
              </Button>
              <Button
                onClick={() => setIsAddingItem(false)}
                size="icon"
                variant="ghost"
                isLoading={loading}>
                <RiCloseLargeLine className="h-4 w-4" color="red" />
                <span className="sr-only">Cancel</span>
              </Button>
            </form>
          </div>
          <div className="flex items-center space-x-2 w-full">
            {errors.title && (
              <p className="mt-4 text-xs text-red-600 dark:text-red-500">
                <span className="font-medium">{errors.title.message}</span>
              </p>
            )}
          </div>
        </div>
      ) : (
        <Button
          onClick={() => setIsAddingItem(true)}
          variant="secondary"
          className="w-full"
          isLoading={loading}
        >
          <RiAddLine className="h-4 w-4 mr-2" />
          Add New List
        </Button>
      )}
    </div>
  )
}
