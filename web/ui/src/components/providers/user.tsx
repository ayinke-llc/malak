"use client";

import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import useWorkspacesStore from "@/store/workspace";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserProvider({
	children,
}: { children: React.ReactNode }) {
	const { token, setUser, isAuthenticated, logout } = useAuthStore.getState();

	const { setWorkspaces, setCurrent } = useWorkspacesStore.getState();

	const router = useRouter();

	client.instance.interceptors.request.use(
		async (config) => {
			if (isAuthenticated()) {
				config.headers.Authorization = `Bearer ${token}`;
			}
			return config;
		},
		(error) => Promise.reject(error),
	);

	client.instance.interceptors.response.use(
		(response) => response,
		(error) => {
			if (error.response && error.response.status === 401) {
				logout();
				router.push("/login");
			}

			return Promise.reject(error);
		},
	);

	useEffect(() => {
		if (!isAuthenticated()) {
			logout();
			router.push("/login");
			return;
		}
	}, [token]);

	useEffect(() => {
		if (isAuthenticated()) {
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
		}
	}, [token]);

	return children;
}
