/* eslint-disable @next/next/no-img-element */
"use client"

import { ServerAPIStatus } from "@/client/Api"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from "@/components/ui/card"
import Skeleton from "@/components/ui/custom/loader/skeleton"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select, SelectContent,
  SelectItem, SelectTrigger,
  SelectValue
} from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import client from "@/lib/client"
import {
  FETCH_PREFERENCES, UPDATE_WORKSPACE,
  UPDATE_WORKSPACE_PREFERENCES, UPLOAD_IMAGE
} from "@/lib/query-constants"
import timezoneMap from "@/lib/timezone"
import useWorkspacesStore from "@/store/workspace"
import { yupResolver } from "@hookform/resolvers/yup"
import { useMutation, useQuery } from "@tanstack/react-query"
import { AxiosError } from "axios"
import { useEffect } from "react"
import { Controller, useForm } from "react-hook-form"
import { toast } from "sonner"
import * as yup from "yup"

export default function GeneralSettings() {
  return (
    <div className="grid gap-6 md:grid-cols-2">
      <CompanyUpdateCard />
      <NewsletterCard />
    </div>
  )
}

const schema = yup.object({
  companyName: yup.string().required("Company name is required"),
  website: yup.string().url("Invalid URL").optional(),
  timezone: yup.string().required("Timezone is required"),
  image: yup.string().optional().test(
    "is-valid-url",
    "Invalid URL",
    (value) => {
      if (!value) return true; // Optional field, allow empty values
      try {
        const url = new URL(value);
        // Allow URLs with localhost
        return !!url && (url.hostname === "localhost" || /^[a-zA-Z0-9.-]+$/.test(url.hostname));
      } catch {
        return false;
      }
    }
  )
});

type FormData = yup.InferType<typeof schema>;

