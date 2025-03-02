"use client";

import { Button } from "@/components/ui/button";
import {
  RiPieChartLine,
  RiBarChartBoxLine,
  RiTeamLine,
  RiExchangeDollarLine,
  RiArrowLeftLine,
  RiShareBoxLine,
  RiMessage2Line,
  RiTimeLine
} from "@remixicon/react";
import { useRouter } from "next/navigation";

export default function Page() {
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
              Cap Table Management
            </h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-4">
              We&apos;re building comprehensive cap table management tools to help you track ownership,
              manage equity, and maintain accurate shareholder records with precision.
            </p>
            <p className="text-sm text-muted-foreground max-w-2xl mx-auto mb-12">
              Our team is working hard to bring you these powerful features. Stay tuned for updates!
            </p>
          </div>

          {/* Feature Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mt-16">
            <div className="group hover:scale-[1.02] transition-transform duration-200 ease-in-out">
              <div className="relative p-8 rounded-xl border bg-card/50 backdrop-blur supports-[backdrop-filter]:bg-background/60 hover:shadow-lg transition-all h-full overflow-hidden">
                <div className="absolute -top-3 -right-3 w-24 h-24 rotate-12">
                  <div className="absolute inset-0 bg-primary/10" />
                  <div className="absolute bottom-7 right-7 text-[10px] font-semibold tracking-wider text-primary uppercase transform -rotate-12">
                    Coming Soon
                  </div>
                </div>
                <div className="rounded-lg bg-primary/10 p-3 w-fit mb-6">
                  <RiPieChartLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Equity Tracking</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Monitor and manage all forms of equity including common stock, preferred shares,
                  options, and convertible instruments with our intuitive interface.
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
                <h3 className="text-2xl font-semibold mb-4">Ownership Analytics</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Visualize ownership structure, dilution impacts, and waterfall analysis
                  for various exit scenarios with powerful analytics tools.
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
                  <RiTeamLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Stakeholder Management</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Maintain detailed records of shareholders, option holders, and their respective
                  holdings with full historical tracking and reporting capabilities.
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
                  <RiExchangeDollarLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Round Modeling</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Model future funding rounds, simulate dilution scenarios, and understand
                  the impact of new investments on existing shareholders in real-time.
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
                  <RiShareBoxLine className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Smart Sharing</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Share cap table insights and analytics with stakeholders securely.
                  Control access levels and customize what information each party can view.
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
                  <RiMessage2Line className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-2xl font-semibold mb-4">Shareholder Communications</h3>
                <p className="text-muted-foreground leading-relaxed">
                  Keep stakeholders informed with secure document storage, automated updates,
                  and streamlined communication channels for important company matters.
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
