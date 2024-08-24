"use client"

import Link from "next/link";
import { MalakLink } from "./dashboard";
import { FC, PropsWithChildren } from "react";

interface Props {
  data: MalakLink
}

const SideLink: FC<PropsWithChildren<Props>> = ({ data }: Props) => {
  const defaultClass = "flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary"
  // const activeClass = "flex items-center gap-3 rounded-lg bg-muted px-3 py-2 text-primary transition-all hover:text-primary"

  return (
    <Link
      href={data.link}
      className={defaultClass}
    >
      <data.icon className="h-4 w-4" />
      {data.name}
    </Link >
  )
}


const MobileLink: FC<PropsWithChildren<Props>> = ({ data }: Props) => {

  const defaultClass = `mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground`

  return (
    <Link
      href={data.link}
      className={defaultClass}
    >
      <data.icon className="h-4 w-4" />
      {data.name}
    </Link >
  )
}

export {
  SideLink,
  MobileLink
}
