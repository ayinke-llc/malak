import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  RiEyeLine,
  RiUserLine,
  RiTimeLine,
  RiDownloadLine,
  RiBarChartLine,
  RiMapPinLine,
  RiUserHeartLine,
  RiFilterLine,
  RiArrowLeftLine,
  RiArrowRightLine,
  RiVipCrownFill,
  RiErrorWarningLine,
} from "@remixicon/react";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip as RechartsTooltip } from "recharts";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { ChartContainer, ChartLegend, ChartLegendContent } from "@/components/ui/chart";
import { useQuery } from "@tanstack/react-query";
import client from "@/lib/client";
import { LIST_DECK_SESSIONS } from "@/lib/query-constants";
import { MalakDeckViewerSession, ServerFetchSessionsDeck } from "@/client/Api";

// Feature flag for pro features
const HAS_PRO_ACCESS = false; // Set this based on user's subscription status

// Time range options with pro status
const TIME_RANGE_OPTIONS = [
  { value: "7", label: "Last 7 days", requiresPro: false },
  { value: "14", label: "Last 14 days", requiresPro: true },
  { value: "30", label: "Last 30 days", requiresPro: true },
  { value: "90", label: "Last 90 days", requiresPro: true },
] as const;

interface DeckAnalyticsProps {
  reference: string;
}

interface EngagementTrend {
  date: string;
  views: number;
}

interface GeographicDistribution {
  country: string;
  views: number;
}

interface MalakContact {
  reference: string;
  name?: string;
}

// Add the time formatting function
function formatTimeSpent(seconds: number): string {
  if (!seconds) return "00:00:00";

  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;

  return [hours, minutes, remainingSeconds]
    .map(val => val.toString().padStart(2, '0'))
    .join(':');
}

// Extend ServerFetchSessionsDeck meta to include our analytics data
declare module "@/client/Api" {
  interface ServerMeta {
    total_views: number;
    unique_views: number;
    average_time_spent: string;
    downloads: number;
    max_views: number;
    engagement_trends: EngagementTrend[];
    geographic_distribution: GeographicDistribution[];
  }

  interface MalakDeckViewerSession {
    contact?: MalakContact;
  }
}

// Add these components at the top level, before the DeckAnalytics component
function LoadingState() {
  return (
    <div className="space-y-6 animate-pulse">
      <Card className="p-6">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i}>
              <div className="h-6 w-24 bg-muted rounded mb-2" />
              <div className="h-8 w-16 bg-muted rounded" />
            </div>
          ))}
        </div>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {[...Array(2)].map((_, i) => (
          <Card key={i} className="p-6">
            <div className="h-6 w-32 bg-muted rounded mb-4" />
            <div className="h-[200px] bg-muted rounded" />
          </Card>
        ))}
      </div>

      <Card className="p-6">
        <div className="h-6 w-32 bg-muted rounded mb-6" />
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-12 bg-muted rounded" />
          ))}
        </div>
      </Card>
    </div>
  );
}

function ErrorState({ error }: { error: Error }) {
  return (
    <Card className="p-8">
      <div className="flex flex-col items-center justify-center text-center space-y-4">
        <RiErrorWarningLine className="h-12 w-12 text-destructive" />
        <div className="space-y-2">
          <h3 className="text-lg font-medium">Failed to load analytics</h3>
          <p className="text-sm text-muted-foreground max-w-md">
            {error.message || "An unexpected error occurred while loading the analytics data. Please try again later."}
          </p>
        </div>
        <Button
          variant="outline"
          onClick={() => window.location.reload()}
          className="mt-4"
        >
          Try again
        </Button>
      </div>
    </Card>
  );
}

function NoSessionsState() {
  return (
    <Card className="p-8">
      <div className="flex flex-col items-center justify-center text-center space-y-4">
        <RiUserHeartLine className="h-12 w-12 text-muted-foreground" />
        <div className="space-y-2">
          <h3 className="text-lg font-medium">No viewer sessions yet</h3>
          <p className="text-sm text-muted-foreground max-w-md">
            When people view your deck, their session data will appear here.
          </p>
        </div>
      </div>
    </Card>
  );
}

