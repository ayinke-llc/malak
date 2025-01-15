"use client";
import { overviews } from "@/data/overview-data";
import type { OverviewData } from "@/data/schema";
import React from "react";
import type { DateRange } from "react-day-picker";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar } from "@/components/ui/calendar";
import { addDays, format, isWithinInterval, parseISO, subDays } from "date-fns";
import { 
  ArrowDown, 
  ArrowUp, 
  Users, 
  DollarSign, 
  Building, 
  FileText,
  TrendingUp,
  Presentation,
  BarChart,
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
} from "recharts";

export default function Overview() {
  const [date, setDate] = React.useState<DateRange | undefined>({
    from: subDays(new Date(), 30),
    to: new Date(),
  });

  const filteredData = React.useMemo(() => {
    if (!date?.from || !date?.to) return overviews;
    return overviews.filter((item) => {
      const itemDate = parseISO(item.date);
      return isWithinInterval(itemDate, { start: date.from!, end: date.to! });
    });
  }, [date]);

  const latestMetrics = filteredData[filteredData.length - 1];
  const previousMetrics = filteredData[filteredData.length - 2];

  const getPercentageChange = (current: number, previous: number) => {
    if (!previous) return 0;
    return ((current - previous) / previous) * 100;
  };

  const formatCurrency = (value: number) => {
    if (value >= 1000000) {
      return `$${(value / 1000000).toFixed(1)}M`;
    }
    if (value >= 1000) {
      return `$${(value / 1000).toFixed(1)}K`;
    }
    return `$${value}`;
  };

  const metrics = [
    {
      title: "Active Investors",
      value: latestMetrics?.["Active Investors"] || 0,
      change: getPercentageChange(
        latestMetrics?.["Active Investors"] || 0,
        previousMetrics?.["Active Investors"] || 0
      ),
      icon: Users,
      format: (v: number) => v.toString(),
    },
    {
      title: "Total Funding",
      value: latestMetrics?.["Total Funding"] || 0,
      change: getPercentageChange(
        latestMetrics?.["Total Funding"] || 0,
        previousMetrics?.["Total Funding"] || 0
      ),
      icon: DollarSign,
      format: formatCurrency,
    },
    {
      title: "Company Valuation",
      value: latestMetrics?.["Company Valuation"] || 0,
      change: getPercentageChange(
        latestMetrics?.["Company Valuation"] || 0,
        previousMetrics?.["Company Valuation"] || 0
      ),
      icon: TrendingUp,
      format: formatCurrency,
    },
    {
      title: "Team Size",
      value: latestMetrics?.["Team Members"] || 0,
      change: getPercentageChange(
        latestMetrics?.["Team Members"] || 0,
        previousMetrics?.["Team Members"] || 0
      ),
      icon: Building,
      format: (v: number) => v.toString(),
    },
    {
      title: "New Pitches",
      value: latestMetrics?.["New Pitches"] || 0,
      change: getPercentageChange(
        latestMetrics?.["New Pitches"] || 0,
        previousMetrics?.["New Pitches"] || 0
      ),
      icon: Presentation,
      format: (v: number) => v.toString(),
    },
    {
      title: "Document Updates",
      value: latestMetrics?.["Document Updates"] || 0,
      change: getPercentageChange(
        latestMetrics?.["Document Updates"] || 0,
        previousMetrics?.["Document Updates"] || 0
      ),
      icon: FileText,
      format: (v: number) => v.toString(),
    },
  ];

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Overview</h1>
        <Calendar
          mode="range"
          selected={date}
          onSelect={setDate}
          className="rounded-md border shadow-sm"
        />
      </div>

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
              <div className="flex items-center gap-1 mt-1">
                {metric.change > 0 ? (
                  <ArrowUp className="h-4 w-4 text-green-500" />
                ) : (
                  <ArrowDown className="h-4 w-4 text-red-500" />
                )}
                <span
                  className={`text-sm ${
                    metric.change > 0 ? "text-green-500" : "text-red-500"
                  }`}
                >
                  {Math.abs(metric.change).toFixed(1)}%
                </span>
                <span className="text-sm text-muted-foreground">vs previous day</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="col-span-1">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Funding & Valuation</CardTitle>
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
                    formatter={(value: number, name: string) => [formatCurrency(value), name]}
                    contentStyle={{
                      backgroundColor: "rgba(0, 0, 0, 0.8)",
                      border: "none",
                      borderRadius: "4px",
                      color: "#fff",
                    }}
                  />
                  <Area
                    type="monotone"
                    name="Total Funding"
                    dataKey="Total Funding"
                    stroke="#82ca9d"
                    fill="#82ca9d"
                    fillOpacity={0.2}
                    strokeWidth={2}
                  />
                  <Area
                    type="monotone"
                    name="Company Valuation"
                    dataKey="Company Valuation"
                    stroke="#8884d8"
                    fill="#8884d8"
                    fillOpacity={0.2}
                    strokeWidth={2}
                  />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-1">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Activity Trends</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={filteredData}>
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
                  <Line
                    type="monotone"
                    name="New Pitches"
                    dataKey="New Pitches"
                    stroke="#ffc658"
                    strokeWidth={2}
                    dot={false}
                    activeDot={{ r: 6 }}
                  />
                  <Line
                    type="monotone"
                    name="Document Updates"
                    dataKey="Document Updates"
                    stroke="#ff7300"
                    strokeWidth={2}
                    dot={false}
                    activeDot={{ r: 6 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Growth Overview</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={filteredData}>
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
                  <Line
                    type="monotone"
                    name="Active Investors"
                    dataKey="Active Investors"
                    stroke="#8884d8"
                    strokeWidth={2}
                    dot={false}
                    activeDot={{ r: 6 }}
                  />
                  <Line
                    type="monotone"
                    name="Team Members"
                    dataKey="Team Members"
                    stroke="#82ca9d"
                    strokeWidth={2}
                    dot={false}
                    activeDot={{ r: 6 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
