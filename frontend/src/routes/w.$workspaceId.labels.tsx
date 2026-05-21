import { createFileRoute } from "@tanstack/react-router";
import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
import { useState } from "react";
import { Plus, Trash2, Tag, Pencil, Check, X } from "lucide-react";
import { motion } from "framer-motion";

const PRESET_COLORS = [
  "#6366f1",
  "#ec4899",
  "#14b8a6",
  "#f59e0b",
  "#10b981",
  "#a855f7",
  "#f97316",
  "#0ea5e9",
];

export const Route = createFileRoute("/w/$workspaceId/labels")({
  head: () => ({ meta: [{ title: "Labels — Kanvix" }] }),
  component: LabelsPage,
});

function LabelsPage() {
  const { workspaceId } = Route.useParams();
  const allLabels = useAppStore((s) => s.labels);
  const tasks = useAppStore((s) => s.tasks);
  const labels = allLabels.filter((l) => l.workspaceId === workspaceId);
  const addLabel = useAppStore((s) => s.addLabel);
  const updateLabel = useAppStore((s) => s.updateLabel);
  const deleteLabel = useAppStore((s) => s.deleteLabel);

  const [name, setName] = useState("");
  const [color, setColor] = useState(PRESET_COLORS[0]);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [editColor, setEditColor] = useState(PRESET_COLORS[0]);

  const create = () => {
    if (!name.trim()) return;
    void addLabel(workspaceId, name.trim(), color);
    setName("");
  };

  return (
    <div className="mx-auto max-w-3xl p-8">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold tracking-tight">Labels</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Centrally manage labels available to all tasks in this workspace.
        </p>
      </div>

      <div className="rounded-xl border border-border bg-card p-4">
        <div className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
          Create label
        </div>
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <Input
            placeholder="Label name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && create()}
            className="max-w-xs"
          />
          <div className="flex items-center gap-1">
            {PRESET_COLORS.map((c) => (
              <button
                key={c}
                type="button"
                onClick={() => setColor(c)}
                className={`h-6 w-6 rounded-full border-2 transition-transform ${
                  color === c ? "scale-110 border-foreground" : "border-transparent"
                }`}
                style={{ background: c }}
                aria-label={`Color ${c}`}
              />
            ))}
          </div>
          <Button onClick={create}>
            <Plus className="mr-1.5 h-4 w-4" /> Add
          </Button>
        </div>
      </div>

      <div className="mt-6 rounded-xl border border-border bg-card">
        {labels.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <Tag className="mb-3 h-8 w-8 text-muted-foreground" />
            <p className="text-sm text-muted-foreground">No labels yet.</p>
          </div>
        ) : (
          <ul className="divide-y divide-border">
            {labels.map((l) => {
              const useCount = tasks.filter((t) => t.labels.includes(l.id)).length;
              const editing = editingId === l.id;
              return (
                <motion.li
                  key={l.id}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="flex items-center gap-3 px-4 py-3"
                >
                  {editing ? (
                    <>
                      <div className="flex items-center gap-1">
                        {PRESET_COLORS.map((c) => (
                          <button
                            key={c}
                            type="button"
                            onClick={() => setEditColor(c)}
                            className={`h-5 w-5 rounded-full border-2 ${
                              editColor === c ? "border-foreground" : "border-transparent"
                            }`}
                            style={{ background: c }}
                          />
                        ))}
                      </div>
                      <Input
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                        className="h-8 max-w-xs"
                        autoFocus
                      />
                      <Button
                        size="sm"
                        onClick={() => {
                          if (editName.trim()) {
                            void updateLabel(l.id, { name: editName.trim(), color: editColor });
                          }
                          setEditingId(null);
                        }}
                      >
                        <Check className="h-4 w-4" />
                      </Button>
                      <Button size="sm" variant="ghost" onClick={() => setEditingId(null)}>
                        <X className="h-4 w-4" />
                      </Button>
                    </>
                  ) : (
                    <>
                      <span
                        className="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium"
                        style={{ background: `${l.color}22`, color: l.color }}
                      >
                        <span
                          className="h-1.5 w-1.5 rounded-full"
                          style={{ background: l.color }}
                        />
                        {l.name}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {useCount} {useCount === 1 ? "task" : "tasks"}
                      </span>
                      <div className="ml-auto flex items-center gap-1">
                        <Button
                          size="icon"
                          variant="ghost"
                          className="h-8 w-8"
                          onClick={() => {
                            setEditingId(l.id);
                            setEditName(l.name);
                            setEditColor(l.color);
                          }}
                        >
                          <Pencil className="h-3.5 w-3.5" />
                        </Button>
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button
                              size="icon"
                              variant="ghost"
                              className="h-8 w-8 text-destructive hover:text-destructive"
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Delete label?</AlertDialogTitle>
                              <AlertDialogDescription>
                                Delete label “{l.name}”? It will be removed from {useCount} task
                                {useCount === 1 ? "" : "s"}.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                              <AlertDialogAction asChild>
                                <Button
                                  variant="destructive"
                                  onClick={() => {
                                    void deleteLabel(l.id);
                                  }}
                                >
                                  Delete
                                </Button>
                              </AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </div>
                    </>
                  )}
                </motion.li>
              );
            })}
          </ul>
        )}
      </div>
    </div>
  );
}
