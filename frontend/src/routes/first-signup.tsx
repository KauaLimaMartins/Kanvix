import { createFileRoute, Navigate, useRouter } from "@tanstack/react-router";
import { useEffect, useState } from "react";

import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { motion } from "framer-motion";
import { ClientOnly } from "@/components/ClientOnly";
import { api } from "@/services/api";

export const Route = createFileRoute("/first-signup")({
  head: () => ({
    meta: [{ title: "Create first user — Kanvix" }],
  }),
  component: FirstSignupPage,
});

function FirstSignupPage() {
  return (
    <ClientOnly
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
          Loading…
        </div>
      }
    >
      <FirstSignupInner />
    </ClientOnly>
  );
}

function FirstSignupInner() {
  const authStatus = useAppStore((s) => s.authStatus);
  const firstSignup = useAppStore((s) => s.firstSignup);
  const isLoading = useAppStore((s) => s.isLoading);
  const error = useAppStore((s) => s.error);
  const router = useRouter();

  const [checked, setChecked] = useState(false);
  const [allowed, setAllowed] = useState(false);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  useEffect(() => {
    void (async () => {
      try {
        const res = await api.auth.setup();
        setAllowed(res.needsFirstSignup);
      } finally {
        setChecked(true);
      }
    })();
  }, []);

  if (authStatus === "authed") return <Navigate to="/workspaces" />;
  if (checked && !allowed) return <Navigate to="/login" />;

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !email.trim() || !password) return;
    await firstSignup(name.trim(), email.trim(), password);
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
            <div className="text-xs text-muted-foreground">First-time setup</div>
          </div>
        </div>

        <div className="rounded-2xl border border-border bg-card p-6 shadow-xl shadow-black/5">
          <h1 className="text-lg font-semibold">Create first user</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            This user will become an administrator automatically.
          </p>

          <form onSubmit={submit} className="mt-5 space-y-3">
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">Name</label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Your name"
                disabled={isLoading || !checked || !allowed}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">Email</label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@company.com"
                disabled={isLoading || !checked || !allowed}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">Password</label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isLoading || !checked || !allowed}
              />
            </div>
            <Button
              type="submit"
              className="w-full"
              disabled={isLoading || !checked || !allowed}
            >
              {isLoading ? "Creating…" : "Create account"}
            </Button>
          </form>

          {error && (
            <div className="mt-3 rounded-md border border-border bg-muted/30 px-3 py-2 text-xs text-muted-foreground">
              {error}
            </div>
          )}
        </div>
      </motion.div>
    </div>
  );
}

