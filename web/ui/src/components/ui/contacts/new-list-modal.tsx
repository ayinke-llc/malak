import type {
  MalakContactList,
  ServerAPIStatus,
  ServerFetchContactListsResponse,
} from "@/client/Api";
import { Button } from "@/components/Button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/Dialog";
import { Input } from "@/components/Input";
import client from "@/lib/client";
import {
  CREATE_CONTACT_LIST,
  LIST_CONTACT_LISTS,
  UPDATE_CONTACT_LIST,
} from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RiAddLine,
  RiCheckLine,
  RiCloseLargeLine,
  RiDeleteBinLine,
  RiEyeLine,
  RiPencilLine,
} from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { AxiosError, AxiosResponse } from "axios";
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
  const { data, error } = useQuery({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: () => client.contacts.fetchContactLists(),
  });

  if (error) {
    toast.error("an error occurred while fetching your contact lists");
  }

  const [loading, setLoading] = useState<boolean>(false);
  const [hasOpenDialog, setHasOpenDialog] = useState(false);
  const [referenceToDelete, setReferenceToDelete] = useState("");
  const [editingID, setEditingId] = useState<string | null>(null);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const queryClient = useQueryClient();

  const onNewListAdded = (list: MalakContactList | string) => {
    queryClient.setQueryData(
      [LIST_CONTACT_LISTS],
      (old: AxiosResponse<ServerFetchContactListsResponse | undefined>) => {
        if (!old) {
          return { data: { lists: [list] } };
        }

        let listExists = false;
        let updatedLists: MalakContactList[] | undefined;

        if (typeof list === 'string') {
          updatedLists = old?.data?.lists?.filter((existingList) => existingList.reference !== list);
        } else {
          updatedLists = old?.data?.lists?.map((existingList) => {
            if (list == null) {
              return existingList
            }
            if (existingList.id === list.id) {
              listExists = true;
              return list;
            }
            return existingList;
          });

          if (!listExists) {
            updatedLists?.unshift(list);
          }
        }

        return {
          ...old,
          data: {
            ...old.data,
            lists: updatedLists,
          },
        };
      },
    );
  };

  const isItemBeingEdited = (item: MalakContactList): boolean =>
    item.id == editingID;

  const handleEdit = (id: string) => {
    setEditingId(id);
  };

  const handleSave = () => {
    setEditingId(null);
  };

  const handleDelete = (reference: string) => {
    setReferenceToDelete(reference);
    setIsDeleteModalOpen(true);
  };

  const confirmDelete = () => {
    mutation.mutate(referenceToDelete);
  };

  const handleDialogItemOpenChange = (open: boolean) => {
    setHasOpenDialog(open);
  };

  const { reset } = useForm({
    resolver: yupResolver(schema),
    defaultValues: {
      title: "",
    },
  });

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTACT_LIST],
    onMutate: () => setLoading(true),
    mutationFn: (reference: string) =>
      client.contacts.deleteContactList(reference),
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setIsDeleteModalOpen(false);
      onNewListAdded(referenceToDelete)
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
            <CreateNewContactList onCreate={onNewListAdded} />
            {data?.data?.lists?.map((item) => (
              <div>
                <div
                  key={item.id}
                  className="flex items-center space-x-2 p-3 rounded-lg border-2 border-gray-200"
                >
                  {isItemBeingEdited(item) ? (
                    <EditList
                      list={item}
                      onEdited={(item) => {
                        handleSave();
                        onNewListAdded(item);
                      }}
                    />
                  ) : (
                    <>
                      <span className="flex-grow">{item.title}</span>
                      <div className="flex space-x-2">
                        <Button
                          type="button"
                          onClick={() => handleEdit(item?.id as string)}
                          size="icon"
                          variant="ghost"
                        >
                          <RiPencilLine className="h-4 w-4" />
                        </Button>
                        <Button
                          type="button"
                          onClick={() =>
                            handleDelete(item?.reference as string)
                          }
                          size="icon"
                          variant="ghost"
                        >
                          <RiDeleteBinLine className="h-4 w-4" color="red" />
                        </Button>
                      </div>
                    </>
                  )}
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
            <Button
              variant="secondary"
              onClick={() => setIsDeleteModalOpen(false)}
              isLoading={loading}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={confirmDelete}
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

const CreateNewContactList = ({
  onCreate,
}: { onCreate: (value: MalakContactList) => void }) => {
  const [isAddingItem, setIsAddingItem] = useState(false);

  const [loading, setLoading] = useState<boolean>(false);

  const mutation = useMutation({
    mutationKey: [CREATE_CONTACT_LIST],
    onMutate: () => setLoading(true),
    mutationFn: (data: CreateContactListInput) =>
      client.contacts.createContactList({ name: data.title }),
    onSuccess: ({ data }) => {
      onCreate(data?.list as MalakContactList);
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
    mutation.mutate(data);
  };

  return (
    <div className="w-full">
      {isAddingItem ? (
        <div>
          <div className="flex items-center space-x-2 w-full">
            <form
              onSubmit={handleSubmit(onSubmit)}
              className="flex items-center space-x-2 w-full"
            >
              <Input
                placeholder="Name of your list"
                className="flex-grow"
                {...register("title")}
              />
              <Button
                type="submit"
                size="icon"
                variant="ghost"
                isLoading={loading}
              >
                <RiCheckLine className="h-4 w-4" color="green" />
                <span className="sr-only">Add item</span>
              </Button>
              <Button
                onClick={() => setIsAddingItem(false)}
                size="icon"
                variant="ghost"
                isLoading={loading}
              >
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
  );
};

const EditList = ({
  onEdited,
  list,
}: { onEdited: (item: MalakContactList) => void; list: MalakContactList }) => {
  const [loading, setLoading] = useState<boolean>(false);

  const {
    register,
    formState: { errors, isDirty },
    handleSubmit,
    reset,
  } = useForm({
    resolver: yupResolver(schema),
    mode: "onChange",
    defaultValues: {
      title: list?.title as string,
    },
  });

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTACT_LIST],
    onMutate: () => setLoading(true),
    mutationFn: (data: CreateContactListInput) =>
      client.contacts.editContactList(list.reference as string, {
        name: data.title,
      }),
    onSuccess: ({ data }) => {
      onEdited(data?.list as MalakContactList);
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

  const onSubmit: SubmitHandler<CreateContactListInput> = (data) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="w-full">
      <div className="flex items-center space-x-2 w-full">
        <Input className="flex-grow" {...register("title")} />
        <Button
          type="submit"
          size="icon"
          variant="ghost"
          isLoading={loading}
          disabled={!isDirty}
        >
          <RiCheckLine className="h-4 w-4" color="green" />
          <span className="sr-only">Save</span>
        </Button>
        <Button type="button" size="icon" variant="ghost">
          <RiDeleteBinLine className="h-4 w-4" color="red" />
        </Button>
      </div>

      {errors.title && (
        <p className="mt-2 text-xs text-red-600 dark:text-red-500">
          <span className="font-medium">{errors.title.message}</span>
        </p>
      )}
    </form>
  );
};
