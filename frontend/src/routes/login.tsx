import { createFileRoute, useRouter, Navigate } from "@tanstack/react-router";
import { useState } from "react";
import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { motion } from "framer-motion";
import { ClientOnly } from "@/components/ClientOnly";

export const Route = createFileRoute("/login")({
  head: () => ({
    meta: [
      { title: "Sign in — Kanvix" },
      { name: "description", content: "Sign in to your Kanvix workspace." },
    ],
  }),
  component: LoginPage,
});

function LoginPage() {
  return (
    <ClientOnly
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
          Loading…
        </div>
      }
    >
      <LoginInner />
    </ClientOnly>
  );
}

function LoginInner() {
  const authStatus = useAppStore((s) => s.authStatus);
  const login = useAppStore((s) => s.login);
  const isLoading = useAppStore((s) => s.isLoading);
  const error = useAppStore((s) => s.error);
  const router = useRouter();
  const [email, setEmail] = useState("you@kanvix.app");
  const [password, setPassword] = useState("••••••••");

  if (authStatus === "authed") return <Navigate to="/workspaces" />;

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) return;
    await login(email.trim(), password);
    router.navigate({ to: "/workspaces" });
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background px-4">
      <div className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute left-1/2 top-0 h-[480px] w-[820px] -translate-x-1/2 rounded-full bg-indigo-500/15 blur-3xl" />
        <div className="absolute bottom-0 right-1/3 h-[360px] w-[520px] rounded-full bg-fuchsia-500/15 blur-3xl" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, ease: [0.2, 0, 0, 1] }}
        className="w-full max-w-sm"
      >
        <div className="mb-8 flex items-center gap-2.5">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-gradient-to-br from-indigo-500 to-fuchsia-500 text-white shadow-lg shadow-indigo-500/25">
            <span className="text-base font-bold">K</span>
          </div>
          <div>
            <div className="text-xl font-semibold tracking-tight">Kanvix</div>
            <div className="text-xs text-muted-foreground">Boards for builders.</div>
          </div>
        </div>

        <div className="rounded-2xl border border-border bg-card p-6 shadow-xl shadow-black/5">
          <h1 className="text-lg font-semibold">Welcome back</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Sign in to continue. Any email works in this demo.
          </p>

          <form onSubmit={submit} className="mt-5 space-y-3">
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">Email</label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@company.com"
                disabled={isLoading}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">Password</label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isLoading}
              />
            </div>
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Signing in…" : "Sign in"}
            </Button>
          </form>

          {error && (
            <div className="mt-3 rounded-md border border-border bg-muted/30 px-3 py-2 text-xs text-muted-foreground">
              {error}
            </div>
          )}

          <div className="mt-4 text-center text-[11px] text-muted-foreground">
            Demo only — no real authentication is performed.
          </div>
        </div>
      </motion.div>
    </div>
  );
}
