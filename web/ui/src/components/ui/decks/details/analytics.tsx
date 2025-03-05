import { Card } from "@/components/ui/card";
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
} from "@remixicon/react";
import { ServerFetchDeckResponse } from "@/client/Api";
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
import { Button } from "@/components/ui/button";

// Feature flag for pro features
const HAS_PRO_ACCESS = false; // Set this based on user's subscription status

// Time range options with pro status
const TIME_RANGE_OPTIONS = [
  { value: "7", label: "Last 7 days", requiresPro: false },
  { value: "14", label: "Last 14 days", requiresPro: true },
  { value: "30", label: "Last 30 days", requiresPro: true },
  { value: "90", label: "Last 90 days", requiresPro: true },
] as const;

interface ViewerSession {
  email: string | null;
  name?: string;
  contactId?: string;
  viewCount: number;
  lastViewed: string;
  averageTimeSpent: string;
  location: string;
  deviceInfo: string;
  sessionId?: string;
}

interface AnalyticsData {
  totalViews: number;
  uniqueViews: number;
  averageTimeSpent: string;
  downloads: number;
  engagementTrends: Array<{
    date: string;
    views: number;
  }>;
  geographicDistribution: Array<{
    country: string;
    views: number;
  }>;
  viewerSessions: ViewerSession[];
}

// Generate more mock session data
const generateMockSessions = (): ViewerSession[] => {
  const sessions: ViewerSession[] = [
    {
      email: "sarah.smith@example.com",
      name: "Sarah Smith",
      contactId: "c123",
      viewCount: 12,
      lastViewed: "2024-04-15T14:30:00Z",
      averageTimeSpent: "06:45",
      location: "United States",
      deviceInfo: "Chrome / macOS",
    },
  ];

  // Add more mock data
  const names = ["Alex Johnson", "Emma Wilson", "Michael Chen", "Sofia Rodriguez", "James Miller"];
  const locations = ["Canada", "Australia", "Japan", "Brazil", "India", "France"];
  const devices = [
    "Firefox / Windows",
    "Chrome / macOS",
    "Safari / iOS",
    "Chrome / Android",
    "Edge / Windows",
  ];

  // Generate 20 more sessions
  for (let i = 0; i < 20; i++) {
    const isAnonymous = Math.random() > 0.6; // 40% chance of being anonymous
    const date = new Date();
    date.setDate(date.getDate() - Math.floor(Math.random() * 30)); // Random date within last 30 days

    if (isAnonymous) {
      sessions.push({
        email: null,
        viewCount: Math.floor(Math.random() * 10) + 1,
        lastViewed: date.toISOString(),
        averageTimeSpent: `0${Math.floor(Math.random() * 10)}:${Math.floor(Math.random() * 60)}`,
        location: locations[Math.floor(Math.random() * locations.length)],
        deviceInfo: devices[Math.floor(Math.random() * devices.length)],
        sessionId: `anon_${Math.floor(Math.random() * 1000)}`,
      });
    } else {
      const name = names[Math.floor(Math.random() * names.length)];
      sessions.push({
        email: name.toLowerCase().replace(" ", ".") + "@example.com",
        name,
        contactId: `c${Math.floor(Math.random() * 1000)}`,
        viewCount: Math.floor(Math.random() * 20) + 1,
        lastViewed: date.toISOString(),
        averageTimeSpent: `0${Math.floor(Math.random() * 10)}:${Math.floor(Math.random() * 60)}`,
        location: locations[Math.floor(Math.random() * locations.length)],
        deviceInfo: devices[Math.floor(Math.random() * devices.length)],
      });
    }
  }

  return sessions;
};

// Mock data for analytics
const mockAnalytics: AnalyticsData = {
  totalViews: 1247,
  uniqueViews: 856,
  averageTimeSpent: "05:32",
  downloads: 124,
  engagementTrends: [
    { date: "2024-01", views: 150 },
    { date: "2024-02", views: 280 },
    { date: "2024-03", views: 420 },
    { date: "2024-04", views: 397 },
  ],
  geographicDistribution: [
    { country: "United States", views: 450 },
    { country: "United Kingdom", views: 230 },
    { country: "Germany", views: 180 },
    { country: "Canada", views: 165 },
    { country: "Others", views: 222 },
  ],
  viewerSessions: generateMockSessions(),
};

interface AnalyticsProps {
  data: ServerFetchDeckResponse;
}

