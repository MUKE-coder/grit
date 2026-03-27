import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build & Distribution — Grit Desktop Course',
  description: 'Compile your Grit desktop app into a native binary, target multiple platforms, customize app icons, and distribute to users.',
}

export default function BuildCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-desktop" className="hover:text-foreground transition-colors">Grit Desktop</Link>
          <span>/</span>
          <span className="text-foreground">Build & Distribution</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 5 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build & Distribution
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this final course, you will compile your desktop app into a native binary,
            learn about platform targets, customize the app icon, optimize the file size,
            and understand how to distribute your app to users.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Development vs Production ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Development vs Production</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Until now, you{"'"}ve been running your app in development mode with <Code>grit start</Code>.
            Development mode includes tools that help you build the app but are not needed by end users:
          </p>

          <Definition term="Development Build">
            A version of your app optimized for developers. Includes hot reload (instant updates when
            you save a file), detailed error messages, debug tools (DevTools, console logs), and
            source maps. Slower and larger than a production build, but much easier to work with.
          </Definition>

          <Definition term="Production Build">
            A version of your app optimized for end users. The Go code is compiled into a single
            native binary. The React code is bundled, minified, and embedded inside that binary.
            No source code is visible, no debug tools, no hot reload — just a fast, small, distributable
            executable file.
          </Definition>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Development</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Production</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Hot reload</td>
                  <td className="px-4 py-3">Yes</td>
                  <td className="px-4 py-3">No</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">DevTools</td>
                  <td className="px-4 py-3">Available</td>
                  <td className="px-4 py-3">Disabled</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">File size</td>
                  <td className="px-4 py-3">Large (many files)</td>
                  <td className="px-4 py-3">Small (single binary)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Speed</td>
                  <td className="px-4 py-3">Slower (interpreted)</td>
                  <td className="px-4 py-3">Fast (compiled)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Requirements</td>
                  <td className="px-4 py-3">Go, Node, pnpm, Wails</td>
                  <td className="px-4 py-3">Nothing (standalone binary)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        {/* ═══ Building Your App ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Building Your App</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To compile your app into a production binary, run:
          </p>

          <CodeBlock filename="Terminal">
{`grit compile`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Or use the Wails command directly:
          </p>

          <CodeBlock filename="Terminal">
{`wails build`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is what happens during the build:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> <strong className="text-foreground">Vite builds React</strong> — bundles, minifies, and optimizes all JavaScript, CSS, and assets</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> <strong className="text-foreground">Go embeds the frontend</strong> — the built React files are embedded into the Go binary using <Code>go:embed</Code></li>
            <li className="flex gap-2"><span className="text-primary">3.</span> <strong className="text-foreground">Go compiles to native code</strong> — the entire application (backend + embedded frontend) becomes a single executable</li>
          </ul>

          <Definition term="Native Binary">
            A compiled executable file that runs directly on the operating system without needing
            Go, Node, or any runtime installed. The user downloads one file and double-clicks it.
            On Windows it{"'"}s a <Code>.exe</Code>, on macOS it{"'"}s a <Code>.app</Code> bundle, and on Linux it{"'"}s
            an ELF binary.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The output binary is placed in the <Code>build/bin/</Code> directory:
          </p>

          <CodeBlock filename="Build output">
{`build/
└── bin/
    └── myapp.exe     # Windows
    └── myapp         # macOS / Linux`}
          </CodeBlock>

          <Challenge number={1} title="Build Your App">
            <p>Run <Code>grit compile</Code> in your project directory. How long does it take? Where is the output file? What is the file size? Try double-clicking the binary — does the app open?</p>
          </Challenge>
        </section>

        {/* ═══ Platform Targets ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Platform Targets</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Wails can build for different operating systems and architectures:
          </p>

          <CodeBlock filename="Terminal">
{`# Build for your current platform (default)
wails build

# Build for Windows (64-bit)
wails build -platform windows/amd64

# Build for macOS (Intel)
wails build -platform darwin/amd64

# Build for macOS (Apple Silicon / M1, M2, M3)
wails build -platform darwin/arm64

# Build for Linux (64-bit)
wails build -platform linux/amd64`}
          </CodeBlock>

          <Definition term="Cross-platform">
            Software that runs on multiple operating systems. Wails apps can target Windows, macOS,
            and Linux from the same Go + React codebase. You write your code once and compile it
            for each platform.
          </Definition>

          <Note>
            Cross-compilation (building for a different OS than your development machine) may require
            additional tools. Building a macOS app from Windows requires a macOS machine or CI service.
            Building a Windows app from macOS works with the right toolchain installed. Building for
            your current platform always works out of the box.
          </Note>

          <Challenge number={2} title="Build for Your Platform">
            <p>Build for your current platform with <Code>grit compile</Code>. Find the output binary in <Code>build/bin/</Code>. Run it directly (double-click on Windows, <Code>./myapp</Code> on macOS/Linux). Does it work the same as dev mode?</p>
          </Challenge>
        </section>

        {/* ═══ App Icon ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">App Icon</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every desktop app needs an icon — it appears in the taskbar, dock, file explorer, and
            the window title bar. Wails uses icon files in the <Code>build/</Code> directory:
          </p>

          <CodeBlock filename="Icon files">
{`build/
├── appicon.png        # 1024x1024 PNG (used for all platforms)
├── windows/
│   └── icon.ico       # Windows icon (multi-size)
├── darwin/
│   └── appicon.icns   # macOS icon (multi-size)
└── linux/
    └── appicon.png    # Linux icon`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To change the icon:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> Create a 1024x1024 pixel PNG image for your app icon</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> Replace <Code>build/appicon.png</Code> with your new image</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> Rebuild with <Code>grit compile</Code> — Wails converts it to platform-specific formats</li>
          </ul>

          <Tip>
            Use a simple, recognizable design for your icon. It will be displayed at sizes from 16x16
            to 1024x1024 pixels, so avoid fine details that disappear at small sizes. Tools like
            Figma, Canva, or icon generators can help you create one.
          </Tip>

          <Challenge number={3} title="Change the App Icon">
            <p>Find the default app icon file in the <Code>build/</Code> directory. Replace it with a custom PNG (1024x1024). Rebuild and check — does the new icon appear in the taskbar and window title?</p>
          </Challenge>
        </section>

        {/* ═══ Window Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Window Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>wails.json</Code> file controls how the production window appears:
          </p>

          <CodeBlock filename="wails.json (window settings)">
{`{
  "name": "My App",
  "width": 1024,
  "height": 768,
  "minWidth": 800,
  "minHeight": 600,
  "frameless": true,
  "fullscreen": false,
  "resizable": true
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>name</Code> field is what appears in the operating system{"'"}s task manager and window list.
            The size settings define the default dimensions and minimum size.
          </p>

          <Challenge number={4} title="Configure Window Size">
            <p>Change the default window size to 1200x800 in <Code>wails.json</Code>. Rebuild with <Code>grit compile</Code> and run the binary. Is the window larger? Try resizing it below the minimum dimensions — does it stop at the minimum?</p>
          </Challenge>
        </section>

        {/* ═══ File Size ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">File Size</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A typical Grit desktop app binary is <strong className="text-foreground">15-30 MB</strong>.
            This includes the entire Go backend, all React code, CSS, and embedded assets. For comparison,
            an Electron app with similar features would be 100-200 MB.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can reduce the binary size by stripping debug symbols:
          </p>

          <CodeBlock filename="Terminal">
{`# Build with stripped debug symbols (smaller binary)
wails build -ldflags "-s -w"

# -s removes the symbol table
# -w removes DWARF debugging information`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For even smaller binaries, you can use UPX compression (a tool that compresses executables):
          </p>

          <CodeBlock filename="Terminal">
{`# Install UPX (if not already installed)
# Windows: scoop install upx
# macOS: brew install upx
# Linux: apt install upx

# Compress the binary
upx --best build/bin/myapp.exe`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Typical size reduction:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Build Type</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Typical Size</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Default build</td>
                  <td className="px-4 py-3">25-30 MB</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Stripped (-s -w)</td>
                  <td className="px-4 py-3">15-20 MB</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Stripped + UPX</td>
                  <td className="px-4 py-3">8-12 MB</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={5} title="Optimize Binary Size">
            <p>Build your app normally and note the file size. Then rebuild with <Code>wails build -ldflags {'"'}-s -w{'"'}</Code>. How much smaller is it? If you have UPX installed, compress it further. What{"'"}s the final size?</p>
          </Challenge>
        </section>

        {/* ═══ Distribution ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Distribution</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once your app is compiled, you need to get it to your users. There are several distribution methods:
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Direct Download</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The simplest method — upload the binary to your website, GitHub Releases, or a file hosting service.
            Users download the file and run it. No installer needed.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Windows Installer (NSIS)</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Wails can generate an NSIS installer for Windows. This creates a <Code>.exe</Code> installer
            that adds your app to the Start Menu, creates desktop shortcuts, and supports uninstallation:
          </p>

          <CodeBlock filename="Terminal">
{`wails build -nsis`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3">macOS DMG</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            On macOS, Wails builds a <Code>.app</Code> bundle. You can package it into a <Code>.dmg</Code> disk
            image for distribution — the standard macOS app distribution format.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Linux AppImage</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            For Linux, you can package the binary as an AppImage — a single executable that works on
            most Linux distributions without installation.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Code Signing</h3>

          <Definition term="Code Signing">
            Digitally signing your application with a certificate so the operating system trusts it.
            Without code signing, Windows shows a {'"'}Windows protected your PC{'"'} warning, and macOS
            shows an {'"'}unidentified developer{'"'} warning. Code signing tells the OS that the app
            comes from a verified developer.
          </Definition>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Windows</strong> — requires a code signing certificate ($200-400/year) from a Certificate Authority</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">macOS</strong> — requires an Apple Developer account ($99/year) and notarization through Apple{"'"}s servers</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Linux</strong> — code signing is optional; most Linux users install from package managers or trust direct downloads</li>
          </ul>

          <Note>
            Code signing is important for production apps distributed to many users, but it{"'"}s not
            required for personal projects, internal tools, or early development. You can distribute
            unsigned apps — users just need to click through an extra warning dialog.
          </Note>

          <Challenge number={6} title="Build an Installer">
            <p>If you{"'"}re on Windows, try building with <Code>wails build -nsis</Code> to generate an installer. Run the installer — does it create a Start Menu entry and desktop shortcut? Can you uninstall it from Settings?</p>
          </Challenge>
        </section>

        {/* ═══ GitHub Releases ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Publishing with GitHub Releases</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The most common way to distribute desktop apps for open-source projects is GitHub Releases.
            You build for each platform and upload the binaries:
          </p>

          <CodeBlock filename="Terminal">
{`# Build for your platform
grit compile

# Create a GitHub release (using gh CLI)
gh release create v1.0.0 build/bin/myapp.exe --title "v1.0.0" --notes "Initial release"

# Or upload multiple platform binaries
gh release create v1.0.0 \
  build/bin/myapp-windows-amd64.exe \
  build/bin/myapp-darwin-amd64 \
  build/bin/myapp-linux-amd64 \
  --title "v1.0.0" --notes "Initial release"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Users visit your GitHub repository{"'"}s Releases page and download the binary for their platform.
          </p>

          <Challenge number={7} title="Prepare for Release">
            <p>Build your app for your current platform. Create a <Code>CHANGELOG.md</Code> file documenting what your app does. If you have a GitHub repository, create a release and upload the binary.</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The difference between development builds and production builds</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to compile your app with <Code>grit compile</Code> or <Code>wails build</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Platform targets — building for Windows, macOS, and Linux</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to customize the app icon</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Window configuration in wails.json</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Optimizing binary size with <Code>-ldflags {'"'}-s -w{'"'}</Code> and UPX</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Distribution methods — direct download, NSIS installer, DMG, AppImage</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> What code signing is and when you need it</li>
          </ul>

          <Challenge number={8} title="Full Build Challenge">
            <p>Build your notes app for your current platform. Test the binary — create some data, export to PDF, verify everything works the same as dev mode. Check the file size.</p>
          </Challenge>

          <Challenge number={9} title="Customize and Rebuild">
            <p>Change the window title in <Code>wails.json</Code>. Replace the app icon with a custom one. Rebuild with <Code>-ldflags {'"'}-s -w{'"'}</Code>. Compare the file size to the previous build. Does the new title and icon appear?</p>
          </Challenge>

          <Challenge number={10} title="Share Your App">
            <p>Send the compiled binary to a friend (or test on a different computer). Does it run without installing Go, Node, or Wails? Does the SQLite database get created automatically on first launch? This is the power of a native binary — zero dependencies for the end user.</p>
          </Challenge>
        </section>

        {/* ═══ Course Complete ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Course Complete</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Congratulations — you{"'"}ve completed the entire Grit Desktop course track. You now know how to:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Scaffold a desktop app with <Code>grit new-desktop</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> Generate full CRUD resources with Wails bindings</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Customize the UI with frameless windows, themes, and shadcn/ui</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Export data to PDF, Excel, and CSV</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Compile and distribute a native binary</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Next steps: build something real. A personal finance tracker, a note-taking app, a project
            management tool, a CRM for your business — the possibilities are endless. Every app follows
            the same pattern: scaffold, generate resources, customize the UI, export data, and ship.
          </p>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-desktop/export', label: 'Prev: PDF & Excel Export' }}
            next={{ href: '/courses/grit-desktop', label: 'Back to All Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
