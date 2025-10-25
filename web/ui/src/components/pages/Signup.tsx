"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import Link from "next/link"
import Image from "next/image"
import { ServerCreatedUserResponse } from "@/client/Api"
import client from "@/lib/client"
import useAuthStore from "@/store/auth"
import { useMutation } from "@tanstack/react-query"
import { AxiosResponse } from "axios"
import { useRouter } from "next/navigation"
import { usePostHog } from "posthog-js/react"
import { toast } from "sonner"
import { AnalyticsEvent } from "@/lib/events"
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup"
import { useForm, Controller } from "react-hook-form"
import {
  Form,
  FormControl,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

const signupSchema = yup.object().shape({
  full_name: yup.string()
    .required("Name is required")
    .min(3, "Name must be at least 3 characters")
    .max(30, "Name must not exceed 30 characters"),

  email: yup.string().email().required("Email address is required"),
  password: yup.string().required().min(8),
});

type SignupFormData = yup.InferType<typeof signupSchema>;

export const SignupForm = ({
  className,
  ...props
}: React.ComponentProps<"div">) => {

  const router = useRouter();
  const posthog = usePostHog();

  const { setUser, setToken } = useAuthStore();

  const {
    control,
    handleSubmit,
    watch,
    reset,
    formState: { errors } } = useForm<SignupFormData>({
      resolver: yupResolver(signupSchema) as any,
    });

  const mutation = useMutation({
    mutationFn: (data: SignupFormData) => {
      return client.auth.registerCreate(data);
    },
    gcTime: 0,
    onError: (err: AxiosResponse<ServerCreatedUserResponse>): void => {
      toast.error(err.data.message);
    },
    onSuccess: (resp: AxiosResponse<ServerCreatedUserResponse>) => {
      reset()
      posthog?.identify(resp.data.user.id);
      setToken(resp.data.token);
      setUser(resp.data.user);
      router.push("/");
    },
  });

  const signupHandler = (data: SignupFormData) => {
    posthog?.capture(AnalyticsEvent.SignupButtonClicked, {});
    mutation.mutate(data)
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
                Sign into or create your account to continue
              </p>
            </div>

            <div className="grid gap-4">
              <div className={cn("flex flex-col gap-6", className)} {...props}>
                <Card>
                  <CardHeader className="text-center">
                    <CardTitle className="text-xl">Create your account</CardTitle>
                    <CardDescription>
                      Enter your email below to create your account
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <form onSubmit={handleSubmit(signupHandler)}>
                      <FieldGroup>
                        <Field>
                          <Controller
                            control={control}
                            render={({ field }) => {
                              return (
                                <>
                                  <FieldLabel htmlFor="full_name">Full Name</FieldLabel>
                                  <Input {...field} id="full_name" type="text" placeholder="Lanre Adelowo" required />

                                  {errors.full_name && (
                                    <p className="text-sm text-red-500">{errors.full_name.message}</p>
                                  )}
                                </>
                              )
                            }}
                            name="full_name"
                          />
                        </Field>
                        <Field>
                          <FieldLabel htmlFor="email">Email</FieldLabel>
                          <FormControl>
                            <Input
                              id="email"
                              type="email"
                              placeholder="lanre@malak.vc"
                              required
                            />
                          </FormControl>
                        </Field>
                        <Field>
                          <Field>
                            <FieldLabel htmlFor="password">Password</FieldLabel>
                            <FormControl>
                              <Input id="password" type="password" required />
                            </FormControl>
                          </Field>
                          <FieldDescription>
                            Must be at least 8 characters long.
                          </FieldDescription>
                        </Field>
                        <Field>
                          <Button type="submit">Create Account</Button>
                          <FieldDescription className="text-center">
                            Already have an account? <Link href="/login">Sign in</Link>
                          </FieldDescription>
                        </Field>
                      </FieldGroup>
                    </form>
                  </CardContent>
                </Card>
                <FieldDescription className="px-6 text-center">
                  By clicking continue, you agree to our <a href="#">Terms of Service</a>{" "}
                  and <a href="#">Privacy Policy</a>.
                </FieldDescription>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
