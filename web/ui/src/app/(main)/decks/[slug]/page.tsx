"use client";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { RiFileCopyLine, RiArrowLeftLine, RiEyeLine, RiTimeLine, RiDownloadLine, RiUserLine, RiSettings4Line, RiPushpin2Line, RiPushpin2Fill, RiExternalLinkLine } from "@remixicon/react";
import { format } from "date-fns";
import { useState, useMemo } from "react";
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
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useQuery } from "@tanstack/react-query";
import client from "@/lib/client";
import { useRouter } from "next/navigation";
import { FETCH_DECK } from "@/lib/query-constants";
import { DECKS_DOMAIN } from "@/lib/config";
import DeleteDeck from "@/components/ui/decks/details/delete";

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
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isPinned, setIsPinned] = useState(false);
  const [isPinning, setIsPinning] = useState(false);

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
    setIsSubmitting(true);
    try {
      // TODO: API call to update settings
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulated API call
      console.log("Form data:", formData);
      toast.success("Settings updated successfully");
    } catch (error) {
      toast.error("Failed to update settings");
    } finally {
      setIsSubmitting(false);
    }
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
    setIsPinning(true);
    try {
      // TODO: API call to toggle pin status
      await new Promise(resolve => setTimeout(resolve, 500)); // Simulated API call
      setIsPinned(!isPinned);
      toast.success(isPinned ? "Deck unpinned" : "Deck pinned", {
        description: isPinned ? "Deck removed from pinned items" : "Deck added to pinned items",
      });
    } catch (error) {
      toast.error("Failed to update pin status");
    } finally {
      setIsPinning(false);
    }
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
          className="inline-flex items-center text-sm text-zinc-400 hover:text-zinc-300"
        >
          <RiArrowLeftLine className="mr-1 h-4 w-4" />
          Back to decks
        </Link>

        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            className="text-zinc-400 hover:text-zinc-300"
            onClick={() => window.open(data?.data?.deck?.short_link, '_blank')}
          >
            <RiExternalLinkLine className="h-5 w-5" />
          </Button>

          <Button
            variant="ghost"
            size="icon"
            className={`${isPinned ? 'text-blue-400 hover:text-blue-300' : 'text-zinc-400 hover:text-zinc-300'
              } ${isPinning ? 'opacity-50 cursor-not-allowed' : ''}`}
            onClick={handleTogglePin}
            disabled={isPinning}
          >
            {isPinning ? (
              <div className="h-5 w-5 animate-spin rounded-full border-2 border-current border-t-transparent" />
            ) : isPinned ? (
              <RiPushpin2Fill className="h-5 w-5" />
            ) : (
              <RiPushpin2Line className="h-5 w-5" />
            )}
          </Button>

          <Dialog>
            <DialogTrigger>
              <div className="text-zinc-400 hover:text-zinc-300 cursor-pointer p-2">
                <RiSettings4Line className="h-5 w-5" />
              </div>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle className="text-zinc-100">Deck settings</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 py-4">
                {/* Sharing Settings */}
                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-zinc-100">Sharing settings</h3>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <Label htmlFor="enable-downloading" className="text-zinc-100">Enable downloading</Label>
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
                      <Label htmlFor="require-email" className="text-zinc-100">Require email to view</Label>
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

                <Separator className="bg-zinc-800" />

                {/* Advanced Settings */}
                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-zinc-100">Advanced settings</h3>
                  <div className="space-y-4">
                    <div>
                      <div className="flex items-center justify-between mb-4">
                        <Label htmlFor="password-protection" className="text-zinc-100">Password protection</Label>
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
                                className="bg-zinc-800 border-zinc-700 text-zinc-100 placeholder:text-zinc-500"
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
                    className="text-zinc-400 hover:text-zinc-300"
                    onClick={handleReset}
                    disabled={isSubmitting}
                  >
                    Reset
                  </Button>
                  <Button
                    type="submit"
                    className="bg-zinc-100 text-zinc-900 hover:bg-zinc-200 disabled:opacity-50"
                    disabled={isSubmitting}
                  >
                    {isSubmitting ? (
                      <div className="flex items-center gap-2">
                        <div className="h-4 w-4 animate-spin rounded-full border-2 border-zinc-900 border-t-transparent" />
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
        <Card className="bg-zinc-900/50 p-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div>
              <div className="flex items-center gap-2 text-zinc-400 mb-1">
                <RiEyeLine className="h-4 w-4" />
                <span className="text-sm">Total views</span>
              </div>
              <p className="text-2xl font-medium text-zinc-100">0</p>
            </div>
            <div>
              <div className="flex items-center gap-2 text-zinc-400 mb-1">
                <RiUserLine className="h-4 w-4" />
                <span className="text-sm">Unique views</span>
              </div>
              <p className="text-2xl font-medium text-zinc-100">0</p>
            </div>
            <div>
              <div className="flex items-center gap-2 text-zinc-400 mb-1">
                <RiTimeLine className="h-4 w-4" />
                <span className="text-sm">Time spent (avg)</span>
              </div>
              <p className="text-2xl font-medium text-zinc-100">00:00</p>
            </div>
            <div>
              <div className="flex items-center gap-2 text-zinc-400 mb-1">
                <RiDownloadLine className="h-4 w-4" />
                <span className="text-sm">Downloads</span>
              </div>
              <p className="text-2xl font-medium text-zinc-100">0</p>
            </div>
          </div>
        </Card>

        {/* Main Content Card */}
        <Card className="bg-zinc-900/50 p-6">
          <div className="space-y-8">
            {/* Header Section */}
            <div>
              <h1 className="text-xl font-medium text-zinc-100 mb-2">
                {data?.data?.deck?.title}
              </h1>
              <p className="text-sm text-zinc-400">
                {/* TODO: Add description field to deck */}
                {data?.data?.deck?.title}
              </p>
            </div>

            {/* Details Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-zinc-400 mb-1">Uploaded</h3>
                  <p className="text-zinc-100">
                    {data?.data?.deck?.created_at ? (
                      format(new Date(data.data.deck.created_at), "MMMM d, yyyy 'at' h:mm a")
                    ) : (
                      "-"
                    )}
                  </p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-zinc-400 mb-1">File Size</h3>
                  <p className="text-zinc-100">-</p>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-zinc-400 mb-1">Share URL</h3>
                  <div className="flex items-center gap-2 max-w-md">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <div
                          className="block rounded bg-zinc-800 px-3 py-2 text-sm text-zinc-100 truncate cursor-pointer w-full"
                          onClick={() => {
                            if (data?.data?.deck?.short_link) {
                              copyToClipboard(`${DECKS_DOMAIN}/${data.data.deck.short_link}`);
                            }
                          }}
                        >
                          <code className="text-zinc-100">
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
                      className="shrink-0 h-9 w-9 p-0 text-zinc-400 hover:text-zinc-100"
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

    </div>
  );
}
