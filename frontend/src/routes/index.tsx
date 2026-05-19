import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useAppStore } from "@/store/useAppStore";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const authStatus = useAppStore((s) => s.authStatus);
  if (authStatus === "unknown") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
        Loading…
      </div>
    );
  }
  if (authStatus !== "authed") return <Navigate to="/login" />;
  return <Navigate to="/workspaces" />;
}
