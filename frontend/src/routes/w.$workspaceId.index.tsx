import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { FolderKanban, Plus } from "lucide-react";
import { Link, createFileRoute, useRouter } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/services/api";
import { useAppStore } from "@/store/useAppStore";
import { useState } from "react";

export const Route = createFileRoute("/w/$workspaceId/")({
  component: WorkspacePage,
});

function WorkspacePage() {
  const { workspaceId } = Route.useParams();
  const workspace = useAppStore((s) => s.workspaces.find((w) => w.id === workspaceId));
  const userRole = useAppStore((s) => s.userRole);
  const allProjects = useAppStore((s) => s.projects);
  const tasks = useAppStore((s) => s.tasks);
  const projects = allProjects.filter((p) => p.workspaceId === workspaceId);
  const addProject = useAppStore((s) => s.addProject);
  const router = useRouter();

  const statsQuery = useQuery({
    queryKey: ["workspaceStats", workspaceId],
    queryFn: () => api.stats.get(workspaceId),
    staleTime: 30_000,
    retry: 1,
  });

  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");

  if (!workspace) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
        Workspace not found.
      </div>
    );
  }

  const canCreateProjects = (workspace.role ?? userRole) === "admin";

  return (
    <div className="p-8">
      <div className="mb-6 flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-semibold">{workspace.name}</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {projects.length} {projects.length === 1 ? "project" : "projects"}
            {" • "}
            {statsQuery.isLoading
              ? "… tasks"
              : `${statsQuery.data?.taskCount ?? tasks.filter((t) => projects.some((p) => p.id === t.projectId)).length} tasks`}
          </p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button disabled={!canCreateProjects}>
              <Plus className="mr-2 h-4 w-4" /> New project
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>New project</DialogTitle>
            </DialogHeader>
            <Input
              autoFocus
              placeholder="Project name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            <DialogFooter>
              <Button variant="ghost" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button
                onClick={() => {
                  if (!canCreateProjects) return;
                  if (!name.trim()) return;
                  void (async () => {
                    const p = await addProject(workspaceId, name.trim());
                    setName("");
                    setOpen(false);
                    router.navigate({
                      to: "/w/$workspaceId/p/$projectId",
                      params: { workspaceId, projectId: p.id },
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

      {projects.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <FolderKanban className="mb-3 h-8 w-8 text-muted-foreground" />
          <p className="text-sm text-muted-foreground">No projects yet.</p>
          <Button className="mt-4" onClick={() => setOpen(true)}>
            <Plus className="mr-2 h-4 w-4" /> Create your first project
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {projects.map((p, i) => {
            const count =
              statsQuery.data?.tasksByProject?.[p.id] ?? tasks.filter((t) => t.projectId === p.id).length;
            return (
              <motion.div
                key={p.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.25, delay: i * 0.04 }}
                whileHover={{ y: -3 }}
              >
                <Link
                  to="/w/$workspaceId/p/$projectId"
                  params={{ workspaceId, projectId: p.id }}
                  className="group block rounded-xl border border-border bg-card p-5 shadow-sm transition-shadow hover:shadow-lg"
                >
                  <div className="flex items-center gap-2">
                    <FolderKanban className="h-4 w-4 text-muted-foreground" />
                    <h2 className="font-medium tracking-tight">{p.name}</h2>
                  </div>
                  {p.description && (
                    <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                      {p.description}
                    </p>
                  )}
                  <p className="mt-3 text-xs text-muted-foreground">
                    {count} {count === 1 ? "task" : "tasks"}
                  </p>
                </Link>
              </motion.div>
            );
          })}
        </div>
      )}
    </div>
  );
}
