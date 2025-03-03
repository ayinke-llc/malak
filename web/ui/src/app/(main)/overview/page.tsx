"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  RiMailLine,
  RiFileTextLine,
  RiPresentationLine, RiTimeLine,
  RiEditLine,
  RiSendPlaneLine,
  RiArrowRightSLine
} from "@remixicon/react";
import Link from "next/link";

export default function Overview() {
  const activityLog = [
    {
      type: "deck",
      title: "Seed Pitch Deck",
      action: "opened",
      user: "Emma Thompson",
      time: "May 20, 2024"
    },
    {
      type: "deck",
      title: "Seed Pitch Deck",
      action: "opened",
      user: "Emma Thompson and 1 other",
      time: "Apr 9, 2024"
    },
    {
      type: "deck",
      title: "Example Deck",
      action: "opened",
      user: "James Wilson",
      time: "Jan 10, 2024"
    },
    {
      type: "update",
      title: "Q1 Investor Update",
      action: "opened via email",
      user: "Oliver Parker",
      time: "Jun 20, 2023"
    },
    {
      type: "update",
      title: "Investor Update",
      action: "opened via email",
      user: "Sophie Anderson",
      time: "Jun 15, 2023"
    }
  ];

  const recentDrafts = [
    {
      title: "The Malak Starter Plan Investor Update",
      lastEdit: "11 days ago"
    },
    {
      title: "The Malak Starter Plan Investor Update",
      lastEdit: "2 months ago"
    },
    {
      title: "Untitled",
      lastEdit: "2 months ago"
    },
    {
      title: "The Malak Starter Plan Investor Update",
      lastEdit: "2 months ago"
    }
  ];

  const recentlySent = [
    {
      title: "Harrison's Co Monthly Update",
      recipients: 1,
      sentTime: "10 months ago"
    },
    {
      title: "Parker's Co Monthly Update",
      recipients: 1,
      sentTime: "10 months ago"
    },
    {
      title: "MSP Monthly Update",
      recipients: 1,
      sentTime: "10 months ago"
    },
    {
      title: "Q1 Investor Update",
      recipients: 1,
      group: "Investors",
      sentTime: "2 years ago"
    }
  ];

  return (
    <div className="space-y-8">
      {/* Header Section */}
      <div className="flex flex-col gap-1">
        <h1 className="text-3xl font-bold tracking-tight">Overview</h1>
        <p className="text-muted-foreground">
          Recent activity and updates
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Activity Log */}
        <Card className="lg:col-span-1">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <RiTimeLine className="h-5 w-5 text-muted-foreground" />
              Activity Log
            </CardTitle>
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
                    <div className="flex items-center text-sm text-muted-foreground">
                      <span>{activity.time}</span>
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
                  <div className="p-2 rounded-full bg-orange-100 text-orange-600">
                    <RiFileTextLine className="h-4 w-4" />
                  </div>
                  <div className="flex-1 space-y-1">
                    <p className="text-sm font-medium leading-none">{draft.title}</p>
                    <p className="text-sm text-muted-foreground">Last edit {draft.lastEdit}</p>
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
                    <p className="text-sm font-medium leading-none">{update.title}</p>
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
