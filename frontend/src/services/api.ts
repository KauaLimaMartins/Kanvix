import type { Column, Label, Project, Subtask, Task, User, Workspace } from "@/types";

export type BootstrapResponse = {
  user: User & { email: string; role: string };
  workspaces: Workspace[];
  projects: Project[];
  columns: Column[];
  tasks: Task[];
  labels: Label[];
  users: User[];
};

type Json = Record<string, unknown> | unknown[] | string | number | boolean | null;

export class ApiError extends Error {
  status: number;
  payload?: unknown;

  constructor(message: string, status: number, payload?: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

const API_BASE = ((import.meta as unknown as { env?: Record<string, string> }).env
  ?.VITE_API_BASE_URL ||
  "/api") as string;

function joinUrl(base: string, path: string) {
  if (!base) return path;
  const b = base.endsWith("/") ? base.slice(0, -1) : base;
  return `${b}${path}`;
}

async function request<T>(path: string, init: RequestInit & { json?: Json } = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.json !== undefined) headers.set("content-type", "application/json");

  const res = await fetch(joinUrl(API_BASE, path), {
    ...init,
    headers,
    credentials: "include",
    body: init.json !== undefined ? JSON.stringify(init.json) : init.body,
  });

  const contentType = res.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");
  const payload = isJson
    ? await res.json().catch(() => undefined)
    : await res.text().catch(() => "");

  if (!res.ok) {
    const msg =
      (payload &&
        typeof payload === "object" &&
        "error" in (payload as Record<string, unknown>) &&
        typeof (payload as Record<string, unknown>).error === "string" &&
        ((payload as Record<string, unknown>).error as string)) ||
      `Request failed (${res.status})`;
    throw new ApiError(msg, res.status, payload);
  }

  return payload as T;
}

export const api = {
  auth: {
    setup: () => request<{ needsFirstSignup: boolean }>("/auth/setup", { method: "GET" }),
    firstSignup: (email: string, password: string, name: string) =>
      request<{ user: User & { email?: string; role?: string } }>("/auth/first-signup", {
        method: "POST",
        json: { email, password, name },
      }),
    login: (email: string, password?: string) =>
      request<{ user: User & { email?: string; role?: string } }>("/auth/login", {
        method: "POST",
        json: { email, password: password ?? "" },
      }),
    me: () => request<{ user: User & { email?: string; role?: string } }>("/auth/me", { method: "GET" }),
    logout: () => request<{ ok: true }>("/auth/logout", { method: "POST" }),
  },

  bootstrap: () => request<BootstrapResponse>("/bootstrap", { method: "GET" }),

  workspaces: {
    list: () => request<{ workspaces: Workspace[] }>("/workspaces", { method: "GET" }),
    create: (name: string) =>
      request<{ workspace: Workspace }>("/workspaces", { method: "POST", json: { name } }),
    update: (id: string, patch: Partial<Workspace>) =>
      request<{ workspace: Workspace }>(`/workspaces/${id}`, { method: "PATCH", json: patch }),
    delete: (id: string) => request<{ ok: true }>(`/workspaces/${id}`, { method: "DELETE" }),
  },

  projects: {
    listByWorkspace: (workspaceId: string) =>
      request<{ projects: Project[] }>(`/workspaces/${workspaceId}/projects`, { method: "GET" }),
    create: (workspaceId: string, name: string) =>
      request<{ project: Project; columns: Column[] }>(`/workspaces/${workspaceId}/projects`, {
        method: "POST",
        json: { name },
      }),
    update: (id: string, patch: Partial<Project>) =>
      request<{ project: Project }>(`/projects/${id}`, { method: "PATCH", json: patch }),
    delete: (id: string) => request<{ ok: true }>(`/projects/${id}`, { method: "DELETE" }),
  },

  columns: {
    listByProject: (projectId: string) =>
      request<{ columns: Column[] }>(`/projects/${projectId}/columns`, { method: "GET" }),
    create: (projectId: string, name: string) =>
      request<{ column: Column }>(`/projects/${projectId}/columns`, {
        method: "POST",
        json: { name },
      }),
    update: (id: string, patch: Partial<Column>) =>
      request<{ column: Column }>(`/columns/${id}`, { method: "PATCH", json: patch }),
    delete: (id: string) => request<{ ok: true }>(`/columns/${id}`, { method: "DELETE" }),
  },

  tasks: {
    listByProject: (projectId: string) =>
      request<{ tasks: Task[] }>(`/projects/${projectId}/tasks`, { method: "GET" }),
    create: (projectId: string, columnId: string, title: string) =>
      request<{ task: Task }>(`/projects/${projectId}/tasks`, {
        method: "POST",
        json: { columnId, title },
      }),
    get: (id: string) => request<{ task: Task }>(`/tasks/${id}`, { method: "GET" }),
    update: (id: string, patch: Partial<Task>) => {
      const json: Record<string, unknown> = { ...patch };
      if ("dueDate" in patch && patch.dueDate === undefined) json.dueDate = null;
      if ("assigneeId" in patch && patch.assigneeId === undefined) json.assigneeId = null;
      return request<{ task: Task }>(`/tasks/${id}`, { method: "PATCH", json });
    },
    delete: (id: string) => request<{ ok: true }>(`/tasks/${id}`, { method: "DELETE" }),
    move: (id: string, toColumnId: string, toIndex: number) =>
      request<{ ok: true }>(`/tasks/${id}/move`, {
        method: "POST",
        json: { toColumnId, toIndex },
      }),
  },

  subtasks: {
    listByTask: (taskId: string) =>
      request<{ subtasks: Subtask[] }>(`/tasks/${taskId}/subtasks`, { method: "GET" }),
    create: (taskId: string, title: string) =>
      request<{ subtask: Subtask }>(`/tasks/${taskId}/subtasks`, { method: "POST", json: { title } }),
    setDone: (subtaskId: string, done: boolean) =>
      request<{ subtask: Subtask }>(`/subtasks/${subtaskId}`, { method: "PATCH", json: { done } }),
  },

  labels: {
    listByWorkspace: (workspaceId: string) =>
      request<{ labels: Label[] }>(`/workspaces/${workspaceId}/labels`, { method: "GET" }),
    create: (workspaceId: string, name: string, color: string) =>
      request<{ label: Label }>(`/workspaces/${workspaceId}/labels`, {
        method: "POST",
        json: { name, color },
      }),
    update: (id: string, patch: Partial<Label>) =>
      request<{ label: Label }>(`/labels/${id}`, { method: "PATCH", json: patch }),
    delete: (id: string) => request<{ ok: true }>(`/labels/${id}`, { method: "DELETE" }),
  },

  stats: {
    get: (workspaceId: string) =>
      request<{
        workspaceId: string;
        projectCount: number;
        taskCount: number;
        tasksByProject: Record<string, number>;
      }>(`/workspaces/${workspaceId}/stats`, { method: "GET" }),
  },

  search: {
    query: (workspaceId: string, q: string, limit?: number) => {
      const params = new URLSearchParams();
      params.set("q", q);
      if (limit !== undefined) params.set("limit", String(limit));
      return request<{
        query: string;
        results: Array<
          | { type: "project"; id: string; name: string; workspaceId: string }
          | { type: "task"; id: string; title: string; projectId: string; workspaceId: string }
        >;
      }>(`/workspaces/${workspaceId}/search?${params.toString()}`, { method: "GET" });
    },
  },

  users: {
    listByWorkspace: (workspaceId: string) =>
      request<{
        users: Array<User & { email: string; role: string }>;
      }>(`/workspaces/${workspaceId}/users`, { method: "GET" }),
    createInWorkspace: (workspaceId: string, u: { email: string; name: string; password: string; role: string }) =>
      request<{ user: User & { email: string; role: string } }>(`/workspaces/${workspaceId}/users`, {
        method: "POST",
        json: u,
      }),
    updateRoleInWorkspace: (workspaceId: string, userId: string, role: string) =>
      request<{ ok: true }>(`/workspaces/${workspaceId}/users/${userId}`, {
        method: "PATCH",
        json: { role },
      }),
  },
};
