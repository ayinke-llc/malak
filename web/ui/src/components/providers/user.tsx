"use client";

import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import useWorkspacesStore from "@/store/workspace";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserProvider({
  children,
}: { children: React.ReactNode }) {
  const { token, setUser, isAuthenticated, logout } = useAuthStore();
  const { setWorkspaces, setCurrent } = useWorkspacesStore();
  const router = useRouter();

  useEffect(() => {
    const requestInterceptor = client.instance.interceptors.request.use(
      async (config) => {
        if (isAuthenticated()) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error),
    );

    const responseInterceptor = client.instance.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response && error.response.status === 401) {
          logout();
          router.push("/login");
        }
        return Promise.reject(error);
      },
    );

    return () => {
      client.instance.interceptors.request.eject(requestInterceptor);
      client.instance.interceptors.response.eject(responseInterceptor);
    };
  }, [token, isAuthenticated, logout, router]);

  useEffect(() => {
    if (!isAuthenticated()) {
      logout();
      router.push("/login");
      return;
    }

    client.user
      .userList()
      .then((res) => {
        setUser(res.data.user);
        if (res.data.current_workspace !== undefined) {
          setCurrent(res.data.current_workspace);
        }
        setWorkspaces(res.data.workspaces);
      })
      .catch((err) => {
        console.log(err, "authenticate user");
      });
  }, [token, isAuthenticated]);

  return children;
}
