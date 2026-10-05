"use client";

import { useState } from "react";
import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { makeQueryClient } from "@/lib/query-client";

export function Providers({ children }: { children: React.ReactNode }) {
  // One client per mount: per request on the server, once in the browser.
  const [queryClient] = useState(makeQueryClient);

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {/*
        richColors gives us sonner's tinted success/error/warning/info
        surfaces. We override sonner's default palette via CSS vars in
        globals.css so green = Grit success (#00b894), red = danger
        (#ff6b6b), etc. — matching the rest of the design system.
      */}
      <Toaster
        richColors
        position="bottom-right"
        theme="dark"
        toastOptions={{
          classNames: {
            toast: "grit-toast",
          },
        }}
      />
    </QueryClientProvider>
  );
}
