import Next from "next";

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      NEXT_PUBLIC_GOOGLE_CLIENT_ID: string;
      NEXT_PUBLIC_MALAK_TERMS_CONDITION_LINK: string
      NEXT_PUBLIC_MALAK_PRIVACY_POLICY_LINK: string
      // use posthog to track analytics or not
      NEXT_PUBLIC_ENABLE_POSTHOG: boolean
    }
  }
}
