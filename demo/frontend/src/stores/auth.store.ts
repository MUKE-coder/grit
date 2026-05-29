import { useState, useEffect } from 'react'
import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import type { User, TokenPair, BusinessInfo, WorkspaceAccess } from '@/types'

/**
 * Workspace = which side of the business the user is currently working in.
 * Drives the sidebar contents AND the home dashboard. Persisted so users
 * don't have to re-pick on every reload.
 */
export type Workspace = 'loans' | 'spares'

interface AuthState {
  user: User | null
  tokens: TokenPair | null
  businesses: BusinessInfo[]
  currentBusinessID: number | null
  currentBranchID: number | null
  currentRole: string | null
  currentWorkspace: Workspace

  login: (user: User, tokens: TokenPair, businesses: BusinessInfo[]) => void
  logout: () => void
  switchBusiness: (businessID: number) => void
  setCurrentBranch: (branchID: number) => void
  setTokens: (tokens: TokenPair) => void
  setWorkspace: (w: Workspace) => void
  isAuthenticated: () => boolean
  getWorkspaceAccess: () => WorkspaceAccess
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      tokens: null,
      businesses: [],
      currentBusinessID: null,
      currentBranchID: null,
      currentRole: null,
      currentWorkspace: 'loans',

      login: (user, tokens, businesses) => {
        const firstBiz = businesses[0]
        set({
          user,
          tokens,
          businesses,
          currentBusinessID: firstBiz?.id ?? null,
          currentRole: firstBiz?.role ?? null,
          currentBranchID: firstBiz?.branch_id ?? null,
          // If this user can only see one workspace, auto-snap to it. Otherwise
          // default to loans (which is also the historic default).
          currentWorkspace:
            firstBiz?.workspace_access === 'spares' ? 'spares' :
            firstBiz?.workspace_access === 'loans'  ? 'loans'  :
            (get().currentWorkspace || 'loans'),
        })
      },

      logout: () => {
        set({
          user: null,
          tokens: null,
          businesses: [],
          currentBusinessID: null,
          currentBranchID: null,
          currentRole: null,
        })
      },

      switchBusiness: (businessID) => {
        const biz = get().businesses.find((b) => b.id === businessID)
        if (biz) {
          const ws =
            biz.workspace_access === 'spares' ? 'spares' :
            biz.workspace_access === 'loans'  ? 'loans'  :
            get().currentWorkspace
          set({
            currentBusinessID: biz.id,
            currentRole: biz.role,
            currentBranchID: biz.branch_id,
            currentWorkspace: ws,
          })
        }
      },

      setCurrentBranch: (branchID) => set({ currentBranchID: branchID }),
      setTokens: (tokens) => set({ tokens }),
      // Guard: don't set a workspace the user doesn't have access to.
      setWorkspace: (w) => {
        const access = get().getWorkspaceAccess()
        if (access === 'both' || access === w) set({ currentWorkspace: w })
      },

      isAuthenticated: () => {
        const { tokens } = get()
        if (!tokens) return false
        return tokens.expires_at * 1000 > Date.now()
      },

      /** Returns the workspace_access for the currently selected business, defaulting to "both". */
      getWorkspaceAccess: (): WorkspaceAccess => {
        const { businesses, currentBusinessID } = get()
        const biz = businesses.find((b) => b.id === currentBusinessID)
        return (biz?.workspace_access ?? 'both') as WorkspaceAccess
      },
    }),
    {
      name: 'grit-auth',
      storage: createJSONStorage(() => localStorage),
    }
  )
)

// Hook to check if Zustand has finished hydrating from localStorage
export function useAuthHydrated(): boolean {
  const [hydrated, setHydrated] = useState(useAuthStore.persist.hasHydrated())

  useEffect(() => {
    const unsub = useAuthStore.persist.onFinishHydration(() => {
      setHydrated(true)
    })
    return unsub
  }, [])

  return hydrated
}
