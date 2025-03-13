import {
  MalakDeckDailyEngagement,
  MalakDeckGeographicStat,
  MalakDeckViewerSession, MalakPlan,
  ServerFetchEngagementsResponse,
  ServerFetchSessionsDeck
} from "@/client/Api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ChartContainer, ChartLegend } from "@/components/ui/chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import client from "@/lib/client";
import { LIST_DECK_ENGAGEMENTS, LIST_DECK_SESSIONS } from "@/lib/query-constants";
import useWorkspacesStore from "@/store/workspace";
import {
  RiArrowLeftLine,
  RiArrowRightLine,
  RiBarChartLine, RiErrorWarningLine, RiFilterLine,
  RiMapPinLine, RiUserHeartLine, RiVipCrownFill
} from "@remixicon/react";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Pie, PieChart, Tooltip as RechartsTooltip } from "recharts";

// Time range options with pro status
const TIME_RANGE_OPTIONS = [
  { value: "7", label: "Last 7 days", requiresPro: false },
  { value: "14", label: "Last 14 days", requiresPro: true },
  { value: "30", label: "Last 30 days", requiresPro: true }
] as const;

interface DeckAnalyticsProps {
  reference: string;
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
      <div className="space-y-6">
        <div className="flex flex-col items-center justify-center text-center space-y-4">
          <RiUserHeartLine className="h-12 w-12 text-muted-foreground" />
          <div className="space-y-2">
            <h3 className="text-lg font-medium">No viewer sessions yet</h3>
            <p className="text-sm text-muted-foreground max-w-md">
              When people view your deck, their session data will appear here.
            </p>
          </div>
        </div>
      </div>
    </Card>
  );
}

interface DeckEngagementsProps {
  reference: string;
}


