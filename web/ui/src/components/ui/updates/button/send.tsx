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
  FETCH_SINGLE_UPDATE,
  LIST_CONTACT_LISTS,
  SEND_UPDATE
} from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RiCloseLargeLine,
  RiMailSendLine,
} from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as EmailValidator from "email-validator";
import { Option } from "lucide-react";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";
import type { ButtonProps } from "./props";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import CreatableSelect, { OptionType } from "@/components/ui/multi-select";
import { ServerAPIStatus, ServerSendUpdateRequest } from "@/client/Api";
import { AxiosError } from "axios";
import { usePostHog } from "posthog-js/react";

interface Option {
  readonly label: string;
  readonly value: string;
}

const schema = yup
  .object({})
  .required();

const SendUpdateButton = ({ reference, isSent }: ButtonProps & { isSent: boolean }) => {

  const posthog = usePostHog()

  const queryClient = useQueryClient();

  const [loading, setLoading] = useState<boolean>(false);
  const [showAllRecipients, setShowAllRecipients] = useState<boolean>(false);

  const { data } = useQuery({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: () => {
      return client.contacts.fetchContactLists({
        include_emails: true,
      })
    },
  });

  const options: OptionType[] = data?.data?.lists?.map(({ list, mappings }) => {
    return {
      emails: mappings?.map((mapping) => {
        return {
          email: mapping?.email as string,
          reference: mapping?.reference as string
        }
      }) ?? [],
      label: list?.title as string,
      value: list?.reference as string
    }
  }) ?? []

  const [values, setValues] = useState<Option[]>([]);

  const addNewContacts = (...inputValues: string[]) => {
    // Normalize all inputs to lowercase
    const normalizedInputs = inputValues.map((input) => input.toLowerCase());

    setLoading(true);

    // Track invalid emails for error reporting
    const invalidEmails: string[] = [];

    const newOptions: Option[] = [];

    normalizedInputs.forEach((inputValue) => {
      if (!EmailValidator.validate(inputValue)) {
        invalidEmails.push(inputValue);
        return; // Skip invalid emails
      }

      if (!values.some((item) => item.value === inputValue)) {
        // Only add if it's not already present
        newOptions.push({
          value: inputValue,
          label: inputValue,
        });
      }
    });

    if (invalidEmails.length > 0) {
      toast.error(
        `The following are not valid email addresses: ${invalidEmails.join(", ")}`
      );
    }

    if (newOptions.length > 0) {
      setValues((prev) => [...prev, ...newOptions]);
      toast.success("added email")
    }

    setLoading(false);
  };

  const removeContact = (index: number) => {
    setValues(values.filter((_, i) => i !== index));
  };

  const toggleShowAllRecipientState = () => setShowAllRecipients(!showAllRecipients);

  const handleOnChange = (opts: OptionType[]) => {
    opts.map((opt) => {
      addNewContacts(...opt.emails.map((value) => value.email))
    })
  }

  const mutation = useMutation({
    mutationKey: [SEND_UPDATE],
    mutationFn: (data: ServerSendUpdateRequest) => {
      return client.workspaces.sendUpdate(reference, data)
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
      queryClient.invalidateQueries({ queryKey: [FETCH_SINGLE_UPDATE, reference] })
      posthog?.capture(AnalyticsEvent.SendUpdate, {
        liveRecipient: true
      })
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
    onMutate: () => setLoading(true),
  });

  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  const onSubmit: SubmitHandler<{}> = () => {
    mutation.mutate({
      emails: values.map((value) => value.value)
    })
  };

  const {
    handleSubmit,
  } = useForm({
    resolver: yupResolver(schema),
  });

  return (
    <>
      <div className="flex justify-center">
        <Dialog>
          <DialogTrigger asChild>
            <Button type="submit" size="lg"
              variant="default"
              className="gap-1">
              <RiMailSendLine size={18} />
              {isSent ? "Add recipient" : "Send"}
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
                    onCustomInputEnter={addNewContacts}
                    onChange={handleOnChange}
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
