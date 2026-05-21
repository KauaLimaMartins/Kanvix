import { useMemo, useState } from "react";
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragOverEvent,
  type DragStartEvent,
} from "@dnd-kit/core";
import { Column } from "./Column";
import { Card } from "./Card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";
import { useAppStore } from "@/store/useAppStore";
import { motion } from "framer-motion";

export function Board({
  projectId,
  onOpenTask,
}: {
  projectId: string;
  onOpenTask: (id: string) => void;
}) {
  const allColumns = useAppStore((s) => s.columns);
  const allTasks = useAppStore((s) => s.tasks);
  const columns = allColumns
    .filter((c) => c.projectId === projectId)
    .sort((a, b) => a.order - b.order);
  const tasks = allTasks.filter((t) => t.projectId === projectId);
  const moveTask = useAppStore((s) => s.moveTask);
  const addColumn = useAppStore((s) => s.addColumn);

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));
  const [activeId, setActiveId] = useState<string | null>(null);
  const [isAddingColumn, setIsAddingColumn] = useState(false);
  const [newColumnName, setNewColumnName] = useState("");

  const tasksByColumn = useMemo(() => {
    const map: Record<string, typeof tasks> = {};
    for (const c of columns) map[c.id] = [];
    for (const t of tasks) (map[t.columnId] ??= []).push(t);
    for (const k of Object.keys(map)) map[k].sort((a, b) => a.order - b.order);
    return map;
  }, [columns, tasks]);

  const activeTask = activeId ? (tasks.find((t) => t.id === activeId) ?? null) : null;

  const onDragStart = (e: DragStartEvent) => setActiveId(String(e.active.id));

  const findContainerOf = (id: string): string | null => {
    if (id.startsWith("col-")) return id.slice(4);
    const t = tasks.find((t) => t.id === id);
    return t ? t.columnId : null;
  };

  const onDragOver = (_e: DragOverEvent) => {
    // no-op; we commit on drag end for simplicity
  };

  const onDragEnd = (e: DragEndEvent) => {
    setActiveId(null);
    const activeIdStr = String(e.active.id);
    const overId = e.over?.id ? String(e.over.id) : null;
    if (!overId) return;

    const toCol = findContainerOf(overId);
    if (!toCol) return;

    const targetList = tasksByColumn[toCol] ?? [];
    let toIndex: number;
    if (overId.startsWith("col-")) {
      toIndex = targetList.length;
    } else {
      toIndex = targetList.findIndex((t) => t.id === overId);
      if (toIndex === -1) toIndex = targetList.length;
    }
    void moveTask(activeIdStr, toCol, toIndex);
  };

  const submitNewColumn = () => {
    const name = newColumnName.trim();
    if (!name) return;
    void addColumn(projectId, name);
    setIsAddingColumn(false);
    setNewColumnName("");
  };

  return (
    <DndContext
      sensors={sensors}
      onDragStart={onDragStart}
      onDragOver={onDragOver}
      onDragEnd={onDragEnd}
      onDragCancel={() => setActiveId(null)}
    >
      <motion.div
        className="flex h-full gap-3 overflow-x-auto overflow-y-hidden p-4"
        initial={{ opacity: 0, scale: 0.99 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.22, ease: [0.2, 0, 0, 1] }}
      >
        {columns.map((c) => (
          <Column key={c.id} column={c} tasks={tasksByColumn[c.id] ?? []} onOpenTask={onOpenTask} />
        ))}
        <div className="shrink-0">
          {isAddingColumn ? (
            <div className="w-72 rounded-md border border-border bg-muted/20 p-2">
              <Input
                value={newColumnName}
                onChange={(e) => setNewColumnName(e.target.value)}
                placeholder="Column name"
                autoFocus
                onKeyDown={(e) => {
                  if (e.key === "Enter") submitNewColumn();
                  if (e.key === "Escape") {
                    setIsAddingColumn(false);
                    setNewColumnName("");
                  }
                }}
              />
              <div className="mt-2 flex gap-2">
                <Button size="sm" onClick={submitNewColumn}>
                  Create
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => {
                    setIsAddingColumn(false);
                    setNewColumnName("");
                  }}
                >
                  Cancel
                </Button>
              </div>
            </div>
          ) : (
            <Button
              variant="ghost"
              className="h-9 w-72 justify-start gap-1.5 border border-dashed border-border bg-muted/20 text-muted-foreground transition-all hover:-translate-y-0.5 hover:bg-muted/40 active:translate-y-0"
              onClick={() => {
                setIsAddingColumn(true);
                setNewColumnName("New column");
              }}
            >
              <Plus className="h-4 w-4" /> Add column
            </Button>
          )}
        </div>
      </motion.div>
      <DragOverlay>{activeTask && <Card task={activeTask} onOpen={() => {}} />}</DragOverlay>
    </DndContext>
  );
}
