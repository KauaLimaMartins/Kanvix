import type { Column as ColType, Task } from "@/types";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { MoreHorizontal, Plus, Trash2 } from "lucide-react";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";

import { Button } from "@/components/ui/button";
import { Card } from "./Card";
import { Input } from "@/components/ui/input";
import { useAppStore } from "@/store/useAppStore";
import { useDroppable } from "@dnd-kit/core";
import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";

export function Column({
  column,
  tasks,
  onOpenTask,
}: {
  column: ColType;
  tasks: Task[];
  onOpenTask: (id: string) => void;
}) {
  const addTask = useAppStore((s) => s.addTask);
  const renameColumn = useAppStore((s) => s.renameColumn);
  const deleteColumn = useAppStore((s) => s.deleteColumn);
  const lastCreatedColumnId = useAppStore((s) => s.lastCreatedColumnId);
  const { setNodeRef, isOver } = useDroppable({
    id: `col-${column.id}`,
    data: { type: "column", columnId: column.id },
  });

  const [adding, setAdding] = useState(false);
  const [newTitle, setNewTitle] = useState("");
  const [editingName, setEditingName] = useState(false);
  const [name, setName] = useState(column.name);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const ids = tasks.map((t) => t.id);

  return (
    <div
      className={`flex w-72 shrink-0 flex-col rounded-lg bg-muted/40 ${
        column.id === lastCreatedColumnId ? "kanvix-ring" : ""
      }`}
    >
      <div className="flex items-center justify-between gap-2 px-3 py-2">
        {editingName ? (
          <Input
            autoFocus
            value={name}
            onChange={(e) => setName(e.target.value)}
            onBlur={() => {
              void renameColumn(column.id, name.trim() || column.name);
              setEditingName(false);
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                void renameColumn(column.id, name.trim() || column.name);
                setEditingName(false);
              }
              if (e.key === "Escape") {
                setName(column.name);
                setEditingName(false);
              }
            }}
            className="h-7 text-sm"
          />
        ) : (
          <button
            onClick={() => setEditingName(true)}
            className="flex items-center gap-2 text-sm font-medium"
          >
            {column.name}
            <span className="rounded-full bg-background px-1.5 py-0.5 text-[10px] text-muted-foreground">
              {tasks.length}
            </span>
          </button>
        )}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button size="icon" variant="ghost" className="h-6 w-6">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => setEditingName(true)}>Rename</DropdownMenuItem>
            <DropdownMenuItem
              className="text-destructive"
              onClick={() => {
                setDeleteOpen(true);
              }}
            >
              <Trash2 className="mr-2 h-4 w-4" /> Delete column
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <div
        ref={setNodeRef}
        className={`flex flex-1 flex-col gap-2 overflow-y-auto px-2 pb-2 transition-colors ${
          isOver ? "bg-accent/30" : ""
        }`}
      >
        <SortableContext items={ids} strategy={verticalListSortingStrategy}>
          {tasks.map((t) => (
            <Card key={t.id} task={t} onOpen={() => onOpenTask(t.id)} />
          ))}
        </SortableContext>

        <AnimatePresence initial={false} mode="wait">
          {adding ? (
            <motion.div
              key="add-task-form"
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.18, ease: [0.2, 0, 0, 1] }}
              className="overflow-hidden"
            >
              <div className="space-y-2 rounded-md border border-border bg-card p-2">
                <Input
                  autoFocus
                  placeholder="Task title"
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" && newTitle.trim()) {
                      void addTask(column.projectId, column.id, newTitle.trim());
                      setNewTitle("");
                    }
                    if (e.key === "Escape") {
                      setAdding(false);
                      setNewTitle("");
                    }
                  }}
                  className="h-8"
                />
                <div className="flex gap-1">
                  <Button
                    size="sm"
                    className="h-7"
                    onClick={() => {
                      if (newTitle.trim()) {
                        void addTask(column.projectId, column.id, newTitle.trim());
                        setNewTitle("");
                      }
                      setAdding(false);
                    }}
                  >
                    Add
                  </Button>
                  <Button
                    size="sm"
                    variant="ghost"
                    className="h-7"
                    onClick={() => {
                      setAdding(false);
                      setNewTitle("");
                    }}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </motion.div>
          ) : (
            <motion.button
              key="add-task-button"
              type="button"
              onClick={() => setAdding(true)}
              className="flex items-center gap-1.5 rounded-md px-2 py-1.5 text-left text-xs text-muted-foreground transition-all duration-150 hover:-translate-y-0.5 hover:bg-accent/60 hover:text-foreground active:translate-y-0"
              whileTap={{ scale: 0.98 }}
            >
              <Plus className="h-3.5 w-3.5" /> Add task
            </motion.button>
          )}
        </AnimatePresence>
      </div>

      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete column?</AlertDialogTitle>
            <AlertDialogDescription>
              Delete column “{column.name}” and its {tasks.length} task
              {tasks.length === 1 ? "" : "s"}?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button
                variant="destructive"
                onClick={() => {
                  void deleteColumn(column.id);
                  setDeleteOpen(false);
                }}
              >
                Delete
              </Button>
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
