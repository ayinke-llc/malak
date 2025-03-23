"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  RiMailLine,
  RiFileTextLine,
  RiPresentationLine,
  RiEyeLine,
  RiArrowRightSLine,
  RiTeamLine,
  RiTimeLine,
  RiDashboardLine,
} from "@remixicon/react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { FETCH_OVERVIEW_DATA } from "@/lib/query-constants";
import client from "@/lib/client";
import { Skeleton } from "@/components/ui/skeleton";
import { ServerWorkspaceOverviewResponse } from "@/client/Api";
import { formatDistanceToNow } from "date-fns";

interface ActivityLog {
  id: string;
  reference: string;
  title: string;
  type: "deck" | "update" | "dashboard";
  action: string;
  date: string;
  email: string;
  first_name: string;
  last_name: string;
}

interface RecentUpdate {
  id: string;
  title: string;
  date: string;
  reference: string;
}

function MetricsCard({
  title,
  value,
  description,
  icon: Icon,
  href
}: {
  title: string;
  value: number;
  description: string;
  icon: React.ElementType;
  href: string;
}) {
  const content = (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        <p className="text-xs text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  );

  if (!href) {
    return content;
  }

  return (
    <Link href={href} className="block transition-transform hover:scale-[1.02]">
      {content}
    </Link>
  );
}

function MetricsCardSkeleton() {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-4 w-4" />
      </CardHeader>
      <CardContent>
        <Skeleton className="h-8 w-16 mb-1" />
        <Skeleton className="h-3 w-32" />
      </CardContent>
    </Card>
  );
}

function ActivityLogItem({ item }: { item: ActivityLog }) {
  const getIconAndColor = (type: string) => {
    switch (type) {
      case "update":
        return {
          Icon: RiFileTextLine,
          bgColor: "bg-blue-100",
          textColor: "text-blue-600",
          href: `/updates/${item.reference}`
        };
      case "deck":
        return {
          Icon: RiPresentationLine,
          bgColor: "bg-purple-100",
          textColor: "text-purple-600",
          href: `/decks/${item.reference}`
        };
      case "dashboard":
        return {
          Icon: RiDashboardLine,
          bgColor: "bg-green-100",
          textColor: "text-green-600",
          href: `/dashboards/${item.reference}`
        };
      default:
        return {
          Icon: RiFileTextLine,
          bgColor: "bg-gray-100",
          textColor: "text-gray-600",
          href: "#"
        };
    }
  };

  const { Icon, bgColor, textColor, href } = getIconAndColor(item.type);

  const displayName = !item.first_name || item.first_name.toLowerCase() === 'investor' ? item.email : `${item.first_name} ${item.last_name}`;

  return (
    <Link href={href} className="block">
      <div className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50">
        <div className={`p-2 rounded-full ${bgColor} ${textColor}`}>
          <Icon className="h-4 w-4" />
        </div>
        <div className="flex-1">
          <div className="flex justify-between items-center">
            <p className="text-sm font-medium">{item.title}</p>
            <span className="text-xs text-muted-foreground">
              {formatDistanceToNow(new Date(item.date), { addSuffix: true })}
            </span>
          </div>
          <p className="text-sm text-muted-foreground">
            Shared with {displayName}
          </p>
        </div>
      </div>
    </Link>
  );
}

function ActivityLogSkeleton() {
  return (
    <div className="flex items-center gap-4 p-2">
      <Skeleton className="h-8 w-8 rounded-full" />
      <div className="flex-1">
        <Skeleton className="h-4 w-48 mb-2" />
        <Skeleton className="h-3 w-32" />
      </div>
    </div>
  );
}

function RecentUpdateItem({ update }: { update: RecentUpdate }) {
  return (
    <Link
      href={`/updates/${update.reference}`}
      className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50"
    >
      <div className="p-2 rounded-full bg-blue-100 text-blue-600">
        <RiFileTextLine className="h-4 w-4" />
      </div>
      <div className="flex-1">
        <div className="flex justify-between items-center">
          <p className="text-sm font-medium">{update.title}</p>
          <span className="text-xs text-muted-foreground">
            {formatDistanceToNow(new Date(update.date), { addSuffix: true })}
          </span>
        </div>
      </div>
    </Link>
  );
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-8 text-center">
      <RiFileTextLine className="h-12 w-12 text-muted-foreground mb-4" />
      <p className="text-muted-foreground">{message}</p>
    </div>
  );
}

