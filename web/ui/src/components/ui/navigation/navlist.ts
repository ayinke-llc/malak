import {
  RiArchiveStackLine,
  RiBook3Line,
  RiChatNewLine,
  RiContactsLine,
  RiHome2Line,
  RiLinkM,
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
    items: [
      {
        title: "All updates",
        url: "/updates"
      },
      {
        title: "View your last update",
        url: "/updates/latest",
        icon: RiChatNewLine,
      },
    ],
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
    items: [
      {
        title: "View contacts",
        url: "/contacts",
        icon: RiChatNewLine,
      },
      {
        title: "View contact lists",
        url: "/contacts"
      }
    ],
  },
  {
    title: "Fundraising",
    url: "#",
    icon: RiContactsLine,
  },
  {
    title: "Captable",
    url: "#",
    icon: RiContactsLine,
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

