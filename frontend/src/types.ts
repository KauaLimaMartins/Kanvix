export type ID = string;

export interface User {
  id: ID;
  name: string;
  avatarColor: string;
}

export interface Workspace {
  id: ID;
  name: string;
  icon?: string; // lucide icon name
  color?: string; // hex color
}

export interface Project {
  id: ID;
  workspaceId: ID;
  name: string;
  description?: string;
}

export interface Column {
  id: ID;
  projectId: ID;
  name: string;
  order: number;
}

export interface Task {
  id: ID;
  projectId: ID;
  columnId: ID;
  title: string;
  description: string; // HTML from TipTap
  labels: string[]; // Label IDs
  dueDate?: string;
  assigneeId?: ID;
  order: number;
  createdAt: string;
}

export interface Label {
  id: ID;
  workspaceId: ID;
  name: string;
  color: string; // hex
}