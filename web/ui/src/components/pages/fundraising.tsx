"use client";

import { Button } from "@/components/ui/button";
import {
  RiUserLine,
  RiFileListLine,
  RiBarChartBoxLine,
  RiCalendarTodoLine,
  RiArrowLeftLine,
  RiTimeLine
} from "@remixicon/react";
import { useRouter } from "next/navigation";

export default function Fundraising() {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <div className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-background z-0" />

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-24 z-10">
          <div className="text-center">
            <div className="inline-flex items-center justify-center gap-2 mb-6 px-4 py-2 rounded-full bg-primary/10 border border-primary/20 shadow-sm">
              <RiTimeLine className="w-5 h-5 text-primary animate-[pulse_2s_ease-in-out_infinite]" />
              <span className="text-sm font-semibold text-primary tracking-wide uppercase">Coming Soon</span>
            </div>
            <h1 className="text-5xl font-bold tracking-tight bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent mb-6">
              Fundraising Features
            </h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-4">
              We&apos;re building powerful fundraising tools to help you manage and track
              your fundraising efforts more effectively.
            </p>
            <p className="text-sm text-muted-foreground max-w-2xl mx-auto mb-12">
              Our team is working hard to bring you these powerful features. Stay tuned for updates!
            </p>
          </div>

          {/* Feature Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-16">
            <div className="group hover:scale-[1.02] transition-transform duration-200 ease-in-out">
              <div className="relative p-8 rounded-xl border bg-card/50 backdrop-blur supports-[backdrop-filter]:bg-background/60 hover:shadow-lg transition-all h-full overflow-hidden">
                <div className="absolute -top-3 -right-3 w-24 h-24 rotate-12">
                  <div className="absolute inset-0 bg-primary/10" />
                  <div className="absolute bottom-7 right-7 text-[10px] font-semibold tracking-wider text-primary uppercase transform -rotate-12">
                    Coming Soon
                  </div>
                </div>
                <div className="rounded-lg bg-primary/10 p-3 w-fit mb-6">
                  <RiUserLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Investor Management</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Track potential investors, manage communications, and maintain detailed
                  records of all fundraising interactions.
                </p>
              </div>
            </div>

            <div className="group hover:scale-[1.02] transition-transform duration-200 ease-in-out">
              <div className="relative p-8 rounded-xl border bg-card/50 backdrop-blur supports-[backdrop-filter]:bg-background/60 hover:shadow-lg transition-all h-full overflow-hidden">
                <div className="absolute -top-3 -right-3 w-24 h-24 rotate-12">
                  <div className="absolute inset-0 bg-primary/10" />
                  <div className="absolute bottom-7 right-7 text-[10px] font-semibold tracking-wider text-primary uppercase transform -rotate-12">
                    Coming Soon
                  </div>
                </div>
                <div className="rounded-lg bg-primary/10 p-3 w-fit mb-6">
                  <RiFileListLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Deal Room</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Create secure virtual data rooms for sharing sensitive documents
                  with investors and tracking engagement.
                </p>
              </div>
            </div>

            <div className="group hover:scale-[1.02] transition-transform duration-200 ease-in-out">
              <div className="relative p-8 rounded-xl border bg-card/50 backdrop-blur supports-[backdrop-filter]:bg-background/60 hover:shadow-lg transition-all h-full overflow-hidden">
                <div className="absolute -top-3 -right-3 w-24 h-24 rotate-12">
                  <div className="absolute inset-0 bg-primary/10" />
                  <div className="absolute bottom-7 right-7 text-[10px] font-semibold tracking-wider text-primary uppercase transform -rotate-12">
                    Coming Soon
                  </div>
                </div>
                <div className="rounded-lg bg-primary/10 p-3 w-fit mb-6">
                  <RiBarChartBoxLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Fundraising Analytics</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Get insights into your fundraising pipeline, conversion rates,
                  and investor engagement metrics.
                </p>
              </div>
            </div>

            <div className="group hover:scale-[1.02] transition-transform duration-200 ease-in-out">
              <div className="relative p-8 rounded-xl border bg-card/50 backdrop-blur supports-[backdrop-filter]:bg-background/60 hover:shadow-lg transition-all h-full overflow-hidden">
                <div className="absolute -top-3 -right-3 w-24 h-24 rotate-12">
                  <div className="absolute inset-0 bg-primary/10" />
                  <div className="absolute bottom-7 right-7 text-[10px] font-semibold tracking-wider text-primary uppercase transform -rotate-12">
                    Coming Soon
                  </div>
                </div>
                <div className="rounded-lg bg-primary/10 p-3 w-fit mb-6">
                  <RiCalendarTodoLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Campaign Management</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Plan and execute fundraising campaigns with tools for goal setting,
                  progress tracking, and investor communications.
                </p>
              </div>
            </div>
          </div>

          {/* CTA Section */}
          <div className="mt-16 text-center">
            <Button
              variant="outline"
              onClick={() => router.back()}
              className="group text-base"
            >
              <RiArrowLeftLine className="mr-2 h-4 w-4 transition-transform group-hover:-translate-x-1" />
              Go Back
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
} 