export default function DeckAnalytics({ reference }: DeckAnalyticsProps) {
  const router = useRouter();
  const [timeFilter, setTimeFilter] = useState("7"); // Default to 7 days for non-pro users
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  const { data: sessionsData, isLoading, error } = useQuery<ServerFetchSessionsDeck>({
    queryKey: [LIST_DECK_SESSIONS, reference, timeFilter, currentPage],
    queryFn: () => client.decks.sessionsDetail(
      reference,
      {
        page: currentPage,
        per_page: itemsPerPage,
        days: parseInt(timeFilter),
      }
    ).then(res => res.data),
  });

  if (isLoading) {
    return <LoadingState />;
  }

  if (error) {
    return <ErrorState error={error as Error} />;
  }

  const handleTimeFilterChange = (value: string) => {
    setTimeFilter(value);
    setCurrentPage(1);
  };

  const handleUpgradeClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    router.push("/settings?tab=billing");
  };

  const sessions = sessionsData?.sessions || [];
  const totalSessions = sessionsData?.meta?.paging?.total || 0;
  const totalPages = Math.ceil(totalSessions / itemsPerPage);

  if (totalSessions === 0) {
    return <NoSessionsState />;
  }

  return (
    <div className="space-y-6">
      {/* Basic Metrics Overview */}
      {false && (<Card className="p-6">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiEyeLine className="h-4 w-4" />
              <span className="text-sm">Total views</span>
            </div>
            <p className="text-2xl font-medium">{sessionsData?.meta?.total_views || 0}</p>
          </div>
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiUserLine className="h-4 w-4" />
              <span className="text-sm">Unique views</span>
            </div>
            <p className="text-2xl font-medium">{sessionsData?.meta?.unique_views || 0}</p>
          </div>
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiTimeLine className="h-4 w-4" />
              <span className="text-sm">Time spent (avg)</span>
            </div>
            <p className="text-2xl font-medium">{sessionsData?.meta?.average_time_spent || "0:00"}</p>
          </div>
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiDownloadLine className="h-4 w-4" />
              <span className="text-sm">Downloads</span>
            </div>
            <p className="text-2xl font-medium">{sessionsData?.meta?.downloads || 0}</p>
          </div>
        </div>
      </Card>)}

      {/* Detailed Analytics */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Engagement Trends */}
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <RiBarChartLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Engagement Trends</h3>
          </div>
          <div className="h-[200px] flex flex-col gap-2">
            {sessionsData?.meta?.engagement_trends?.map((trend: EngagementTrend) => (
              <div key={trend.date} className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">{trend.date}</span>
                <div className="flex-1 mx-4">
                  <div
                    className="bg-primary h-2 rounded"
                    style={{
                      width: `${(trend.views / (sessionsData.meta.max_views || 1)) * 100}%`,
                    }}
                  />
                </div>
                <span className="text-sm font-medium">{trend.views}</span>
              </div>
            )) || []}
          </div>
        </Card>

        {/* Geographic Distribution */}
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <RiMapPinLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Geographic Distribution</h3>
          </div>
          <div className="flex flex-col items-center">
            <ChartContainer
              config={{
                United_States: {
                  label: "United States",
                  color: "hsl(217 91% 60%)",
                },
                United_Kingdom: {
                  label: "United Kingdom",
                  color: "hsl(271 91% 65%)",
                },
                Germany: {
                  label: "Germany",
                  color: "hsl(292 84% 61%)",
                },
                Canada: {
                  label: "Canada",
                  color: "hsl(316 70% 50%)",
                },
                Others: {
                  label: "Others",
                  color: "hsl(322 75% 46%)",
                },
              }}
              className="mx-auto aspect-[4/3] w-full max-w-lg"
            >
              <PieChart margin={{ top: 20, right: 0, bottom: 20, left: 0 }}>
                <Pie
                  data={sessionsData?.meta?.geographic_distribution?.map((item: GeographicDistribution) => ({
                    name: item.country,
                    value: item.views,
                    fill: `var(--color-${item.country.replace(/\s+/g, "_")})`,
                  })) || []}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  innerRadius={40}
                  outerRadius={80}
                  paddingAngle={2}
                />
                <ChartLegend
                  verticalAlign="bottom"
                  content={<ChartLegendContent />}
                  className="flex flex-wrap justify-center gap-4 pt-4"
                />
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="rounded-lg border bg-background p-2 shadow-md">
                          <p className="text-sm font-medium">{data.name}</p>
                          <p className="text-sm text-muted-foreground">{data.value.toLocaleString()} views</p>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
              </PieChart>
            </ChartContainer>
          </div>
        </Card>
      </div>

      {/* Viewer Sessions */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <RiUserHeartLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Viewer Sessions</h3>
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <RiFilterLine className="h-4 w-4 text-muted-foreground" />
              <Select value={timeFilter} onValueChange={handleTimeFilterChange}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Select time range" />
                </SelectTrigger>
                <SelectContent>
                  {TIME_RANGE_OPTIONS.map((option) => (
                    <Tooltip key={option.value}>
                      <TooltipTrigger asChild>
                        <div className="flex items-center justify-between w-full px-2 py-1.5 relative">
                          <SelectItem
                            value={option.value}
                            disabled={option.requiresPro && !HAS_PRO_ACCESS}
                            className="w-full"
                          >
                            {option.label}
                          </SelectItem>
                          {option.requiresPro && !HAS_PRO_ACCESS && (
                            <button
                              onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                                router.push("/settings?tab=billing");
                              }}
                              className="absolute right-2 top-1/2 -translate-y-1/2 group"
                            >
                              <RiVipCrownFill
                                className="h-4 w-4 text-amber-400 group-hover:text-amber-500 transition-colors shrink-0"
                              />
                            </button>
                          )}
                        </div>
                      </TooltipTrigger>
                      {option.requiresPro && !HAS_PRO_ACCESS && (
                        <TooltipContent
                          side="right"
                          className="w-72 p-4 bg-popover border rounded-md shadow-md"
                          sideOffset={5}
                        >
                          <div className="space-y-3">
                            <div className="flex items-center gap-2 text-foreground">
                              <RiVipCrownFill className="h-5 w-5 text-amber-500" />
                              <p className="font-semibold">Unlock Extended Analytics</p>
                            </div>
                            <p className="text-sm text-muted-foreground leading-relaxed">
                              Get access to longer historical data, advanced filters, and more detailed insights
                              with our Pro plan.
                            </p>
                            <Button
                              size="sm"
                              variant="default"
                              className="w-full font-medium"
                              onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                                router.push("/settings?tab=billing");
                              }}
                            >
                              Upgrade to Pro
                            </Button>
                          </div>
                        </TooltipContent>
                      )}
                    </Tooltip>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="text-sm text-muted-foreground">
              {totalSessions} sessions
            </div>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left border-b">
                <th className="pb-2 font-medium text-muted-foreground">Viewer</th>
                <th className="pb-2 font-medium text-muted-foreground">Views</th>
                <th className="pb-2 font-medium text-muted-foreground">Time Spent</th>
                <th className="pb-2 font-medium text-muted-foreground">Location</th>
                <th className="pb-2 font-medium text-muted-foreground">Device</th>
                <th className="pb-2 font-medium text-muted-foreground">Last Viewed</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {sessions.map((session: MalakDeckViewerSession) => (
                <tr key={session.id} className="hover:bg-muted/50">
                  <td className="py-3">
                    {session.contact ? (
                      <Link
                        href={`/contacts/${session.contact.reference}`}
                        className="text-sm text-primary hover:underline"
                      >
                        {session?.contact?.first_name || session?.contact?.email}
                      </Link>
                    ) : (
                      <span className="text-sm text-muted-foreground">
                        Anonymous {session.session_id ? `(${session.session_id})` : ""}
                      </span>
                    )}
                  </td>
                  <td className="py-3">
                    <span className="text-sm font-medium">1</span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm">
                      {formatTimeSpent(session.time_spent_seconds || 0)}
                    </span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm text-muted-foreground">
                      {session.country || 'N/A'}
                    </span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm text-muted-foreground">
                      {session.os && session.device_info 
                        ? `${session.os} / ${session.device_info}`
                        : session.os || session.device_info || 'N/A'}
                    </span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm text-muted-foreground">
                      {new Date(session.viewed_at || session.created_at || '').toLocaleDateString()}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex items-center justify-between mt-4 pt-4 border-t">
          <div className="text-sm text-muted-foreground">
            Showing {((currentPage - 1) * itemsPerPage) + 1}-{Math.min(currentPage * itemsPerPage, totalSessions)} of{" "}
            {totalSessions} sessions
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setCurrentPage(currentPage - 1)}
              disabled={currentPage === 1}
              className="p-2 rounded hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <RiArrowLeftLine className="h-4 w-4" />
            </button>
            <div className="text-sm">
              Page {currentPage} of {totalPages}
            </div>
            <button
              onClick={() => setCurrentPage(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="p-2 rounded hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <RiArrowRightLine className="h-4 w-4" />
            </button>
          </div>
        </div>
      </Card>
    </div>
  );
} 
