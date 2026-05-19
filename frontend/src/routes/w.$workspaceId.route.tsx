import { createFileRoute, Outlet, Navigate } from "@tanstack/react-router";
import { AppShell } from "@/components/layout/AppShell";
import { useAppStore } from "@/store/useAppStore";
import { ClientOnly } from "@/components/ClientOnly";

export const Route = createFileRoute("/w/$workspaceId")({
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  return (
    <ClientOnly
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
          Loading…
        </div>
      }
    >
      <Inner />
    </ClientOnly>
  );
}

function Inner() {
  const authStatus = useAppStore((s) => s.authStatus);
  if (authStatus !== "authed") return <Navigate to="/login" />;
  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
}
