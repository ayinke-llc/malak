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
import { AnalyticsEvent } from "@/lib/events";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Search } from "lucide-react";

interface ContactListOption extends OptionType {
  emails: Array<{ email: string; reference: string }>;
}

const schema = yup
  .object({})
  .required();

const SendUpdateButton = ({ reference, isSent }: ButtonProps & { isSent: boolean }) => {
  const posthog = usePostHog()
  const queryClient = useQueryClient();

  const [loading, setLoading] = useState<boolean>(false);
  const [showAllRecipients, setShowAllRecipients] = useState<boolean>(false);
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [selectedTab, setSelectedTab] = useState<"lists" | "emails">("lists");
  const [values, setValues] = useState<OptionType[]>([]);
  const [emailOptions, setEmailOptions] = useState<OptionType[]>([
    { value: "john.doe@example.com", label: "john.doe@example.com", emails: [] },
    { value: "jane.smith@example.com", label: "jane.smith@example.com", emails: [] },
    { value: "dev.team@company.com", label: "dev.team@company.com", emails: [] },
  ]);

  const { data } = useQuery({
    queryKey: [LIST_CONTACT_LISTS],
    queryFn: () => {
      return client.contacts.fetchContactLists({
        include_emails: true,
      })
    },
  });

  const options: ContactListOption[] = data?.data?.lists?.map(({ list, mappings }) => {
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
  }) ?? [];

  // Filter both contact lists and individual emails based on search term
  const filteredOptions = options.filter(option => 
    option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
    option.emails.some(email => 
      email.email.toLowerCase().includes(searchTerm.toLowerCase())
    )
  );

  const filteredEmailOptions = emailOptions.filter(option =>
    option.label.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const addNewContacts = (...inputValues: string[]) => {
    // Normalize all inputs to lowercase
    const normalizedInputs = inputValues.map((input) => input.toLowerCase());

    setLoading(true);

    // Track invalid emails for error reporting
    const invalidEmails: string[] = [];
    const newOptions: OptionType[] = [];

    normalizedInputs.forEach((inputValue) => {
      if (!EmailValidator.validate(inputValue)) {
        invalidEmails.push(inputValue);
        return; // Skip invalid emails
      }

      if (!values.some((item) => item.value === inputValue)) {
        // Only add if it's not already present
        const newOption: OptionType = {
          value: inputValue,
          label: inputValue,
          emails: []
        };
        newOptions.push(newOption);
        // Also add to email options for future use
        if (!emailOptions.some(opt => opt.value === inputValue)) {
          setEmailOptions(prev => [...prev, newOption]);
        }
      }
    });

    if (invalidEmails.length > 0) {
      toast.error(
        `The following are not valid email addresses: ${invalidEmails.join(", ")}`
      );
    }

    if (newOptions.length > 0) {
      setValues((prev) => [...prev, ...newOptions]);
      toast.success("Email added successfully");
    }

    setLoading(false);
  };

  const removeContact = (index: number) => {
    setValues(values.filter((_, i) => i !== index));
  };

  const toggleShowAllRecipientState = () => setShowAllRecipients(!showAllRecipients);

  const handleOnChange = (opts: ContactListOption[]) => {
    opts.map((opt) => {
      addNewContacts(...opt.emails.map((value) => value.email))
    })
  }

  const handleEmailSelection = (selectedOptions: OptionType[]) => {
    setValues(selectedOptions);
  };

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
      emails: values.map((value) => String(value.value))
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
          <DialogContent className="sm:max-w-[600px]">
            <form onSubmit={handleSubmit(onSubmit)} autoComplete="off">
              <DialogHeader>
                <DialogTitle>Send this update</DialogTitle>
                <DialogDescription className="mt-1 text-sm leading-6">
                  Select recipients from your contact lists or enter email addresses directly.
                </DialogDescription>
              </DialogHeader>

              <div className="mt-4">
                <Tabs defaultValue="lists" value={selectedTab} onValueChange={(v) => setSelectedTab(v as "lists" | "emails")}>
                  <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="lists">Contact Lists</TabsTrigger>
                    <TabsTrigger value="emails">Individual Emails</TabsTrigger>
                  </TabsList>
                  
                  <div className="mt-4 mb-4">
                    <div className="relative">
                      <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input
                        placeholder={selectedTab === "lists" ? "Search contact lists..." : "Search or enter email addresses..."}
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="pl-8"
                      />
                    </div>
                  </div>

                  <TabsContent value="lists" className="mt-0">
                    <ScrollArea className="h-[200px] rounded-md border p-4">
                      <div className="space-y-2">
                        {filteredOptions.map((option) => (
                          <div
                            key={option.value}
                            className="flex items-center justify-between rounded-lg border p-3 hover:bg-muted cursor-pointer"
                            onClick={() => handleOnChange([option])}
                          >
                            <div>
                              <h4 className="font-medium">{option.label}</h4>
                              <p className="text-sm text-muted-foreground">
                                {option.emails.length} recipients
                              </p>
                            </div>
                            <Badge variant="secondary">
                              {values.some(v => option.emails.some(email => email.email === v.value)) ? 'Selected' : 'Select'}
                            </Badge>
                          </div>
                        ))}
                        {filteredOptions.length === 0 && (
                          <p className="text-center text-sm text-muted-foreground">No contact lists found</p>
                        )}
                      </div>
                    </ScrollArea>
                  </TabsContent>

                  <TabsContent value="emails" className="mt-0">
                    <div className="space-y-4">
                      <CreatableSelect
                        placeholder="Enter email addresses..."
                        isMulti
                        options={filteredEmailOptions}
                        value={values}
                        allowCustomInput={true}
                        onCustomInputEnter={addNewContacts}
                        onChange={handleEmailSelection}
                        className="min-h-[200px]"
                      />
                    </div>
                  </TabsContent>
                </Tabs>
              </div>

              {values.length > 0 && (
                <div className="mt-6">
                  <h4 className="text-sm font-medium mb-2">Selected Recipients ({values.length})</h4>
                  <ScrollArea className="h-[100px] rounded-md border p-2">
                    <div className="flex flex-wrap gap-2">
                      {values.map((recipient, index) => (
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
                    </div>
                  </ScrollArea>
                </div>
              )}

              <DialogFooter className="mt-6">
                <DialogClose asChild>
                  <Button
                    type="button"
                    variant="outline"
                    loading={loading}
                  >
                    Cancel
                  </Button>
                </DialogClose>
                <Button
                  type="submit"
                  variant="default"
                  loading={loading}
                  disabled={values.length === 0}
                >
                  Send to {values.length} {values.length === 1 ? 'recipient' : 'recipients'}
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