export default function Overview() {
  const { data, isLoading, error } = useQuery<ServerWorkspaceOverviewResponse>({
    queryKey: [FETCH_OVERVIEW_DATA],
    queryFn: async () => {
      const response = await client.workspaces.overviewList();
      return response.data;
    },
  });

  if (error) {
    return (
      <div className="flex items-center justify-center h-[50vh]">
        <p className="text-destructive">Failed to load overview data. Please try again later.</p>
      </div>
    );
  }

  const metrics = {
    totalViews: data?.decks.total_viewer_sessions ?? 0,
    activeDecks: data?.decks.total_decks ?? 0,
    totalUpdates: data?.updates.total ?? 0,
    totalContacts: data?.contacts.total_contacts ?? 0,
  };

  return (
    <div className="space-y-6">
      {/* Header Section */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Your activity logs, investor updates and pitch decks at a glance
        </p>
      </div>

      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {isLoading ? (
          <>
            <MetricsCardSkeleton />
            <MetricsCardSkeleton />
            <MetricsCardSkeleton />
            <MetricsCardSkeleton />
          </>
        ) : (
          <>
            <MetricsCard
              title="Active Decks"
              value={metrics.activeDecks}
              description="Currently shared decks"
              icon={RiPresentationLine}
              href="/decks"
            />
            <MetricsCard
              title="Total Updates"
              value={metrics.totalUpdates}
              description="Investor updates sent"
              icon={RiFileTextLine}
              href="/updates"
            />
            <MetricsCard
              title="Total Views"
              value={metrics.totalViews}
              description="Across your decks"
              icon={RiEyeLine}
              href=""
            />
            <MetricsCard
              title="Total Contacts"
              value={metrics.totalContacts}
              description="Active recipients"
              icon={RiTeamLine}
              href="/contacts"
            />
          </>
        )}
      </div>

      {/* Activity Cards */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Activity Log */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <RiTimeLine className="h-5 w-5 text-muted-foreground" />
              Recent Activity
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="space-y-4">
                {Array(5).fill(0).map((_, i) => (
                  <ActivityLogSkeleton key={i} />
                ))}
              </div>
            ) : data?.shares.recent_shares && data.shares.recent_shares.length > 0 ? (
              <div className="space-y-4">
                {data.shares.recent_shares.map((item) => (
                  <ActivityLogItem
                    key={item.id}
                    item={{
                      id: item.id ?? "",
                      reference: item.item_reference ?? "",
                      title: item.title ?? "",
                      type: item.item_type ?? "update",
                      action: "shared",
                      date: item.shared_at ?? "",
                      email: item.email ?? "",
                      first_name: item.first_name ?? "",
                      last_name: item.last_name ?? "",
                    }}
                  />
                ))}
              </div>
            ) : (
              <EmptyState message="No recent activity to show" />
            )}
          </CardContent>
        </Card>

        {/* Recently Sent Updates */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <RiMailLine className="h-5 w-5 text-muted-foreground" />
              Recently Sent Updates
            </CardTitle>
            <Link
              href="/updates"
              className="text-sm text-muted-foreground hover:text-primary flex items-center gap-1 transition-colors"
            >
              View all
              <RiArrowRightSLine className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="space-y-4">
                {Array(3).fill(0).map((_, i) => (
                  <ActivityLogSkeleton key={i} />
                ))}
              </div>
            ) : data?.updates.last_updates && data.updates.last_updates.length > 0 ? (
              <div className="space-y-4">
                {data.updates.last_updates.map((update) => (
                  <RecentUpdateItem
                    key={update.id}
                    update={{
                      id: update.id ?? "",
                      title: update.title ?? "",
                      date: update.sent_at ?? update.created_at ?? "",
                      reference: update.reference ?? "",
                    }}
                  />
                ))}
              </div>
            ) : (
              <EmptyState message="No updates sent yet" />
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 
