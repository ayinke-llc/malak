import {
  RiArchiveStackLine,
  RiBook3Line,
  RiContactsLine,
  RiDashboardHorizontalLine,
  RiHome2Line,
  RiMoneyDollarCircleLine,
  RiPieChartLine,
  RiPlug2Line,
  RiSettingsLine
} from "@remixicon/react";

export const links = [
  {
    title: "Overview",
    url: "/overview",
    icon: RiHome2Line,
  },
  {
    title: "Updates",
    url: "/updates",
    icon: RiArchiveStackLine,
  },
  {
    title: "Decks",
    url: "/decks",
    icon: RiBook3Line,
  },
  {
    title: "Contacts",
    url: "/contacts",
    icon: RiContactsLine,
  },
  {
    title: "Integrations",
    url: "/integrations",
    icon: RiPlug2Line
  },
  {
    title: "Metrics Dashboards",
    url: "/dashboards",
    icon: RiDashboardHorizontalLine
  },
  {
    title: "Fundraising",
    url: "/fundraising",
    icon: RiMoneyDollarCircleLine,
  },
  {
    title: "Captable",
    url: "/captable",
    icon: RiPieChartLine,
  },
  {
    title: "Settings",
    url: "/settings",
    icon: RiSettingsLine
  },
];
