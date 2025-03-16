"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  RiMailLine,
  RiFileTextLine,
  RiPresentationLine,
  RiTimeLine,
  RiEditLine,
  RiSendPlaneLine,
  RiArrowRightSLine,
  RiEyeLine,
  RiTeamLine,
  RiBarChartLine,
  RiAddLine,
} from "@remixicon/react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";

export default function Overview() {
  // Mock data - replace with real API calls
  const metrics = {
    totalContacts: 156,
    activeDecks: 8,
    totalViews: 1243,
    engagementRate: 68,
    monthlyGrowth: 12.5
  };

  const activityLog = [
    {
      type: "deck",
      title: "Seed Pitch Deck",
      action: "opened",
      user: "Emma Thompson",
      time: "May 20, 2024",
      details: "Viewed 8 slides • 5 min duration"
    },
    {
      type: "deck",
      title: "Seed Pitch Deck",
      action: "opened",
      user: "Emma Thompson and 1 other",
      time: "Apr 9, 2024",
      details: "Viewed all slides • 12 min duration"
    },
    {
      type: "deck",
      title: "Example Deck",
      action: "opened",
      user: "James Wilson",
      time: "Jan 10, 2024",
      details: "Downloaded deck"
    },
    {
      type: "update",
      title: "Q1 Investor Update",
      action: "opened via email",
      user: "Oliver Parker",
      time: "Jun 20, 2023",
      details: "Read time: 3 min"
    },
    {
      type: "update",
      title: "Investor Update",
      action: "opened via email",
      user: "Sophie Anderson",
      time: "Jun 15, 2023",
      details: "Read time: 5 min"
    }
  ];

  const recentDrafts = [
    {
      title: "The Malak Starter Plan Investor Update",
      lastEdit: "11 days ago",
      progress: 80,
      type: "update"
    },
    {
      title: "Series A Pitch Deck",
      lastEdit: "2 months ago",
      progress: 45,
      type: "deck"
    },
    {
      title: "Q2 2024 Investor Update",
      lastEdit: "2 months ago",
      progress: 25,
      type: "update"
    },
    {
      title: "Product Roadmap Presentation",
      lastEdit: "2 months ago",
      progress: 90,
      type: "deck"
    }
  ];

  const recentlySent = [
    {
      title: "Harrison's Co Monthly Update",
      recipients: 12,
      sentTime: "10 months ago",
      openRate: 85,
      engagement: "High"
    },
    {
      title: "Parker's Co Monthly Update",
      recipients: 8,
      sentTime: "10 months ago",
      openRate: 75,
      engagement: "Medium"
    },
    {
      title: "MSP Monthly Update",
      recipients: 15,
      sentTime: "10 months ago",
      openRate: 92,
      engagement: "High"
    },
    {
      title: "Q1 Investor Update",
      recipients: 20,
      group: "Investors",
      sentTime: "2 years ago",
      openRate: 88,
      engagement: "High"
    }
  ];

  return (
    <div className="space-y-8">
      {/* Header Section */}
      <div className="flex flex-col gap-1">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Overview</h1>
            <p className="text-muted-foreground">
              Recent activity and updates
            </p>
          </div>
        </div>
      </div>

      {/* Metrics Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Contacts</CardTitle>
            <RiTeamLine className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.totalContacts}</div>
            <p className="text-xs text-muted-foreground">
              +{metrics.monthlyGrowth}% from last month
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Decks</CardTitle>
            <RiPresentationLine className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.activeDecks}</div>
            <p className="text-xs text-muted-foreground">
              {metrics.activeDecks} decks shared this month
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Views</CardTitle>
            <RiEyeLine className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.totalViews}</div>
            <p className="text-xs text-muted-foreground">
              Across all decks and updates
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Activity Log */}
        <Card className="lg:col-span-1">
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
            <div className="space-y-6">
              {activityLog.map((activity, index) => (
                <div
                  key={index}
                  className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50"
                >
                  <div className={`p-2 rounded-full 
                    ${activity.type === "update" ? "bg-blue-100 text-blue-600" : ""}
                    ${activity.type === "deck" ? "bg-purple-100 text-purple-600" : ""}
                  `}>
                    {activity.type === "update" && <RiMailLine className="h-4 w-4" />}
                    {activity.type === "deck" && <RiPresentationLine className="h-4 w-4" />}
                  </div>
                  <div className="flex-1 space-y-1">
                    <p className="text-sm font-medium leading-none">
                      {activity.user} {activity.action} {activity.title}
                    </p>
                    <div className="flex items-center text-sm text-muted-foreground gap-2">
                      <span>{activity.time}</span>
                      <span>•</span>
                      <span>{activity.details}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Drafts */}
        <Card className="lg:col-span-1">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <RiEditLine className="h-5 w-5 text-muted-foreground" />
              Recent Drafts
            </CardTitle>
            <Link
              href="/updates/drafts"
              className="text-sm text-muted-foreground hover:text-primary flex items-center gap-1 transition-colors"
            >
              View all
              <RiArrowRightSLine className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentDrafts.map((draft, index) => (
                <div
                  key={index}
                  className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50"
                >
                  <div className={`p-2 rounded-full ${
                    draft.type === "update" ? "bg-orange-100 text-orange-600" : "bg-purple-100 text-purple-600"
                  }`}>
                    {draft.type === "update" ? (
                      <RiFileTextLine className="h-4 w-4" />
                    ) : (
                      <RiPresentationLine className="h-4 w-4" />
                    )}
                  </div>
                  <div className="flex-1 space-y-2">
                    <div className="flex justify-between items-center">
                      <p className="text-sm font-medium leading-none">{draft.title}</p>
                      <span className="text-xs text-muted-foreground">{draft.lastEdit}</span>
                    </div>
                    <Progress value={draft.progress} className="h-1.5" />
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recently Sent Updates */}
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <RiSendPlaneLine className="h-5 w-5 text-muted-foreground" />
              Recently Sent Updates
            </CardTitle>
            <Link
              href="/updates/sent"
              className="text-sm text-muted-foreground hover:text-primary flex items-center gap-1 transition-colors"
            >
              View all
              <RiArrowRightSLine className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentlySent.map((update, index) => (
                <div
                  key={index}
                  className="flex items-center gap-4 p-2 rounded-lg transition-colors hover:bg-muted/50"
                >
                  <div className="p-2 rounded-full bg-green-100 text-green-600">
                    <RiMailLine className="h-4 w-4" />
                  </div>
                  <div className="flex-1 space-y-1">
                    <div className="flex justify-between items-center">
                      <p className="text-sm font-medium leading-none">{update.title}</p>
                      <span className={`text-xs px-2 py-1 rounded-full ${
                        update.engagement === "High" 
                          ? "bg-green-100 text-green-700"
                          : update.engagement === "Medium"
                          ? "bg-yellow-100 text-yellow-700"
                          : "bg-red-100 text-red-700"
                      }`}>
                        {update.openRate}% Open Rate
                      </span>
                    </div>
                    <div className="flex items-center text-sm text-muted-foreground gap-2">
                      <span>Sent to {update.recipients} {update.recipients === 1 ? 'person' : 'people'}</span>
                      {update.group && (
                        <>
                          <span>•</span>
                          <span>{update.group}</span>
                        </>
                      )}
                      <span>•</span>
                      <span>Sent {update.sentTime}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 