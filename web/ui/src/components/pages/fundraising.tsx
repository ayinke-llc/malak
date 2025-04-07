"use client";

import KanbanBoard from "./KanbanBoard";

interface FundraisingProps {
  slug: string;
}

export default function Fundraising({ slug }: FundraisingProps) {
  return <KanbanBoard slug={slug} />;
} 
