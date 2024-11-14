import { siteConfig } from "@/app/siteConfig";
import { Button } from "@/components/ui/button";
import { ArrowAnimated } from "@/components/ui/icons/ArrowAnimated";
import { TremorPlaceholder } from "@/components/ui/icons/TremorPlaceholder";

export default function Settings() {
  return (
    <>
      <div className="mt-4 sm:mt-6 lg:mt-10">
        <div className="my-40 flex w-full flex-col items-center justify-center">
          <TremorPlaceholder className="size-20 shrink-0" aria-hidden="true" />
          <h2 className="mt-6 text-lg font-semibold sm:text-xl">
            This feature is coming soon
          </h2>
          <p className="mt-3 max-w-md text-center text-gray-500">
            Coming soon. You will be able to configure your workspace settings
            here
          </p>
          <Button className="group mt-6" variant="secondary" asChild>
            <a href={siteConfig.externalLink.blocks}>
              Contact us
              <ArrowAnimated
                className="stroke-gray-900 dark:stroke-gray-50"
                aria-hidden="true"
              />
            </a>
          </Button>
        </div>
      </div>
    </>
  );
}
