import Next from "next";

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      NEXT_PUBLIC_GOOGLE_CLIENT_ID: string;
      NEXT_PUBLIC_MALAK_TERMS_CONDITION_LINK: string;
      NEXT_PUBLIC_MALAK_PRIVACY_POLICY_LINK: string;
      // use posthog to track analytics or not
      NEXT_PUBLIC_MALAK_ENABLE_POSTHOG: boolean;
      NEXT_PUBLIC_MALAK_POSTHOG_KEY: string;
      NEXT_PUBLIC_MALAK_POSTHOG_HOST: string;
      NEXT_PUBLIC_SENTRY_DSN?: string;
      NEXT_PUBLIC_DECKS_DOMAIN: string,
      NEXT_PUBLIC_SUPPORT_EMAIL: string,

      // Integrations
      NEXT_PUBLIC_INTEGRATION_STRIPE_CLIENT_ID: string
    }
  }
}
