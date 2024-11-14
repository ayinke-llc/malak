import type { ServerAPIStatus } from "@/client/Api";
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
import client from "@/lib/client";
import {
  CREATE_CONTACT_MUTATION,
  LIST_CONTACT_LISTS,
} from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RiCloseLargeLine,
  RiMailSendLine,
} from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import * as EmailValidator from "email-validator";
import { Option } from "lucide-react";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
// import CreatableSelect from "react-select/creatable";
import { toast } from "sonner";
import * as yup from "yup";
import type { ButtonProps } from "./props";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import CreatableSelect from "@/components/ui/multi-select"

interface Option {
  readonly label: string;
  readonly value: string;
}

type SendUpdateInput = {
  email: string;
  link?: boolean;
  recipients?: Option[];
};

const schema = yup
  .object({
    email: yup.string().min(5).max(50).required(),
    link: yup.boolean().optional(),
  })
  .required();

type ListOption = Option & {
  emails?: string[];
};

const SendUpdateButton = ({ }: ButtonProps) => {

  const [loading, setLoading] = useState<boolean>(false);
  const [showAllRecipients, setShowAllRecipients] = useState<boolean>(false);

  const [options, setOptions] = useState<ListOption[]>([
    { value: "oops", label: "oops" },
    { value: "test", label: "test" },
  ]);

  const { data, error } = useQuery({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: () => client.contacts.fetchContactLists(),
  });

  const [values, setValues] = useState<Option[]>([]);

  const contactMutation = useMutation({
    mutationKey: [CREATE_CONTACT_MUTATION],
    mutationFn: (data: { email: string }) =>
      client.contacts.contactsCreate(data),
    onSuccess: ({ data }) => {
      toast.info(`${data.contact.email} has been added as a contact now`);

      const newOption = {
        value: data.contact.id,
        label: data.contact.email,
      } as Option;

      setValues((prev) => [...prev, newOption]);
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

  const createNewContact = (inputValue: string) => {
    inputValue = inputValue.toLowerCase();

    setLoading(true);

    if (!EmailValidator.validate(inputValue)) {
      toast.error("you can only add an email address as a new recipient");
      setLoading(false);
      return;
    }

    const newOption = {
      value: inputValue,
      label: inputValue,
    } as Option;

    // probably just use a set here
    if (values.some((item) => item.value === inputValue)) {
      setLoading(false);
      return;
    }

    setValues((prev) => [...prev, newOption]);
    setLoading(false);
  };

  const removeContact = (index: number) => {
    setValues(values.filter((_, i) => i !== index));
  };

  const toggleShowAllRecipientState = () =>
    setShowAllRecipients(!showAllRecipients);

  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm({
    resolver: yupResolver(schema),
  });

  const onSubmit: SubmitHandler<SendUpdateInput> = (data) => {
    setLoading(true);
  };

  return (
    <>
      <div className="flex justify-center">
        <Dialog>
          <DialogTrigger asChild>
            <Button type="submit" size="lg"
              variant="default"
              className="gap-1">
              <RiMailSendLine size={18} />
              Send
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <form onSubmit={handleSubmit(onSubmit)}>
              <DialogHeader>
                <DialogTitle>Send this update</DialogTitle>
                <DialogDescription className="mt-1 text-sm leading-6">
                  An email will be sent to all selected contacts immediately.
                  Please re-verify your content is ready and good to go
                </DialogDescription>

                <div className="mt-4">
                  <CreatableSelect placeholder="Select a list or add an email"
                    isMulti
                    options={options}
                    allowCustomInput={true}
                    onCustomInputEnter={createNewContact}
                  />
                </div>

                {values.length > 0 && (
                  <div className="flex-1 pt-5">
                    <div
                      className={cn(
                        "w-full rounded-md border bg-background p-2",
                        showAllRecipients ? "h-[100px]" : "h-full",
                        "overflow-y-auto"
                      )}
                    >
                      <div className="flex flex-wrap gap-2">
                        {values
                          .slice(0, showAllRecipients ? values.length : 5)
                          .map((recipient, index) => (
                            <Badge
                              key={index}
                              variant="secondary"
                              className="flex items-center gap-1 pr-1"
                            >
                              <span className="text-sm">{recipient.label}</span>
                              <button
                                onClick={() => removeContact(index)}
                                className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                              >
                                <RiCloseLargeLine className="h-3 w-3 text-muted-foreground hover:text-foreground" />
                                <span className="sr-only">Remove recipient</span>
                              </button>
                            </Badge>
                          ))}
                        {values.length > 5 && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={toggleShowAllRecipientState}
                            className="h-8"
                          >
                            {showAllRecipients
                              ? "Show less"
                              : `+${values.length - 5} more`}
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                )}
              </DialogHeader>
              <DialogFooter className="mt-6">
                <DialogClose asChild>
                  <Button
                    type={"button"}
                    className="mt-2 w-full sm:mt-0 sm:w-fit"
                    variant="secondary"
                    loading={loading}
                  >
                    Cancel
                  </Button>
                </DialogClose>
                <Button
                  type="submit"
                  className="w-full sm:w-fit"
                  loading={loading}
                >
                  Send
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </>
  );
};

export default SendUpdateButton;
