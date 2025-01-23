"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  RiFileCopyLine, RiArrowLeftLine, RiEyeLine,
  RiTimeLine, RiDownloadLine, RiUserLine, RiSettings4Line,
  RiPushpin2Line, RiPushpin2Fill, RiExternalLinkLine
} from "@remixicon/react";
import { format } from "date-fns";
import { useState, useMemo, useEffect } from "react";
import { toast } from "sonner";
import Link from "next/link";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useForm, Controller } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  Tooltip, TooltipContent,
  TooltipTrigger
} from "@/components/ui/tooltip";
import { useMutation, useQuery } from "@tanstack/react-query";
import client from "@/lib/client";
import { useRouter } from "next/navigation";
import { FETCH_DECK } from "@/lib/query-constants";
import { DECKS_DOMAIN } from "@/lib/config";
import DeleteDeck from "@/components/ui/decks/details/delete";
import { ServerAPIStatus, ServerFetchDeckResponse } from "@/client/Api";
import { AxiosError, AxiosResponse } from "axios";

type SettingsFormData = {
  enableDownloading: boolean;
  requireEmail: boolean;
  passwordProtection: boolean;
  password?: string | undefined;
};

const settingsSchema = yup.object().shape({
  enableDownloading: yup.boolean().defined(),
  requireEmail: yup.boolean().defined(),
  passwordProtection: yup.boolean().defined(),
  password: yup.string().when('passwordProtection', {
    is: true,
    then: (schema) => schema.required('Password is required when protection is enabled')
      .min(6, 'Password must be at least 6 characters'),
    otherwise: (schema) => schema.optional().nullable(),
  }),
}) satisfies yup.ObjectSchema<SettingsFormData>;