export default function DeckAnalytics({ data }: AnalyticsProps) {
  const router = useRouter();
  const [timeFilter, setTimeFilter] = useState("7"); // Default to 7 days for non-pro users
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  // Filter sessions based on selected time range
  const filteredSessions = mockAnalytics.viewerSessions.filter((session) => {
    const sessionDate = new Date(session.lastViewed);
    const now = new Date();
    const diffDays = Math.floor((now.getTime() - sessionDate.getTime()) / (1000 * 60 * 60 * 24));
    return diffDays <= parseInt(timeFilter);
  });

  // Calculate pagination
  const totalPages = Math.ceil(filteredSessions.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const paginatedSessions = filteredSessions.slice(startIndex, startIndex + itemsPerPage);

  // Reset to first page when filter changes
  const handleTimeFilterChange = (value: string) => {
    setTimeFilter(value);
    setCurrentPage(1);
  };

  const handleUpgradeClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    router.push("/settings?tab=billing");
  };

  return (
    <div className="space-y-6">
      {/* Basic Metrics Overview */}
      <Card className="p-6">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiEyeLine className="h-4 w-4" />
              <span className="text-sm">Total views</span>
            </div>
            <p className="text-2xl font-medium">{mockAnalytics.totalViews}</p>
          </div>
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiUserLine className="h-4 w-4" />
              <span className="text-sm">Unique views</span>
            </div>
            <p className="text-2xl font-medium">{mockAnalytics.uniqueViews}</p>
          </div>
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiTimeLine className="h-4 w-4" />
              <span className="text-sm">Time spent (avg)</span>
            </div>
            <p className="text-2xl font-medium">{mockAnalytics.averageTimeSpent}</p>
          </div>
          <div>
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <RiDownloadLine className="h-4 w-4" />
              <span className="text-sm">Downloads</span>
            </div>
            <p className="text-2xl font-medium">{mockAnalytics.downloads}</p>
          </div>
        </div>
      </Card>

      {/* Detailed Analytics */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Engagement Trends */}
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <RiBarChartLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Engagement Trends</h3>
          </div>
          <div className="h-[200px] flex flex-col gap-2">
            {mockAnalytics.engagementTrends.map((trend) => (
              <div key={trend.date} className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">{trend.date}</span>
                <div className="flex-1 mx-4">
                  <div
                    className="bg-primary h-2 rounded"
                    style={{
                      width: `${(trend.views / 420) * 100}%`,
                    }}
                  />
                </div>
                <span className="text-sm font-medium">{trend.views}</span>
              </div>
            ))}
          </div>
        </Card>

        {/* Geographic Distribution */}
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <RiMapPinLine className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-medium">Geographic Distribution</h3>
          </div>
          <div className="h-[200px] flex flex-col gap-2">
            {mockAnalytics.geographicDistribution.map((geo) => (
              <div key={geo.country} className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground w-24">{geo.country}</span>
                <div className="flex-1 mx-4">
                  <div
                    className="bg-primary h-2 rounded"
                    style={{
                      width: `${(geo.views / 450) * 100}%`,
                    }}
                  />
                </div>
                <span className="text-sm font-medium">{geo.views}</span>
              </div>
            ))}
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
              {filteredSessions.length} sessions
            </div>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left border-b">
                <th className="pb-2 font-medium text-muted-foreground">Viewer</th>
                <th className="pb-2 font-medium text-muted-foreground">Views</th>
                <th className="pb-2 font-medium text-muted-foreground">Avg. Time</th>
                <th className="pb-2 font-medium text-muted-foreground">Location</th>
                <th className="pb-2 font-medium text-muted-foreground">Device</th>
                <th className="pb-2 font-medium text-muted-foreground">Last Viewed</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {paginatedSessions.map((session, index) => (
                <tr key={index} className="hover:bg-muted/50">
                  <td className="py-3">
                    {session.contactId ? (
                      <Link
                        href={`/contacts/${session.contactId}`}
                        className="text-sm text-primary hover:underline"
                      >
                        {session.name}
                      </Link>
                    ) : (
                      <span className="text-sm text-muted-foreground">
                        Anonymous {session.sessionId ? `(${session.sessionId})` : ""}
                      </span>
                    )}
                  </td>
                  <td className="py-3">
                    <span className="text-sm font-medium">{session.viewCount}</span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm">{session.averageTimeSpent}</span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm">{session.location}</span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm text-muted-foreground">{session.deviceInfo}</span>
                  </td>
                  <td className="py-3">
                    <span className="text-sm text-muted-foreground">
                      {new Date(session.lastViewed).toLocaleDateString()}
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
            Showing {startIndex + 1}-{Math.min(startIndex + itemsPerPage, filteredSessions.length)} of{" "}
            {filteredSessions.length} sessions
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