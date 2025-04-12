"use client";

import type { ServerCreatedUserResponse } from "@/client/Api";
import { Button } from "@/components/ui/button";
import client from "@/lib/client";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { Input } from "@/components/ui/input";
import {
  MALAK_PRIVACY_POLICY_LINK,
  MALAK_TERMS_CONDITION_LINK,
} from "@/lib/config";
import { AnalyticsEvent } from "@/lib/events";
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
import { useState } from "react";

const loginSchema = yup.object({
  email: yup.string().email("Invalid email").required("Email is required"),
  password: yup.string().required("Password is required"),
});

const signupSchema = yup.object({
  email: yup.string().email("Invalid email").required("Email is required"),
  firstName: yup.string().required("First name is required"),
  lastName: yup.string().required("Last name is required"),
  password: yup
    .string()
    .required("Password is required")
    .min(8, "Password must be at least 8 characters"),
});

type LoginFormData = yup.InferType<typeof loginSchema>;
type SignupFormData = yup.InferType<typeof signupSchema>;

export default function LoginPage() {
  const [isLogin, setIsLogin] = useState(true);
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
    posthog?.capture(AnalyticsEvent.LoginButtonClicked, {
      auth: "google",
    });
    googleLogin();
  };

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: yupResolver(loginSchema),
  });

  const {
    register: registerSignup,
    handleSubmit: handleSignupSubmit,
    formState: { errors: signupErrors, isSubmitting: isSignupSubmitting },
  } = useForm<SignupFormData>({
    resolver: yupResolver(signupSchema),
  });

  const emailPasswordLogin = async (data: LoginFormData) => {
    return new Promise((resolve) => {
      setTimeout(() => {
        setUser({ id: "1", email: data.email, name: "Test User" });
        setToken("fake-token");
        router.push("/");
        resolve(true);
      }, 1000);
    });
  };

  const handleSignup = async (data: SignupFormData) => {
    return new Promise((resolve) => {
      setTimeout(() => {
        setUser({
          id: "1",
          email: data.email,
          name: `${data.firstName} ${data.lastName}`,
        });
        setToken("fake-token");
        router.push("/");
        resolve(true);
      }, 1000);
    });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-violet-950/5 to-muted">
      <div className="container relative h-screen flex-col items-center justify-center grid lg:max-w-none lg:grid-cols-2 lg:px-0">
        <div className="relative hidden h-full flex-col bg-muted p-10 text-white lg:flex dark:border-r">
          <div className="absolute inset-0">
            <Image
              src="/overview-malak.png"
              alt="Malak Background"
              fill
              className="object-cover opacity-60 transition-opacity duration-500 hover:opacity-70"
              priority
              quality={100}
            />
            <div className="absolute inset-0 bg-gradient-to-br from-violet-950/40 via-background/50 to-background/70" />
            <div className="absolute inset-0 bg-gradient-to-t from-background/90 via-transparent" />
            <div className="absolute inset-0 bg-grid-white/[0.04]" />
          </div>
          <div className="relative z-20 flex items-center text-xl font-medium">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="mr-2 h-7 w-7"
            >
              <path d="M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3" />
            </svg>
            Malak
          </div>
          <div className="relative z-20 mt-auto">
            <blockquote className="space-y-3">
              <p className="text-2xl font-medium text-white drop-shadow-sm">
                Streamline your communication and stay organized with our
                powerful platform
              </p>
              <footer className="text-base text-white/90">
                Your all-in-one collaboration solution
              </footer>
            </blockquote>
          </div>
        </div>
        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-8 sm:w-[380px]">
            <div className="flex flex-col space-y-3 text-center">
              <h1 className="text-3xl font-semibold tracking-tight">
                {isLogin ? "Welcome back" : "Create an account"}
              </h1>
              <p className="text-sm text-muted-foreground">
                {isLogin ? (
                  <>
                    Don't have an account?{" "}
                    <button
                      onClick={() => setIsLogin(false)}
                      className="text-primary hover:underline"
                    >
                      Create one instead
                    </button>
                  </>
                ) : (
                  <>
                    Already have an account?{" "}
                    <button
                      onClick={() => setIsLogin(true)}
                      className="text-primary hover:underline"
                    >
                      Sign in
                    </button>
                  </>
                )}
              </p>
            </div>

            {isLogin ? (
              <form onSubmit={handleSubmit(emailPasswordLogin)} className="space-y-5">
                <div className="space-y-2">
                  <Input
                    type="email"
                    placeholder="Email"
                    className="h-11"
                    {...register("email")}
                  />
                  {errors.email && (
                    <p className="text-sm text-red-500 ml-1">{errors.email.message}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Input
                    type="password"
                    placeholder="Password"
                    className="h-11"
                    {...register("password")}
                  />
                  {errors.password && (
                    <p className="text-sm text-red-500 ml-1">{errors.password.message}</p>
                  )}
                </div>
                <Button
                  type="submit"
                  className="w-full h-11 text-base"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Signing in..." : "Sign in"}
                </Button>
              </form>
            ) : (
              <form onSubmit={handleSignupSubmit(handleSignup)} className="space-y-5">
                <div className="space-y-2">
                  <Input
                    type="email"
                    placeholder="Email"
                    className="h-11"
                    {...registerSignup("email")}
                  />
                  {signupErrors.email && (
                    <p className="text-sm text-red-500 ml-1">{signupErrors.email.message}</p>
                  )}
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Input
                      placeholder="First Name"
                      className="h-11"
                      {...registerSignup("firstName")}
                    />
                    {signupErrors.firstName && (
                      <p className="text-sm text-red-500 ml-1">{signupErrors.firstName.message}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Input
                      placeholder="Last Name"
                      className="h-11"
                      {...registerSignup("lastName")}
                    />
                    {signupErrors.lastName && (
                      <p className="text-sm text-red-500 ml-1">{signupErrors.lastName.message}</p>
                    )}
                  </div>
                </div>
                <div className="space-y-2">
                  <Input
                    type="password"
                    placeholder="Password"
                    className="h-11"
                    {...registerSignup("password")}
                  />
                  {signupErrors.password && (
                    <p className="text-sm text-red-500 ml-1">{signupErrors.password.message}</p>
                  )}
                </div>
                <Button
                  type="submit"
                  className="w-full h-11 text-base"
                  disabled={isSignupSubmitting}
                >
                  {isSignupSubmitting ? "Creating account..." : "Create account"}
                </Button>
              </form>
            )}

            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-muted-foreground/20" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-3 text-muted-foreground">
                  Or continue with
                </span>
              </div>
            </div>
            <div className="grid gap-4">
              <Button
                type="button"
                variant="outline"
                className="relative overflow-hidden h-11"
                onClick={loginWithGoogle}
              >
                <div className="flex items-center justify-center gap-3">
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
              <Button
                variant="outline"
                className="relative overflow-hidden h-11"
                disabled
              >
                <div className="flex items-center justify-center gap-3">
                  <RiMicrosoftLine className="h-5 w-5" />
                  <span>Continue with Microsoft</span>
                  <span className="absolute right-4 text-xs bg-muted px-2 py-0.5 rounded-full text-muted-foreground">
                    Soon
                  </span>
                </div>
              </Button>
            </div>
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                href={MALAK_TERMS_CONDITION_LINK}
                className="underline underline-offset-4 hover:text-primary transition-colors"
                target="_blank"
              >
                Terms of Service
              </Link>{" "}
              and{" "}
              <Link
                href={MALAK_PRIVACY_POLICY_LINK}
                className="underline underline-offset-4 hover:text-primary transition-colors"
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
