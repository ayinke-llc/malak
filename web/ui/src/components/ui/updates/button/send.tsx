"use client"

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
  LIST_ALL_CONTACTS,
  LIST_CONTACT_LISTS,
  SEND_UPDATE
} from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiMailSendLine, RiCloseLine, RiUserLine, RiTeamLine, RiAddLine } from "@remixicon/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as EmailValidator from "email-validator";
import { useState, useEffect } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";
import type { ButtonProps } from "./props";
import { Badge } from "@/components/ui/badge";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ServerAPIStatus, ServerSendUpdateRequest } from "@/client/Api";
import { AxiosError } from "axios";
import { usePostHog } from "posthog-js/react";
import { AnalyticsEvent } from "@/lib/events";

interface ContactList {
  emails: Array<string>;
  label: string;
  value: string;
}

interface FormValues {
  recipients: string[];
}

const schema = yup.object().shape({
  recipients: yup.array().of(yup.string().required()).min(1).required(),
});

const SendUpdateButton = ({ reference, isSent }: ButtonProps & { isSent: boolean }) => {
  const posthog = usePostHog();
  const queryClient = useQueryClient();

  const [loading, setLoading] = useState<boolean>(false);
  const [showAllRecipients, setShowAllRecipients] = useState<boolean>(false);
  const [open, setOpen] = useState(false);
  const [inputValue, setInputValue] = useState("");

  const { data } = useQuery({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: async () => {
      try {
        const response = await client.contacts.fetchContactLists({
          include_emails: true,
        });
        return response?.data ?? { lists: [] };
      } catch (error) {
        return { lists: [] };
      }
    },
    initialData: { lists: [] }
  });

  const { data: contactsData } = useQuery({
    queryKey: [LIST_ALL_CONTACTS],
    queryFn: async () => {
      const response = await client.contacts.listAllContacts();
      return response?.data?.contacts || [];
    },
  });

  const options: ContactList[] = (data?.lists || []).map(({ list, mappings }) => ({
    emails: (mappings || []).map((m) => m?.email ?? "").filter(Boolean),
    label: list?.title ?? "Untitled List",
    value: list?.reference ?? "",
  }));

  const {
    handleSubmit,
    setValue,
    formState: { isSubmitting }
  } = useForm<FormValues>({
    resolver: yupResolver(schema),
    defaultValues: {
      recipients: []
    }
  });

  const [values, setValues] = useState<string[]>([]);

  useEffect(() => {
    setValue('recipients', values);
  }, [values, setValue]);

  const handleOnChange = (opts: ContactList[]) => {
    if (!Array.isArray(opts)) {
      return
    }

    const newEmails = opts.flatMap(opt =>
      Array.isArray(opt?.emails) ? opt.emails.filter(e => EmailValidator.validate(e)) : []
    );

    setValues(prevValues => {
      const combinedEmails = [...(prevValues || []), ...newEmails];
      const uniqueEmails = [...new Set(combinedEmails)];
      setValue('recipients', uniqueEmails);
      return uniqueEmails;
    });
  };

  const addNewContacts = (...inputValues: string[]) => {
    const normalizedInputs = inputValues
      .map((input) => input?.toLowerCase()?.trim() ?? '')
      .filter(Boolean);

    const invalidEmails: string[] = [];
    const newOptions: string[] = [];

    normalizedInputs.forEach((inputValue) => {
      if (!EmailValidator.validate(inputValue)) {
        invalidEmails.push(inputValue);
        return;
      }
      if (!values.includes(inputValue)) {
        newOptions.push(inputValue);
      }
    });

    if (invalidEmails.length > 0) {
      toast.error(`Invalid email${invalidEmails.length > 1 ? 's' : ''}: ${invalidEmails.join(", ")}`);
      return;
    }

    if (newOptions.length > 0) {
      setValues((prev) => [...prev, ...newOptions]);
      setValue('recipients', [...values, ...newOptions]);
      toast.success(newOptions.length === 1 ? "Email added" : `${newOptions.length} emails added`);
    }
  };

  const removeContact = (index: number) => {
    const newValues = [...values];
    newValues.splice(index, 1);
    setValues(newValues);
    setValue('recipients', newValues);
  };

  const mutation = useMutation({
    mutationKey: [SEND_UPDATE],
    mutationFn: (data: ServerSendUpdateRequest) => {
      return client.workspaces.sendUpdate(reference, data)
    },
    onSuccess: ({ data }) => {
      toast.success(data?.message || "Update sent successfully");
      queryClient.invalidateQueries({ queryKey: [FETCH_SINGLE_UPDATE, reference] });
      posthog?.capture(AnalyticsEvent.SendUpdate, {
        liveRecipient: true
      });
      setValues([]);
      setOpen(false);
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response?.data?.message) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onSettled: () => setLoading(false),
    onMutate: () => setLoading(true),
  });

  const onSubmit: SubmitHandler<FormValues> = async (data) => {
    try {
      await schema.validate(data);
      mutation.mutate({
        emails: data.recipients.filter(Boolean)
      });
    } catch (error) {
      if (error instanceof yup.ValidationError) {
        toast.error(error.message);
      }
    }
  };

  return (
    <div className="flex justify-center">
      <Dialog>
        <DialogTrigger asChild>
          <Button type="submit" size="lg" variant="default" className="gap-1">
            <RiMailSendLine size={18} />
            {isSent ? "Add recipient" : "Send"}
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-2xl">
          <form onSubmit={handleSubmit(onSubmit)}>
            <DialogHeader>
              <DialogTitle className="text-xl">Send this update</DialogTitle>
              <DialogDescription className="mt-3 text-base leading-6">
                An email will be sent to all selected contacts immediately.
                Please re-verify your content is ready and good to go.
              </DialogDescription>

              <div className="mt-4 flex flex-col gap-4">
                <Popover open={open} onOpenChange={setOpen}>
                  <PopoverTrigger asChild>
                    <Button
                      variant="outline"
                      role="combobox"
                      aria-expanded={open}
                      className="w-full justify-between"
                    >
                      Select contacts or enter email...
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-[520px] p-0" side="bottom" align="start">
                    <Command className="w-full rounded-lg border shadow-md">
                      <CommandInput
                        value={inputValue}
                        onValueChange={setInputValue}
                        placeholder="Search contacts or enter email..."
                        className="h-11 px-4"
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            e.preventDefault();
                            const selectedItem = document.querySelector('[data-selected="true"]');
                            if (selectedItem) {
                              selectedItem?.click();
                            } else if (inputValue.trim()) {
                              addNewContacts(inputValue.trim());
                              setInputValue("");
                              setOpen(false);
                            }
                          }
                        }}
                      />
                      <CommandList className="max-h-[300px] overflow-y-auto">
                        <CommandEmpty className="py-6 text-center text-sm">
                          {!inputValue ? (
                            "Type an email address to add directly."
                          ) : (
                            <>Press Enter to add "{inputValue}" as an email address</>
                          )}
                        </CommandEmpty>

                        {options.length > 0 && (
                          <>
                            <CommandGroup heading="Contact Lists">
                              {options.map((option) => (
                                <CommandItem
                                  key={option.value || Math.random().toString(36).substring(2, 9)}
                                  onSelect={() => {
                                    handleOnChange([option]);
                                    setOpen(false);
                                  }}
                                  className="flex items-center gap-2 px-4 py-2 hover:bg-accent cursor-pointer"
                                >
                                  <RiTeamLine className="h-4 w-4" />
                                  <span>{option.label}</span>
                                  <span className="ml-auto text-xs text-muted-foreground">
                                    {option.emails.length} contacts
                                  </span>
                                </CommandItem>
                              ))}
                            </CommandGroup>
                            <CommandSeparator />
                          </>
                        )}

                        {contactsData && contactsData.length > 0 && (
                          <>
                            <CommandGroup heading="Emails">
                              {contactsData.map((contact) => (
                                <CommandItem
                                  key={contact.email}
                                  onSelect={() => {
                                    addNewContacts(contact.email as string);
                                    setOpen(false);
                                  }}
                                  className="flex items-center gap-2 px-4 py-2 hover:bg-accent cursor-pointer"
                                >
                                  <RiUserLine className="h-4 w-4" />
                                  <span>{contact.email}</span>
                                </CommandItem>
                              ))}
                            </CommandGroup>
                            <CommandSeparator />
                          </>
                        )}

                        <CommandGroup heading="Actions">
                          <CommandItem
                            onSelect={() => {
                              if (inputValue.trim()) {
                                addNewContacts(inputValue.trim());
                                setInputValue("");
                                setOpen(false);
                              }
                            }}
                            className="flex items-center gap-2 px-4 py-2"
                          >
                            <RiAddLine className="h-4 w-4" />
                            <span>Add new contact</span>
                          </CommandItem>
                        </CommandGroup>
                      </CommandList>
                    </Command>
                  </PopoverContent>
                </Popover>

                {values.length > 0 && (
                  <div className="flex-1">
                    <div className="w-full rounded-md border bg-background p-3">
                      <div className="flex flex-wrap gap-2">
                        {(showAllRecipients ? values : values.slice(0, 5)).map((recipient, index) => (
                          <Badge
                            key={`${recipient}-${index}`}
                            variant="secondary"
                            className="flex items-center gap-1.5 px-2 py-1"
                          >
                            <RiUserLine className="h-3 w-3 text-muted-foreground" />
                            <span className="text-sm">{recipient}</span>
                            <button
                              type="button"
                              onClick={() => removeContact(index)}
                              className="ml-1 rounded-full outline-none hover:bg-secondary-foreground/10 p-0.5 transition-colors"
                            >
                              <RiCloseLine className="h-3 w-3 text-muted-foreground hover:text-foreground" />
                              <span className="sr-only">Remove {recipient}</span>
                            </button>
                          </Badge>
                        ))}
                        {values.length > 5 && (
                          <button
                            type="button"
                            onClick={() => setShowAllRecipients(!showAllRecipients)}
                            className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground transition-colors"
                          >
                            {showAllRecipients ? (
                              "Show less"
                            ) : (
                              <span className="flex items-center gap-1">
                                +{values.length - 5} more recipients
                              </span>
                            )}
                          </button>
                        )}
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </DialogHeader>
            <DialogFooter className="mt-6">
              <DialogClose asChild>
                <Button
                  type="button"
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
                loading={loading || isSubmitting}
                disabled={values.length === 0}
              >
                Send
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default SendUpdateButton;
