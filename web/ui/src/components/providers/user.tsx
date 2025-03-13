"use client";

import { ServerAPIStatus } from "@/client/Api";
import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import useWorkspacesStore from "@/store/workspace";
import { AxiosError } from "axios";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function UserProvider({
  children,
}: { children: React.ReactNode }) {
  const [loading, setLoading] = useState(true);

  const token = useAuthStore((state) => state.token);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const setUser = useAuthStore((state) => state.setUser);
  const logout = useAuthStore((state) => state.logout);
  const isRehydrated = useAuthStore(state => state.isRehydrated);

  const clear = useWorkspacesStore(state => state.clear);
  const { setWorkspaces, setCurrent } = useWorkspacesStore();
  const router = useRouter();

  const handleLogout = () => {
    clear();
    logout();
    setLoading(false);
    router.push("/login");
  };

  useEffect(() => {
    if (!isRehydrated) {
      return;
    }

    if (!isAuthenticated()) {
      handleLogout();
      return;
    }

    const requestInterceptor = client.instance.interceptors.request.use(
      async (config) => {
        config.headers.Authorization = `Bearer ${token}`;
        return config;
      },
      (error) => Promise.reject(error),
    );

    const responseInterceptor = client.instance.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          if (error.response.status === 401) {
            handleLogout();
          }

          if (error.response.status === 402) {
            router.push("/settings?tab=billing");
          }
        }

        return Promise.reject(error);
      }
    );

    client.user
      .userList()
      .then((res) => {
        setUser(res.data.user);

        if (res.data.current_workspace !== undefined) {
          setCurrent(res.data.current_workspace);
        }

        setWorkspaces(res.data.workspaces);
        setLoading(false);
      })
      .catch((err: AxiosError<ServerAPIStatus>) => {
        if (err?.response?.status === 402) {
          setLoading(false);
          router.push("/settings?tab=billing");
          return;
        }

        toast.error(err?.response?.data?.message);
        handleLogout();
      });

    return () => {
      client.instance.interceptors.request.eject(requestInterceptor);
      client.instance.interceptors.response.eject(responseInterceptor);
    };
  }, [token, isRehydrated, clear, isAuthenticated, logout, router, setCurrent, setUser, setWorkspaces]);

  if (loading) {
    return (
      <>
        <div className="fixed inset-0 bg-background flex items-center justify-center">
          <div className="flex flex-col items-center gap-2">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        </div>
      </>
    );
  }

  return <>{children}</>;
}

