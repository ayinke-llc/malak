import { MalakUser, MalakWorkspace } from '@/client/Api';
import create from 'zustand';
import { persist } from 'zustand/middleware';

type UserState = {
  token: string | null
  user: MalakUser | null
  // this is the currently active workspace for the user
  workspace: MalakWorkspace | null
}

type Actions = {
  isAuthenticated: () => boolean
  setUser: (user: MalakUser) => void
  setWorkspace: (workspace: MalakWorkspace) => void
  setToken: (token: string) => void
  logout: () => void
}

const useAuthStore = create(
  persist<UserState & Actions>(
    (set, get) => ({
      user: null,
      token: null,
      workspace: null,
      isAuthenticated: (): boolean => {
        const { user, token } = get()
        return user !== null && token !== null
      },
      setUser: (user: MalakUser) => set({ user }),
      setWorkspace: (workspace: MalakWorkspace) => set({ workspace }),
      setToken: (token: string) => set({ token }),
      logout: (): void => set({ user: null, token: null, workspace: null })
    }), {
    "name": "auth",
  })
)

export default useAuthStore;
