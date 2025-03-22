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
} from "@remixicon/react";
import Link from "next/link";

export default function Overview() {
  // Simplified mock data
  const metrics = {
    totalViews: 1243,
    activeDecks: 8,
    totalUpdates: 24,
    totalContacts: 156,
  };

  const activityLogs = [
    {
      title: "Series A Pitch Deck",
      type: "deck",
      action: "shared",
      date: "2 hours ago",
      recipient: "John Smith"
    },
    {
      title: "Q2 2024 Investor Update",
      type: "update",
      action: "viewed",
      date: "1 day ago",
      recipient: "Emma Wilson"
    },
    {
      title: "Product Roadmap Deck",
      type: "deck",
      action: "downloaded",
      date: "1 day ago",
      recipient: "Sarah Chen"
    },
    {
      title: "Monthly Metrics Update",
      type: "update",
      action: "shared",
      date: "2 days ago",
      recipient: "David Kumar"
    },
    {
      title: "Seed Round Deck",
      type: "deck",
      action: "downloaded",
      date: "2 days ago",
      recipient: "Mike Johnson"
    },
    {
      title: "Growth Strategy Presentation",
      type: "deck",
      action: "viewed",
      date: "3 days ago",
      recipient: "Lisa Rodriguez"
    },
    {
      title: "Q1 2024 Financial Update",
      type: "update",
      action: "shared",
      date: "4 days ago",
      recipient: "Alex Thompson"
    },
    {
      title: "Market Analysis Deck",
      type: "deck",
      action: "viewed",
      date: "5 days ago",
      recipient: "James Wilson"
    },
    {
      title: "Team Growth Update",
      type: "update",
      action: "downloaded",
      date: "6 days ago",
      recipient: "Rachel Kim"
    },
    {
      title: "Partnership Proposal Deck",
      type: "deck",
      action: "shared",
      date: "1 week ago",
      recipient: "Tom Martinez"
    }
  ];

  const recentUpdates = [
    {
      id: "update-1",
      title: "The Malak Starter Plan Update",
      date: "11 days ago",
      views: 24,
      recipients: 12
    },
    {
      id: "update-2",
      title: "Q2 2024 Investor Update",
      date: "2 months ago",
      views: 45,
      recipients: 18
    },
    {
      id: "update-3",
      title: "Monthly Progress Report",
      date: "2 months ago",
      views: 32,
      recipients: 15
    }
  ];

  return (
    <div className="space-y-6">
      {/* Header Section */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Your investor updates and pitch decks at a glance
        </p>
      </div>

      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Link href="/decks" className="block transition-transform hover:scale-[1.02]">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Decks</CardTitle>
              <RiPresentationLine className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics.activeDecks}</div>
              <p className="text-xs text-muted-foreground">
                Currently shared decks
              </p>
            </CardContent>
          </Card>
        </Link>
        <Link href="/updates" className="block transition-transform hover:scale-[1.02]">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Updates</CardTitle>
              <RiFileTextLine className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics.totalUpdates}</div>
              <p className="text-xs text-muted-foreground">
                Investor updates sent
              </p>
            </CardContent>
          </Card>
        </Link>
        <Link href="/analytics" className="block transition-transform hover:scale-[1.02]">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Views</CardTitle>
              <RiEyeLine className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics.totalViews}</div>
              <p className="text-xs text-muted-foreground">
                Across all content
              </p>
            </CardContent>
          </Card>
        </Link>
        <Link href="/contacts" className="block transition-transform hover:scale-[1.02]">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Contacts</CardTitle>
              <RiTeamLine className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics.totalContacts}</div>
              <p className="text-xs text-muted-foreground">
                Active recipients
              </p>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Activity Cards */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Activity Log */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <RiTimeLine className="h-5 w-5 text-muted-foreground" />
              Activity Log
            </CardTitle>
            <Link
              href="/activity"
              className="text-sm text-muted-foreground hover:text-primary flex items-center gap-1 transition-colors"
            >
              View all
              <RiArrowRightSLine className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {activityLogs.map((item, index) => (
                <div
                  key={index}
                  className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50"
                >
                  <div className={`p-2 rounded-full ${
                    item.type === "update" ? "bg-blue-100 text-blue-600" : "bg-purple-100 text-purple-600"
                  }`}>
                    {item.type === "update" ? (
                      <RiFileTextLine className="h-4 w-4" />
                    ) : (
                      <RiPresentationLine className="h-4 w-4" />
                    )}
                  </div>
                  <div className="flex-1">
                    <div className="flex justify-between items-center">
                      <p className="text-sm font-medium">{item.title}</p>
                      <span className="text-xs text-muted-foreground">{item.date}</span>
                    </div>
                    <p className="text-sm text-muted-foreground">
                      {item.action} by {item.recipient}
                    </p>
                  </div>
                </div>
              ))}
            </div>
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
            <div className="space-y-4">
              {recentUpdates.map((update, index) => (
                <Link
                  key={index}
                  href={`/updates/${update.id}`}
                  className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50"
                >
                  <div className="p-2 rounded-full bg-blue-100 text-blue-600">
                    <RiFileTextLine className="h-4 w-4" />
                  </div>
                  <div className="flex-1">
                    <div className="flex justify-between items-center">
                      <p className="text-sm font-medium">{update.title}</p>
                      <span className="text-xs text-muted-foreground">{update.date}</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                      <span>{update.views} views</span>
                      <span>•</span>
                      <span>{update.recipients} recipients</span>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 