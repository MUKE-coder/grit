import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/mobile/building')

export default function MobileBuildingPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Mobile</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Building &amp; Publishing
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Expo Go is for development. To ship to the App Store and Google Play you use{' '}
                <strong>EAS Build</strong> — Expo&apos;s cloud service that compiles native
                binaries. This page covers configuring EAS, building, submitting to the stores,
                and pushing JavaScript-only updates over the air.
              </p>
              <LaneFlow
                id="mob-build"
                lanes={['Your app', 'EAS Build', 'Ship']}
                nodes={[
                  { id: 'code', lane: 0, row: 1, title: 'Expo app', sub: 'JS + native', tone: 'primary' },
                  { id: 'eas', lane: 1, row: 1, title: 'EAS Build', sub: 'cloud compile', tone: 'cyan' },
                  { id: 'ios', lane: 2, row: 0, title: 'App Store', sub: '.ipa', tone: 'blue' },
                  { id: 'android', lane: 2, row: 1, title: 'Play Store', sub: '.aab', tone: 'green' },
                  { id: 'ota', lane: 2, row: 2, title: 'OTA update', sub: 'JS-only', tone: 'amber' },
                ]}
                edges={[
                  { from: 'code', to: 'eas', label: 'build', tone: 'cyan' },
                  { from: 'eas', to: 'ios', tone: 'blue' },
                  { from: 'eas', to: 'android', tone: 'green' },
                  { from: 'code', to: 'ota', label: 'push', dashed: true, tone: 'amber' },
                ]}
                legend={[{ tone: 'cyan', label: 'Cloud build' }, { tone: 'amber', label: 'Over-the-air' }]}
                caption="EAS compiles native binaries in the cloud; ship to the stores or push JS-only OTA updates"
              />
            </div>

            {/* Two kinds of change */}
            <div className="prose-grit mb-4">
              <h2>Two ways to ship a change</h2>
              <p>
                Understanding the split up front saves a lot of time. A change either needs a new
                native binary or it doesn&apos;t:
              </p>
            </div>
            <CodeBlock language="text" filename="ship decision" code={`Did you add/upgrade a native module, or change app.json native config?
   │
   ├─ YES → EAS Build → submit to the stores (review required)
   │        (new SDK, new permission, new icon/splash, version bump)
   │
   └─ NO  → EAS Update (OTA) → users get it on next launch (no review)
            (JS/TS changes: bug fixes, UI tweaks, new screens)`} className="mb-10" />

            {/* app.json */}
            <div className="prose-grit mb-4">
              <h2>app.json — shipped in the scaffold</h2>
              <p>
                Grit generates <code>apps/expo/app.json</code> for you. It holds the app name,
                slug, scheme, icons, splash, and version. This is the file you edit to rename the
                app, bump the version, or change the bundle identifier before a store build.
              </p>
            </div>
            <CodeBlock language="json" filename="apps/expo/app.json (shape)" code={`{
  "expo": {
    "name": "myapp",
    "slug": "myapp",
    "scheme": "myapp",
    "version": "1.0.0",
    "icon": "./assets/icon.png",
    "splash": { "image": "./assets/splash.png" },
    "ios": { "bundleIdentifier": "com.yourco.myapp" },
    "android": { "package": "com.yourco.myapp" }
  }
}`} className="mb-4" />
            <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mb-10">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <strong className="text-foreground">Note.</strong> <code className="text-xs font-mono">app.json</code>{' '}
                ships in the scaffold, but <code className="text-xs font-mono">eas.json</code>{' '}
                does not — it is created for you by <code className="text-xs font-mono">eas build:configure</code>{' '}
                (below). That keeps EAS opt-in: you only get build config once you decide to build.
              </p>
            </div>

            {/* EAS Build */}
            <div className="prose-grit mb-4">
              <h2>EAS Build</h2>
              <p>
                Install the EAS CLI, log in to your Expo account, then configure and build. The
                first <code>eas build:configure</code> writes <code>eas.json</code> with the
                default build profiles (development, preview, production).
              </p>
            </div>
            <CodeBlock terminal code={`# Install the EAS CLI and log in
npm install -g eas-cli
eas login

# From the Expo app directory
cd apps/expo

# Creates eas.json + links the project to your Expo account
eas build:configure

# Build native binaries in the cloud
eas build --platform ios        # requires an Apple Developer account
eas build --platform android
eas build --platform all`} className="mb-4" />
            <div className="prose-grit mb-10">
              <p>
                Builds run on Expo&apos;s servers, so you can produce an iOS build from Windows or
                Linux without a Mac. When a build finishes, EAS gives you a link to the{' '}
                <code>.ipa</code> / <code>.aab</code> (or an installable dev/preview build you can
                open on a device).
              </p>
            </div>

            {/* Submit */}
            <div className="prose-grit mb-4">
              <h2>Submit to the stores</h2>
              <p>
                Once you have a production build, <code>eas submit</code> uploads it to App Store
                Connect or Google Play Console. You&apos;ll need the respective developer accounts
                and store listings set up first.
              </p>
            </div>
            <CodeBlock terminal code={`eas submit --platform ios
eas submit --platform android`} className="mb-4" />
            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">Store</th>
                    <th className="text-left p-3 font-medium">You need</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  <tr>
                    <td className="p-3">Apple App Store</td>
                    <td className="p-3">Apple Developer Program ($99/yr), a bundle identifier, an App Store Connect listing</td>
                  </tr>
                  <tr>
                    <td className="p-3">Google Play</td>
                    <td className="p-3">Google Play Developer account ($25 one-time), a package name, a Play Console listing</td>
                  </tr>
                </tbody>
              </table>
            </div>

            {/* OTA */}
            <div className="prose-grit mb-4">
              <h2>OTA updates with EAS Update</h2>
              <p>
                For JavaScript-only changes — bug fixes, copy tweaks, new screens, restyled
                components — you can push directly to users without a store review. Configure
                updates once (this installs <code>expo-updates</code> and wires the runtime), then
                publish to a branch:
              </p>
            </div>
            <CodeBlock terminal code={`# One-time: adds expo-updates + update config
eas update:configure

# Publish a JS bundle to the production branch
eas update --branch production --message "Fix login screen layout"`} className="mb-4" />
            <div className="prose-grit mb-10">
              <p>
                Installed apps pointed at that branch pick up the update on their next launch. OTA
                updates cannot change native code — anything that touches native modules or{' '}
                <code>app.json</code> native config still requires a fresh EAS Build and store
                submission. Keep the update&apos;s runtime version compatible with the installed
                binary.
              </p>
            </div>

            {/* API deployment */}
            <div className="prose-grit mb-4">
              <h2>Deploy the API separately</h2>
              <p>
                The mobile app is only half the story — it needs the Go API running somewhere
                public. The API deploys exactly like the API-only architecture: as a Docker
                container. Point the app at it with <code>EXPO_PUBLIC_API_URL</code> so release
                builds hit production instead of a LAN IP.
              </p>
            </div>
            <CodeBlock terminal code={`# Build and push the API image
cd apps/api
docker build -t myapp-api .
docker push registry.example.com/myapp-api:latest

# In the Expo build environment
EXPO_PUBLIC_API_URL=https://api.myapp.com`} className="mb-4" />
            <div className="rounded-lg border border-border/40 bg-accent/20 p-4 mb-10">
              <p className="text-sm text-muted-foreground leading-relaxed">
                Set <code className="text-xs font-mono">EXPO_PUBLIC_API_URL</code> as an
                environment variable / EAS secret for your production build profile. Grit&apos;s
                API client uses it verbatim (appending <code className="text-xs font-mono">/api</code>),
                so the shipped app never falls back to a development host.
              </p>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/resource-generation" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Resource Generation
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/offline" className="gap-1.5">
                  Offline &amp; Sync
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
