import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { AccountSettings } from "./account";
import { GeneralSettings } from "./general";
import { NotificationSettings } from "./notification";
import { PrivacySettings } from "./privacy";

export default function Settings() {

  return (
    <>
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3
                id="company-settings"
                className="text-lg font-medium"
              >
                Company Preferences
              </h3>
              <p className="text-sm text-muted-foreground">
                View and manage your company's preferences
              </p>
            </div>

            <div>
            </div>
          </div>
        </section>

        <section className="mt-10">
          <Tabs defaultValue="general" className="space-y-6">
            <TabsList className="w-full justify-start border-b pb-px mb-4">
              <TabsTrigger value="general">General</TabsTrigger>
              <TabsTrigger value="account">Account</TabsTrigger>
              <TabsTrigger value="notifications">Notifications</TabsTrigger>
              <TabsTrigger value="privacy">Privacy</TabsTrigger>
              <TabsTrigger value="api">API Key & Webhooks</TabsTrigger>
            </TabsList>
            <TabsContent value="general">
              <GeneralSettings />
            </TabsContent>
            <TabsContent value="account">
              <AccountSettings />
            </TabsContent>
            <TabsContent value="notifications">
              <NotificationSettings />
            </TabsContent>
            <TabsContent value="privacy">
              <PrivacySettings />
            </TabsContent>
            <TabsContent value="api">
              <PrivacySettings />
            </TabsContent>
          </Tabs>
        </section>
      </div>
    </>
  );
}
