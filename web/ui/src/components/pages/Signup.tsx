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

export const SignupForm = ({
  className,
  ...props
}: React.ComponentProps<"div">) => {
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
                    <form>
                      <FieldGroup>
                        <Field>
                          <FieldLabel htmlFor="name">Full Name</FieldLabel>
                          <Input id="name" type="text" placeholder="Lanre Adelowo" required />
                        </Field>
                        <Field>
                          <FieldLabel htmlFor="email">Email</FieldLabel>
                          <Input
                            id="email"
                            type="email"
                            placeholder="lanre@malak.vc"
                            required
                          />
                        </Field>
                        <Field>
                          <Field>
                            <FieldLabel htmlFor="password">Password</FieldLabel>
                            <Input id="password" type="password" required />
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
