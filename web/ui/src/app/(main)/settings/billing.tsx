"use client"

import {
  Card,
  CardContent, CardDescription,
  CardHeader, CardTitle
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { SUPPORT_EMAIL } from "@/lib/config"
import { useMutation } from "@tanstack/react-query"
import client from "@/lib/client"
import { useRouter } from "next/navigation"
import { FETCH_BILLING_PORTAL_URL } from "@/lib/query-constants"
import { ServerAPIStatus } from "@/client/Api"
import { AxiosError } from "axios"
import { toast } from "sonner"
import useWorkspacesStore from "@/store/workspace"
import { RiAlertLine, RiExternalLinkLine, RiMailSendLine } from "@remixicon/react"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"

export function BillingPage() {

  const router = useRouter()
  const current = useWorkspacesStore(state => state.current)

  const mutation = useMutation({
    mutationKey: [FETCH_BILLING_PORTAL_URL],
    mutationFn: () => client.workspaces.billingCreate(),
    onSuccess: ({ data }) => {
      router.push(data.link)
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred while creating Stripe billing link");
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  const handleBillingPortal = () => mutation.mutate()

  return (
    <div className="flex justify-start">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <CardTitle>Billing</CardTitle>
          <CardDescription>Manage your subscription and get support</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {!current?.is_subscription_active && (
            <Alert variant="destructive">
              <RiAlertLine className="h-4 w-4" />
              <AlertTitle>No active subscription</AlertTitle>
              <AlertDescription>
                You currently don't have an active subscription. Subscribe to get access to Malak
              </AlertDescription>
            </Alert>
          )}

          <div className="space-y-4">
            <div className="flex flex-col space-y-2">
              <h3 className="text-lg font-semibold">Current Plan</h3>
              <p className="text-sm text-muted-foreground">
                {current?.is_subscription_active ? (
                  <p>You are currently on the
                    <span className="font-semibold"> {current?.plan?.plan_name}</span>plan.</p>
                ) : (
                  <span>You don't have an active subscription.</span>
                )}
              </p>
            </div>
            <Separator className="mb-5" />
            <div className="flex flex-col space-y-2">
              <h3 className="text-lg font-semibold">Billing Portal</h3>
              <p className="text-sm text-muted-foreground">View your invoices and manage your payment method.</p>
              <Button onClick={handleBillingPortal}
                variant="outline"
                loading={mutation.isPending}
                className="w-full sm:w-auto">
                Go to Billing Portal
                <RiExternalLinkLine className="ml-2 h-4 w-4" />
              </Button>
            </div>
            <Separator className="mb-5" />
            <div className="flex flex-col space-y-2">
              <h3 className="text-lg font-semibold">Need Help?</h3>
              <p className="text-sm text-muted-foreground">
                If you need any further help with billing, our support team are here to help.
              </p>
              <Button variant="secondary" className="w-full sm:w-auto">
                <a href={`mailto:${SUPPORT_EMAIL}`}>
                  <span className="flex items-center">
                    Contact Support
                    <RiMailSendLine className="ml-2 h-4 w-4" />
                  </span>
                </a>
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div >
  )
}
