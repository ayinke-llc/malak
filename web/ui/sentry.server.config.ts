// This file configures the initialization of Sentry on the server.
// The config you add here will be used whenever the server handles a request.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

import { MALAK_SENTRY_DSN, MALAK_SENTRY_ENABLED } from "@/lib/config";
import * as Sentry from "@sentry/nextjs";

if (MALAK_SENTRY_ENABLED) {
  Sentry.init({
    dsn: MALAK_SENTRY_DSN,
    tracesSampleRate: 1,
    debug: false,
  });
}
