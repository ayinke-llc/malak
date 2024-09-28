import {
  RiArchiveStackLine,
  RiBook3Line,
  RiContactsLine,
  RiHome2Line,
  RiLinkM,
} from "@remixicon/react";

export const links = [
  {
    name: "Home",
    href: "/",
    icon: RiHome2Line,
  },
  {
    name: "Updates",
    href: "/updates",
    icon: RiArchiveStackLine,
  },
  {
    name: "Decks",
    href: "/decks",
    icon: RiBook3Line,
  },
  {
    name: "Contacts",
    href: "/contacts",
    icon: RiContactsLine,
  },
] as const;

export const shortcuts = [
  {
    name: "Send a new update",
    href: "/updates/new",
    icon: RiLinkM,
  },
  {
    name: "Analytics of your last update",
    href: "#",
    icon: RiLinkM,
  },
  {
    name: "Fundraising",
    href: "#",
    icon: RiLinkM,
  },
  {
    name: "Captable",
    href: "#",
    icon: RiLinkM,
  },
] as const;
