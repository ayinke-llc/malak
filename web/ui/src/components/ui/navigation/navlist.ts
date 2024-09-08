import { RemixiconComponentType, RiArrowDownLine, RiBook3Line, RiHome2Line, RiLinkM } from "@remixicon/react";

export const links = [
  {
    name: "Home",
    href: "/",
    icon: RiHome2Line
  },
  {
    name: "Updates",
    href: "/updates",
    icon: RiArrowDownLine
  },
  {
    name: "Decks",
    href: "/decks",
    icon: RiBook3Line
  }
] as const

export const shortcuts = [
  {
    name: "Add new user",
    href: "#",
    icon: RiLinkM,
  },
  {
    name: "Workspace usage",
    href: "#",
    icon: RiLinkM,
  },
  {
    name: "Cost spend control",
    href: "#",
    icon: RiLinkM,
  },
  {
    name: "Analytics – Most recent Analytics",
    href: "#",
    icon: RiLinkM,
  },
] as const

