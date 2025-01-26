"use client";

import { ServerAPIStatus } from "@/client/Api";
import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import useWorkspacesStore from "@/store/workspace";
import { AxiosError, AxiosResponse } from "axios";
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

  const { setWorkspaces, setCurrent } = useWorkspacesStore();
  const router = useRouter();

  useEffect(() => {
    if (!isRehydrated) {
      return
    }

    if (!isAuthenticated()) {
      logout();
      router.push("/login");
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
        if (error.response && error.response.status === 401) {
          logout();
          router.push("/login");
        }
        return Promise.reject(error);
      },
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
        console.log(err, token, "THERE")
        toast.error(err?.response?.data?.message)
        setLoading(true)
      });

    return () => {
      client.instance.interceptors.request.eject(requestInterceptor);
      client.instance.interceptors.response.eject(responseInterceptor);
    };
  }, [token, isRehydrated]);


  if (loading) {
    return <div>Loading...</div>;
  }

  return children;
}
