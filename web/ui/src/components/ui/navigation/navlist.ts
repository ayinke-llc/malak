import {
  RiArchiveStackLine,
  RiBook3Line,
  RiChatNewLine,
  RiContactsLine,
  RiHome2Line,
  RiHomeGearFill,
  RiHomeGearLine,
  RiLinkM,
  RiPlug2Line,
  RiSettings6Line,
  RiSettingsLine,
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
    title: "Fundraising",
    url: "/fundraising",
    icon: RiContactsLine,
  },
  {
    title: "Captable",
    url: "/captable",
    icon: RiContactsLine,
  },
  {
    title: "Integrations",
    url: "/integrations",
    icon: RiPlug2Line
  },
  {
    title: "Settings",
    url: "/settings",
    icon: RiSettingsLine
  },
];

export const shortcuts = [
  {
    title: "Analytics of your last update",
    url: "#",
    icon: RiLinkM,
  },
  {
    title: "Fundraising",
    url: "#",
    icon: RiLinkM,
  },
  {
    title: "Captable",
    url: "#",
    icon: RiLinkM,
  },
] as const;

