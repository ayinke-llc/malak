"use client";

import { ServerAPIStatus } from "@/client/Api";
import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import useWorkspacesStore from "@/store/workspace";
import { AxiosError } from "axios";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";

// Setup interceptors outside component to ensure they're always available
client.instance.interceptors.request.use(
  async (config) => {
    // Skip adding auth header for shared routes or when no token is available
    if (config.url?.startsWith('/shared')) {
      return config;
    }

    const token = useAuthStore.getState().token;
    if (!token) {
      return config;
    }

    config.headers.Authorization = `Bearer ${token}`;
    return config;
  },
  (error) => Promise.reject(error),
);

client.instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const pathname = window.location.pathname;
      // Skip auth checks for shared routes
      if (!pathname?.startsWith('/shared')) {
        if (error.response.status === 401) {
          useAuthStore.getState().logout();
          window.location.href = '/login';
        }

        if (error.response.status === 402) {
          window.location.href = "/settings?tab=billing";
        }
      }
    }

    return Promise.reject(error);
  }
);

export default function UserProvider({
  children,
}: { children: React.ReactNode }) {
  const [loading, setLoading] = useState(true);
  const pathname = usePathname();

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
    // Skip authentication for shared routes
    if (pathname?.startsWith('/shared')) {
      setLoading(false);
      return;
    }

    if (!isRehydrated) {
      return;
    }

    // Redirect from login to home if authenticated
    if (pathname === '/login' && isAuthenticated()) {
      router.push('/');
      return;
    }

    if (!isAuthenticated()) {
      handleLogout();
      return;
    }

    // Validate token exists before making the request
    if (!token) {
      handleLogout();
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
  }, [token, isRehydrated, clear, isAuthenticated, logout, router, setCurrent, setUser, setWorkspaces, pathname]);

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

