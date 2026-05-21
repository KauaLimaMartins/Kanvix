import { Briefcase, Building2, LayoutGrid, LogOut, Plus, Sparkles } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Navigate, createFileRoute, useRouter } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import { ClientOnly } from "@/components/ClientOnly";
import { Input } from "@/components/ui/input";
import { motion } from "framer-motion";
import { useAppStore } from "@/store/useAppStore";
import { useState } from "react";

export const Route = createFileRoute("/workspaces")({
  head: () => ({
    meta: [
      { title: "Choose a workspace — Kanvix" },
      { name: "description", content: "Pick a workspace to continue." },
    ],
  }),
  component: WorkspacesPage,
});

const ICON_MAP: Record<string, React.ComponentType<{ className?: string }>> = {
  Building2,
  Briefcase,
  Sparkles,
  LayoutGrid,
};

function WorkspacesPage() {
  return (
    <ClientOnly
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
          Loading…
        </div>
      }
    >
      <WorkspacesInner />
    </ClientOnly>
  );
}

function WorkspacesInner() {
  const authStatus = useAppStore((s) => s.authStatus);
  const workspaces = useAppStore((s) => s.workspaces);
  const projects = useAppStore((s) => s.projects);
  const addWorkspace = useAppStore((s) => s.addWorkspace);
  const logout = useAppStore((s) => s.logout);
  const userEmail = useAppStore((s) => s.userEmail);
  const userRole = useAppStore((s) => s.userRole);
  const router = useRouter();

  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");

  if (authStatus !== "authed") return <Navigate to="/login" />;

  return (
    <div className="min-h-screen bg-background">
      <header className="flex items-center justify-between border-b border-border px-6 py-4">
        <div className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-indigo-500 to-fuchsia-500 text-white">
            <span className="text-sm font-bold">K</span>
          </div>
          <span className="font-semibold tracking-tight">Kanvix</span>
        </div>
        <div className="flex items-center gap-3">
          {userEmail && <span className="text-xs text-muted-foreground">{userEmail}</span>}
          <Button
            variant="ghost"
            size="sm"
            onClick={async () => {
              await logout();
              router.navigate({ to: "/login" });
            }}
          >
            <LogOut className="mr-2 h-4 w-4" />
            Sign out
          </Button>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-16">
        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
          className="mb-10"
        >
          <h1 className="text-3xl font-semibold tracking-tight">Choose a workspace</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Pick where you want to work today. You can switch anytime.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {workspaces.map((w, i) => {
            const Icon = ICON_MAP[w.icon ?? "LayoutGrid"] ?? LayoutGrid;
            const count = projects.filter((p) => p.workspaceId === w.id).length;
            return (
              <motion.button
                key={w.id}
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: i * 0.05 }}
                whileHover={{ y: -3 }}
                whileTap={{ scale: 0.98 }}
                onClick={() =>
                  router.navigate({ to: "/w/$workspaceId", params: { workspaceId: w.id } })
                }
                className="group relative overflow-hidden rounded-2xl border border-border bg-card p-5 text-left shadow-sm transition-shadow hover:shadow-lg"
              >
                <div
                  className="absolute inset-x-0 top-0 h-1"
                  style={{ background: w.color ?? "#6366f1" }}
                />
                <div
                  className="flex h-11 w-11 items-center justify-center rounded-xl text-white shadow-sm"
                  style={{ background: w.color ?? "#6366f1" }}
                >
                  <Icon className="h-5 w-5" />
                </div>
                <div className="mt-4 font-semibold tracking-tight">{w.name}</div>
                <div className="mt-1 text-xs text-muted-foreground">
                  {count} {count === 1 ? "project" : "projects"}
                </div>
              </motion.button>
            );
          })}

          {userRole === "admin" && (
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <motion.button
                  initial={{ opacity: 0, y: 12 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3, delay: workspaces.length * 0.05 }}
                  whileHover={{ y: -3 }}
                  whileTap={{ scale: 0.98 }}
                  className="flex h-full min-h-[160px] flex-col items-center justify-center rounded-2xl border border-dashed border-border bg-muted/20 p-5 text-sm text-muted-foreground transition-colors hover:bg-muted/40"
                >
                  <Plus className="mb-2 h-5 w-5" />
                  New workspace
                </motion.button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>New workspace</DialogTitle>
                </DialogHeader>
                <Input
                  autoFocus
                  placeholder="Workspace name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
                <DialogFooter>
                  <Button variant="ghost" onClick={() => setOpen(false)}>
                    Cancel
                  </Button>
                  <Button
                    onClick={async () => {
                      if (!name.trim()) return;
                      const w = await addWorkspace(name.trim());
                      setName("");
                      setOpen(false);
                      router.navigate({ to: "/w/$workspaceId", params: { workspaceId: w.id } });
                    }}
                  >
                    Create
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          )}
        </div>
      </main>
    </div>
  );
}
