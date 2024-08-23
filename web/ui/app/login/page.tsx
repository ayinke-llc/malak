"use client"

import Image from "next/image"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { useGoogleLogin } from "@react-oauth/google"
import { MALAK_PRIVACY_POLICY_LINK, MALAK_TERMS_CONDITION_LINK } from "@/lib/config"
import client from "@/lib/client"
import { useMutation } from "@tanstack/react-query"
import { useToast } from "@/components/ui/use-toast"
import { HttpResponse, ServerAPIStatus, ServerCreatedUserResponse } from "@/client/Api"
import { redirect } from "next/navigation"

export default function Login() {

  const { toast } = useToast();

  const mutation = useMutation({
    mutationFn: ({ code }: { code: string }) => {
      return client.auth.connectCreate("google", {
        code,
      })
    },
    gcTime: 0,
    onError: (err: HttpResponse<Response, ServerAPIStatus>): void => {
      toast({
        variant: "destructive",
        title: err.error.message,
      })
    },
    onSuccess: (resp: HttpResponse<ServerCreatedUserResponse>) => {
      redirect("/")
    }
  })

  const googleLogin = useGoogleLogin({
    flow: 'auth-code',
    onSuccess: async (codeResponse) => {
      mutation.mutate({ code: codeResponse.code })
    },
    onError: errorResponse => console.log(errorResponse),
  });

  const loginWithGoogle = () => {
    googleLogin()
  }

  return (
    <div className="w-screen h-screen lg:grid lg:grid-cols-2">
      <div className="flex items-center justify-center py-12">
        <div className="mx-auto grid w-[350px] gap-6">
          <div className="grid gap-2 text-center">
            <h1 className="text-3xl font-bold">Login</h1>
            <p className="text-balance text-muted-foreground">
              Use either your Google or Github account to authenticate into your dashboard
            </p>
          </div>
          <div className="grid gap-4">
            <Button type="button" className="w-full" onClick={loginWithGoogle}>
              <svg className="w-4 h-4 me-2" xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" /><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" /><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" /><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" /><path d="M1 1h22v22H1z" fill="none" /></svg>
              Login with Google
            </Button>
            <Button variant="secondary" className="w-full" disabled={true}>
              <svg className="w-4 h-4 me-2" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 .333A9.911 9.911 0 0 0 6.866 19.65c.5.092.678-.215.678-.477 0-.237-.01-1.017-.014-1.845-2.757.6-3.338-1.169-3.338-1.169a2.627 2.627 0 0 0-1.1-1.451c-.9-.615.07-.6.07-.6a2.084 2.084 0 0 1 1.518 1.021 2.11 2.11 0 0 0 2.884.823c.044-.503.268-.973.63-1.325-2.2-.25-4.516-1.1-4.516-4.9A3.832 3.832 0 0 1 4.7 7.068a3.56 3.56 0 0 1 .095-2.623s.832-.266 2.726 1.016a9.409 9.409 0 0 1 4.962 0c1.89-1.282 2.717-1.016 2.717-1.016.366.83.402 1.768.1 2.623a3.827 3.827 0 0 1 1.02 2.659c0 3.807-2.319 4.644-4.525 4.889a2.366 2.366 0 0 1 .673 1.834c0 1.326-.012 2.394-.012 2.72 0 .263.18.572.681.475A9.911 9.911 0 0 0 10 .333Z" clipRule="evenodd" />
              </svg>
              Login with Github ( coming soon )
            </Button>
          </div>
          <div className="mt-4 text-center text-sm">
            By clicking continue, you agree to our
            {" "}<Link className="underline" href={MALAK_TERMS_CONDITION_LINK} target="_blank">Terms of service</Link> and
            <Link className="underline" href={MALAK_PRIVACY_POLICY_LINK} target="_blank"> Privacy Policy</Link>
          </div>
        </div>
      </div>
      <div className="hidden bg-muted lg:block">
        <Image
          src="/placeholder.svg"
          alt="Image"
          width="1920"
          height="1080"
          className="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
      </div>
    </div>
  )
}
