"use client"

import { useInitializeUserData } from "@/hooks/user";

export default function UserProvider({ children }: { children: React.ReactNode }) {
  useInitializeUserData()
  return children;
}
