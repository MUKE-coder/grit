package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDesktopFrontendConfig(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "frontend", "package.json"):       desktopPackageJSON(opts),
		filepath.Join(root, "frontend", "vite.config.ts"):     desktopViteConfig(),
		filepath.Join(root, "frontend", "tailwind.config.ts"): desktopTailwindConfig(),
		filepath.Join(root, "frontend", "postcss.config.js"):  desktopPostcssConfig(),
		filepath.Join(root, "frontend", "tsconfig.json"):      desktopTSConfig(),
		filepath.Join(root, "frontend", "index.html"):         desktopIndexHTML(),
		filepath.Join(root, "frontend", "src", "main.tsx"):    desktopMainTSX(),
		filepath.Join(root, "frontend", "src", "index.css"):   desktopIndexCSS(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopPackageJSON(opts DesktopOptions) string {
	return fmt.Sprintf(`{
  "name": "%s-frontend",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "@tanstack/react-router": "^1.95.0",
    "@tanstack/react-query": "^5.62.0",
    "react-hook-form": "^7.54.0",
    "@hookform/resolvers": "^4.1.0",
    "zod": "^3.24.0",
    "lucide-react": "^0.468.0",
    "clsx": "^2.1.0",
    "tailwind-merge": "^2.6.0",
    "date-fns": "^4.1.0",
    "sonner": "^1.7.0"
  },
  "devDependencies": {
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@vitejs/plugin-react": "^4.3.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.7.0",
    "vite": "^6.0.0",
    "@tanstack/router-plugin": "^1.95.0"
  }
}
`, opts.ProjectName)
}

func desktopViteConfig() string {
	return `import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import path from "path";

export default defineConfig({
  plugins: [TanStackRouterVite(), react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
`
}

func desktopTailwindConfig() string {
	return `import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        "bg-elevated": "var(--bg-elevated)",
        "bg-hover": "var(--bg-hover)",
        border: "var(--border)",
        foreground: "var(--foreground)",
        "text-secondary": "var(--text-secondary)",
        "text-muted": "var(--text-muted)",
        accent: "var(--accent)",
        success: "var(--success)",
        danger: "var(--danger)",
      },
      fontFamily: {
        sans: ["Onest", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
    },
  },
  plugins: [],
};

export default config;
`
}

func desktopPostcssConfig() string {
	return `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}

func desktopTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": false,
    "noUnusedParameters": false,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"]
}
`
}

func desktopIndexHTML() string {
	return `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Grit Desktop</title>
</head>
<body class="bg-background text-foreground antialiased min-h-screen">
  <div id="root"></div>
  <script type="module" src="/src/main.tsx"></script>
</body>
</html>
`
}

func desktopMainTSX() string {
	return `import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createRouter, createHashHistory } from "@tanstack/react-router";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./lib/query-client";
import { AuthProvider } from "./hooks/use-auth";
import { routeTree } from "./routeTree.gen";
import "./index.css";

const hashHistory = createHashHistory();
const router = createRouter({ routeTree, history: hashHistory });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <AuthProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </AuthProvider>
  </React.StrictMode>
);
`
}

func desktopIndexCSS() string {
	return `@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --background: #0a0a0f;
  --bg-elevated: #22222e;
  --bg-hover: #2a2a38;
  --border: #2a2a3a;
  --foreground: #e8e8f0;
  --text-secondary: #9090a8;
  --text-muted: #606078;
  --accent: #6c5ce7;
  --accent-hover: #7d6ff0;
  --success: #00b894;
  --danger: #ff6b6b;
  --warning: #fdcb6e;
}

.light {
  --background: #fafafa;
  --bg-elevated: #ffffff;
  --bg-hover: #f0f0f5;
  --border: #e2e2ea;
  --foreground: #1a1a2e;
  --text-secondary: #64648c;
  --text-muted: #9090a8;
  --accent: #6c5ce7;
  --accent-hover: #5a4bd6;
  --success: #00b894;
  --danger: #ff6b6b;
  --warning: #fdcb6e;
}

body {
  font-family: "Onest", system-ui, -apple-system, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--background);
}

::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-muted);
}
`
}
