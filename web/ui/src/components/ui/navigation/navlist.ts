import {
  RiArchiveStackLine,
  RiBook3Line,
  RiContactsLine,
  RiDashboardHorizontalLine,
  RiHome2Line,
  RiLineChartLine,
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
    title: "Dashboards",
    url: "/dashboards",
    icon: RiDashboardHorizontalLine
  },
  {
    title: "Metrics",
    url: "/metrics",
    icon: RiLineChartLine
  },
  {
    title: "Fundraising",
    url: "/fundraising",
    icon: RiMoneyDollarCircleLine,
    comingSoon: true,
  },
  {
    title: "Captable",
    url: "/captable",
    icon: RiPieChartLine,
    comingSoon: true,
  },
  {
    title: "Settings",
    url: "/settings",
    icon: RiSettingsLine
  },
];
