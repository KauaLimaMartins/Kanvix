import { createFileRoute, useNavigate, useSearch } from "@tanstack/react-router";

import { Board } from "@/components/kanban/Board";
import { TaskDrawer } from "@/components/task/TaskDrawer";
import { useAppStore } from "@/store/useAppStore";
import { z } from "zod";

const searchSchema = z.object({
  task: z.string().optional(),
});

export const Route = createFileRoute("/w/$workspaceId/p/$projectId")({
  validateSearch: searchSchema,
  component: ProjectPage,
});

function ProjectPage() {
  const { workspaceId, projectId } = Route.useParams();
  const { task: taskId } = useSearch({ from: Route.id });
  const navigate = useNavigate({ from: Route.fullPath });

  const project = useAppStore((s) => s.projects.find((p) => p.id === projectId));

  if (!project) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
        Project not found.
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      <div className="border-b border-border px-6 py-3">
        <h1 className="text-lg font-semibold">{project.name}</h1>
        {project.description && (
          <p className="text-xs text-muted-foreground">{project.description}</p>
        )}
      </div>
      <div className="min-h-0 flex-1">
        <Board
          projectId={projectId}
          onOpenTask={(id) =>
            navigate({
              params: { workspaceId, projectId },
              search: { task: id },
            })
          }
        />
      </div>
      <TaskDrawer
        taskId={taskId ?? null}
        onClose={() =>
          navigate({
            params: { workspaceId, projectId },
            search: { task: undefined },
          })
        }
      />
    </div>
  );
}
