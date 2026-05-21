import { ApiError, api } from "@/services/api";
import type { Column, ID, Label, Project, Subtask, Task, User, Workspace } from "@/types";

import { create } from "zustand";
import { toast } from "sonner";

function errorMessage(e: unknown, fallback: string) {
  if (e instanceof ApiError) return e.message;
  if (e instanceof Error) return e.message;
  return fallback;
}

interface AppState {
  workspaces: Workspace[];
  projects: Project[];
  columns: Column[];
  tasks: Task[];
  subtasks: Subtask[];
  labels: Label[];
  users: User[];
  lastCreatedColumnId: string | null;
  lastCreatedTaskId: string | null;
  lastCreatedSubtaskId: string | null;
  theme: "dark" | "light";
  authStatus: "unknown" | "guest" | "authed";
  userEmail: string | null;
  userRole: "admin" | "member" | null;
  needsFirstSignup: boolean | null;
  isLoading: boolean;
  error: string | null;
  hydrate: () => Promise<void>;
  // auth
  login: (email: string, password?: string) => Promise<void>;
  firstSignup: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  // workspaces
  addWorkspace: (name: string) => Promise<Workspace>;
  renameWorkspace: (id: ID, name: string) => Promise<void>;
  updateWorkspace: (id: ID, patch: Partial<Workspace>) => Promise<void>;
  deleteWorkspace: (id: ID) => Promise<void>;
  // projects
  addProject: (workspaceId: ID, name: string) => Promise<Project>;
  renameProject: (id: ID, name: string) => Promise<void>;
  deleteProject: (id: ID) => Promise<void>;
  // columns
  addColumn: (projectId: ID, name: string) => Promise<Column>;
  renameColumn: (id: ID, name: string) => Promise<void>;
  deleteColumn: (id: ID) => Promise<void>;
  // tasks
  addTask: (projectId: ID, columnId: ID, title: string) => Promise<Task>;
  updateTask: (id: ID, patch: Partial<Task>) => Promise<void>;
  deleteTask: (id: ID) => Promise<void>;
  moveTask: (taskId: ID, toColumnId: ID, toIndex: number) => Promise<void>;
  // subtasks
  fetchSubtasks: (taskId: ID) => Promise<void>;
  addSubtask: (taskId: ID, title: string) => Promise<Subtask>;
  setSubtaskDone: (subtaskId: ID, done: boolean) => Promise<void>;
  // labels
  addLabel: (workspaceId: ID, name: string, color: string) => Promise<Label>;
  updateLabel: (id: ID, patch: Partial<Label>) => Promise<void>;
  deleteLabel: (id: ID) => Promise<void>;
  // theme
  setTheme: (t: "dark" | "light") => void;
}

