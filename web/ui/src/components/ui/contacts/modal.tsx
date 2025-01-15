import type { ServerAPIStatus, ServerCreateContactRequest, ServerFetchContactResponse } from "@/client/Api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import client from "@/lib/client";
import { CREATE_CONTACT_MUTATION } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiAddLine } from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { AxiosError, AxiosResponse } from "axios";
import { useEffect, useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";

type ContactFormInput = {
  email: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  company?: string;
};

interface CreateContactModalProps {
  mode?: 'create' | 'edit';
  reference?: string;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  triggerButton?: boolean;
}

const schema = yup
  .object({
    email: yup.string().required().email(),
    first_name: yup.string(),
    last_name: yup.string(),
    phone: yup.string(),
    company: yup.string(),
  })
  .required();

export default function CreateContactModal({ 
  mode = 'create', 
  reference,
  open: controlledOpen,
  onOpenChange: controlledOnOpenChange,
  triggerButton = true
}: CreateContactModalProps) {
  const [loading, setLoading] = useState<boolean>(false);
  const [hasOpenDialog, setHasOpenDialog] = useState(false);

  const isControlled = controlledOpen !== undefined;
  const isOpen = isControlled ? controlledOpen : hasOpenDialog;
  const onOpenChange = isControlled ? controlledOnOpenChange : setHasOpenDialog;

  const handleDialogItemOpenChange = (open: boolean) => {
    if (onOpenChange) {
      onOpenChange(open);
    }
  };

  const {
    register,
    formState: { errors },
    handleSubmit,
    reset,
    setValue,
  } = useForm<ContactFormInput>({
    resolver: yupResolver(schema),
    defaultValues: {
      email: "",
      first_name: "",
      last_name: "",
      phone: "",
      company: "",
    },
  });

  // Fetch contact data if in edit mode
  const { data: contactData } = useQuery({
    queryKey: ['contact', reference],
    queryFn: async () => {
      if (!reference) return null;
      const response = await fetch(`/api/contacts/${reference}`);
      if (!response.ok) throw new Error('Failed to fetch contact');
      return response.json() as Promise<ServerFetchContactResponse>;
    },
    enabled: mode === 'edit' && !!reference,
  });

  useEffect(() => {
    if (mode === 'edit' && contactData?.contact) {
      const contact = contactData.contact;
      setValue('email', contact.email || '');
      setValue('first_name', contact.first_name || '');
      setValue('last_name', contact.last_name || '');
      setValue('phone', contact.phone || '');
      setValue('company', contact.company || '');
    }
  }, [mode, contactData, setValue]);

  const createMutation = useMutation<
    AxiosResponse<ServerFetchContactResponse>,
    AxiosError<ServerAPIStatus>,
    ServerCreateContactRequest
  >({
    mutationKey: [CREATE_CONTACT_MUTATION],
    mutationFn: (data) => client.contacts.contactsCreate(data),
    onSuccess: ({ data }) => {
      toast.success(data.message);
      handleDialogItemOpenChange(false);
      reset();
    },
    onError(err) {
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

  const updateMutation = useMutation<
    AxiosResponse<ServerFetchContactResponse>,
    AxiosError<ServerAPIStatus>,
    ServerCreateContactRequest
  >({
    mutationKey: ['UPDATE_CONTACT'],
    mutationFn: async (data) => {
      if (!reference) throw new Error('No reference');
      const response = await fetch(`/api/contacts/${reference}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!response.ok) throw new Error('Failed to update contact');
      return response.json();
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      handleDialogItemOpenChange(false);
    },
    onError(err) {
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

  const onSubmit: SubmitHandler<ContactFormInput> = (data) => {
    setLoading(true);
    if (mode === 'edit') {
      updateMutation.mutate(data);
    } else {
      createMutation.mutate(data);
    }
  };

  return (
    <Dialog onOpenChange={handleDialogItemOpenChange} open={isOpen}>
      {triggerButton && (
        <DialogTrigger asChild>
          <div className="w-full text-right">
            <Button
              type="button"
              variant="default"
              className="whitespace-nowrap gap-1"
            >
              <RiAddLine />
              Add User
            </Button>
          </div>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-lg">
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col gap-y-4"
        >
          <DialogHeader>
            <DialogTitle>{mode === 'edit' ? 'Edit contact' : 'Add a new contact'}</DialogTitle>
            <DialogDescription className="mt-1 text-sm leading-6">
              {mode === 'edit' 
                ? 'Update contact information'
                : 'Get started with connecting and building relationships with a specific investor'
              }
            </DialogDescription>
          </DialogHeader>

          <div>
            <Label htmlFor="email">Email address</Label>
            <Input
              id="email"
              placeholder="yo@lanre.wtf"
              className="mt-2"
              type="email"
              {...register("email")}
            />
            {errors.email && (
              <p className="mt-1 text-xs text-red-600 dark:text-red-500">
                <span className="font-medium">{errors.email.message}</span>
              </p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label htmlFor="first_name">First name</Label>
              <Input
                id="first_name"
                placeholder="John"
                className="mt-2"
                {...register("first_name")}
              />
            </div>
            <div>
              <Label htmlFor="last_name">Last name</Label>
              <Input
                id="last_name"
                placeholder="Doe"
                className="mt-2"
                {...register("last_name")}
              />
            </div>
          </div>

          <div>
            <Label htmlFor="phone">Phone number</Label>
            <Input
              id="phone"
              placeholder="+1 (555) 123-4567"
              className="mt-2"
              type="tel"
              {...register("phone")}
            />
          </div>

          <div>
            <Label htmlFor="company">Company</Label>
            <Input
              id="company"
              placeholder="Acme Inc."
              className="mt-2"
              {...register("company")}
            />
          </div>

          <DialogFooter className="mt-6">
            <DialogClose asChild>
              <Button
                className="mt-2 w-full sm:mt-0 sm:w-fit"
                variant="secondary"
              >
                Cancel
              </Button>
            </DialogClose>
            <Button
              type="submit"
              className="w-full sm:w-fit"
              disabled={loading}
            >
              {mode === 'edit' ? 'Save changes' : 'Add Contact'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
