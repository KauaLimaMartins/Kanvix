import { useAppStore } from "@/store/useAppStore";
import { Button } from "@/components/ui/button";
import { Moon, Sun, LogOut } from "lucide-react";
import { useRouter } from "@tanstack/react-router";
import { WorkspaceSearch } from "@/components/search/WorkspaceSearch";

export function Topbar() {
  const theme = useAppStore((s) => s.theme);
  const setTheme = useAppStore((s) => s.setTheme);
  const logout = useAppStore((s) => s.logout);
  const userEmail = useAppStore((s) => s.userEmail);
  const router = useRouter();

  return (
    <header className="flex h-12 items-center justify-end gap-2 border-b border-border bg-background/80 px-4 backdrop-blur">
      {userEmail && (
        <span className="mr-auto text-xs text-muted-foreground">Signed in as {userEmail}</span>
      )}
      <WorkspaceSearch />
      <Button
        variant="ghost"
        size="icon"
        onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        aria-label="Toggle theme"
      >
        {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
      </Button>
      <Button
        variant="ghost"
        size="icon"
        aria-label="Sign out"
        onClick={async () => {
          await logout();
          router.navigate({ to: "/login" });
        }}
      >
        <LogOut className="h-4 w-4" />
      </Button>
    </header>
  );
}
