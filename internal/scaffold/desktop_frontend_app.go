package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDesktopFrontendAppFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		// Route files
		filepath.Join(root, "frontend", "src", "routes", "__root.tsx"):   desktopRootRoute(),
		filepath.Join(root, "frontend", "src", "routes", "login.tsx"):    desktopLoginRoute(),
		filepath.Join(root, "frontend", "src", "routes", "register.tsx"): desktopRegisterRoute(),
		// Utilities and hooks
		filepath.Join(root, "frontend", "src", "lib", "utils.ts"):        desktopUtilsTS(),
		filepath.Join(root, "frontend", "src", "lib", "query-client.ts"): desktopQueryClientTS(),
		filepath.Join(root, "frontend", "src", "hooks", "use-auth.tsx"):  desktopUseAuthHook(),
		filepath.Join(root, "frontend", "src", "hooks", "use-theme.ts"):  desktopUseThemeHook(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// ── Root Route (__root.tsx) ─────────────────────────────────────────────────

func desktopRootRoute() string {
	return `import { createRootRoute, Outlet } from "@tanstack/react-router";
import { Toaster } from "sonner";

export const Route = createRootRoute({
  component: RootLayout,
});

function RootLayout() {
  return (
    <>
      <Outlet />
      <Toaster
        theme="dark"
        position="bottom-right"
        toastOptions={{
          style: {
            background: "var(--bg-elevated)",
            border: "1px solid var(--border)",
            color: "var(--foreground)",
          },
        }}
      />
    </>
  );
}
`
}

// ── Login Route ─────────────────────────────────────────────────────────────

func desktopLoginRoute() string {
	return `import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { useState } from "react";
import { toast } from "sonner";
import { useAuth } from "../hooks/use-auth";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login(email, password);
      navigate({ to: "/" });
    } catch (err: any) {
      toast.error(err?.message || "Invalid credentials");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-background">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-foreground">Welcome Back</h1>
          <p className="text-sm text-text-secondary mt-1">Sign in to your account</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 bg-bg-elevated border border-border rounded-xl p-6">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="admin@example.com"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter password"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-accent hover:bg-accent/90 text-white font-medium py-2.5 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>
        <p className="text-center text-sm text-text-secondary mt-4">
          {"Don't have an account? "}
          <Link to="/register" className="text-accent hover:underline">Create one</Link>
        </p>
        <p className="text-center text-xs text-text-muted mt-2">Default: admin@example.com / admin123</p>
      </div>
    </div>
  );
}
`
}

// ── Register Route ──────────────────────────────────────────────────────────

func desktopRegisterRoute() string {
	return `import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { useState } from "react";
import { toast } from "sonner";
import { useAuth } from "../hooks/use-auth";

export const Route = createFileRoute("/register")({
  component: RegisterPage,
});

function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }
    if (password.length < 6) {
      toast.error("Password must be at least 6 characters");
      return;
    }
    setLoading(true);
    try {
      await register(name, email, password);
      navigate({ to: "/" });
    } catch (err: any) {
      toast.error(err?.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-background">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-foreground">Create Account</h1>
          <p className="text-sm text-text-secondary mt-1">Get started with your new account</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 bg-bg-elevated border border-border rounded-xl p-6">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your name"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="At least 6 characters"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Confirm Password</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Repeat password"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-accent hover:bg-accent/90 text-white font-medium py-2.5 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Creating account..." : "Create Account"}
          </button>
        </form>
        <p className="text-center text-sm text-text-secondary mt-4">
          {"Already have an account? "}
          <Link to="/login" className="text-accent hover:underline">Sign in</Link>
        </p>
      </div>
    </div>
  );
}
`
}

// ── Utilities (unchanged) ───────────────────────────────────────────────────

func desktopUtilsTS() string {
	return `import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(date: string | Date): string {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
`
}

func desktopQueryClientTS() string {
	return `import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});
`
}

func desktopUseAuthHook() string {
	return `import { createContext, useContext, useState, useEffect, ReactNode } from "react";

// These are auto-generated by Wails — they will exist when wails dev runs
// @ts-ignore
import { Login as WailsLogin, Register as WailsRegister, GetCurrentUser } from "../../wailsjs/go/main/App";

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const stored = localStorage.getItem("grit-user");
    if (stored) {
      try {
        const parsed = JSON.parse(stored);
        if (parsed && parsed.id) {
          GetCurrentUser(parsed.id)
            .then((u: any) => {
              if (u) setUser(u);
              else localStorage.removeItem("grit-user");
            })
            .catch(() => localStorage.removeItem("grit-user"))
            .finally(() => setIsLoading(false));
          return;
        }
      } catch {
        // Invalid JSON in localStorage, ignore
      }
    }
    setIsLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    const response = await WailsLogin(email, password);
    setUser(response.user);
    localStorage.setItem("grit-user", JSON.stringify(response.user));
  };

  const register = async (name: string, email: string, password: string) => {
    const response = await WailsRegister(name, email, password);
    setUser(response.user);
    localStorage.setItem("grit-user", JSON.stringify(response.user));
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem("grit-user");
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within AuthProvider");
  return context;
}
`
}

func desktopUseThemeHook() string {
	return `import { useState, useEffect } from "react";

export function useTheme() {
  const [theme, setTheme] = useState<"dark" | "light">(() => {
    if (typeof window !== "undefined") {
      return (localStorage.getItem("grit-theme") as "dark" | "light") || "dark";
    }
    return "dark";
  });

  useEffect(() => {
    const root = document.documentElement;
    if (theme === "light") {
      root.classList.add("light");
      root.classList.remove("dark");
    } else {
      root.classList.add("dark");
      root.classList.remove("light");
    }
    localStorage.setItem("grit-theme", theme);
  }, [theme]);

  const toggleTheme = () => setTheme((t) => (t === "dark" ? "light" : "dark"));

  return { theme, setTheme, toggleTheme };
}
`
}