export default function DeckDetails({ params }: { params: { slug: string } }) {

  const router = useRouter();

  const [isPinned, setIsPinned] = useState<boolean>(false);

  const { data, isLoading, error } = useQuery({
    queryKey: [FETCH_DECK],
    queryFn: () => client.decks.decksDetail(params.slug),
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  const defaultValues = useMemo(() => ({
    enableDownloading: data?.data?.deck?.preferences?.enable_downloading ?? true,
    requireEmail: data?.data?.deck?.preferences?.require_email ?? false,
    passwordProtection: data?.data?.deck?.preferences?.password?.enabled ?? false,
    password: data?.data?.deck?.preferences?.password?.password,
  }), [data]);

  const pinMutation = useMutation({
    mutationFn: () => {
      return client.decks.togglePin(params.slug, {})
    },
    gcTime: 0,
    onError: (err: AxiosError<ServerAPIStatus>): void => {
      toast.error(err?.response?.data?.message || "an error occurred while updating pinned status");
    },
    onSuccess: (resp: AxiosResponse<ServerFetchDeckResponse>) => {
      setIsPinned(resp.data?.deck?.is_pinned as boolean)
      toast.success(resp.data.message)
    },
  });

  const mutation = useMutation({
    mutationFn: (data: SettingsFormData) => {
      return client.decks.preferencesUpdate(params.slug, {
        enable_downloading: data.enableDownloading,
        password_protection: {
          enabled: data.passwordProtection,
          value: data.password
        },
        require_email: data.requireEmail
      })
    },
    gcTime: 0,
    onError: (err: AxiosError<ServerAPIStatus>): void => {
      toast.error(err?.response?.data?.message || "an error occurred while updating preferences");
    },
    onSuccess: (resp: AxiosResponse<ServerFetchDeckResponse>) => {
      toast.success(resp.data.message)
    },
  });

  const {
    control,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<SettingsFormData>({
    resolver: yupResolver(settingsSchema),
    defaultValues,
  });

  const handleReset = () => {
    reset(defaultValues);
  };

  const passwordProtection = watch("passwordProtection");

  const onSubmit = async (formData: SettingsFormData) => {
    mutation.mutate(formData)
  };

  const copyToClipboard = (text: string) => {
    try {
      navigator.clipboard.writeText(text);
      toast.success("Link copied to clipboard");
    } catch (err) {
      toast.error("Failed to copy link");
    }
  };

  const handleTogglePin = async () => {
    pinMutation.mutate()
  };

  if (error || isLoading) {
    return (
      <div className="flex items-center justify-center h-[calc(100vh-4rem)]">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-zinc-100 border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="pt-6">
      <div className="mb-8 flex items-center justify-between">
        <Link
          href="/decks"
          className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground"
        >
          <RiArrowLeftLine className="mr-1 h-4 w-4" />
          Back to decks
        </Link>

        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            className="text-muted-foreground hover:text-foreground"
            onClick={() => window.open(data?.data?.deck?.short_link, '_blank')}
          >
            <RiExternalLinkLine className="h-5 w-5" />
          </Button>

          <Button
            variant="ghost"
            size="icon"
            className={`${isPinned ? 'text-blue-600 hover:text-blue-700' : 'text-muted-foreground hover:text-foreground'
              } ${pinMutation.isPending ? 'opacity-50 cursor-not-allowed' : ''}`}
            onClick={handleTogglePin}
            disabled={pinMutation.isPending}
          >
            {pinMutation.isPending ? (
              <div className="h-5 w-5 animate-spin rounded-full border-2 border-current border-t-transparent" />
            ) : isPinned ? (
              <RiPushpin2Fill className="h-5 w-5" />
            ) : (
              <RiPushpin2Line className="h-5 w-5" />
            )}
          </Button>

          <Dialog>
            <DialogTrigger>
              <div className="text-muted-foreground hover:text-foreground cursor-pointer p-2">
                <RiSettings4Line className="h-5 w-5" />
              </div>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Deck settings</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 py-4">
                {/* Sharing Settings */}
                <div className="space-y-4">
                  <h3 className="text-sm font-medium">Sharing settings</h3>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <Label htmlFor="enable-downloading">Enable downloading</Label>
                      <Controller
                        name="enableDownloading"
                        control={control}
                        render={({ field: { value, onChange } }) => (
                          <Switch
                            id="enable-downloading"
                            checked={value}
                            onCheckedChange={onChange}
                          />
                        )}
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <Label htmlFor="require-email">Require email to view</Label>
                      <Controller
                        name="requireEmail"
                        control={control}
                        render={({ field: { value, onChange } }) => (
                          <Switch
                            id="require-email"
                            checked={value}
                            onCheckedChange={onChange}
                          />
                        )}
                      />
                    </div>
                  </div>
                </div>

                <Separator />

                {/* Advanced Settings */}
                <div className="space-y-4">
                  <h3 className="text-sm font-medium">Advanced settings</h3>
                  <div className="space-y-4">
                    <div>
                      <div className="flex items-center justify-between mb-4">
                        <Label htmlFor="password-protection">Password protection</Label>
                        <Controller
                          name="passwordProtection"
                          control={control}
                          render={({ field: { value, onChange } }) => (
                            <Switch
                              id="password-protection"
                              checked={value}
                              onCheckedChange={onChange}
                            />
                          )}
                        />
                      </div>
                      {passwordProtection && (
                        <div className="mt-2 space-y-2">
                          <Controller
                            name="password"
                            control={control}
                            render={({ field }) => (
                              <Input
                                {...field}
                                type="password"
                                placeholder="Enter password"
                                className="placeholder:text-muted-foreground"
                              />
                            )}
                          />
                          {errors.password && (
                            <p className="text-sm text-red-500">{errors.password.message}</p>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                <div className="flex justify-end gap-2 pt-4">
                  <Button
                    type="button"
                    variant="ghost"
                    className="text-muted-foreground hover:text-foreground"
                    onClick={handleReset}
                    disabled={mutation.isPending}
                  >
                    Reset
                  </Button>
                  <Button
                    type="submit"
                    disabled={mutation.isPending}
                  >
                    {mutation.isPending ? (
                      <div className="flex items-center gap-2">
                        <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                        Saving...
                      </div>
                    ) : (
                      "Save changes"
                    )}
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
          <DeleteDeck reference={params.slug} />
        </div>
      </div>

      <div className="space-y-6">
        {/* Metrics Overview */}
        <Card className="p-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div>
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <RiEyeLine className="h-4 w-4" />
                <span className="text-sm">Total views</span>
              </div>
              <p className="text-2xl font-medium">0</p>
            </div>
            <div>
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <RiUserLine className="h-4 w-4" />
                <span className="text-sm">Unique views</span>
              </div>
              <p className="text-2xl font-medium">0</p>
            </div>
            <div>
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <RiTimeLine className="h-4 w-4" />
                <span className="text-sm">Time spent (avg)</span>
              </div>
              <p className="text-2xl font-medium">00:00</p>
            </div>
            <div>
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <RiDownloadLine className="h-4 w-4" />
                <span className="text-sm">Downloads</span>
              </div>
              <p className="text-2xl font-medium">0</p>
            </div>
          </div>
        </Card>

        {/* Main Content Card */}
        <Card className="p-6">
          <div className="space-y-8">
            {/* Header Section */}
            <div>
              <h1 className="text-xl font-medium mb-2">
                {data?.data?.deck?.title}
              </h1>
              <p className="text-sm text-muted-foreground">
                {/* TODO: Add description field to deck */}
                {data?.data?.deck?.title}
              </p>
            </div>

            {/* Details Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-muted-foreground mb-1">Uploaded</h3>
                  <p>
                    {data?.data?.deck?.created_at ? (
                      format(new Date(data.data.deck.created_at), "MMMM d, yyyy 'at' h:mm a")
                    ) : (
                      "-"
                    )}
                  </p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-muted-foreground mb-1">File Size</h3>
                  <p>-</p>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-muted-foreground mb-1">Share URL</h3>
                  <div className="flex items-center gap-2 max-w-md">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <div
                          className="block rounded border px-3 py-2 text-sm truncate cursor-pointer w-full"
                          onClick={() => {
                            if (data?.data?.deck?.short_link) {
                              copyToClipboard(`${DECKS_DOMAIN}/${data.data.deck.short_link}`);
                            }
                          }}
                        >
                          <code>
                            {DECKS_DOMAIN}/{data?.data?.deck?.short_link || "-"}
                          </code>
                        </div>
                      </TooltipTrigger>
                      <TooltipContent side="top">
                        <p className="text-sm">Click to copy</p>
                      </TooltipContent>
                    </Tooltip>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="shrink-0 h-9 w-9 p-0 text-muted-foreground hover:text-foreground"
                      onClick={() => {
                        if (data?.data?.deck?.short_link) {
                          copyToClipboard(`${DECKS_DOMAIN}/${data.data.deck.short_link}`);
                        }
                      }}
                      disabled={!data?.data?.deck?.short_link}
                    >
                      <RiFileCopyLine className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Card>
      </div>

    </div >
  );
}
