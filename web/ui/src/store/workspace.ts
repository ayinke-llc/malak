import type { MalakWorkspace } from "@/client/Api";
import create from "zustand";
import { persist } from "zustand/middleware";

type WorkspaceState = {
  current: MalakWorkspace | null;
  workspaces: MalakWorkspace[];
};

type Actions = {
  setCurrent: (_workspace: MalakWorkspace) => void;
  setWorkspaces: (_workspaces: MalakWorkspace[]) => void;
  clear: () => void;
  appendWorkspaceAfterCreation: (_workspace: MalakWorkspace) => void
};

const useWorkspacesStore = create(
  persist<WorkspaceState & Actions>(
    (set) => ({
      current: null,
      workspaces: [],
      setCurrent: (workspace: MalakWorkspace) => set({ current: workspace }),
      setWorkspaces: (workspaces: MalakWorkspace[]) => set({ workspaces: workspaces || [] }),
      appendWorkspaceAfterCreation: (workspace: MalakWorkspace) =>
        set((state) => ({ workspaces: [...(state.workspaces || []), workspace] })),
      clear: () => set({ current: null, workspaces: [] }),
    }),
    {
      name: "workspace",
    },
  ),
);

export default useWorkspacesStore;
