"use client";

import { createContext, useContext } from "react";
import type { User } from "@repo/shared/types";
import { useMe } from "@/hooks/use-auth";

// User shape comes from packages/shared/types/user.ts — same shape the
// Go API serialises, same shape the admin imports. Never inline it.

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextValue>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data: user, isLoading } = useMe();

  return (
    <AuthContext.Provider
      value={{
        user: user ?? null,
        isAuthenticated: !!user,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  return useContext(AuthContext);
}
