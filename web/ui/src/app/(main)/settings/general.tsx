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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import timezoneMap from "@/lib/timezone"
import useWorkspacesStore from "@/store/workspace"

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

const CompanyUpdateCard = () => {

  const current = useWorkspacesStore(state => state.current)

  const getTimezoneOffset = (timezone: string) => {
    const date = new Date();
    const formatter = new Intl.DateTimeFormat("en-US", {
      timeZone: timezone,
      timeZoneName: "shortOffset",
    });
    const parts = formatter.formatToParts(date);

    const offsetPart = parts.find(part => part.type === "timeZoneName");
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
        <div className="space-y-2">
          <Label htmlFor="name">Company name:</Label>
          <Input
            id="name"
            placeholder={current?.workspace_name}
            className="mt-2"
          />
        </div>

        <div className="flex space-x-4">
          <div className="flex-1 space-y-2">
            <Label htmlFor="website">Website</Label>
            <Input
              id="website"
              placeholder={"https://google.com"}
              className="mt-2"
            />
          </div>

          <div className="flex-1 space-y-2">
            <Label htmlFor="timezone">Timezone</Label>
            <Select defaultValue="GMT">
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
          </div>

        </div>

      </CardContent>
      <CardFooter>
        <div className="space-x-3">
          <Button>Save Preferences</Button>
          <Button variant={"destructive"}>Delete Workspace</Button>
        </div>
      </CardFooter>
    </Card>
  )
}

