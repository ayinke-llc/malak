import type { MalakUser } from "@/client/Api";
import { create } from "zustand";
import { persist } from "zustand/middleware";

type UserState = {
  token: string | null;
  user: MalakUser | null;
};

type Actions = {
  isAuthenticated: () => boolean;
  setUser: (user: MalakUser) => void;
  setToken: (token: string) => void;
  logout: () => void;
  // Flag to track rehydration
  isRehydrated: boolean,
};

const useAuthStore = create(
  persist<UserState & Actions>(
    (set, get) => ({
      user: null,
      token: null,
      workspace: null,
      isAuthenticated: (): boolean => {
        const { user, token } = get();
        return user !== null && token !== null;
      },
      setUser: (user: MalakUser) => set({ user }),
      setToken: (token: string) => set({ token }),
      logout: (): void => set({ user: null, token: null }),
      isRehydrated: false,
    }),
    {
      name: "auth",
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.isRehydrated = true;
        }
      },
    },
  ),
);

export default useAuthStore;
