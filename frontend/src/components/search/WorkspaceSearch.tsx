import { useEffect, useMemo, useState } from "react";
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/services/api";
import { useRouter } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";

export function WorkspaceSearch() {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");

  const match = router.state.matches.find((m) => m.params && "workspaceId" in m.params);
  const workspaceId =
    (match?.params as { workspaceId?: string } | undefined)?.workspaceId ?? null;

  const trimmed = query.trim();

  const searchQuery = useQuery({
    queryKey: ["workspaceSearch", workspaceId, trimmed],
    queryFn: () => api.search.query(workspaceId as string, trimmed, 20),
    enabled: open && !!workspaceId && trimmed.length > 0,
    staleTime: 15_000,
    retry: 1,
  });

  const grouped = useMemo(() => {
    const projects: Array<{ id: string; name: string; workspaceId: string }> = [];
    const tasks: Array<{ id: string; title: string; projectId: string; workspaceId: string }> = [];
    for (const r of searchQuery.data?.results ?? []) {
      if (r.type === "project") projects.push(r);
      if (r.type === "task") tasks.push(r);
    }
    return { projects, tasks };
  }, [searchQuery.data?.results]);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setOpen(true);
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  useEffect(() => {
    if (!open) setQuery("");
  }, [open]);

  return (
    <>
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setOpen(true)}
        disabled={!workspaceId}
      >
        <Search className="mr-2 h-4 w-4" />
        Search
      </Button>

      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput
          placeholder={workspaceId ? "Search projects and tasks…" : "Select a workspace to search"}
          value={query}
          onValueChange={setQuery}
        />
        <CommandList>
          {trimmed.length === 0 ? (
            <CommandEmpty>Type to search.</CommandEmpty>
          ) : searchQuery.isLoading ? (
            <CommandEmpty>Searching…</CommandEmpty>
          ) : searchQuery.isError ? (
            <CommandEmpty>Search failed.</CommandEmpty>
          ) : (
            <CommandEmpty>No results.</CommandEmpty>
          )}

          {grouped.projects.length > 0 && (
            <CommandGroup heading="Projects">
              {grouped.projects.map((p) => (
                <CommandItem
                  key={p.id}
                  value={`project:${p.name}`}
                  onSelect={() => {
                    setOpen(false);
                    router.navigate({
                      to: "/w/$workspaceId/p/$projectId",
                      params: { workspaceId: p.workspaceId, projectId: p.id },
                    });
                  }}
                >
                  {p.name}
                </CommandItem>
              ))}
            </CommandGroup>
          )}

          {grouped.tasks.length > 0 && (
            <CommandGroup heading="Tasks">
              {grouped.tasks.map((t) => (
                <CommandItem
                  key={t.id}
                  value={`task:${t.title}`}
                  onSelect={() => {
                    setOpen(false);
                    router.navigate({
                      to: "/w/$workspaceId/p/$projectId",
                      params: { workspaceId: t.workspaceId, projectId: t.projectId },
                      search: { task: t.id },
                    });
                  }}
                >
                  {t.title}
                </CommandItem>
              ))}
            </CommandGroup>
          )}
        </CommandList>
      </CommandDialog>
    </>
  );
}

