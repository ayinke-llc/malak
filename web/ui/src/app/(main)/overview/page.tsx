"use client";
import { overviews } from "@/data/overview-data";
import type { OverviewData } from "@/data/schema";
import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { format, parseISO } from "date-fns";
import { 
  ArrowDown, 
  ArrowUp, 
  Users, 
  Mail,
  Building, 
  FileText,
  TrendingUp,
  Presentation,
  Eye,
  Clock,
  Lock,
  Share2,
  BarChart4,
  Activity
} from "lucide-react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  AreaChart,
  Area,
  BarChart,
  Bar,
} from "recharts";

export default function Overview() {
  const filteredData = overviews;
  const latestMetrics = filteredData[filteredData.length - 1];
  const previousMetrics = filteredData[filteredData.length - 2];

  const getPercentageChange = (current: number, previous: number) => {
    if (!previous) return 0;
    return ((current - previous) / previous) * 100;
  };

  const metrics = [
    {
      title: "Active Investors",
      value: latestMetrics?.["Active Investors"] || 0,
      icon: Users,
      format: (v: number) => v.toString(),
    },
    {
      title: "Update Opens",
      value: latestMetrics?.["Update Opens"] || 0,
      icon: Mail,
      format: (v: number) => v.toString(),
    },
    {
      title: "Deck Views",
      value: latestMetrics?.["Deck Views"] || 0,
      icon: Eye,
      format: (v: number) => v.toString(),
    }
  ];

  const contentMetrics = [
    {
      title: "Active Decks",
      value: latestMetrics?.["Active Decks"] || 0,
      icon: Presentation,
      format: (v: number) => v.toString(),
    },
    {
      title: "Protected Content",
      value: latestMetrics?.["Protected Content"] || 0,
      icon: Lock,
      format: (v: number) => v.toString(),
    },
    {
      title: "Shared Links",
      value: latestMetrics?.["Shared Links"] || 0,
      icon: Share2,
      format: (v: number) => v.toString(),
    }
  ];

  const recentActivity = [
    {
      type: "update",
      title: "Q4 2023 Investor Update",
      metric: "85% open rate",
      time: "2 hours ago"
    },
    {
      type: "deck",
      title: "Series A Pitch Deck",
      metric: "12 new views",
      time: "5 hours ago"
    },
    {
      type: "investor",
      title: "New Investor Added",
      metric: "Total: 45 investors",
      time: "1 day ago"
    },
    {
      type: "share",
      title: "Financial Model Shared",
      metric: "3 accesses",
      time: "2 days ago"
    }
  ];

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Investor Relations Overview</h1>
      </div>

      <h2 className="text-xl font-semibold">Engagement Metrics</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {metrics.map((metric) => (
          <Card key={metric.title} className="hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium flex items-center gap-2">
                <metric.icon className="h-4 w-4" />
                {metric.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metric.format(metric.value)}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <h2 className="text-xl font-semibold mt-8">Content Overview</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {contentMetrics.map((metric) => (
          <Card key={metric.title} className="hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium flex items-center gap-2">
                <metric.icon className="h-4 w-4" />
                {metric.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metric.format(metric.value)}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
        <Card className="col-span-1">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Investor Engagement Trends</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={filteredData}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-50" />
                  <XAxis
                    dataKey="date"
                    tickFormatter={(value) => format(parseISO(value), "MMM dd")}
                    stroke="#888888"
                  />
                  <YAxis stroke="#888888" />
                  <Tooltip
                    labelFormatter={(value) =>
                      format(parseISO(value as string), "MMM dd, yyyy")
                    }
                    contentStyle={{
                      backgroundColor: "rgba(0, 0, 0, 0.8)",
                      border: "none",
                      borderRadius: "4px",
                      color: "#fff",
                    }}
                  />
                  <Area
                    type="monotone"
                    name="Update Opens"
                    dataKey="Update Opens"
                    stroke="#8884d8"
                    fill="#8884d8"
                    fillOpacity={0.2}
                    strokeWidth={2}
                  />
                  <Area
                    type="monotone"
                    name="Deck Views"
                    dataKey="Deck Views"
                    stroke="#82ca9d"
                    fill="#82ca9d"
                    fillOpacity={0.2}
                    strokeWidth={2}
                  />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-1">
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentActivity.map((activity, index) => (
                <div key={index} className="flex items-center gap-4 p-2 hover:bg-gray-50 rounded-lg">
                  {activity.type === "update" && <Mail className="h-5 w-5 text-blue-500" />}
                  {activity.type === "deck" && <Presentation className="h-5 w-5 text-purple-500" />}
                  {activity.type === "investor" && <Users className="h-5 w-5 text-green-500" />}
                  {activity.type === "share" && <Share2 className="h-5 w-5 text-orange-500" />}
                  <div className="flex-1">
                    <p className="font-medium">{activity.title}</p>
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-gray-500">{activity.metric}</span>
                      <span className="text-xs text-gray-400">•</span>
                      <span className="text-sm text-gray-500">{activity.time}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Content Performance</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={filteredData}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-50" />
                  <XAxis
                    dataKey="date"
                    tickFormatter={(value) => format(parseISO(value), "MMM dd")}
                    stroke="#888888"
                  />
                  <YAxis stroke="#888888" />
                  <Tooltip
                    labelFormatter={(value) =>
                      format(parseISO(value as string), "MMM dd, yyyy")
                    }
                    contentStyle={{
                      backgroundColor: "rgba(0, 0, 0, 0.8)",
                      border: "none",
                      borderRadius: "4px",
                      color: "#fff",
                    }}
                  />
                  <Bar
                    name="Active Decks"
                    dataKey="Active Decks"
                    fill="#8884d8"
                    radius={[4, 4, 0, 0]}
                  />
                  <Bar
                    name="Protected Content"
                    dataKey="Protected Content"
                    fill="#82ca9d"
                    radius={[4, 4, 0, 0]}
                  />
                  <Bar
                    name="Shared Links"
                    dataKey="Shared Links"
                    fill="#ffc658"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
