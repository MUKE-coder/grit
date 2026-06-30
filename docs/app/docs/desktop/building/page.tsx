import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from "@/config/docs-metadata";
import { Download } from "lucide-react";

export const metadata = getDocMetadata("/docs/desktop/building");

export default function DesktopBuildingPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Desktop (Wails)
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Building & Distribution
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Compile your desktop app into a native executable and distribute
                it to users.
              </p>
            </div>

            {/* Compile */}
            <div className="prose-grit mb-10">
              <h2>Compile</h2>
              <p>
                Build a production-ready native executable with a single command:
              </p>
            </div>

            <CodeBlock terminal code="grit compile" />

            <div className="prose-grit mb-10 mt-6">
              <p>
                This auto-detects your desktop project and runs{" "}
                <code>wails build</code>. The output binary is placed in{" "}
                <code>build/bin/</code>.
              </p>
              <p>
                You can also run <code>wails build</code> directly for more
                control over build flags.
              </p>
            </div>

            {/* Output */}
            <div className="prose-grit mb-10">
              <h2>Build Output</h2>
            </div>

            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">Platform</th>
                    <th className="text-left p-3 font-medium">Output</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  <tr>
                    <td className="p-3 font-medium text-foreground">Windows</td>
                    <td className="p-3 font-mono text-xs">
                      build/bin/myapp.exe
                    </td>
                  </tr>
                  <tr>
                    <td className="p-3 font-medium text-foreground">macOS</td>
                    <td className="p-3 font-mono text-xs">
                      build/bin/myapp.app
                    </td>
                  </tr>
                  <tr>
                    <td className="p-3 font-medium text-foreground">Linux</td>
                    <td className="p-3 font-mono text-xs">build/bin/myapp</td>
                  </tr>
                </tbody>
              </table>
            </div>

            {/* Cross-Platform */}
            <div className="prose-grit mb-10">
              <h2>Cross-Platform Builds</h2>
              <p>
                Build for a specific platform using the{" "}
                <code>-platform</code> flag:
              </p>
            </div>

            <CodeBlock
              terminal
              code={`# Windows
wails build -platform windows/amd64

# macOS (Intel)
wails build -platform darwin/amd64

# macOS (Apple Silicon)
wails build -platform darwin/arm64

# Linux
wails build -platform linux/amd64`}
            />

            {/* NSIS Installer */}
            <div className="prose-grit mb-10 mt-10">
              <h2>Windows Installer (NSIS)</h2>
              <p>
                Create a Windows installer with NSIS (requires{" "}
                <a
                  href="https://nsis.sourceforge.io/"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline"
                >
                  NSIS
                </a>{" "}
                installed):
              </p>
            </div>

            <CodeBlock terminal code="wails build -nsis" />

            <div className="prose-grit mb-10 mt-6">
              <p>
                This produces a <code>.exe</code> installer that handles
                installation, Start Menu shortcuts, and uninstallation.
              </p>
            </div>

            {/* Tips */}
            <div className="prose-grit mb-10">
              <h2>Production Tips</h2>
            </div>

            <div className="space-y-4 mb-10">
              <div className="rounded-lg border border-border bg-card p-4">
                <h3 className="font-semibold text-foreground mb-1">
                  App Icon
                </h3>
                <p className="text-sm text-muted-foreground">
                  Place your icon at <code>build/appicon.png</code> (1024x1024).
                  Wails automatically converts it for each platform during
                  build.
                </p>
              </div>
              <div className="rounded-lg border border-border bg-card p-4">
                <h3 className="font-semibold text-foreground mb-1">
                  Environment Config
                </h3>
                <p className="text-sm text-muted-foreground">
                  The <code>.env</code> file is read at runtime from the working
                  directory. For production, either bundle it alongside the
                  executable or use OS environment variables.
                </p>
              </div>
              <div className="rounded-lg border border-border bg-card p-4">
                <h3 className="font-semibold text-foreground mb-1">
                  Embedded Assets
                </h3>
                <p className="text-sm text-muted-foreground">
                  The React frontend is embedded into the Go binary via{" "}
                  <code>//go:embed all:frontend/dist</code>. No separate files
                  need to be distributed — everything is in one executable.
                </p>
              </div>
              <div className="rounded-lg border border-border bg-card p-4">
                <h3 className="font-semibold text-foreground mb-1">
                  Database Location
                </h3>
                <p className="text-sm text-muted-foreground">
                  SQLite creates <code>app.db</code> in the working directory.
                  For a better user experience, consider using{" "}
                  <code>os.UserConfigDir()</code> to store the database in the
                  user&apos;s config folder.
                </p>
              </div>
            </div>

            {/* wails.json */}
            <div className="prose-grit mb-10">
              <h2>Wails Configuration</h2>
              <p>
                The <code>wails.json</code> file controls build settings:
              </p>
            </div>

            <CodeBlock
              code={`{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "Your Name"
  }
}`}
            />

            {/* Desktop Handbook */}
            <div className="mt-10">
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfpiJDPD3QgNG9hYzVFo5iLR0yrDPTJedWnBH7"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-3 rounded-lg border border-primary/25 bg-primary/5 px-5 py-4 hover:bg-primary/10 transition-colors group"
              >
                <Download className="h-5 w-5 text-primary/70 group-hover:text-primary shrink-0" />
                <div>
                  <span className="text-sm font-semibold text-foreground block">Download Desktop Handbook (PDF)</span>
                  <span className="text-xs text-muted-foreground/60">Complete offline reference for Grit Desktop development</span>
                </div>
              </a>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
