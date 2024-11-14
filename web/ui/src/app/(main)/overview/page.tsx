"use client";
import { overviews } from "@/data/overview-data";
import type { OverviewData } from "@/data/schema";
import React from "react";
import type { DateRange } from "react-day-picker";

export type PeriodValue = "previous-period" | "last-year" | "no-comparison";

const categories: {
  title: keyof OverviewData;
  type: "currency" | "unit";
}[] = [
    {
      title: "Rows read",
      type: "unit",
    },
    {
      title: "Rows written",
      type: "unit",
    },
    {
      title: "Queries",
      type: "unit",
    },
    {
      title: "Payments completed",
      type: "currency",
    },
    {
      title: "Sign ups",
      type: "unit",
    },
    {
      title: "Logins",
      type: "unit",
    },
  ];

export type KpiEntry = {
  title: string;
  percentage: number;
  current: number;
  allowed: number;
  unit?: string;
};

const data: KpiEntry[] = [
  {
    title: "Updates viewed",
    percentage: 48.1,
    current: 48.1,
    allowed: 100,
    unit: "M",
  },
  {
    title: "Updates reacted to",
    percentage: 78.3,
    current: 78.3,
    allowed: 100,
    unit: "M",
  },
  {
    title: "CRM",
    percentage: 26,
    current: 5.2,
    allowed: 20,
    unit: "GB",
  },
];

const data2: KpiEntry[] = [
  {
    title: "Weekly active users",
    percentage: 21.7,
    current: 21.7,
    allowed: 100,
    unit: "%",
  },
  {
    title: "Total users",
    percentage: 70,
    current: 28,
    allowed: 40,
  },
  {
    title: "Uptime",
    percentage: 98.3,
    current: 98.3,
    allowed: 100,
    unit: "%",
  },
];

export type KpiEntryExtended = Omit<
  KpiEntry,
  "current" | "allowed" | "unit"
> & {
  value: string;
  color: string;
};

const data3: KpiEntryExtended[] = [
  {
    title: "Base tier",
    percentage: 68.1,
    value: "$200",
    color: "bg-indigo-600 dark:bg-indigo-500",
  },
  {
    title: "On-demand charges",
    percentage: 20.8,
    value: "$61.1",
    color: "bg-purple-600 dark:bg-purple-500",
  },
  {
    title: "Caching",
    percentage: 11.1,
    value: "$31.9",
    color: "bg-gray-400 dark:bg-gray-600",
  },
];


export default function Overview() {

  return null
}
