import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Trash2, X } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { LabelPicker } from "@/components/labels/LabelPicker";
import { TaskEditor } from "./TaskEditor";
import { useAppStore } from "@/store/useAppStore";

export function TaskDrawer({ taskId, onClose }: { taskId: string | null; onClose: () => void }) {
  const task = useAppStore((s) => s.tasks.find((t) => t.id === taskId) ?? null);
  const project = useAppStore((s) =>
    task ? (s.projects.find((p) => p.id === task.projectId) ?? null) : null,
  );
  const columns = useAppStore((s) => s.columns);
  const allLabels = useAppStore((s) => s.labels);
  const users = useAppStore((s) => s.users);
  const updateTask = useAppStore((s) => s.updateTask);
  const deleteTask = useAppStore((s) => s.deleteTask);
  const moveTask = useAppStore((s) => s.moveTask);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description);
    }
  }, [task?.id]);

  if (!task || !project) return null;

  const taskColumns = columns
    .filter((c) => c.projectId === task.projectId)
    .sort((a, b) => a.order - b.order);

  const taskLabelObjs = task.labels
    .map((id) => allLabels.find((l) => l.id === id))
    .filter((l): l is NonNullable<typeof l> => Boolean(l));

  const commit = () => {
    if (title !== task.title || description !== task.description) {
      updateTask(task.id, { title, description });
    }
  };

  return (
    <Sheet
      open={!!task}
      onOpenChange={(o) => {
        if (!o) {
          commit();
          onClose();
        }
      }}
    >
      <SheetContent className="w-full overflow-y-auto p-0 sm:max-w-xl">
        <SheetHeader className="border-b border-border px-6 py-4">
          <SheetTitle className="sr-only">Task details</SheetTitle>
          <Input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            onBlur={commit}
            placeholder="Task title"
            className="border-0 bg-transparent px-0 text-lg font-semibold shadow-none focus-visible:ring-0"
          />
        </SheetHeader>

        <div className="space-y-5 px-6 py-5">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs">Status</Label>
              <Select value={task.columnId} onValueChange={(v) => moveTask(task.id, v, 0)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {taskColumns.map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs">Assignee</Label>
              <Select
                value={task.assigneeId ?? "none"}
                onValueChange={(v) =>
                  updateTask(task.id, { assigneeId: v === "none" ? undefined : v })
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Unassigned</SelectItem>
                  {users.map((u) => (
                    <SelectItem key={u.id} value={u.id}>
                      {u.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs">Due date</Label>
              <Input
                type="date"
                value={task.dueDate ?? ""}
                onChange={(e) => updateTask(task.id, { dueDate: e.target.value || undefined })}
              />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs">Created</Label>
              <div className="rounded-md border border-input px-3 py-2 text-sm text-muted-foreground">
                {new Date(task.createdAt).toLocaleDateString()}
              </div>
            </div>
          </div>

          <div className="space-y-1.5">
            <Label className="text-xs">Labels</Label>
            <div className="flex flex-wrap items-center gap-1.5">
              {taskLabelObjs.map((l) => (
                <span
                  key={l.id}
                  className="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium"
                  style={{ background: `${l.color}22`, color: l.color }}
                >
                  <span className="h-1.5 w-1.5 rounded-full" style={{ background: l.color }} />
                  {l.name}
                  <button
                    type="button"
                    onClick={() =>
                      updateTask(task.id, {
                        labels: task.labels.filter((x) => x !== l.id),
                      })
                    }
                    className="ml-0.5 opacity-70 hover:opacity-100"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </span>
              ))}
              <LabelPicker
                workspaceId={project.workspaceId}
                selectedIds={task.labels}
                onChange={(ids) => updateTask(task.id, { labels: ids })}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label className="text-xs">Description</Label>
            <TaskEditor value={description} onChange={setDescription} />
          </div>

          <div className="flex justify-between pt-2">
            <Button
              variant="ghost"
              size="sm"
              className="text-destructive hover:text-destructive"
              onClick={() => {
                if (confirm("Delete this task?")) {
                  deleteTask(task.id);
                  onClose();
                }
              }}
            >
              <Trash2 className="mr-2 h-4 w-4" /> Delete
            </Button>
            <Button
              size="sm"
              onClick={() => {
                commit();
                onClose();
              }}
            >
              Done
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
