import { Check, Plus, Tag } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";

import { Button } from "@/components/ui/button";
import { useAppStore } from "@/store/useAppStore";
import { useState } from "react";

export function LabelPicker({
  workspaceId,
  selectedIds,
  onChange,
}: {
  workspaceId: string;
  selectedIds: string[];
  onChange: (ids: string[]) => void;
}) {
  const allLabels = useAppStore((s) => s.labels);
  const labels = allLabels.filter((l) => l.workspaceId === workspaceId);
  const [open, setOpen] = useState(false);

  const toggle = (id: string) => {
    onChange(selectedIds.includes(id) ? selectedIds.filter((x) => x !== id) : [...selectedIds, id]);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm" className="h-7 gap-1.5 text-xs">
          <Plus className="h-3 w-3" /> Add label
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-60 p-1" align="start">
        {labels.length === 0 ? (
          <div className="p-3 text-center text-xs text-muted-foreground">
            <Tag className="mx-auto mb-2 h-4 w-4 opacity-50" />
            No labels yet. Create one in the Labels page.
          </div>
        ) : (
          <div className="max-h-64 overflow-y-auto">
            {labels.map((l) => {
              const checked = selectedIds.includes(l.id);
              return (
                <button
                  key={l.id}
                  type="button"
                  onClick={() => toggle(l.id)}
                  className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-sm hover:bg-accent"
                >
                  <span className="h-3 w-3 shrink-0 rounded-sm" style={{ background: l.color }} />
                  <span className="flex-1 truncate">{l.name}</span>
                  {checked && <Check className="h-3.5 w-3.5 opacity-70" />}
                </button>
              );
            })}
          </div>
        )}
      </PopoverContent>
    </Popover>
  );
}