const CompanyUpdateCard = () => {

  const current = useWorkspacesStore((state) => state.current);

  const mutation = useMutation({
    mutationKey: [UPDATE_WORKSPACE],
    mutationFn: (data: FormData) => {
      return client.workspaces.workspacesPartialUpdate({
        workspace_name: data.companyName,
        timezone: data.timezone,
        website: data.website,
        logo: data.image,
      })
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message || "An error occurred while updating workspace");
    },
  });

  const uploadMutation = useMutation({
    mutationKey: [UPLOAD_IMAGE],
    mutationFn: (file: File) => client.uploads.uploadImage({ image_body: file }),
    onSuccess: ({ data }) => {
      setValue("image", data.url);
      toast.success("logo uploaded successfully");
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err?.response?.data?.message || "an error occurred while uploading your logo");
    }
  });

  const onFileChange = async (file: File | null) => {
    if (file) {
      uploadMutation.mutate(file)
    }
  };

  const {
    control,
    handleSubmit,
    formState: { errors },
    setValue,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: {
      companyName: current?.workspace_name || "",
      website: current?.website || "",
      timezone: current?.timezone || "UTC",
      image: current?.logo_url || ""
    },
  });

  const onSubmit = (data: FormData) => mutation.mutate(data)

  const getTimezoneOffset = (timezone: string) => {
    const date = new Date();
    const formatter = new Intl.DateTimeFormat("en-US", {
      timeZone: timezone,
      timeZoneName: "shortOffset",
    });
    const parts = formatter.formatToParts(date);

    const offsetPart = parts.find((part) => part.type === "timeZoneName");
    if (offsetPart) {
      return offsetPart.value;
    }

    return "UTC";
  };


  return (
    <Card>
      <CardHeader>
        <CardTitle>Company details</CardTitle>
        <CardDescription>Manage your account settings and preferences.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <Label htmlFor="companyName">Company name:</Label>
            <Controller
              name="companyName"
              control={control}
              render={({ field }) => (
                <Input
                  {...field}
                  id="companyName"
                  placeholder="Enter company name"
                  className="mt-2"
                />
              )}
            />
            {errors.companyName && (
              <p className="text-sm text-red-500">{errors.companyName.message}</p>
            )}
          </div>

          <div className="flex space-x-4 mt-4">
            <div className="flex-1 space-y-2">
              <Label htmlFor="website">Website</Label>
              <Controller
                name="website"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    id="website"
                    placeholder="https://example.com"
                    className="mt-2"
                  />
                )}
              />
              {errors.website && (
                <p className="text-sm text-red-500">{errors.website.message}</p>
              )}
            </div>

            <div className="flex-1 space-y-2">
              <Label htmlFor="timezone">Timezone</Label>
              <Controller
                name="timezone"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <SelectTrigger id="timezone">
                      <SelectValue placeholder="Select Timezone" />
                    </SelectTrigger>
                    <SelectContent>
                      {Object.entries(timezoneMap).map(([timezone, label], idx) => (
                        <SelectItem value={timezone} key={idx}>
                          {`${label} (${getTimezoneOffset(timezone)})`}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.timezone && (
                <p className="text-sm text-red-500">{errors.timezone.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-2 mt-4">
            <Label htmlFor="image">Company Logo</Label>
            <Controller
              name="image"
              control={control}
              render={({ field }) => (
                <div>
                  <Input
                    type="file"
                    id="image"
                    accept="image/png, image/jpeg"
                    onChange={(e) => onFileChange(e.target.files?.[0] || null)}
                  />
                  {field.value && (
                    <img
                      src={field.value}
                      alt="Uploaded Preview"
                      className="mt-2 w-32 h-32 object-cover rounded"
                    />
                  )}
                </div>
              )}
            />
            {errors.image && (
              <p className="text-sm text-red-500">{errors.image.message}</p>
            )}
          </div>

          <CardFooter className="mt-6 p-0">
            <div className="space-x-3">
              <Button type="submit">Save Preferences</Button>
              <Button variant={"destructive"}>Delete Workspace</Button>
            </div>
          </CardFooter>
        </form>
      </CardContent>
    </Card>
  );
};

const marketingSchema = yup.object({
  marketingEmails: yup.boolean().required("This field is required."),
  productUpdates: yup.boolean().required("This field is required."),
});

type CommunicationFormData = yup.InferType<typeof marketingSchema>;

const NewsletterCard = () => {

  const { data, isLoading, error } = useQuery({
    queryKey: [FETCH_PREFERENCES],
    queryFn: () => client.workspaces.preferencesList(),
  });

  useEffect(() => {
    if (error) {
      toast.error("Error occurred while fetching communication preferences");
    }
  }, [error]);

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm({
    defaultValues: {
      marketingEmails: data?.data?.preferences?.communication?.enable_marketing as boolean,
      productUpdates: data?.data?.preferences?.communication?.enable_product_updates as boolean,
    },
    resolver: yupResolver(marketingSchema),
  });

  useEffect(() => {
    if (data?.data?.preferences?.communication) {
      reset({
        marketingEmails: data.data.preferences.communication.enable_marketing,
        productUpdates: data.data.preferences.communication.enable_product_updates,
      });
    }
  }, [data, reset]);

  const mutation = useMutation({
    mutationKey: [UPDATE_WORKSPACE_PREFERENCES],
    mutationFn: (data: CommunicationFormData) => {
      return client.workspaces.preferencesUpdate({
        preferences: {
          billing: {},
          newsletter: {
            enable_marketing: data.marketingEmails,
            enable_product_updates: data.productUpdates,
          }
        }
      })
    },
    onSuccess: ({ data }) => {
      toast.success(data.message);
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message || "An error occurred while updating workspace preferences");
    },
  });

  const onSubmit = (data: CommunicationFormData) => mutation.mutate(data)

  return (
    <Card>
      <form onSubmit={handleSubmit(onSubmit)}>
        <CardHeader>
          <CardTitle>Communication Settings</CardTitle>
          <CardDescription>Manage your communication preferences.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {isLoading ? <Skeleton count={5} /> : (
            <>
              <div className="flex items-center justify-between">
                <Label htmlFor="marketing-emails">Receive Marketing Emails?</Label>
                <Controller
                  name="marketingEmails"
                  control={control}
                  render={({ field }) => (
                    <Switch
                      id="marketing-emails"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  )}
                />
              </div>
              {errors.marketingEmails && (
                <p className="text-red-500 text-sm">{errors.marketingEmails.message}</p>
              )}

              <div className="flex items-center justify-between">
                <Label htmlFor="product-updates">Product Update Notifications?</Label>
                <Controller
                  name="productUpdates"
                  control={control}
                  render={({ field }) => (
                    <Switch
                      id="product-updates"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  )}
                />
              </div>
              {errors.productUpdates && (
                <p className="text-red-500 text-sm">{errors.productUpdates.message}</p>
              )}
            </>
          )}
        </CardContent>
        <CardFooter className="mt-6">
          <Button type="submit">Update Communication Settings</Button>
        </CardFooter>
      </form>
    </Card>
  );
}; 
