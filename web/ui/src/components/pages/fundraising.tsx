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
            <p className="text-sm text-muted-foreground max-w-2xl mx-auto mb-2">
              Our team is working hard to bring you these powerful features. Stay tuned for updates!
            </p>
            <p className="text-sm font-medium text-primary mb-12">
              Expected Release Date: March 31st, 2025
            </p>
          </div>

          {/* Preview Image */}
          <div className="mt-12 mb-16 rounded-xl overflow-hidden border shadow-lg">
            <div className="relative aspect-[16/9] w-full">
              <img
                src="/fundraising-crm.png"
                alt="Fundraising CRM Preview"
                className="object-cover w-full h-full"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-background/20 to-transparent" />
            </div>
            <div className="bg-card/50 backdrop-blur supports-[backdrop-filter]:bg-background/60 p-4 text-center space-y-1">
              <p className="text-sm text-muted-foreground">
                Preview of the upcoming Investor Pipeline feature
              </p>
              <p className="text-xs font-medium text-primary">
                Available March 31st, 2025
              </p>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
} 