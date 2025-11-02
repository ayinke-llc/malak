"use client"

import useAuthStore from "@/store/auth"
import { Alert, AlertDescription, AlertTitle } from "./ui/alert"
import { MdMarkEmailUnread } from "react-icons/md"
import { useMutation } from "@tanstack/react-query"
import { toast } from "sonner"
import { Button } from "./ui/button"
import { useState } from "react"
import client from "@/lib/client"
import { ServerAPIStatus } from "@/client/Api"
import { AxiosError } from "axios"


const EmailVerificationBadge = () => {
  const [resendDisabled, setResendDisabled] = useState(false)
  const verification_date = useAuthStore(state => state.user?.email_verified_at)

  const resendMutation = useMutation({
    mutationKey: ['resend-verification-email'],
    mutationFn: async () => {
      return client.user.resendVerificationCreate()
    },
    onSuccess: () => {
      toast.success("verification email sent successfully!")
      setResendDisabled(true)

      setTimeout(() => {
        setResendDisabled(false)
      }, 60000)
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      toast.error(err.response?.data.message ?? "An error occurred while sending verification email");
    },
  })

  if (verification_date && verification_date !== "") {
    return null
  }

  return (
    <div className="grid w-full max-w-xl items-start gap-4">
      <Alert variant={"destructive"}>
        <MdMarkEmailUnread className="h-4 w-4" />
        <AlertTitle>Email verification required</AlertTitle>
        <AlertDescription className="flex items-center justify-between gap-4">
          <span>Please verify your email address to access all features.</span>
          <Button
            variant="default"
            size="sm"
            onClick={() => resendMutation.mutate()}
            disabled={resendMutation.isPending || resendDisabled}
          >
            {resendMutation.isPending ? "Sending..." : resendDisabled ? "Sent" : "Resend Email"}
          </Button>
        </AlertDescription>
      </Alert>
    </div>
  )
}

export default EmailVerificationBadge;
