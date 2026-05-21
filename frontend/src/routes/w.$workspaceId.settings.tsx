import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useState, useMemo } from "react";
import {
  LayoutGrid,
  Briefcase,
  Rocket,
  Heart,
  Star,
  Folder,
  Target,
  Sparkles,
  Flag,
  Zap,
  Globe,
  Code,
  type LucideIcon,
} from "lucide-react";

const ICONS: Record<string, LucideIcon> = {
  LayoutGrid,
  Briefcase,
  Rocket,
  Heart,
  Star,
  Folder,
  Target,
  Sparkles,
  Flag,
  Zap,
  Globe,
  Code,
};

const COLORS = [
  "#6366f1",
  "#ec4899",
  "#14b8a6",
  "#f59e0b",
  "#10b981",
  "#a855f7",
  "#f97316",
  "#0ea5e9",
  "#ef4444",
  "#84cc16",
];

export const Route = createFileRoute("/w/$workspaceId/settings")({
  head: () => ({ meta: [{ title: "Workspace settings — Kanvix" }] }),
  component: SettingsPage,
});

function SettingsPage() {
  const { workspaceId } = Route.useParams();
  const router = useRouter();
  const workspaces = useAppStore((s) => s.workspaces);
  const userRole = useAppStore((s) => s.userRole);
  const updateWorkspace = useAppStore((s) => s.updateWorkspace);
  const deleteWorkspace = useAppStore((s) => s.deleteWorkspace);
  const ws = useMemo(() => workspaces.find((w) => w.id === workspaceId), [workspaces, workspaceId]);

  const [name, setName] = useState(ws?.name ?? "");
  const [color, setColor] = useState(ws?.color ?? COLORS[0]);
  const [icon, setIcon] = useState(ws?.icon ?? "LayoutGrid");
  const [saved, setSaved] = useState(false);

  if (!ws) {
    return <div className="p-8 text-sm text-muted-foreground">Workspace not found.</div>;
  }

  const canEdit = (ws.role ?? userRole) === "admin";
  if (!canEdit) {
    return <div className="p-8 text-sm text-muted-foreground">Forbidden.</div>;
  }

  const PreviewIcon = ICONS[icon] ?? LayoutGrid;

  const save = () => {
    if (!name.trim()) return;
    void updateWorkspace(workspaceId, { name: name.trim(), color, icon });
    setSaved(true);
    setTimeout(() => setSaved(false), 1500);
  };

  return (
    <div className="mx-auto max-w-3xl p-8">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold tracking-tight">Workspace settings</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Customize how this workspace appears across Kanvix.
        </p>
      </div>

      <div className="rounded-xl border border-border bg-card p-6">
        <div className="mb-6 flex items-center gap-3">
          <div
            className="flex h-12 w-12 items-center justify-center rounded-lg text-white shadow-sm"
            style={{ background: color }}
          >
            <PreviewIcon className="h-6 w-6" />
          </div>
          <div>
            <div className="text-sm font-semibold">{name || "Untitled workspace"}</div>
            <div className="text-xs text-muted-foreground">Live preview</div>
          </div>
        </div>

        <div className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="ws-name">Title</Label>
            <Input
              id="ws-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Workspace name"
              className="max-w-md"
            />
          </div>

          <div className="space-y-2">
            <Label>Color</Label>
            <div className="flex flex-wrap gap-2">
              {COLORS.map((c) => (
                <button
                  key={c}
                  type="button"
                  onClick={() => setColor(c)}
                  className={`h-8 w-8 rounded-full border-2 transition-transform ${
                    color === c ? "scale-110 border-foreground" : "border-transparent"
                  }`}
                  style={{ background: c }}
                  aria-label={`Color ${c}`}
                />
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <Label>Icon</Label>
            <div className="flex flex-wrap gap-2">
              {Object.entries(ICONS).map(([key, Ico]) => {
                const selected = key === icon;
                return (
                  <button
                    key={key}
                    type="button"
                    onClick={() => setIcon(key)}
                    className={`flex h-10 w-10 items-center justify-center rounded-lg border transition-colors ${
                      selected
                        ? "border-foreground bg-accent text-foreground"
                        : "border-border text-muted-foreground hover:bg-accent/60"
                    }`}
                    aria-label={key}
                  >
                    <Ico className="h-4.5 w-4.5" />
                  </button>
                );
              })}
            </div>
          </div>
        </div>

        <div className="mt-6 flex items-center gap-2">
          <Button onClick={save}>Save changes</Button>
          {saved && <span className="text-xs text-muted-foreground">Saved ✓</span>}
        </div>
      </div>

      <div className="mt-6 rounded-xl border border-destructive/30 bg-card p-6">
        <h2 className="text-sm font-semibold text-destructive">Danger zone</h2>
        <p className="mt-1 text-xs text-muted-foreground">
          Deleting a workspace removes all its projects, columns, tasks, and labels.
        </p>
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant="destructive" className="mt-3">
              Delete workspace
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Delete workspace?</AlertDialogTitle>
              <AlertDialogDescription>
                Delete workspace “{ws.name}”? This cannot be undone.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction asChild>
                <Button
                  variant="destructive"
                  onClick={() => {
                    void deleteWorkspace(workspaceId);
                    router.navigate({ to: "/workspaces" });
                  }}
                >
                  Delete
                </Button>
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  );
}
