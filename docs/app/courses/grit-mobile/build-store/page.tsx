import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build & App Store — Grit Mobile Course',
  description: 'Learn EAS Build, TestFlight, Play Store submission, OTA updates, and how to deploy your Go API for production mobile apps.',
}

export default function BuildStoreCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-mobile" className="hover:text-foreground transition-colors">Grit Mobile</Link>
          <span>/</span>
          <span className="text-foreground">Build & App Store</span>
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
            Build & App Store
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Your app works in development — now it{"'"}s time to ship it. In this course, you will learn
            to build production binaries with EAS, distribute through TestFlight and the Play Store,
            prepare your store listing, and push updates without going through app review.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Development vs Production Build ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Development vs Production Build</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            During development, you run the app through Expo Go — a shell app that loads your
            JavaScript code over the network. This is fast for iteration but not distributable.
            For the App Store, you need a <strong className="text-foreground">standalone binary</strong> — a
            real <Code>.ipa</Code> (iOS) or <Code>.apk</Code>/<Code>.aab</Code> (Android) file.
          </p>

          <Definition term="Development Build">
            Running your app through Expo Go or a development client. The JavaScript code is loaded
            over the network from your computer. Fast to iterate but cannot be submitted to the
            App Store or Play Store.
          </Definition>

          <Definition term="Production Build">
            A standalone binary (<Code>.ipa</Code> for iOS, <Code>.aab</Code> for Android) that
            contains your app{"'"}s JavaScript bundle and all native dependencies. This is what you
            submit to the App Store or Play Store. Users download and install it like any other app.
          </Definition>

          <Challenge number={1} title="Understand the Difference">
            <p>Explain in your own words why you cannot submit an Expo Go-based app to the App Store. What is the difference between loading JavaScript over the network and bundling it into a standalone binary?</p>
          </Challenge>
        </section>

        {/* ═══ EAS Build ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">EAS Build</h2>

          <Definition term="EAS (Expo Application Services)">
            A cloud build service by Expo. EAS builds your app on Expo{"'"}s servers so you don{"'"}t
            need Xcode (macOS only) or Android Studio installed locally. You run a single command,
            and EAS handles compilation, signing, and packaging.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Install the EAS CLI and configure your project:
          </p>

          <CodeBlock language="bash">{`# Install EAS CLI globally
npm install -g eas-cli

# Log in to your Expo account
eas login

# Initialize EAS in your project
cd apps/expo && eas build:configure`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates an <Code>eas.json</Code> file with build profiles:
          </p>

          <CodeBlock language="json">{`{
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal"
    },
    "preview": {
      "distribution": "internal"
    },
    "production": {}
  }
}`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To build for a specific platform:
          </p>

          <CodeBlock language="bash">{`# Build for iOS
eas build --platform ios

# Build for Android
eas build --platform android

# Build for both
eas build --platform all`}</CodeBlock>

          <Tip>
            Use the <Code>preview</Code> profile for internal testing (distributes via QR code or link)
            and <Code>production</Code> for store submissions.
          </Tip>

          <Challenge number={2} title="Set Up EAS">
            <p>Install <Code>eas-cli</Code> with <Code>npm install -g eas-cli</Code>. Run <Code>eas whoami</Code> to check if you{"'"}re logged in. If not, run <Code>eas login</Code> to create an account and log in.</p>
          </Challenge>
        </section>

        {/* ═══ iOS Build ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">iOS Build</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Building for iOS requires an <strong className="text-foreground">Apple Developer account</strong> ($99/year).
            Apple requires all iOS apps to be signed with certificates and provisioning profiles.
            EAS handles this automatically — it creates and manages certificates for you.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Requirements for an iOS build:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Apple Developer account ($99/year)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> App ID registered in Apple Developer portal (EAS does this automatically)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Signing certificates and provisioning profiles (EAS manages these)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Bundle identifier (e.g., <Code>com.yourcompany.fitness</Code>)</li>
          </ul>

          <Definition term="TestFlight">
            Apple{"'"}s official beta testing platform. After building your iOS app, you upload it to
            TestFlight through App Store Connect. You can then invite up to 10,000 testers who
            install the app via the TestFlight app. This lets you test with real users before
            submitting to the public App Store.
          </Definition>

          <CodeBlock language="bash">{`# Build for iOS production
eas build --platform ios --profile production

# Submit to App Store Connect (TestFlight)
eas submit --platform ios`}</CodeBlock>

          <Note>
            You do <strong>not</strong> need a Mac to build for iOS with EAS. The build happens on
            Expo{"'"}s cloud servers (which run macOS). You only need a Mac if you want to build locally.
          </Note>

          <Challenge number={3} title="iOS Requirements">
            <p>What do you need to build for iOS? List all the requirements. Which ones does EAS handle automatically, and which ones do you need to set up yourself?</p>
          </Challenge>
        </section>

        {/* ═══ Android Build ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Android Build</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Android builds require a <strong className="text-foreground">signing keystore</strong> — a file
            that proves you are the developer of the app. EAS can generate one for you automatically
            on your first build.
          </p>

          <Definition term="APK (Android Package)">
            An installable file you can share directly with anyone. They open the file on their
            Android phone and install it. No Play Store required. Good for testing and internal
            distribution.
          </Definition>

          <Definition term="AAB (Android App Bundle)">
            An optimized format required by the Google Play Store. Google processes the AAB and
            generates optimized APKs for each device configuration (screen size, CPU architecture).
            This results in smaller downloads for users.
          </Definition>

          <CodeBlock language="bash">{`# Build an APK (for direct sharing/testing)
eas build --platform android --profile preview

# Build an AAB (for Play Store submission)
eas build --platform android --profile production

# Submit to Google Play Store
eas submit --platform android`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For the Play Store, you need a <strong className="text-foreground">Google Play Developer account</strong> ($25
            one-time fee). Unlike Apple, the Android build process works on any operating system.
          </p>

          <Challenge number={4} title="APK vs AAB">
            <p>What{"'"}s the difference between an APK and an AAB? When would you use each? If you wanted to send a test build to a friend who has an Android phone, which format would you use?</p>
          </Challenge>
        </section>

        {/* ═══ App Store Assets ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">App Store Assets</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Both the Apple App Store and Google Play Store require specific assets before you can
            publish your app:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Asset</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Requirement</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">App Icon</td>
                  <td className="px-4 py-3">1024x1024 PNG, no transparency (iOS), 512x512 PNG (Android)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Screenshots</td>
                  <td className="px-4 py-3">Various sizes for different devices (iPhone, iPad, Android phones/tablets)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Description</td>
                  <td className="px-4 py-3">Short description (~80 chars) + full description (~4000 chars)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Keywords</td>
                  <td className="px-4 py-3">Up to 100 characters of comma-separated keywords (iOS only)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Privacy Policy</td>
                  <td className="px-4 py-3">A URL to your privacy policy (required by both stores)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Category</td>
                  <td className="px-4 py-3">Primary and optional secondary category (e.g., Health & Fitness)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Tip>
            Use a tool like <a href="https://hotpot.ai/icon-resizer" target="_blank" rel="noreferrer" className="text-primary hover:underline">Hotpot Icon Resizer</a> to
            generate all required icon sizes from a single 1024x1024 source image.
          </Tip>

          <Challenge number={5} title="Write a Store Description">
            <p>Write a 3-sentence App Store description for your fitness app. It should explain what the app does, who it{"'"}s for, and what makes it useful. Keep it clear and compelling.</p>
          </Challenge>

          <Challenge number={6} title="Plan Your Screenshots">
            <p>List 5 screenshots you would take for your app{"'"}s store listing. For each, describe the screen shown and what feature it highlights. The first screenshot is the most important — it{"'"}s what users see when browsing.</p>
          </Challenge>
        </section>

        {/* ═══ OTA Updates ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">OTA Updates</h2>

          <Definition term="OTA (Over-The-Air) Updates">
            Pushing JavaScript updates directly to users{"'"} devices without going through the App Store
            or Play Store review process. When the user opens the app, it checks for updates and
            downloads the new JavaScript bundle in the background. The next time they open the app,
            they see the updated version.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            OTA updates are one of the biggest advantages of React Native with Expo. App Store review
            can take 1-3 days. With OTA, your fix reaches users in minutes.
          </p>

          <CodeBlock language="bash">{`# Push an OTA update to all users on the production channel
eas update --branch production --message "Fix workout timer bug"

# Push to a preview channel for testing first
eas update --branch preview --message "Test new feature"`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            There are important limitations to OTA updates:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> <strong className="text-foreground">Can update:</strong> JavaScript code, images, fonts — anything in the JS bundle</li>
            <li className="flex gap-2"><span className="text-red-400 shrink-0 mt-0.5">✕</span> <strong className="text-foreground">Cannot update:</strong> Native code, new native modules, app icon, splash screen</li>
          </ul>

          <Note>
            If you add a new native dependency (e.g., a camera library), you need a full build
            and store submission. OTA only works for JavaScript-level changes.
          </Note>

          <Challenge number={7} title="OTA vs Full Build">
            <p>Explain why OTA updates are useful. Give two examples of changes that can be shipped via OTA and two examples that require a full build and store submission.</p>
          </Challenge>
        </section>

        {/* ═══ Deploying the API ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Deploying the API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your mobile app needs a live API to talk to. During development, the app connects to
            <Code>localhost</Code>. In production, it connects to your deployed API server.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Deploy the Go API using <Code>grit deploy</Code> (covered in the Grit Web Deploy course)
            to a VPS with a domain and HTTPS. Then update the mobile app{"'"}s API URL:
          </p>

          <CodeBlock language="tsx">{`// lib/api.ts
const API_URL = __DEV__
  ? "http://localhost:8080"       // Development
  : "https://api.fitness.com"     // Production`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>__DEV__</Code> variable is <Code>true</Code> in development and <Code>false</Code> in
            production builds. This lets you switch between local and remote APIs automatically.
          </p>

          <Tip>
            Always use HTTPS for your production API. Mobile platforms are increasingly strict about
            insecure connections — iOS blocks HTTP by default, and Android 9+ does the same.
          </Tip>

          <Challenge number={8} title="Update the API URL">
            <p>If your API is deployed at <Code>api.fitness.com</Code>, what changes would you make in the mobile app{"'"}s API client? Where is the base URL configured? How do you handle the difference between development and production?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what you learned in this course:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> EAS Build compiles your app in the cloud — no Xcode or Android Studio needed</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> iOS requires an Apple Developer account; Android requires a Google Play Developer account</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> TestFlight distributes iOS betas; APKs distribute Android test builds directly</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> OTA updates ship JS changes instantly — no store review required</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Deploy the Go API to a VPS with HTTPS for production mobile apps</li>
          </ul>

          <Challenge number={9} title="Write a Launch Checklist">
            <p>Write a complete launch checklist for your fitness app with these sections: (1) API deployment, (2) iOS build with TestFlight testing, (3) Android build with APK for internal testing, (4) App Store and Play Store listing assets, (5) OTA update strategy for post-launch fixes.</p>
          </Challenge>

          <Challenge number={10} title="Ship It">
            <p>If you have developer accounts, run a real build with <Code>eas build --platform all</Code>. If not, run <Code>eas build --platform android --profile preview</Code> (free, no account required for APK) and install the APK on an Android device or emulator. Take a screenshot of your app running as a standalone build.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-mobile/notifications', label: 'Prev: Push Notifications' }}
            next={{ href: '/courses/grit-mobile', label: 'Back to All Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
