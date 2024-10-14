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
import { CREATE_CONTACT_MUTATION } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiDeleteBinLine, RiEyeLine, RiPencilLine } from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";

type ListItem = {
  id: number
  text: string
}

type CreateContactListInput = {
  title: string;
};

const schema = yup
  .object({
    title: yup.string().required().min(3).max(50),
  })
  .required();

export default function CreateNewListModal() {

  const [loading, setLoading] = useState<boolean>(false);
  const [hasOpenDialog, setHasOpenDialog] = useState(false);

  const [items, setItems] = useState<ListItem[]>([
    { id: 1, text: 'Item 1' },
    { id: 2, text: 'Item 2' },
    { id: 3, text: 'Item 3' },
    { id: 4, text: 'Item 4' },
    { id: 5, text: 'Item 5' },
  ])
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editText, setEditText] = useState('')
  const [deleteId, setDeleteId] = useState<number | null>(null)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)

  const handleEdit = (id: number, text: string) => {
    setEditingId(id)
    setEditText(text)
  }

  const handleSave = () => {
    setItems(items.map(item =>
      item.id === editingId ? { ...item, text: editText } : item
    ))
    setEditingId(null)
  }

  const handleDelete = (id: number) => {
    setDeleteId(id)
    setIsDeleteModalOpen(true)
  }

  const confirmDelete = () => {
    setItems(items.filter(item => item.id !== deleteId))
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
            <DialogTitle>Item List</DialogTitle>
            <DialogDescription>
              View, edit, or delete items in the list.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {items.map((item) => (
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
                  <span className="flex-grow">{item.text}</span>
                )}
                <div className="flex space-x-2">
                  {editingId === item.id ? (
                    <Button onClick={handleSave} size="sm">Save</Button>
                  ) : (
                    <Button
                      onClick={() => handleEdit(item.id, item.text)}
                      size="icon"
                      variant="ghost"
                    >
                      <RiPencilLine className="h-4 w-4" />
                    </Button>
                  )}
                  <Button
                    onClick={() => handleDelete(item.id)}
                    size="icon"
                    variant="ghost"
                  >
                    <RiDeleteBinLine className="h-4 w-4" />
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
