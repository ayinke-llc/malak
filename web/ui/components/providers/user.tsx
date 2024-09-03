"use client"

import { useInitializeUserData } from "@/hooks/user";
import useAuthStore from "@/store/auth";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserProvider({ children }: { children: React.ReactNode }) {
  useInitializeUserData()
  const path = usePathname()
  const { isAuthenticated, user } = useAuthStore()
  const router = useRouter()

  useEffect(() => {
    console.log(isAuthenticated(), "HERE")
    if (!isAuthenticated() && path !== "/login") {
      router.push("/login")
      return
    }
  }, [user])

  return children;
}
