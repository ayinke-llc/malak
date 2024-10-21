import type { ServerAPIStatus } from "@/client/Api";
import { Button } from "@/components/Button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/Dialog";
import { Label } from "@/components/Label";
import { Switch } from "@/components/Switch";
import client from "@/lib/client";
import { CREATE_CONTACT_MUTATION, LIST_CONTACT_LISTS } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiCloseLargeLine, RiCloseLine, RiMailSendLine, RiMarkupLine, RiTwitterXLine } from "@remixicon/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import * as EmailValidator from "email-validator";
import { Option } from "lucide-react";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import CreatableSelect from "react-select/creatable";
import { toast } from "sonner";
import * as yup from "yup";
import type { ButtonProps } from "./props";
import { Badge } from "@/components/Badge";
import { cx } from "@/lib/utils";

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
  const [showAllRecipients, setShowAllRecipients] = useState<boolean>(false)

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
    inputValue = inputValue.toLowerCase()

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
      return
    }

    setValues((prev) => [...prev, newOption]);
    setLoading(false)
  };

  const removeContact = (index: number) => {
    setValues(values.filter((_, i) => i !== index))
  }

  const toggleShowAllRecipientState = () => setShowAllRecipients(!showAllRecipients)

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
            <Button type="submit" size="lg" variant="primary" className="gap-1">
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
                  <CreatableSelect
                    isDisabled={loading}
                    isLoading={loading}
                    onCreateOption={createNewContact}
                    onChange={(value) => {
                      createNewContact("oopsoops@gmail.com")
                    }}
                    options={options}
                  />
                </div>

                {values.length > 0 && <div className="flex-1 mt-5">
                  <div
                    className={cx(
                      showAllRecipients ? "h-[100px]" : "h-full",
                      "w-full rounded-md border p-2 overflow-y-auto",
                    )}>
                    <div className="flex flex-wrap justify-start gap-3">
                      {values.
                        slice(0, showAllRecipients ? values.length : 5).
                        map((recipient, index) => (
                          <Badge key={index} color="gray"
                            className="flex items-center space-x-1 gap-3 mt-1"
                            variant="neutral">
                            <span>{recipient.label}</span>
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-4 w-4 p-0"
                              onClick={() => removeContact(index)}
                            >
                              <RiCloseLargeLine className="h-3 w-3" color="red" />
                              <span className="sr-only">Remove recipient</span>
                            </Button>
                          </Badge>
                        ))}
                      {values.length > 5 && (
                        <Button
                          size="sm"
                          variant="light"
                          onClick={toggleShowAllRecipientState}
                        >
                          {showAllRecipients ? "hide recipients" : `+${values.length - 5} more`}
                        </Button>
                      )}
                    </div>
                  </div>
                </div>}

                <div className="mt-4 gap-10">
                  <Switch disabled id="r3" {...register("link")} />
                  <Label disabled htmlFor="r3">
                    Coming soon. Generate a public viewable link for this update
                  </Label>
                </div>
              </DialogHeader>
              <DialogFooter className="mt-6">
                <DialogClose asChild>
                  <Button
                    type={"button"}
                    className="mt-2 w-full sm:mt-0 sm:w-fit"
                    variant="secondary"
                    isLoading={loading}
                  >
                    Cancel
                  </Button>
                </DialogClose>
                <Button
                  type="submit"
                  className="w-full sm:w-fit"
                  isLoading={loading}
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
