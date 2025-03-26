import { cn } from "@/lib/utils";

export function BetaBanner() {
  return (
    <div className="relative isolate flex items-center gap-x-6 overflow-hidden bg-gray-50 px-6 py-2.5 sm:px-3.5 sm:before:flex-1 dark:bg-gray-900">
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2">
        <p className="text-sm leading-6 text-gray-900 dark:text-gray-100">
          <strong className="font-semibold">Malak is in beta</strong>
          <svg viewBox="0 0 2 2" className="mx-2 inline h-0.5 w-0.5 fill-current" aria-hidden="true">
            <circle cx="1" cy="1" r="1" />
          </svg>
          We're actively developing and improving the platform.
        </p>
      </div>
      <div className="flex flex-1 justify-end">
        <a
          href="https://github.com/ayinke-llc/malak"
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm font-semibold leading-6 text-gray-900 dark:text-gray-100 hover:text-gray-700 dark:hover:text-gray-300"
        >
          Learn more <span aria-hidden="true">&rarr;</span>
        </a>
      </div>
    </div>
  );
} 