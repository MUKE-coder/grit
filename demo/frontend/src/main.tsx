import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import { routeTree } from './routeTree.gen'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

const router = createRouter({
  routeTree,
  context: { queryClient },
  // Preload route chunks on mount, not just on hover. Combined with React
  // Query's 5-minute staleTime above, this makes sidebar nav feel instant
  // after app boot — by the time the user clicks, both the JS chunk and
  // any cached query data are already in memory.
  defaultPreload: 'intent',
  defaultPreloadDelay: 0,
  defaultPreloadStaleTime: 5 * 60 * 1000,
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 3000,
          style: {
            background: '#FFFFFF',
            color: '#0F172A',
            borderRadius: '12px',
            border: '1px solid #E2E8F0',
            boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
          },
          success: {
            iconTheme: { primary: '#16A34A', secondary: '#FFFFFF' },
          },
          error: {
            iconTheme: { primary: '#EF4444', secondary: '#FFFFFF' },
          },
        }}
      />
    </QueryClientProvider>
  </React.StrictMode>,
)