function DeckEngagements({ reference }: DeckEngagementsProps) {
  const { data: engagementsData, isLoading, error } = useQuery({
    queryKey: [LIST_DECK_ENGAGEMENTS, reference],
    queryFn: async () => {
      const response = await client.decks.analyticsDetail(reference);
      return response.data as ServerFetchEngagementsResponse;
    },
  });

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <div className="flex items-center gap-2 mb-4">
              <div className="h-5 w-5 rounded bg-muted animate-pulse" />
              <div className="h-5 w-32 rounded bg-muted animate-pulse" />
            </div>
            <div className="space-y-3">
              <div className="h-8 w-full rounded bg-muted animate-pulse" />
              <div className="h-8 w-3/4 rounded bg-muted animate-pulse" />
              <div className="h-8 w-1/2 rounded bg-muted animate-pulse" />
            </div>
          </div>
          <div>
            <div className="flex items-center gap-2 mb-4">
              <div className="h-5 w-5 rounded bg-muted animate-pulse" />
              <div className="h-5 w-32 rounded bg-muted animate-pulse" />
            </div>
            <div className="flex items-center justify-center h-[200px]">
              <div className="h-40 w-40 rounded-full bg-muted animate-pulse" />
            </div>
          </div>
        </div>
      </Card>
    );
  }

  if (error) {
    return <ErrorState error={error as Error} />;
  }

  const { daily_engagements, geographic_stats } = engagementsData?.engagements ?? {};

  if (!daily_engagements?.length && !geographic_stats?.length) {
    return <NoSessionsState />;
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      {/* Engagement Trends */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <RiBarChartLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Engagement Trends</h3>
          </div>
          <span className="text-sm text-muted-foreground">Data may be delayed by a few hours</span>
        </div>
        <div className="h-[200px] flex flex-col gap-2">
          {daily_engagements?.map((trend: MalakDeckDailyEngagement) => (
            <div key={trend.engagement_date} className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">{new Date(trend.engagement_date || '').toLocaleDateString()}</span>
              <div className="flex-1 mx-4">
                <div
                  className="bg-primary h-2 rounded"
                  style={{
                    width: `${(trend.engagement_count || 0) / (Math.max(...daily_engagements.map(d => d.engagement_count || 0)) || 1) * 100}%`,
                  }}
                />
              </div>
              <span className="text-sm font-medium">{trend.engagement_count}</span>
            </div>
          )) || []}
        </div>
      </Card>

      {/* Geographic Distribution */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <RiMapPinLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Geographic Distribution</h3>
          </div>
          <span className="text-sm text-muted-foreground">Data may be delayed by a few hours</span>
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
                data={geographic_stats?.map((item: MalakDeckGeographicStat) => ({
                  name: item.country,
                  value: item.view_count,
                  fill: `var(--color-${item.country?.replace(/\s+/g, "_")})`,
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
  );
}

interface ViewerSessionsProps {
  reference: string;
}

function ViewerSessions({ reference }: ViewerSessionsProps) {
  const router = useRouter();
  const [timeFilter, setTimeFilter] = useState("7");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;
  const plan = useWorkspacesStore(state => state.current?.plan as MalakPlan);

  const { data: sessionsData, isLoading, error } = useQuery({
    queryKey: [LIST_DECK_SESSIONS, reference, timeFilter, currentPage],
    queryFn: async () => {
      const response = await client.decks.sessionsDetail(
        reference,
        {
          page: currentPage,
          per_page: itemsPerPage,
          days: parseInt(timeFilter),
        }
      );
      return response.data as ServerFetchSessionsDeck;
    },
  });

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="h-5 w-5 rounded bg-muted animate-pulse" />
              <div className="h-5 w-32 rounded bg-muted animate-pulse" />
            </div>
            <div className="h-9 w-44 rounded bg-muted animate-pulse" />
          </div>
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex items-center gap-4">
                <div className="h-8 w-32 rounded bg-muted animate-pulse" />
                <div className="h-8 flex-1 rounded bg-muted animate-pulse" />
                <div className="h-8 w-24 rounded bg-muted animate-pulse" />
                <div className="h-8 w-24 rounded bg-muted animate-pulse" />
              </div>
            ))}
          </div>
        </div>
      </Card>
    );
  }

  if (error) {
    return <ErrorState error={error as Error} />;
  }

  const handleTimeFilterChange = (value: string) => {
    setTimeFilter(value);
    setCurrentPage(1);
  };

  const sessions = sessionsData?.sessions ?? [];
  const totalSessions = sessionsData?.meta?.paging?.total ?? 0;
  const totalPages = Math.ceil(totalSessions / itemsPerPage);

  if (totalSessions === 0) {
    return <NoSessionsState />;
  }

  return (
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
                          disabled={option.requiresPro && !plan.metadata?.deck?.analytics?.can_view_historical_sessions}
                          className="w-full"
                        >
                          {option.label}
                        </SelectItem>
                        {option.requiresPro && !plan.metadata?.deck?.analytics?.can_view_historical_sessions && (
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
                    {option.requiresPro && !plan.metadata?.deck?.analytics?.can_view_historical_sessions && (
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
  );
}

// Add this new component before the DeckAnalytics export
interface ChartLegendContentProps {
  payload?: Array<{
    value: number;
    payload: {
      name: string;
      value: number;
      fill: string;
    };
  }>;
}

function ChartLegendContent({ payload }: ChartLegendContentProps) {
  if (!payload) return null;

  return (
    <div className="flex flex-wrap justify-center gap-4 mt-4">
      {payload.map((entry, index) => (
        <div key={`legend-${index}`} className="flex items-center gap-2">
          <div
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: entry.payload.fill }}
          />
          <span className="text-sm text-muted-foreground">
            {entry.payload.name} ({entry.payload.value})
          </span>
        </div>
      ))}
    </div>
  );
}

export default function DeckAnalytics({ reference }: DeckAnalyticsProps) {
  if (!reference) {
    return <ErrorState error={new Error("No deck reference provided")} />;
  }

  return (
    <div className="space-y-6">
      {/* Basic Metrics Overview - Commented out until API supports it */}
      <DeckEngagements reference={reference} />
      <ViewerSessions reference={reference} />
    </div>
  );
} 
