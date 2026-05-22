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
import { ApiError, api } from "@/services/api";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Pencil, Trash2 } from "lucide-react";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useQuery, useQueryClient } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { createFileRoute } from "@tanstack/react-router";
import { toast } from "sonner";
import { useAppStore } from "@/store/useAppStore";
import { useState } from "react";

export const Route = createFileRoute("/w/$workspaceId/users")({
  component: UsersPage,
});

function errorMessage(e: unknown, fallback: string) {
  if (e instanceof ApiError) return e.message;
  if (e instanceof Error) return e.message;
  return fallback;
}

function UsersPage() {
  const { workspaceId } = Route.useParams();
  const workspace = useAppStore((s) => s.workspaces.find((w) => w.id === workspaceId) ?? null);
  const userRole = useAppStore((s) => s.userRole);
  const canManageUsers = (workspace?.role ?? userRole) === "admin";
  const hydrate = useAppStore((s) => s.hydrate);

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
  const [isCreating, setIsCreating] = useState(false);
  const [lastCreatedUserId, setLastCreatedUserId] = useState<string | null>(null);
  const [editingUserId, setEditingUserId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState("");
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [deleteUserId, setDeleteUserId] = useState<string | null>(null);
  const [deleteAction, setDeleteAction] = useState<"unassign" | "reassign" | "disable">("unassign");
  const [reassignToUserId, setReassignToUserId] = useState<string>("none");

  const users = usersQuery.data?.users ?? [];
  const deleteUser = users.find((u) => u.id === deleteUserId) ?? null;
  const reassignCandidates =
    deleteUser == null ? [] : users.filter((u) => u.id !== deleteUser.id && !u.disabled);

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
          <p className="mt-1 text-sm text-muted-foreground">
            Manage users and roles for {workspace.name}.
          </p>
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
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Full name"
                />
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
                disabled={isCreating}
                onClick={() => {
                  setOpen(false);
                }}
              >
                Cancel
              </Button>
              <Button
                disabled={isCreating}
                onClick={() => {
                  const nextEmail = email.trim();
                  const nextName = name.trim();
                  if (!nextEmail || !nextName || !password) {
                    toast.error("Name, email, and password are required.");
                    return;
                  }
                  if (password.length < 8) {
                    toast.error("Password must be at least 8 characters.");
                    return;
                  }
                  void (async () => {
                    setIsCreating(true);
                    try {
                      const res = await api.users.createInWorkspace(workspaceId, {
                        email: nextEmail,
                        name: nextName,
                        password,
                        role,
                      });
                      setLastCreatedUserId(res.user.id);
                      setTimeout(() => setLastCreatedUserId(null), 900);
                      setName("");
                      setEmail("");
                      setPassword("");
                      setRole("member");
                      setOpen(false);
                      await qc.invalidateQueries({ queryKey: ["workspaceUsers", workspaceId] });
                      await qc.invalidateQueries({ queryKey: ["bootstrap"] });
                    } catch (e) {
                      toast.error(errorMessage(e, "Could not create user."));
                    } finally {
                      setIsCreating(false);
                    }
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
            {users.map((u) => {
              const isEditing = editingUserId === u.id;
              return (
                <div
                  key={u.id}
                  className={`flex items-center justify-between gap-3 px-4 py-3 ${
                    u.id === lastCreatedUserId ? "kanvix-ring" : ""
                  }`}
                >
                  <div className="min-w-0 flex-1">
                    {isEditing ? (
                      <Input
                        autoFocus
                        value={editingName}
                        onChange={(e) => setEditingName(e.target.value)}
                        onBlur={() => {
                          const next = editingName.trim();
                          setEditingUserId(null);
                          setEditingName("");
                          if (!next || next === u.name) return;
                          void (async () => {
                            await api.users.updateInWorkspace(workspaceId, u.id, { name: next });
                            await qc.invalidateQueries({
                              queryKey: ["workspaceUsers", workspaceId],
                            });
                            await hydrate();
                          })();
                        }}
                        onKeyDown={(e) => {
                          if (e.key === "Escape") {
                            setEditingUserId(null);
                            setEditingName("");
                          }
                          if (e.key === "Enter") {
                            const next = editingName.trim();
                            setEditingUserId(null);
                            setEditingName("");
                            if (!next || next === u.name) return;
                            void (async () => {
                              await api.users.updateInWorkspace(workspaceId, u.id, { name: next });
                              await qc.invalidateQueries({
                                queryKey: ["workspaceUsers", workspaceId],
                              });
                              await hydrate();
                            })();
                          }
                        }}
                        className="h-8 max-w-sm"
                      />
                    ) : (
                      <div className="flex items-center gap-2">
                        <div className="truncate text-sm font-medium">{u.name}</div>
                        {u.disabled && (
                          <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium text-muted-foreground">
                            Disabled
                          </span>
                        )}
                      </div>
                    )}
                    <div className="truncate text-xs text-muted-foreground">{u.email}</div>
                  </div>

                  <div className="flex items-center gap-2">
                    {!isEditing && (
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-8 w-8"
                        onClick={() => {
                          setEditingUserId(u.id);
                          setEditingName(u.name);
                        }}
                      >
                        <Pencil className="h-3.5 w-3.5" />
                      </Button>
                    )}

                    <Select
                      value={u.role}
                      onValueChange={(v) => {
                        void (async () => {
                          await api.users.updateRoleInWorkspace(workspaceId, u.id, v);
                          await qc.invalidateQueries({ queryKey: ["workspaceUsers", workspaceId] });
                          await hydrate();
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

                    <Button
                      size="icon"
                      variant="ghost"
                      className="h-8 w-8 text-destructive hover:text-destructive"
                      onClick={() => {
                        setDeleteUserId(u.id);
                        setDeleteAction("unassign");
                        setReassignToUserId("none");
                        setDeleteOpen(true);
                      }}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <AlertDialog
        open={deleteOpen}
        onOpenChange={(o) => {
          setDeleteOpen(o);
          if (!o) {
            setDeleteUserId(null);
            setDeleteAction("unassign");
            setReassignToUserId("none");
          }
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {deleteAction === "disable" ? "Disable user?" : "Delete user from workspace?"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              Choose what to do with tasks currently assigned to {deleteUser?.name ?? "this user"}.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="space-y-3">
            <RadioGroup
              value={deleteAction}
              onValueChange={(v) => setDeleteAction(v as "unassign" | "reassign" | "disable")}
              className="gap-3"
            >
              <label className="flex cursor-pointer items-start gap-3 rounded-lg border border-border p-3">
                <RadioGroupItem value="unassign" className="mt-0.5" />
                <div className="min-w-0">
                  <div className="text-sm font-medium">Remove attribution</div>
                  <div className="text-xs text-muted-foreground">
                    Unassign this user from all tasks in this workspace.
                  </div>
                </div>
              </label>

              <label className="flex cursor-pointer items-start gap-3 rounded-lg border border-border p-3">
                <RadioGroupItem value="reassign" className="mt-0.5" />
                <div className="min-w-0 flex-1">
                  <div className="text-sm font-medium">Move attribution</div>
                  <div className="text-xs text-muted-foreground">
                    Reassign all tasks in this workspace to another user.
                  </div>
                  <div className="mt-2">
                    <Select
                      value={reassignToUserId}
                      onValueChange={setReassignToUserId}
                      disabled={deleteAction !== "reassign"}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select user" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="none">Select user</SelectItem>
                        {reassignCandidates.map((u) => (
                          <SelectItem key={u.id} value={u.id}>
                            {u.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </label>

              <label className="flex cursor-pointer items-start gap-3 rounded-lg border border-border p-3">
                <RadioGroupItem value="disable" className="mt-0.5" />
                <div className="min-w-0">
                  <div className="text-sm font-medium">Disable user</div>
                  <div className="text-xs text-muted-foreground">
                    Keep them assigned to all tasks, but prevent sign-in.
                  </div>
                </div>
              </label>
            </RadioGroup>
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button
                variant="destructive"
                disabled={
                  !deleteUser ||
                  (deleteAction === "reassign" &&
                    (reassignToUserId === "none" || reassignCandidates.length === 0))
                }
                onClick={() => {
                  if (!deleteUser) return;
                  const payload =
                    deleteAction === "reassign"
                      ? { action: "reassign" as const, reassignToUserId }
                      : { action: deleteAction };
                  void (async () => {
                    await api.users.deleteFromWorkspace(workspaceId, deleteUser.id, payload);
                    setDeleteOpen(false);
                    await qc.invalidateQueries({ queryKey: ["workspaceUsers", workspaceId] });
                    await hydrate();
                  })();
                }}
              >
                {deleteAction === "disable" ? "Disable" : "Delete"}
              </Button>
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
