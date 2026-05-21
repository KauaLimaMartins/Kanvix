import { createFileRoute } from "@tanstack/react-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import { api } from "@/services/api";
import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export const Route = createFileRoute("/w/$workspaceId/users")({
  component: UsersPage,
});

function UsersPage() {
  const { workspaceId } = Route.useParams();
  const workspace = useAppStore((s) => s.workspaces.find((w) => w.id === workspaceId) ?? null);
  const userRole = useAppStore((s) => s.userRole);
  const canManageUsers = (workspace?.role ?? userRole) === "admin";

  const qc = useQueryClient();
  const usersQuery = useQuery({
    queryKey: ["workspaceUsers", workspaceId],
    queryFn: () => api.users.listByWorkspace(workspaceId),
    enabled: canManageUsers,
    staleTime: 10_000,
    retry: 1,
  });

  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<"admin" | "member">("member");

  if (!workspace) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
        Workspace not found.
      </div>
    );
  }

  if (!canManageUsers) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
        Forbidden.
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Users</h1>
          <p className="mt-1 text-sm text-muted-foreground">Manage users and roles for {workspace.name}.</p>
        </div>

        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button>New user</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>New user</DialogTitle>
            </DialogHeader>
            <div className="space-y-3">
              <div className="space-y-1.5">
                <Label className="text-xs">Name</Label>
                <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Full name" />
              </div>
              <div className="space-y-1.5">
                <Label className="text-xs">Email</Label>
                <Input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="user@company.com"
                />
              </div>
              <div className="space-y-1.5">
                <Label className="text-xs">Password</Label>
                <Input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="At least 8 characters"
                />
              </div>
              <div className="space-y-1.5">
                <Label className="text-xs">Role</Label>
                <Select value={role} onValueChange={(v) => setRole(v as "admin" | "member")}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="member">Member</SelectItem>
                    <SelectItem value="admin">Admin</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="ghost"
                onClick={() => {
                  setOpen(false);
                }}
              >
                Cancel
              </Button>
              <Button
                onClick={() => {
                  if (!email.trim() || !name.trim() || !password) return;
                  void (async () => {
                    await api.users.createInWorkspace(workspaceId, {
                      email: email.trim(),
                      name: name.trim(),
                      password,
                      role,
                    });
                    setName("");
                    setEmail("");
                    setPassword("");
                    setRole("member");
                    setOpen(false);
                    await qc.invalidateQueries({ queryKey: ["workspaceUsers", workspaceId] });
                    await qc.invalidateQueries({ queryKey: ["bootstrap"] });
                  })();
                }}
              >
                Create
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="rounded-xl border border-border bg-card">
        <div className="border-b border-border px-4 py-3 text-sm font-medium">Workspace users</div>
        {usersQuery.isLoading ? (
          <div className="px-4 py-6 text-sm text-muted-foreground">Loading…</div>
        ) : usersQuery.isError ? (
          <div className="px-4 py-6 text-sm text-muted-foreground">Could not load users.</div>
        ) : (usersQuery.data?.users?.length ?? 0) === 0 ? (
          <div className="px-4 py-6 text-sm text-muted-foreground">No users.</div>
        ) : (
          <div className="divide-y divide-border">
            {usersQuery.data!.users.map((u) => (
              <div key={u.id} className="flex items-center justify-between px-4 py-3">
                <div className="min-w-0">
                  <div className="truncate text-sm font-medium">{u.name}</div>
                  <div className="truncate text-xs text-muted-foreground">{u.email}</div>
                </div>
                <Select
                  value={u.role}
                  onValueChange={(v) => {
                    void (async () => {
                      await api.users.updateRoleInWorkspace(workspaceId, u.id, v);
                      await qc.invalidateQueries({ queryKey: ["workspaceUsers", workspaceId] });
                      await qc.invalidateQueries({ queryKey: ["bootstrap"] });
                    })();
                  }}
                >
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="member">Member</SelectItem>
                    <SelectItem value="admin">Admin</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

