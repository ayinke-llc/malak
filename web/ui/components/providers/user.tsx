"use client"

import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserProvider({ children }: { children: React.ReactNode }) {
  const path = usePathname()
  const { setUser, isAuthenticated, user, logout } = useAuthStore()
  const router = useRouter()

  // for some reason, it is not always set except you do it like this
  // TODO: resolve this hack
  const token = useAuthStore.getState().token;

  client.instance.interceptors.request.use(
    async (config) => {
      if (isAuthenticated()) {
        config.headers['Authorization'] = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  client.instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response && error.response.status === 401) {
        logout()
        router.push("/login")
      }

      return Promise.reject(error);
    }
  );

  useEffect(() => {
    if (!isAuthenticated()) {
      logout()
      router.push("/login")
      return
    }
  }, [token])

  useEffect(() => {
    if (isAuthenticated()) {
      client.user.userList().then(res => {
        setUser(res.data.user)
      }).catch((err) => {
        console.log(err, "authenticate user")
      })
    }
  }, [token])

  if (user !== null && user.roles.length === 0) {
    router.push("/workspaces/new")
  }

  return children;
}
