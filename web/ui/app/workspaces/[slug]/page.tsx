"use client"

import { Button } from "@/components/ui/button"
import { useParams } from "next/navigation"

export default function Home() {

  const params = useParams()

  return (
    <div className="flex flex-col items-center gap-1 text-center">
      <h3 className="text-2xl font-bold tracking-tight">
        You have no updates
      </h3>
      <p className="text-sm text-muted-foreground">
        {params.slug}
      </p>
      <Button className="mt-4">Create your first update</Button>
    </div>
  )
}
