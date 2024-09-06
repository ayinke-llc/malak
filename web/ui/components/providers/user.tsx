"use client"

import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserProvider({ children }: { children: React.ReactNode }) {
  const { token, setUser, isAuthenticated, user, logout } = useAuthStore.getState()
  const router = useRouter()

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

  if (user !== null && user.metadata.current_workspace.startsWith("0000")) {
    router.push("/workspaces/new")
  }

  return children;
}
