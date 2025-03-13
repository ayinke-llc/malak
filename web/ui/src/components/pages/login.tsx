"use client";

import type { ServerCreatedUserResponse } from "@/client/Api";
import { Button } from "@/components/ui/button";
import { EVENT_LOGIN_BUTTON_CLICKED } from "@/lib/analytics-constansts";
import client from "@/lib/client";
import {
  MALAK_PRIVACY_POLICY_LINK,
  MALAK_TERMS_CONDITION_LINK,
} from "@/lib/config";
import useAuthStore from "@/store/auth";
import { useGoogleLogin } from "@react-oauth/google";
import { RiMicrosoftLine } from "@remixicon/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosResponse } from "axios";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { usePostHog } from "posthog-js/react";
import { toast } from "sonner";

export default function LoginPage() {
  const router = useRouter();
  const posthog = usePostHog();

  const { setUser, setToken } = useAuthStore();

  const mutation = useMutation({
    mutationFn: ({ code }: { code: string }) => {
      return client.auth.connectCreate("google", {
        code,
      });
    },
    gcTime: 0,
    onError: (err: AxiosResponse<ServerCreatedUserResponse>): void => {
      toast.error(err.data.message);
    },
    onSuccess: (resp: AxiosResponse<ServerCreatedUserResponse>) => {
      posthog?.identify(resp.data.user.id);
      setToken(resp.data.token);
      setUser(resp.data.user);
      router.push("/");
    },
  });

  const googleLogin = useGoogleLogin({
    flow: "auth-code",
    onSuccess: async (codeResponse) => {
      mutation.mutate({ code: codeResponse.code });
    },
    onError: (errorResponse) => {
      toast.error(errorResponse?.error_description);
    },
  });

  const loginWithGoogle = () => {
    posthog?.capture(EVENT_LOGIN_BUTTON_CLICKED, {
      auth: "google",
    });
    googleLogin();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted">
      <div className="container relative h-screen flex-col items-center justify-center grid lg:max-w-none lg:grid-cols-2 lg:px-0">
        <div className="relative hidden h-full flex-col bg-muted p-10 text-white lg:flex dark:border-r">
          <div className="absolute inset-0">
            <Image
              src="/overview-malak.png"
              alt="Malak Background"
              fill
              className="object-cover opacity-50 transition-opacity duration-300 hover:opacity-60"
              priority
              quality={100}
            />
            <div className="absolute inset-0 bg-gradient-to-br from-violet-950/30 via-background/40 to-background/60" />
            <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-transparent" />
            <div className="absolute inset-0 bg-grid-white/[0.03]" />
          </div>
          <div className="relative z-20 flex items-center text-lg font-medium">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="mr-2 h-6 w-6"
            >
              <path d="M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3" />
            </svg>
            Malak
          </div>
          <div className="relative z-20 mt-auto">
            <blockquote className="space-y-2">
              <p className="text-lg font-medium text-white drop-shadow-sm">
                Streamline your communication and stay organized with our powerful platform
              </p>
              <footer className="text-sm text-white/80">
                Your all-in-one collaboration solution
              </footer>
            </blockquote>
          </div>
        </div>
        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div className="flex flex-col space-y-2 text-center">
              <h1 className="text-2xl font-semibold tracking-tight">Welcome back</h1>
              <p className="text-sm text-muted-foreground">
                Sign in to your account to continue
              </p>
            </div>
            <div className="grid gap-4">
              <Button
                type="button"
                className="relative overflow-hidden"
                onClick={loginWithGoogle}
              >
                <div className="flex items-center justify-center gap-2">
                  <svg
                    className="h-5 w-5"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                  >
                    <path
                      d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                      fill="#4285F4"
                    />
                    <path
                      d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                      fill="#34A853"
                    />
                    <path
                      d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                      fill="#FBBC05"
                    />
                    <path
                      d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                      fill="#EA4335"
                    />
                  </svg>
                  <span>Continue with Google</span>
                </div>
              </Button>
              <Button variant="outline" className="relative overflow-hidden" disabled>
                <div className="flex items-center justify-center gap-2">
                  <RiMicrosoftLine className="h-5 w-5" />
                  <span>Continue with Microsoft</span>
                  <span className="absolute right-4 text-xs text-muted-foreground">(Soon)</span>
                </div>
              </Button>
            </div>
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                href={MALAK_TERMS_CONDITION_LINK}
                className="underline underline-offset-4 hover:text-primary"
                target="_blank"
              >
                Terms of Service
              </Link>{" "}
              and{" "}
              <Link
                href={MALAK_PRIVACY_POLICY_LINK}
                className="underline underline-offset-4 hover:text-primary"
                target="_blank"
              >
                Privacy Policy
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
} 
