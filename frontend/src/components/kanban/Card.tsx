import { CSS } from "@dnd-kit/utilities";
import { Calendar, Pencil } from "lucide-react";
import type { Task } from "@/types";
import { useAppStore } from "@/store/useAppStore";
import { useSortable } from "@dnd-kit/sortable";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useEffect, useState } from "react";

export function Card({ task, onOpen }: { task: Task; onOpen: () => void }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: task.id,
    data: { type: "task", columnId: task.columnId },
  });
  const allLabels = useAppStore((s) => s.labels);
  const users = useAppStore((s) => s.users);
  const lastCreatedTaskId = useAppStore((s) => s.lastCreatedTaskId);
  const updateTask = useAppStore((s) => s.updateTask);
  const labels = (Array.isArray(task.labels) ? task.labels : [])
    .map((id) => allLabels.find((l) => l.id === id))
    .filter((l): l is NonNullable<typeof l> => Boolean(l));

  const [renaming, setRenaming] = useState(false);
  const [draftTitle, setDraftTitle] = useState(task.title);

  useEffect(() => {
    if (!renaming) setDraftTitle(task.title);
  }, [renaming, task.title]);

  const style = {
    transform: CSS.Translate.toString(transform),
    transition: transition ?? "transform 180ms cubic-bezier(0.2, 0, 0, 1)",
    opacity: isDragging ? 0.5 : 1,
    boxShadow: isDragging
      ? "0 12px 28px -10px rgba(0,0,0,0.35), 0 4px 10px -4px rgba(0,0,0,0.2)"
      : undefined,
  };
  const assignee = task.assigneeId ? users.find((u) => u.id === task.assigneeId) : null;
  const due = task.dueDate ? new Date(task.dueDate) : null;
  const hasValidDue = due != null && !Number.isNaN(due.getTime());

  const commitRename = () => {
    const next = draftTitle.trim();
    if (!next) {
      setDraftTitle(task.title);
      setRenaming(false);
      return;
    }
    if (next !== task.title) {
      updateTask(task.id, { title: next });
    }
    setRenaming(false);
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...(renaming ? {} : listeners)}
      onClick={() => {
        if (!renaming) onOpen();
      }}
      className={`group relative rounded-lg border border-border bg-card p-3 text-card-foreground shadow-sm transition-all duration-150 hover:-translate-y-0.5 hover:border-border/80 hover:shadow-md ${
        task.id === lastCreatedTaskId ? "kanvix-ring" : ""
      }`}
    >
      <Button
        size="icon"
        variant="ghost"
        className="absolute right-1 top-1 h-7 w-7 opacity-0 transition-opacity group-hover:opacity-100"
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setRenaming(true);
        }}
      >
        <Pencil className="h-3.5 w-3.5" />
      </Button>

      {renaming ? (
        <Input
          autoFocus
          value={draftTitle}
          onChange={(e) => setDraftTitle(e.target.value)}
          onBlur={commitRename}
          onKeyDown={(e) => {
            if (e.key === "Enter") commitRename();
            if (e.key === "Escape") {
              setDraftTitle(task.title);
              setRenaming(false);
            }
          }}
          onClick={(e) => e.stopPropagation()}
          className="h-8 text-sm font-medium leading-snug"
        />
      ) : (
        <div className="pr-7 text-sm font-medium leading-snug">{task.title}</div>
      )}
      {labels.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-1">
          {labels.map((l) => (
            <span
              key={l.id}
              className="inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 text-[10px] font-medium"
              style={{ background: `${l.color}22`, color: l.color }}
            >
              <span className="h-1 w-1 rounded-full" style={{ background: l.color }} />
              {l.name}
            </span>
          ))}
        </div>
      )}
      <div className="mt-3 flex items-center justify-between text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          {hasValidDue && (
            <>
              <Calendar className="h-3 w-3" />
              {due!.toLocaleDateString(undefined, {
                month: "short",
                day: "numeric",
              })}
            </>
          )}
        </span>
        {assignee && (
          <span
            className="flex h-6 w-6 items-center justify-center rounded-full text-[10px] font-medium text-white"
            style={{ background: assignee.avatarColor }}
            title={assignee.name}
          >
            {assignee.name
              .split(" ")
              .map((p) => p[0])
              .join("")
              .slice(0, 2)}
          </span>
        )}
      </div>
    </div>
  );
}
