"use client"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select, SelectContent,
  SelectItem, SelectTrigger,
  SelectValue
} from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import timezoneMap from "@/lib/timezone"
import useWorkspacesStore from "@/store/workspace"
import { Controller, useForm } from "react-hook-form"
import * as yup from "yup"
import { yupResolver } from "@hookform/resolvers/yup"

export function GeneralSettings() {
  return (
    <div className="grid gap-6 md:grid-cols-2">

      <CompanyUpdateCard />

      <Card>
        <CardHeader>
          <CardTitle>Communication Settings</CardTitle>
          <CardDescription>Manage your communication preferences.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email-frequency">Email Frequency</Label>
            <Select defaultValue="weekly">
              <SelectTrigger id="email-frequency">
                <SelectValue placeholder="Select Email Frequency" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="daily">Daily</SelectItem>
                <SelectItem value="weekly">Weekly</SelectItem>
                <SelectItem value="monthly">Monthly</SelectItem>
                <SelectItem value="never">Never</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="newsletter">Newsletter Subscription</Label>
            <Input id="newsletter" type="email" placeholder="Enter your email" />
          </div>
          <div className="flex items-center justify-between">
            <Label htmlFor="marketing-emails">Receive Marketing Emails</Label>
            <Switch id="marketing-emails" />
          </div>
          <div className="flex items-center justify-between">
            <Label htmlFor="product-updates">Product Update Notifications</Label>
            <Switch id="product-updates" />
          </div>
        </CardContent>
        <CardFooter>
          <Button>Update Communication Settings</Button>
        </CardFooter>
      </Card>
    </div>
  )
}


const schema = yup.object({
  companyName: yup.string().required("Company name is required"),
  website: yup.string().url("Invalid URL").optional(),
  timezone: yup.string().required("Timezone is required"),
});

type FormData = yup.InferType<typeof schema>;

const CompanyUpdateCard = () => {
  const current = useWorkspacesStore((state) => state.current);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: {
      companyName: current?.workspace_name || "",
      website: current?.website || "",
      timezone: current?.timezone || "UTC",
    },
  });

  const onSubmit = (data: FormData) => {
    console.log("Form data submitted:", data);
  };

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