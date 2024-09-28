import type { MalakWorkspace } from "@/client/Api";
import create from "zustand";
import { persist } from "zustand/middleware";

type WorkspaceState = {
	current: MalakWorkspace | null;
	workspaces: MalakWorkspace[];
};

type Actions = {
	setCurrent: (workspace: MalakWorkspace) => void;
	setWorkspaces: (workspaces: MalakWorkspace[]) => void;
	clear: () => void;
};

const useWorkspacesStore = create(
	persist<WorkspaceState & Actions>(
		(set, get) => ({
			current: null,
			workspaces: [],
			setCurrent: (workspace: MalakWorkspace) => set({ current: workspace }),
			setWorkspaces: (workspaces: MalakWorkspace[]) => set({ workspaces }),
			clear: () => set({ current: null, workspaces: [] }),
		}),
		{
			name: "workspace",
		},
	),
);

export default useWorkspacesStore;