export const useAppStore = create<AppState>()((set, get) => ({
  workspaces: [],
  projects: [],
  columns: [],
  tasks: [],
  subtasks: [],
  labels: [],
  users: [],
  lastCreatedColumnId: null,
  lastCreatedTaskId: null,
  lastCreatedSubtaskId: null,
  theme: "dark",
  authStatus: "unknown",
  userEmail: null,
  userRole: null,
  needsFirstSignup: null,
  isLoading: false,
  error: null,

  hydrate: async () => {
    set({ isLoading: true, error: null });
    try {
      const data = await api.bootstrap();
      set({
        authStatus: "authed",
        userEmail: data.user.email ?? null,
        userRole: (data.user.role as "admin" | "member" | undefined) ?? null,
        needsFirstSignup: false,
        workspaces: data.workspaces,
        projects: data.projects,
        columns: data.columns,
        tasks: data.tasks.map((t) => ({ ...t, labels: Array.isArray(t.labels) ? t.labels : [] })),
        subtasks: [],
        labels: data.labels,
        users: data.users,
        isLoading: false,
      });
    } catch (e) {
      if (e instanceof ApiError && e.status === 401) {
        let needsFirstSignup: boolean | null = null;
        try {
          const res = await api.auth.setup();
          needsFirstSignup = res.needsFirstSignup;
        } catch {
          needsFirstSignup = null;
        }
        set({
          authStatus: "guest",
          userEmail: null,
          userRole: null,
          needsFirstSignup,
          workspaces: [],
          projects: [],
          columns: [],
          tasks: [],
          subtasks: [],
          labels: [],
          users: [],
          isLoading: false,
          error: null,
        });
        return;
      }
      const msg = errorMessage(e, "Request failed");
      set({ isLoading: false, error: msg });
      toast.error(msg);
    }
  },

  login: async (email, password) => {
    set({ isLoading: true, error: null });
    try {
      await api.auth.login(email, password);
      const data = await api.bootstrap();
      set({
        authStatus: "authed",
        userEmail: data.user.email ?? email,
        userRole: (data.user.role as "admin" | "member" | undefined) ?? null,
        needsFirstSignup: false,
        workspaces: data.workspaces,
        projects: data.projects,
        columns: data.columns,
        tasks: data.tasks.map((t) => ({ ...t, labels: Array.isArray(t.labels) ? t.labels : [] })),
        subtasks: [],
        labels: data.labels,
        users: data.users,
        isLoading: false,
      });
    } catch (e) {
      const msg = errorMessage(e, "Login failed");
      set({ isLoading: false, error: msg });
      toast.error(msg);
      throw e;
    }
  },

  firstSignup: async (name, email, password) => {
    set({ isLoading: true, error: null });
    try {
      await api.auth.firstSignup(email, password, name);
      const data = await api.bootstrap();
      set({
        authStatus: "authed",
        userEmail: data.user.email ?? email,
        userRole: (data.user.role as "admin" | "member" | undefined) ?? null,
        needsFirstSignup: false,
        workspaces: data.workspaces,
        projects: data.projects,
        columns: data.columns,
        tasks: data.tasks.map((t) => ({ ...t, labels: Array.isArray(t.labels) ? t.labels : [] })),
        subtasks: [],
        labels: data.labels,
        users: data.users,
        isLoading: false,
      });
    } catch (e) {
      const msg = errorMessage(e, "Signup failed");
      set({ isLoading: false, error: msg });
      toast.error(msg);
      throw e;
    }
  },

  logout: async () => {
    set({ isLoading: true, error: null });
    try {
      await api.auth.logout();
    } finally {
      set({
        authStatus: "guest",
        userEmail: null,
        userRole: null,
        needsFirstSignup: null,
        workspaces: [],
        projects: [],
        columns: [],
        tasks: [],
        subtasks: [],
        labels: [],
        users: [],
        isLoading: false,
      });
    }
  },

  addWorkspace: async (name) => {
    try {
      const res = await api.workspaces.create(name);
      set((s) => ({ workspaces: [...s.workspaces, res.workspace] }));
      return res.workspace;
    } catch (e) {
      toast.error(errorMessage(e, "Could not create workspace"));
      throw e;
    }
  },
  renameWorkspace: async (id, name) => {
    try {
      const res = await api.workspaces.update(id, { name });
      set((s) => ({ workspaces: s.workspaces.map((w) => (w.id === id ? res.workspace : w)) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not update workspace"));
      throw e;
    }
  },
  updateWorkspace: async (id, patch) => {
    try {
      const res = await api.workspaces.update(id, patch);
      set((s) => ({ workspaces: s.workspaces.map((w) => (w.id === id ? res.workspace : w)) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not update workspace"));
      throw e;
    }
  },
  deleteWorkspace: async (id) => {
    try {
      await api.workspaces.delete(id);
      set((s) => {
        const projectIds = s.projects.filter((p) => p.workspaceId === id).map((p) => p.id);
        return {
          workspaces: s.workspaces.filter((w) => w.id !== id),
          projects: s.projects.filter((p) => p.workspaceId !== id),
          columns: s.columns.filter((c) => !projectIds.includes(c.projectId)),
          tasks: s.tasks.filter((t) => !projectIds.includes(t.projectId)),
          labels: s.labels.filter((l) => l.workspaceId !== id),
        };
      });
    } catch (e) {
      toast.error(errorMessage(e, "Could not delete workspace"));
      throw e;
    }
  },

  addProject: async (workspaceId, name) => {
    try {
      const res = await api.projects.create(workspaceId, name);
      set((s) => ({
        projects: [...s.projects, res.project],
        columns: [...s.columns, ...res.columns],
      }));
      return res.project;
    } catch (e) {
      toast.error(errorMessage(e, "Could not create project"));
      throw e;
    }
  },
  renameProject: async (id, name) => {
    try {
      const res = await api.projects.update(id, { name });
      set((s) => ({ projects: s.projects.map((p) => (p.id === id ? res.project : p)) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not update project"));
      throw e;
    }
  },
  deleteProject: async (id) => {
    try {
      await api.projects.delete(id);
      set((s) => ({
        projects: s.projects.filter((p) => p.id !== id),
        columns: s.columns.filter((c) => c.projectId !== id),
        tasks: s.tasks.filter((t) => t.projectId !== id),
      }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not delete project"));
      throw e;
    }
  },

  addColumn: async (projectId, name) => {
    try {
      const res = await api.columns.create(projectId, name);
      const id = res.column.id;
      set((s) => ({ columns: [...s.columns, res.column], lastCreatedColumnId: id }));
      setTimeout(() => {
        if (get().lastCreatedColumnId === id) set({ lastCreatedColumnId: null });
      }, 900);
      return res.column;
    } catch (e) {
      toast.error(errorMessage(e, "Could not create column"));
      throw e;
    }
  },
  renameColumn: async (id, name) => {
    try {
      const res = await api.columns.update(id, { name });
      set((s) => ({ columns: s.columns.map((c) => (c.id === id ? res.column : c)) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not update column"));
      throw e;
    }
  },
  deleteColumn: async (id) => {
    try {
      await api.columns.delete(id);
      set((s) => ({
        columns: s.columns.filter((c) => c.id !== id),
        tasks: s.tasks.filter((t) => t.columnId !== id),
      }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not delete column"));
      throw e;
    }
  },

  addTask: async (projectId, columnId, title) => {
    try {
      const res = await api.tasks.create(projectId, columnId, title);
      const id = res.task.id;
      set((s) => ({ tasks: [...s.tasks, res.task], lastCreatedTaskId: id }));
      setTimeout(() => {
        if (get().lastCreatedTaskId === id) set({ lastCreatedTaskId: null });
      }, 900);
      return res.task;
    } catch (e) {
      toast.error(errorMessage(e, "Could not create task"));
      throw e;
    }
  },
  updateTask: async (id, patch) => {
    try {
      const res = await api.tasks.update(id, patch);
      set((s) => ({ tasks: s.tasks.map((t) => (t.id === id ? res.task : t)) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not update task"));
      throw e;
    }
  },
  deleteTask: async (id) => {
    try {
      await api.tasks.delete(id);
      set((s) => ({ tasks: s.tasks.filter((t) => t.id !== id) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not delete task"));
      throw e;
    }
  },

  fetchSubtasks: async (taskId) => {
    try {
      const res = await api.subtasks.listByTask(taskId);
      set((s) => ({
        subtasks: [...s.subtasks.filter((st) => st.taskId !== taskId), ...res.subtasks],
      }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not load subtasks"));
      throw e;
    }
  },

  addSubtask: async (taskId, title) => {
    try {
      const res = await api.subtasks.create(taskId, title);
      const id = res.subtask.id;
      set((s) => ({ subtasks: [...s.subtasks, res.subtask], lastCreatedSubtaskId: id }));
      setTimeout(() => {
        if (get().lastCreatedSubtaskId === id) set({ lastCreatedSubtaskId: null });
      }, 900);
      return res.subtask;
    } catch (e) {
      toast.error(errorMessage(e, "Could not create subtask"));
      throw e;
    }
  },

  setSubtaskDone: async (subtaskId, done) => {
    const before = get().subtasks;
    set((s) => ({ subtasks: s.subtasks.map((st) => (st.id === subtaskId ? { ...st, done } : st)) }));
    try {
      const res = await api.subtasks.setDone(subtaskId, done);
      set((s) => ({ subtasks: s.subtasks.map((st) => (st.id === subtaskId ? res.subtask : st)) }));
    } catch (e) {
      set({ subtasks: before });
      toast.error(errorMessage(e, "Could not update subtask"));
      throw e;
    }
  },

  moveTask: async (taskId, toColumnId, toIndex) => {
    const before = get().tasks;
    set((s) => {
      const task = s.tasks.find((t) => t.id === taskId);
      if (!task) return {};
      const fromColumnId = task.columnId;
      const sourceList = s.tasks
        .filter((t) => t.columnId === fromColumnId && t.id !== taskId)
        .sort((a, b) => a.order - b.order);
      const targetList =
        fromColumnId === toColumnId
          ? sourceList
          : s.tasks
              .filter((t) => t.columnId === toColumnId && t.id !== taskId)
              .sort((a, b) => a.order - b.order);

      const inserted = [...targetList];
      const clampedIndex = Math.max(0, Math.min(toIndex, inserted.length));
      inserted.splice(clampedIndex, 0, { ...task, columnId: toColumnId });

      const reorderedTarget = inserted.map((t, i) => ({ ...t, order: i }));
      const reorderedSource =
        fromColumnId === toColumnId ? [] : sourceList.map((t, i) => ({ ...t, order: i }));

      const updatedIds = new Set([
        ...reorderedTarget.map((t) => t.id),
        ...reorderedSource.map((t) => t.id),
      ]);
      const merged = [
        ...s.tasks.filter((t) => !updatedIds.has(t.id)),
        ...reorderedTarget,
        ...reorderedSource,
      ];
      return { tasks: merged };
    });
    try {
      await api.tasks.move(taskId, toColumnId, toIndex);
    } catch (e) {
      set({ tasks: before });
      toast.error(errorMessage(e, "Could not move task"));
      throw e;
    }
  },

  addLabel: async (workspaceId, name, color) => {
    try {
      const res = await api.labels.create(workspaceId, name, color);
      set((s) => ({ labels: [...s.labels, res.label] }));
      return res.label;
    } catch (e) {
      toast.error(errorMessage(e, "Could not create label"));
      throw e;
    }
  },
  updateLabel: async (id, patch) => {
    try {
      const res = await api.labels.update(id, patch);
      set((s) => ({ labels: s.labels.map((l) => (l.id === id ? res.label : l)) }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not update label"));
      throw e;
    }
  },
  deleteLabel: async (id) => {
    try {
      await api.labels.delete(id);
      set((s) => ({
        labels: s.labels.filter((l) => l.id !== id),
        tasks: s.tasks.map((t) => ({ ...t, labels: t.labels.filter((x) => x !== id) })),
      }));
    } catch (e) {
      toast.error(errorMessage(e, "Could not delete label"));
      throw e;
    }
  },

  setTheme: (theme) => set({ theme }),
}));
