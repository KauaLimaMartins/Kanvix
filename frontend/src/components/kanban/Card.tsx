import { CSS } from "@dnd-kit/utilities";
import { Calendar } from "lucide-react";
import type { Task } from "@/types";
import { useAppStore } from "@/store/useAppStore";
import { useSortable } from "@dnd-kit/sortable";

export function Card({ task, onOpen }: { task: Task; onOpen: () => void }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: task.id,
    data: { type: "task", columnId: task.columnId },
  });
  const allLabels = useAppStore((s) => s.labels);
  const users = useAppStore((s) => s.users);
  const labels = task.labels
    .map((id) => allLabels.find((l) => l.id === id))
    .filter((l): l is NonNullable<typeof l> => Boolean(l));

  const style = {
    transform: CSS.Translate.toString(transform),
    transition: transition ?? "transform 180ms cubic-bezier(0.2, 0, 0, 1)",
    opacity: isDragging ? 0.5 : 1,
    boxShadow: isDragging
      ? "0 12px 28px -10px rgba(0,0,0,0.35), 0 4px 10px -4px rgba(0,0,0,0.2)"
      : undefined,
  };
  const assignee = task.assigneeId ? users.find((u) => u.id === task.assigneeId) : null;

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      onClick={onOpen}
      className="group cursor-grab rounded-lg border border-border bg-card p-3 text-card-foreground shadow-sm transition-all duration-150 hover:-translate-y-0.5 hover:border-border/80 hover:shadow-md active:cursor-grabbing"
    >
      <div className="text-sm font-medium leading-snug">{task.title}</div>
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
          {task.dueDate && (
            <>
              <Calendar className="h-3 w-3" />
              {new Date(task.dueDate).toLocaleDateString(undefined, {
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
