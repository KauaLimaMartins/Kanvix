import { Link, useRouter } from "@tanstack/react-router";
import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  ChevronsUpDown,
  Plus,
  FolderKanban,
  Trash2,
  Tag,
  LayoutDashboard,
  Settings,
  Users,
} from "lucide-react";
import { useState } from "react";
import { motion } from "framer-motion";

export function Sidebar() {
  const workspaces = useAppStore((s) => s.workspaces);
  const projects = useAppStore((s) => s.projects);
  const userRole = useAppStore((s) => s.userRole);
  const addWorkspace = useAppStore((s) => s.addWorkspace);
  const addProject = useAppStore((s) => s.addProject);
  const deleteProject = useAppStore((s) => s.deleteProject);
  const router = useRouter();

  const match = router.state.matches.find((m) => m.params && "workspaceId" in m.params);
  const currentWorkspaceId =
    (match?.params as { workspaceId?: string } | undefined)?.workspaceId ?? workspaces[0]?.id;
  const currentWorkspace = workspaces.find((w) => w.id === currentWorkspaceId);
  const currentProjectId = (
    router.state.matches.find((m) => m.params && "projectId" in m.params)?.params as
      | { projectId?: string }
      | undefined
  )?.projectId;
  const pathname = router.state.location.pathname;
  const isLabelsRoute = pathname.endsWith("/labels");
  const isSettingsRoute = pathname.endsWith("/settings");
  const isUsersRoute = pathname.endsWith("/users");
  const canManageUsers = (currentWorkspace?.role ?? userRole) === "admin";

  const workspaceProjects = projects.filter((p) => p.workspaceId === currentWorkspaceId);

  const [newWsOpen, setNewWsOpen] = useState(false);
  const [newWsName, setNewWsName] = useState("");
  const [newProjOpen, setNewProjOpen] = useState(false);
  const [newProjName, setNewProjName] = useState("");

  return (
    <aside className="flex w-64 shrink-0 flex-col border-r border-sidebar-border bg-sidebar text-sidebar-foreground">
      <Link
        to="/workspaces"
        className="flex items-center gap-2 border-b border-sidebar-border px-4 py-3.5"
      >
        <div className="flex h-7 w-7 items-center justify-center rounded-md bg-gradient-to-br from-indigo-500 to-fuchsia-500 text-white shadow-sm">
          <span className="text-[13px] font-bold">K</span>
        </div>
        <span className="text-sm font-semibold tracking-tight">Kanvix</span>
      </Link>

      <div className="px-2 pt-3">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="flex w-full items-center justify-between rounded-md px-2 py-2 text-left transition-colors hover:bg-sidebar-accent">
              <div className="flex min-w-0 items-center gap-2">
                {currentWorkspace?.color && (
                  <span
                    className="h-2 w-2 shrink-0 rounded-full"
                    style={{ background: currentWorkspace.color }}
                  />
                )}
                <div className="min-w-0">
                  <div className="text-[10px] uppercase tracking-wider text-sidebar-foreground/50">
                    Workspace
                  </div>
                  <div className="truncate text-sm font-medium">
                    {currentWorkspace?.name ?? "—"}
                  </div>
                </div>
              </div>
              <ChevronsUpDown className="h-4 w-4 opacity-60" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-56">
            {workspaces.map((w) => (
              <DropdownMenuItem
                key={w.id}
                onClick={() =>
                  router.navigate({ to: "/w/$workspaceId", params: { workspaceId: w.id } })
                }
              >
                <span
                  className="mr-2 h-2 w-2 rounded-full"
                  style={{ background: w.color ?? "#6366f1" }}
                />
                {w.name}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => router.navigate({ to: "/workspaces" })}>
              All workspaces
            </DropdownMenuItem>
            {userRole === "admin" && (
              <DropdownMenuItem
                onSelect={(e) => {
                  e.preventDefault();
                  setNewWsOpen(true);
                }}
              >
                <Plus className="mr-2 h-4 w-4" /> New workspace
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {currentWorkspaceId && (
        <nav className="mt-3 space-y-0.5 px-2">
          <Link
            to="/w/$workspaceId"
            params={{ workspaceId: currentWorkspaceId }}
            className={`flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors ${
              !currentProjectId && !isLabelsRoute && !isSettingsRoute && !isUsersRoute
                ? "bg-sidebar-accent text-sidebar-accent-foreground"
                : "hover:bg-sidebar-accent/60"
            }`}
          >
            <LayoutDashboard className="h-4 w-4 opacity-60" />
            Projects
          </Link>
          <Link
            to="/w/$workspaceId/labels"
            params={{ workspaceId: currentWorkspaceId }}
            className={`flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors ${
              isLabelsRoute
                ? "bg-sidebar-accent text-sidebar-accent-foreground"
                : "hover:bg-sidebar-accent/60"
            }`}
          >
            <Tag className="h-4 w-4 opacity-60" />
            Labels
          </Link>
          {canManageUsers && (
            <Link
              to="/w/$workspaceId/settings"
              params={{ workspaceId: currentWorkspaceId }}
              className={`flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors ${
                isSettingsRoute
                  ? "bg-sidebar-accent text-sidebar-accent-foreground"
                  : "hover:bg-sidebar-accent/60"
              }`}
            >
              <Settings className="h-4 w-4 opacity-60" />
              Settings
            </Link>
          )}
          {canManageUsers && (
            <Link
              to="/w/$workspaceId/users"
              params={{ workspaceId: currentWorkspaceId }}
              className={`flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors ${
                isUsersRoute
                  ? "bg-sidebar-accent text-sidebar-accent-foreground"
                  : "hover:bg-sidebar-accent/60"
              }`}
            >
              <Users className="h-4 w-4 opacity-60" />
              Users
            </Link>
          )}
        </nav>
      )}

      <div className="mt-5 flex items-center justify-between px-4">
        <span className="text-[10px] font-medium uppercase tracking-wider text-sidebar-foreground/50">
          Projects
        </span>
        <Dialog open={newProjOpen} onOpenChange={setNewProjOpen}>
          <DialogTrigger asChild>
            <Button size="icon" variant="ghost" className="h-6 w-6" disabled={!canManageUsers}>
              <Plus className="h-4 w-4" />
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>New project</DialogTitle>
            </DialogHeader>
            <Input
              autoFocus
              placeholder="Project name"
              value={newProjName}
              onChange={(e) => setNewProjName(e.target.value)}
            />
            <DialogFooter>
              <Button variant="ghost" onClick={() => setNewProjOpen(false)}>
                Cancel
              </Button>
              <Button
                onClick={() => {
                  if (!newProjName.trim() || !currentWorkspaceId) return;
                  void (async () => {
                    const p = await addProject(currentWorkspaceId, newProjName.trim());
                    setNewProjName("");
                    setNewProjOpen(false);
                    router.navigate({
                      to: "/w/$workspaceId/p/$projectId",
                      params: { workspaceId: currentWorkspaceId, projectId: p.id },
                    });
                  })();
                }}
              >
                Create
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <nav className="mt-1 flex-1 space-y-0.5 overflow-y-auto px-2 pb-3">
        {workspaceProjects.length === 0 && (
          <p className="px-2 py-4 text-xs text-sidebar-foreground/60">No projects yet.</p>
        )}
        {workspaceProjects.map((p) => {
          const active = p.id === currentProjectId;
          return (
            <motion.div
              key={p.id}
              initial={{ opacity: 0, x: -4 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.15 }}
              className="group flex items-center"
            >
              <Link
                to="/w/$workspaceId/p/$projectId"
                params={{ workspaceId: p.workspaceId, projectId: p.id }}
                className={`flex flex-1 items-center gap-2 truncate rounded-md px-2 py-1.5 text-sm transition-all duration-150 ${
                  active
                    ? "bg-sidebar-accent text-sidebar-accent-foreground"
                    : "hover:bg-sidebar-accent/60 hover:translate-x-0.5"
                }`}
              >
                <FolderKanban className="h-4 w-4 shrink-0 opacity-60" />
                <span className="truncate">{p.name}</span>
              </Link>
              <Button
                size="icon"
                variant="ghost"
                className="h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
                onClick={() => {
                  if (confirm(`Delete project "${p.name}"?`)) {
                    void deleteProject(p.id);
                    if (active && currentWorkspaceId) {
                      router.navigate({
                        to: "/w/$workspaceId",
                        params: { workspaceId: currentWorkspaceId },
                      });
                    }
                  }
                }}
              >
                <Trash2 className="h-3.5 w-3.5" />
              </Button>
            </motion.div>
          );
        })}
      </nav>

      <Dialog open={newWsOpen} onOpenChange={setNewWsOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>New workspace</DialogTitle>
          </DialogHeader>
          <Input
            autoFocus
            placeholder="Workspace name"
            value={newWsName}
            onChange={(e) => setNewWsName(e.target.value)}
          />
          <DialogFooter>
            <Button variant="ghost" onClick={() => setNewWsOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => {
                if (!newWsName.trim()) return;
                void (async () => {
                  const w = await addWorkspace(newWsName.trim());
                  setNewWsName("");
                  setNewWsOpen(false);
                  router.navigate({ to: "/w/$workspaceId", params: { workspaceId: w.id } });
                })();
              }}
            >
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </aside>
  );
}
