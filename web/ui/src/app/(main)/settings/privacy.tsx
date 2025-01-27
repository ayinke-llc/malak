import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"

export function PrivacySettings() {
  return (
    <div className="flex justify-center">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <CardTitle>Privacy Settings</CardTitle>
          <CardDescription>Manage your privacy and data sharing preferences.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="profile-visibility" className="text-base">
                  Public Profile Visibility
                </Label>
                <p className="text-sm text-muted-foreground">Allow others to see your public profile</p>
              </div>
              <Switch id="profile-visibility" />
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="data-collection" className="text-base">
                  Data Collection
                </Label>
                <p className="text-sm text-muted-foreground">Allow us to collect data to improve our services</p>
              </div>
              <Switch id="data-collection" />
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="third-party-sharing" className="text-base">
                  Third-Party Data Sharing
                </Label>
                <p className="text-sm text-muted-foreground">Allow sharing of your data with trusted partners</p>
              </div>
              <Switch id="third-party-sharing" />
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="personalized-ads" className="text-base">
                  Personalized Ads
                </Label>
                <p className="text-sm text-muted-foreground">Allow us to show you personalized advertisements</p>
              </div>
              <Switch id="personalized-ads" />
            </div>
          </div>
        </CardContent>
        <CardFooter>
          <Button className="w-full">Save Privacy Settings</Button>
        </CardFooter>
      </Card>
    </div>
  )
}

