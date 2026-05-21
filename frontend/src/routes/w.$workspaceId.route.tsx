import { createFileRoute, Outlet, Navigate, useRouterState } from "@tanstack/react-router";
import { AppShell } from "@/components/layout/AppShell";
import { useAppStore } from "@/store/useAppStore";
import { ClientOnly } from "@/components/ClientOnly";
import { useLayoutEffect, useRef } from "react";

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
  const locationKey = useRouterState({
    select: (s) => s.location.href,
  });
  const pageRef = useRef<HTMLDivElement | null>(null);
  const lastKeyRef = useRef<string | null>(null);

  useLayoutEffect(() => {
    const el = pageRef.current;
    if (!el) return;
    if (lastKeyRef.current === locationKey) {
      return;
    }
    lastKeyRef.current = locationKey;
    el.classList.remove("kanvix-page-enter");
    void el.offsetWidth;
    el.classList.add("kanvix-page-enter");
  }, [locationKey]);
  if (authStatus === "unknown") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
        Loading…
      </div>
    );
  }
  if (authStatus !== "authed") return <Navigate to="/login" />;
  return (
    <AppShell>
      <div ref={pageRef} className="h-full">
        <Outlet />
      </div>
    </AppShell>
  );
}
